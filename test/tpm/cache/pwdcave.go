package cache

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"test/tools"

	"github.com/awnumar/memguard"
)

// Consider another case when user end the shell not the program

// Simple (not really correct way) check for index & normal vault id match
func isMatchIndNor(inds []tools.CredInd, nors []tools.CredSurf) bool {
	return len(inds) == len(nors) // Further check: id == id?
}

// Get Dou, Token by LockedBuffer
func GetRawCred[C tools.Credential](lb *memguard.LockedBuffer) tools.Vault[C] {
	var raw tools.Vault[C]
	if err := json.Unmarshal(lb.Bytes(), &raw); err != nil {
		lb.Destroy()
		panic(err)
	}
	defer func(){
		lb.Melt()
		lb.Scramble() // random overwrite before deallocate
		lb.Destroy()
	}()
	return raw 
}

// Get Pwd
func GetCredRest(tpmDev *tools.TPM, token []byte) *tools.CredRest {
	encodedToken := base64.StdEncoding.EncodeToString(token)
	pwdPrivateBlob := tools.LoadSign(tools.GetSignPath(encodedToken))
	pwdPublicBlob := tools.LoadSign(tools.GetSignPath(encodedToken))

	// Create Parent Storage Key
	psk, err := tpmDev.CreatePSK()
	if err != nil {
		panic(err)
	}
	defer tpmDev.Flush(psk.ObjectHandle)

	// Get pwd for authValue
	salt := tools.LoadSign(tools.GetSignPath(tools.SALTMD5))
	pwd := tools.InputPassword()
	authValue := tools.KDF(pwd, salt)
	bufAuth := memguard.NewBufferFromBytes(authValue)

	// Check if correct password
	userPwd := tpmDev.UnsealObject(psk.ObjectHandle, pwdPrivateBlob, pwdPublicBlob, bufAuth.Bytes()) // Will panic if wrong
	bufPwd := memguard.NewBufferFromBytes(userPwd)
	fmt.Printf("[TPM] Sealed Object Loaded and Unsealed. Successfully obtain the Vault Key for unlocking PWD.\n")

	return &tools.CredRest{
		Token: token,
		Password: bufPwd.Bytes(),
	}
}

// Change Pwd

// Change DoU

// Get Token (Blob Filename)

// Delete Credential (Del Password, DoU, UUID -> Token)

// Need to combine with GetNewCredential()
func AddCredential(
	tpmDev *tools.TPM,
	norBuf *memguard.LockedBuffer,
	indCe *CredentialEnclave,
	resCe *CredentialEnclave,
) (*memguard.LockedBuffer, *CredentialEnclave, *CredentialEnclave){

	creds := tools.GetNewCredential()

	// Decrypt those Vaults
	resLB, err := resCe.GetLockedBuffer()
	if err != nil {
		fmt.Errorf("[Cache] Unable to unlock CredRest enclave.")
		return nil, nil, nil
	}
	indLB, err := indCe.GetLockedBuffer()
	if err != nil {
		fmt.Errorf("[Cache] Unable to unlock CredSurf enclave.")
		return nil, nil, nil
	}
	DestroyLockedBuffer(indLB)
	DestroyLockedBuffer(resLB)

	// Get IndCe and norBuf Raw
	indPlain := GetRawCred[tools.CredInd](indLB)
	norPlain := GetRawCred[tools.CredSurf](norBuf)
	resPlain := GetRawCred[tools.CredRest](resLB)

	for _, cred := range creds {
    switch c := cred.(type) {
			case *tools.CredSurf:
				// Add to Vault[CredSurf].Credentials
				norPlain.Credentials = append(norPlain.Credentials, *c)

			case *tools.CredRest:
				// Add to Vault[CredRest].Credentials | Will be del after 1 minute
				resPlain.Credentials = append(resPlain.Credentials, *c)

				// Get pwd
				bufPwd := memguard.NewBufferFromBytes(c.Password)

				// Create Blobs in blobs dir
				psk, err := tpmDev.CreatePSK()
				if err != nil {
					panic(err)
				}
				defer tpmDev.Flush(psk.ObjectHandle)

				// Get pwd for authValue
				salt := tools.LoadSign(tools.GetSignPath(tools.SALTMD5))
				pwd := tools.InputPassword()
				authValue := tools.KDF(pwd, salt)
				bufAuth := memguard.NewBufferFromBytes(authValue)

				if tools.VerifyPassword(tpmDev, psk, bufAuth) {
					// Create SK Object
					privateBlob, publicBlob, err := tpmDev.CreateSealedObj(
						psk.ObjectHandle,
						bufPwd.Bytes(),
						bufAuth.Bytes(), // Prompt for Pwd
					)
					if err != nil {
						panic(fmt.Errorf("[TPMBlob] Unable to create Private and Public Blob for Pwd. Err: %w", err))
					}
					tools.WriteTmp(tools.MAINMD5, publicBlob)
					tools.WriteTmp(tools.MAINMD5, privateBlob)
				}

				DestroyLockedBuffer(bufPwd)
				DestroyLockedBuffer(bufAuth)

			case *tools.CredInd:
				// Add to Vault[CredInd].Credentials
				indPlain.Credentials = append(indPlain.Credentials, *c)
    }
	}

	// Encrypted once again
	jsonNorVault, err := json.Marshal(norPlain)
	if err != nil {
		memguard.SafePanic(err)
	}
	jsonIndVault, err := json.Marshal(indPlain)
	if err != nil {
		memguard.SafePanic(err)
	}
	jsonResVault, err := json.Marshal(resPlain)
	if err != nil {
		memguard.SafePanic(err)
	}
	tools.WipeVault(&norPlain)
	tools.WipeVault(&indPlain)
	tools.WipeVault(&resPlain)

	newNorBuf := memguard.NewBufferFromBytes(jsonNorVault)
	if err != nil {
		memguard.SafePanic(err)
	}
	newIndCe, err := NewCredentialEnclave(jsonIndVault, false) // need lil math here for = start time - timepass & also time for this = Session time
	if err != nil {
		memguard.SafePanic(err)
	}
	newResCe, err := NewCredentialEnclave(jsonResVault, true) // need lil math here for = start time - timepass
	if err != nil {
		memguard.SafePanic(err)
	}
	return newNorBuf, newIndCe, newResCe
}

