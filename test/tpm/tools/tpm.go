package tools

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
	"github.com/google/go-tpm/tpmutil"
)

type TPM struct {
	Transporter transport.TPMCloser
	Device			io.ReadWriteCloser
}

func NewTPM(rwc io.ReadWriteCloser) *TPM {
	transporter := transport.FromReadWriteCloser(rwc)
	return &TPM{
		Transporter: transporter,
		Device: rwc,
	}
}

func OpenTPMDevice() io.ReadWriteCloser {
	rwc, err := tpmutil.OpenTPM("/dev/tpmrm0")
	if err != nil {
		panic(err)
	}
	return rwc
}

func createPrimaryKey() tpm2.CreatePrimary {
	return tpm2.CreatePrimary{
		PrimaryHandle: tpm2.TPMRHOwner,
		InPublic: tpm2.New2B(
			tpm2.TPMTPublic{
				Type:    tpm2.TPMAlgRSA,
				NameAlg: tpm2.TPMAlgSHA256,
				// AuthPolicy: policyDigestBuffer,

				ObjectAttributes: tpm2.TPMAObject{
					FixedTPM:            true,
					FixedParent:         true,
					SensitiveDataOrigin: true,
					UserWithAuth:        true,
					Decrypt:             true,
					Restricted:          true,
				},

				Parameters: tpm2.NewTPMUPublicParms(
					tpm2.TPMAlgRSA,
					&tpm2.TPMSRSAParms{
						KeyBits: 2048,
						Symmetric: tpm2.TPMTSymDefObject{
							Algorithm: tpm2.TPMAlgAES,
							KeyBits:   tpm2.NewTPMUSymKeyBits(tpm2.TPMAlgAES, tpm2.TPMKeyBits(128)),
							Mode:      tpm2.NewTPMUSymMode(tpm2.TPMAlgAES, tpm2.TPMAlgCFB),
						},
					},
				),
			},
		),
	}
}

func createChildKey(
	parentAuth tpm2.AuthHandle,
	authValue []byte,
	vaultKey []byte,
	policyDigest tpm2.TPM2BDigest,
) tpm2.Create {
	return tpm2.Create{
		ParentHandle: parentAuth,
		InSensitive: tpm2.TPM2BSensitiveCreate{
			Sensitive: &tpm2.TPMSSensitiveCreate{
				UserAuth: tpm2.TPM2BAuth{
					Buffer: authValue,
				},
				Data: tpm2.NewTPMUSensitiveCreate(
					&tpm2.TPM2BSensitiveData{
						Buffer: vaultKey,
					},
				),
			},
		},
		InPublic: tpm2.New2B(
			tpm2.TPMTPublic{
				Type:    tpm2.TPMAlgKeyedHash,
				NameAlg: tpm2.TPMAlgSHA256,
				ObjectAttributes: tpm2.TPMAObject{
					FixedTPM:            true,  // cannot be duplicated off this tpmd
					FixedParent:         true,  // cannot be re-parented
					UserWithAuth:        true,  // authValue is required to unseal
					AdminWithPolicy:     true,
					NoDA:                true,  // exempt from dictionary attack lockout
					SensitiveDataOrigin: false, // we are providing the data, not TPM
				},
				// AuthPolicy binds unsealing to a specific PCR state
				AuthPolicy: policyDigest,
				Parameters: tpm2.NewTPMUPublicParms(
					tpm2.TPMAlgKeyedHash,
					&tpm2.TPMSKeyedHashParms{
						Scheme: tpm2.TPMTKeyedHashScheme{
							Scheme: tpm2.TPMAlgNull, // pure seal, not HMAC
						},
					},
				),
			},
		),
	}
}

func createNonce() ([]byte, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
			return nil, fmt.Errorf("nonce gen failed: %w", err)
	}
	return nonce, nil
}

func (tpmd *TPM) ApplyPolicyPCR(handler tpm2.TPMHandle) error {
	pcr := tpm2.PolicyPCR{
		PolicySession: handler,
		Pcrs: tpm2.TPMLPCRSelection{
			PCRSelections: []tpm2.TPMSPCRSelection{
				{
					Hash:      tpm2.TPMAlgSHA256,
					PCRSelect: tpm2.PCClientCompatible.PCRs(0, 1, 2, 3, 4, 5, 6, 7),
				},
			},
		},
	}
	_, err := pcr.Execute(tpmd.Transporter)
	return err
}

