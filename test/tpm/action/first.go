package action

import (
	"fmt"
	"test/tools"
)

func FirstTCROSS(){
	fmt.Println("==============================")
	fmt.Println("Welcome to TCROSS first login.")
	fmt.Println("==============================")

	if tools.CheckFirstTime() {
		
		fmt.Println("\n[First] Registering...")

		// Input Master Password
		pwd := tools.InputPassword()
		fmt.Println(pwd)

		// KDF
		salt, err := tools.GenerateSalt()
		if err != nil {
			panic(err)
		}
		authValue := tools.KDF(pwd, salt)
	
		// Generate Vault Key
		vaultKey := tools.GenerateKey()

		// Create Empty Vault
		emptyNormVault := &tools.Vault[tools.CredSurf]{}
		emptySensiVault := &tools.Vault[tools.CredRest]{}

		norEv, err := tools.Encrypt(emptyNormVault, vaultKey)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[NorVault] IV:         %s\n", norEv.IV)
		fmt.Printf("[NorVault] Ciphertext: %s...\n", norEv.Ciphertext[:32])

		indEv, err := tools.Encrypt(emptySensiVault, vaultKey)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[IndVault] IV:         %s\n",indEv.IV)
		fmt.Printf("[IndVault] Ciphertext: %s...\n", indEv.Ciphertext[:32])

		// Connect with TPM chip
		rwc := tools.OpenTPMDevice()
		defer rwc.Close()

		tpmDev := tools.NewTPM(rwc)

		// Create Parent Storage Key
		psk, err := tpmDev.CreatePSK()
		if err != nil {
			panic(err)
		}
		defer tpmDev.Flush(psk.ObjectHandle)

		// Create SK Object
		privateBlob, publicBlob, err := tpmDev.CreateSealedObj(
			psk.ObjectHandle,
			vaultKey,
			authValue,
		)
		if err != nil {
			panic(fmt.Errorf("[TPMBlob] Unable to create Private and Public Blob for Object. Err: %w", err))
		}

		// Stored Sign
		tools.CreateNewDir("blobs")
		tools.WriteTmp(tools.MAINMD5, publicBlob)
		tools.WriteTmp(tools.MAINMD5, privateBlob)
		tools.WriteTmp(tools.SALTMD5, salt)
		if err := tools.SaveVaultFile(norEv, tools.NORMD5); err !=nil {
			panic(err)
		}
		if err := tools.SaveVaultFile(indEv, tools.INDMD5); err !=nil {
			panic(err)
		}

	fmt.Println("[Done] Enrollment complete")

	} else {
		fmt.Println("[Done] This is not your first time here.")
	}
}
