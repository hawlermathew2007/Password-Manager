package action

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"test/cache"
	"test/tools"
	"time"

	"github.com/awnumar/memguard"
)

// vault.enc.backup may need to be created during using (if panic -> keep)
// For every panic it should automatically save the changes in enclave and load to .enc
func AccessTCROSS() {
	fmt.Println("============================")
	fmt.Println("Welcome to TCROSS login CLI.")
	fmt.Println("============================")

	memguard.CatchInterrupt()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGHUP)
	signal.Notify(sigCh, syscall.SIGINT)

	// Harden memory
	tools.HardenProcess()
	
	// Input Master Password
	pwd := tools.InputPassword()
	fmt.Println(pwd)

	// Load Salt & Blobs
	salt := tools.LoadSign(tools.GetSignPath(tools.SALTMD5))
	privateBlob := tools.LoadSign(tools.GetSignPath(tools.MAINMD5))
	publicBlob := tools.LoadSign(tools.GetSignPath(tools.MAINMD5))
	fmt.Println("[Load] Loading: salt & sealed Blob.")

	// Get authValue
	authValue := tools.KDF(pwd, salt)
	fmt.Printf("[KDF] AuthValue created.%X\n", authValue)

	// Connect with TPM chip
	rwc := tools.OpenTPMDevice()
	defer rwc.Close()

	tpmDev := tools.NewTPM(rwc)
	fmt.Println("[TPM] Connected to TPM.")

	// Create Parent Storage Key
	psk, err := tpmDev.CreatePSK()
	if err != nil {
		panic(err)
	}
	defer tpmDev.Flush(psk.ObjectHandle)

	// Load Sealed Object
	vaultKey := tpmDev.UnsealObject(
		psk.ObjectHandle,
		privateBlob,
		publicBlob,
		authValue,
	)
	buf := memguard.NewBufferFromBytes(vaultKey)
	fmt.Printf("[TPM] Sealed Object Loaded and Unsealed. Successfully obtain the Vault Key.\n")

	// Decrypt NorVault
	NorEv, err := tools.LoadVaultFile(tools.GetSignPath(tools.NORMD5))
	if err != nil {
		panic(err)
	}
	NorVault, err := tools.Decrypt[tools.CredSurf](NorEv, buf.Bytes())
	if err != nil {
		panic(fmt.Errorf("vault decryption failed: %w", err))
	}
	indEv, err := tools.LoadVaultFile(tools.GetSignPath(tools.INDMD5))
	if err != nil {
		panic(err)
	}
	IndVault, err := tools.Decrypt[tools.CredInd](indEv, buf.Bytes())
	if err != nil {
		panic(fmt.Errorf("vault decryption failed: %w", err))
	}
	fmt.Printf("[Vault] All vaults unlocked. %d credentials loaded.\n", len(NorVault.Credentials))

	// Wipe vaultKey
	cache.DestroyLockedBuffer(buf)

	// Create empty pwd vault
	var ResVault tools.Vault[tools.CredRest]

	// Updates & Re-encrypt Vault
	esc := 0
	jsonNorVault, err := json.Marshal(NorVault)
	if err != nil {
		memguard.SafePanic(err)
	}
	jsonIndVault, err := json.Marshal(IndVault)
	if err != nil {
		memguard.SafePanic(err)
	}
	jsonResVault, err := json.Marshal(ResVault)
	if err != nil {
		memguard.SafePanic(err)
	}
	tools.WipeVault(NorVault)
	tools.WipeVault(IndVault)

	norBuf := memguard.NewBufferFromBytes(jsonNorVault)
	if err != nil {
		memguard.SafePanic(err)
	}
	indCe, err := cache.NewCredentialEnclave(jsonIndVault, false)
	if err != nil {
		memguard.SafePanic(err)
	}
	resCe, err := cache.NewCredentialEnclave(jsonResVault, false)
	if err != nil {
		memguard.SafePanic(err)
	}

	tools.PrintVaultInstruction()
	for esc != 1 {
		opt := tools.TakeVaultOptions()
		switch opt {
			case 0:
				cache.PrintCredentials(tpmDev, norBuf, indCe)
			case 1:
				norBuf, indCe, resCe = cache.AddCredential(tpmDev, norBuf, indCe, resCe)
			case 2:
				fmt.Println("[Nah] Coming soon :)")
				fmt.Println("[Nah] Or never. Because I am lazy.")
			case 3:
				fmt.Println("[Nah] Coming soon :)")
				fmt.Println("[Nah] Or never. Because I am lazy.")
			case 4:
				fmt.Println("[Nah] Coming soon :)")
				fmt.Println("[Nah] Or never. Because I am lazy.")
			case 10:
				esc = 1
			default:
				fmt.Println("[Error] Invalid options.")
		}
	}

	newVaultKey := tools.GenerateKey()

	// Create new blobs for new vault key
	newPrivateBlob, newPublicBlob, err := tpmDev.CreateSealedObj(
		psk.ObjectHandle,
		newVaultKey,
		authValue,
	)
	if err != nil {
		panic(fmt.Errorf("[TPMBlob] Unable to create Private and Public Blob for Object. Err: %w", err))
	}

	// Decrypt Index Vault Enclave into a guarded LockedBuffer (immutable by default after Open)
	indLB, err := indCe.GetLockedBuffer()
	if err != nil {
		fmt.Errorf("[Cache] Unable to unlock enclave.")
		return
	}

	// Decrypt to Plain struct
	var indPlain tools.Vault[tools.CredInd]
	if err := json.Unmarshal(indLB.Bytes(), &indPlain); err != nil {
		fmt.Errorf("Unmarshal ind failed.")
		return
	}
	cache.DestroyLockedBuffer(indLB)
	// defer tools.WipeVault(&indPlain)

	var norPlain tools.Vault[tools.CredSurf]
	if err := json.Unmarshal(norBuf.Bytes(), &norPlain); err != nil {
		norBuf.Destroy()
		panic(err)
	}

	// Encrypt the updated vault
	newEvNor, err := tools.Encrypt(&norPlain, newVaultKey)
	if err != nil {
		panic(fmt.Errorf("encrypt vault failed: %w", err))
	}
	newEvInd, err := tools.Encrypt(&indPlain, newVaultKey)
	if err != nil {
		panic(fmt.Errorf("encrypt vault failed: %w", err))
	}
	
	// Stored Sign
	tools.WriteTmp(tools.MAINMD5, newPublicBlob)
	tools.WriteTmp(tools.MAINMD5, newPrivateBlob)
	if err := tools.SaveVaultFile(newEvNor, tools.GetSignPath(tools.NORMD5)); err !=nil {
		panic(err)
	}
	if err := tools.SaveVaultFile(newEvInd, tools.GetSignPath(tools.INDMD5)); err !=nil {
		panic(err)
	}

	select {
		case sig := <-sigCh:
			fmt.Printf("\n[signal] Received %v — wiping immediately\n", sig)
		case <- time.After(30 * time.Second):
			fmt.Println("[signal] Session end. Exiting Program")
	}

	fmt.Printf("[Vault] Vault saved. %d credentials total.\n", len(norPlain.Credentials))
	fmt.Println("[Done] Accessing complete.")

	memguard.Purge()
	memguard.SafeExit(0)
}