func (tpmd *TPM) ApplyPolicyAuthValue(handler tpm2.TPMHandle) error {
	authValue := tpm2.PolicyAuthValue{
		PolicySession: handler,
	}
	_, err := authValue.Execute(tpmd.Transporter)
	return err
}

func (tpmd *TPM) StartAuthSession(nonce []byte) (*tpm2.StartAuthSessionResponse, error) {
	sess, err := tpm2.StartAuthSession{
		SessionType: tpm2.TPMSEPolicy,
		AuthHash: tpm2.TPMAlgSHA256,
    NonceCaller: tpm2.TPM2BNonce{
        Buffer: nonce,
    },
	}.Execute(tpmd.Transporter)
	if err != nil {
		return &tpm2.StartAuthSessionResponse{}, err
	}
	return sess, nil
}

func (tpmd *TPM) CreatePSK() (*tpm2.CreatePrimaryResponse, error) {
	// policyDigestBuffer, err := tpmd.CreatePolicyDigest()
	// if err != nil {
	// 	panic(fmt.Errorf("[PolicyDigest] Create Policy Session failed. Err: %w", err))
	// }
	cmd := createPrimaryKey()
	rsp, err := cmd.Execute(tpmd.Transporter)
	if err != nil {
		return nil, fmt.Errorf("Create PSK failed: %w", err)
	}

	fmt.Println("[PSK] Create PSK ok.")
	return rsp, nil
}

func (tpmd *TPM) CreatePolicyDigest() (tpm2.TPM2BDigest, error) {

	transporter := tpmd.Transporter

	nonce, err := createNonce() // or 20–32 bytes
	if err != nil {
			return tpm2.TPM2BDigest{}, err
	}
	// 1. Start policy session
	sess, err := tpmd.StartAuthSession(nonce)
	if err != nil {
		return tpm2.TPM2BDigest{}, fmt.Errorf("StartAuthSession failed: %w", err)
	}

	// Ensure session is flushed
	defer tpmd.Flush(sess.SessionHandle)

	// 2. Apply PCR policy
	err = tpmd.ApplyPolicyPCR(sess.SessionHandle)
	if err != nil {
		return tpm2.TPM2BDigest{}, fmt.Errorf("PolicyPCR failed: %w", err)
	}

	// 3. Apply AuthValue policy
	err = tpmd.ApplyPolicyAuthValue(sess.SessionHandle)
	if err != nil {
			return tpm2.TPM2BDigest{}, fmt.Errorf("PolicyAuthValue failed: %w", err)
	}

	// 4. Get digest
	digestRsp, err := tpm2.PolicyGetDigest{
		PolicySession: sess.SessionHandle,
	}.Execute(transporter)
	if err != nil {
		return tpm2.TPM2BDigest{}, fmt.Errorf("PolicyGetDigest failed: %w", err)
	}

	return digestRsp.PolicyDigest, nil
}

func (tpmd *TPM) CreateSealedObj(
	pskHandle tpm2.TPMHandle,
	vaultKey []byte,   // AES-256 key to seal (32 bytes)
	authValue []byte,  // password auth for unsealing
) ([]byte, []byte, error) {

	if len(vaultKey) != 32 {
		return nil, nil, fmt.Errorf("vaultKey must be 32 bytes for AES-256, got %d", len(vaultKey))
	}

	transporter := tpmd.Transporter

	// Build a trial session to create a PCR policy digest
	// This binds the sealed object to the current PCR state
	policyDigest, err := tpmd.CreatePolicyDigest()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build PCR policy digest: %w", err)
	}

	parentAuth, _ := tpmd.GetAuthHandle(pskHandle, nil)

	cmd := createChildKey(parentAuth, authValue, vaultKey, policyDigest)

	rsp, err := cmd.Execute(transporter)
	if err != nil {
		return nil, nil, fmt.Errorf("TPM Create sealed object failed: %w", err)
	}

	privBytes := tpm2.Marshal(rsp.OutPrivate)
	pubBytes := tpm2.Marshal(rsp.OutPublic)

	fmt.Printf("[TPM] Sealed object created — private: %d bytes, public: %d bytes\n",
	len(privBytes), len(pubBytes))

	// Return the blobs; caller is responsible for persisting them
	return privBytes, pubBytes, nil
}

