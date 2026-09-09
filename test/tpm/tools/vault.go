package tools

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"github.com/google/uuid"
	"golang.org/x/term"
)

type CredType uint8

const (
	INDMD5 = "9e29de6e364610eed3cc30630e27a541"
	NORMD5 = "b97bb17879798f27b4d799af06e2028a"
	MAINMD5 = "6b6a517a78bce0c4d59a41cf90f15bac"
	SALTMD5 = "daf3ae813ad102e1014c6a5cb9b3a7fb"
)

const (
	TypeInd CredType = iota
	TypeNorm
	TypeSensi
)

type Credential interface {
	GetType() CredType
}

type CredSurf struct {
	ID			 uuid.UUID	`json:"id"` 
	Domain   []byte  		`json:"domain"` // Should be in bytes somehow
	Username []byte 		`json:"username,omitempty"`
}

type CredRest struct {
	Token		 []byte 		`json:"token"` 
	Password []byte 		`json:"password"`
}

type CredInd struct {
	ID			 uuid.UUID	`json:"id"` 
	Token		 []byte 		`json:"token"` 
}
 
type Vault[C Credential] struct {
	Credentials []C `json:"credentials"`
}

type EncryptedVault struct {
	IV         string `json:"iv"`         // base64-encoded 12-byte nonce
	Ciphertext string `json:"ciphertext"` // base64-encoded AES-GCM output
}

func (c CredSurf) GetType() CredType {
	return TypeNorm
}

func (c CredRest) GetType() CredType {
	return TypeSensi
}

func (c CredInd) GetType() CredType {
	return TypeInd
}

func LoadSign(filepath string) []byte {
	file, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return file
}

func GenerateKey() []byte {
	vaultKey := make([]byte, 32)
	if _, err := rand.Read(vaultKey); err != nil {
		panic(err)
	}
	return vaultKey
}

func PrintVaultInstruction() {
	fmt.Println("\nInstruction] Here are the options of how you can interact with the vault.")
	fmt.Println("[0] Print Credentials")
	fmt.Println("[1] Add 1 new Credential")
	fmt.Println("[2] Update a Credential")
	fmt.Println("[3] Delete a Credential")
	fmt.Println("[4] Update AuthValue (Password)")
	fmt.Println("[10] Quit")
}

func TakeVaultOptions() int {
	var opt int
	fmt.Print("\nYour option (0-4): ")
	fmt.Scanln(&opt)
	fmt.Println()
	return opt
}

func GetNewCredential() []Credential {
	var domain []byte
	var username []byte

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Domain: ")
	domain, _ = reader.ReadBytes('\n')
	domain = bytes.TrimSpace(domain)

	fmt.Print("Username: ")
	username, _ = reader.ReadBytes('\n')
	username = bytes.TrimSpace(username)

	fmt.Print("Password: ")
	pwd, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	userID := uuid.New()
	token := GenerateRandomBytes(16)

	return []Credential{
		&CredSurf{
			ID: userID,
			Domain:   domain,
			Username: username,
		},
		&CredRest{
			Token: token,
			Password: pwd,
		},
		&CredInd{
			ID: userID,
			Token: token,
		},
	}
}

func CreateVault[C Credential]()[]byte {
	emptyVault := Vault[C]{Credentials: []C{}}
	plainJSON, _ := json.MarshalIndent(emptyVault, "", "  ")
	return plainJSON
}

func UpdateCredential() {

}

func DeleteCredential() {

}

func Encrypt[C Credential](v *Vault[C], key []byte) (*EncryptedVault, error) {
	plaintext, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal vault: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	iv := make([]byte, gcm.NonceSize()) // 12 bytes
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("generate iv: %w", err)
	}

	ciphertext := gcm.Seal(nil, iv, plaintext, nil)

	return &EncryptedVault{
		IV:         base64.StdEncoding.EncodeToString(iv),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}, nil
}

// Decrypt decodes and decrypts an EncryptedVault back into a Vault.
func Decrypt[C Credential](ev *EncryptedVault, key []byte) (*Vault[C], error) {
	iv, err := base64.StdEncoding.DecodeString(ev.IV)
	if err != nil {
		return nil, fmt.Errorf("decode iv: %w", err)
	}
	ciphertext, err := base64.StdEncoding.DecodeString(ev.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt (wrong key or corrupted data): %w", err)
	}

	var v Vault[C]
	if err := json.Unmarshal(plaintext, &v); err != nil {
		return nil, fmt.Errorf("unmarshal vault: %w", err)
	}
	return &v, nil
}

func SaveVaultFile(ev *EncryptedVault, vaultFile string) error {
	data, err := json.MarshalIndent(ev, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal encrypted vault: %w", err)
	}
	WriteTmp(vaultFile, data)
	fmt.Printf("[vault] saved to %q\n", vaultFile)
	return nil
}

func LoadVaultFile(vaultFile string) (*EncryptedVault, error) {
	data, err := os.ReadFile(vaultFile)
	if err != nil {
		return nil, fmt.Errorf("read vault file: %w", err)
	}
	var ev EncryptedVault
	if err := json.Unmarshal(data, &ev); err != nil {
		return nil, fmt.Errorf("parse vault file: %w", err)
	}
	return &ev, nil
}