// WithCredential decrypts the enclave into a LockedBuffer, calls fn with the parsed rawCredential, then Scrambles and Destroys the buffer. The LockedBuffer (used mrprotect) is immutable on Open(); we Melt() it only if we needed to modify it (we don't here — just reading). IMPORTANT: fn must NOT retain any pointer into the LockedBuffer after it returns. The memory will be wiped immediately after fn exits. Abandon soon
func (ce *CredentialEnclave) DoSthWithPwd(fn func(password *memguard.LockedBuffer)) error {
	ce.mu.Lock()
	if ce.enclave == nil {
		ce.mu.Unlock()
		return fmt.Errorf("Credential has been wiped (TTL expired or WipeNow called)")
	}
	enc := ce.enclave
	ce.mu.Unlock()

	// Decrypt into a guarded LockedBuffer (immutable by default after Open)
	lb, err := enc.Open()
	if err != nil {
		return fmt.Errorf("open enclave: %w", err)
	}
	defer func() {
		lb.Melt()
		lb.Scramble() // random overwrite before deallocate
		lb.Destroy()
	}()

	// Parse from the locked region.
	// json.Unmarshal will allocate a rawCredential on the heap — we wipe it
	// immediately after extracting what we need.
	var cred tools.CredRest
	if err := json.Unmarshal(lb.Bytes(), &cred); err != nil {
		return fmt.Errorf("unmarshal password: %w", err)
	}

	// Call the user function with copies of the strings, then wipe.
	// We make explicit copies so the caller's local variables are distinct
	// from the rawCredential fields (which share the same backing arrays).
	password := memguard.NewBuffer(len(cred.Password))
	password.Move([]byte(cred.Password)) // copies + wipes the []byte cast

	// Wipe rawCredential fields before GC
	tools.WipeRawCredential(&cred)

	fn(password)

	// Wipe the local string copies (their backing arrays live on Go heap)
	password.Scramble()
	password.Destroy()

	return nil
}

// Print all Credentials as requested
func PrintCredentials(tpmDev *tools.TPM, norBuf *memguard.LockedBuffer, indCe *CredentialEnclave) {

	// Get IndCe LockedBuffer
	indLB, err := indCe.GetLockedBuffer()
	if err != nil {
		fmt.Errorf("[Cache] Unable to unlock enclave.")
		return
	}
	DestroyLockedBuffer(indLB)

	// Get IndCe and norBuf Raw
	indPlain := GetRawCred[tools.CredInd](indLB)
	norPlain := GetRawCred[tools.CredSurf](norBuf)

	// Get Pwd based on UUID from CredInd => Make AddCredentials & PwdCaVe
	// Merge DoU with Pwd and Print

	if len(norPlain.Credentials) == 0 {
		fmt.Println("[Vault] The vault is empty.")
		return
	}

	if !isMatchIndNor(indPlain.Credentials, norPlain.Credentials) {
		panic("The IDs in both Index and Normal Vault does not match.")
	}

	for i:=0;i < len(norPlain.Credentials);i++ {
	   nor := norPlain.Credentials[i]
	   ind := indPlain.Credentials[i]

		 token := ind.Token

		 // Make AddCredential() & GetCredRest() first and go back here
		 // Token -> Find blobs -> Unseal blobs -> Get Pwd
		 pwd := GetCredRest(tpmDev, token)

	   merged := map[string]any{
	       "id":       nor.ID,
	       "domain":   nor.Domain,
	       "username": nor.Username,
	       "password": pwd.Password,
	   }

	   out, _ := json.MarshalIndent(merged, "", "  ")
	   fmt.Println(string(out))
	}
}