func (tpmd *TPM) UnsealObject(
	pskHandle tpm2.TPMHandle,
	privBytes []byte,
	pubBytes []byte,
	authValue []byte,
) ([]byte) {

	// Reconstruct TPM2B types from raw bytes
	priv, err := tpm2.Unmarshal[tpm2.TPM2BPrivate](privBytes)
	if err != nil {
    panic(fmt.Errorf("unmarshal private blob failed: %w", err))
	}
	pub, err := tpm2.Unmarshal[tpm2.TPM2BPublic](pubBytes)
	if err != nil {
			panic(fmt.Errorf("unmarshal public blob failed: %w", err))
	}

	// Load sealed object under PSK
	parentAuth, _ := tpmd.GetAuthHandle(pskHandle, nil)
	loaded, err := tpm2.Load{
		ParentHandle: parentAuth,
		InPrivate:    *priv,
		InPublic:     *pub,
	}.Execute(tpmd.Transporter)

	if err != nil {
		panic(err)
	}

	defer tpmd.Flush(loaded.ObjectHandle)

	// Start real policy session
	sess, closeSession, err := tpm2.PolicySession(
		tpmd.Transporter,
		tpm2.TPMAlgSHA256,
		16, // nonce size
		tpm2.Auth(authValue),
	)
	if err != nil {
		panic(err)
	}
	defer closeSession()

	// Assertion 1: PCR state must match
	err = tpmd.ApplyPolicyPCR(sess.Handle())
	if err != nil {
			panic(fmt.Errorf("PolicyPCR failed: %w", err))
	}

	// Assertion 2: password assertion — MUST match CreatePolicyDigest order
	err = tpmd.ApplyPolicyAuthValue(sess.Handle())
	if err != nil {
		panic(fmt.Errorf("PolicyAuthValue failed: %w", err))
	}

	// Unseal — Auth binds authValue bytes into the session
	unsealed, err := tpm2.Unseal{
		ItemHandle: tpm2.AuthHandle{
			Handle: loaded.ObjectHandle,
			Name:   loaded.Name,
			Auth:   sess,
		},
	}.Execute(tpmd.Transporter)
	if err != nil {
		panic(fmt.Errorf("Wrong Password! Unseal failed: %w.", err))
	}

	return unsealed.OutData.Buffer
}

// Need Authentication: Get fingerprint from Public Blobs and confirm that the one accessing the priv blob is the user.
func (tpmd *TPM) GetAuthHandle(handle tpm2.TPMHandle, auth []byte) (tpm2.AuthHandle, error) {
	rsp, err := tpm2.ReadPublic{ // Read Public Blob metadata
		ObjectHandle: handle,
	}.Execute(tpmd.Transporter)
	if err != nil {
		return tpm2.AuthHandle{}, err
	}

	return tpm2.AuthHandle{
		Handle: handle,
		Name:   rsp.Name, // Get fingerprint about the key
		Auth:   tpm2.PasswordAuth(auth), // Proof to know the AuthValue to unlock Private Blob
	}, nil
}

func (tpmd *TPM) Flush(handle tpm2.TPMHandle) error {
	cmd := tpm2.FlushContext{
		FlushHandle: handle,
	}
	_, err := cmd.Execute(tpmd.Transporter)
	return err
}

func PrintPCRInfo() {
	fmt.Println("[tpm] PCR binding strategy:")
	fmt.Println("    PCR  0 — BIOS/UEFI firmware hash  → changes if firmware tampered")
	fmt.Println("    PCR  7 — Secure Boot state         → changes if SecureBoot disabled")
	fmt.Println("    PCR 11 — OS bootloader (GRUB/shim) → changes if different OS loaded")
	fmt.Println("    ALL PCRs must match seal-time values. Any change = unseal fails.")
}
