package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/awnumar/memguard"
	"github.com/google/go-tpm/tpm2"
	"golang.org/x/term"
)

// Get password safely (avoid Go GC stuff)
func InputPassword() []byte {
	fmt.Print("\nEnter password: ")
	bytePassword, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return bytePassword
}

func VerifyPassword(
	tpmDev *TPM,
	psk *tpm2.CreatePrimaryResponse, 
	bufAuth *memguard.LockedBuffer,
) bool {
	// Get blobs for test
	testPrivateBlob := LoadSign(GetSignPath(MAINMD5))
	testPublicBlob := LoadSign(GetSignPath(MAINMD5))
	vaultKey := tpmDev.UnsealObject(
		psk.ObjectHandle,
		testPrivateBlob,
		testPublicBlob,
		bufAuth.Bytes(),
	) // Will panic if wrong

	bufKey := memguard.NewBufferFromBytes(vaultKey)

	bufKey.Melt()
	bufKey.Scramble()
	bufKey.Destroy()

	return true
}

// Avoid corruption due to system crash
func WriteTmp(name string, data []byte) {
	basePath := GetSignPath("")
	tmp := filepath.Join(basePath, name + ".tmp")
	final := filepath.Join(basePath, name)
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		panic(fmt.Errorf("[File] write temp %s failed: %w", name, err))
	}
	if err := os.Rename(tmp, final); err != nil {
		panic(fmt.Errorf("[File] rename %s failed: %w", name, err))
	}
}

func GetSignPath(filename string) string {
	// exe, _ := os.Executable()
	// return filepath.Join(filepath.Dir(exe), filename)
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Join(pwd, "sign", filename)
}

func CreateNewDir(dirname string){
	basePath := GetSignPath("")
	dirPath := filepath.Join(basePath, dirname)
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		panic(err)
	}
	fmt.Println("[File] Create a new directory successfully.")
}

func CheckFileExists(filepath string) bool {
	fmt.Println(filepath)
	if _, err := os.Stat(filepath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	} else {
		return true
	}
}

func CheckFirstTime() bool {

	priBlobPath := GetSignPath(MAINMD5)
	pubBlobPath := GetSignPath(MAINMD5)
	saltPath := GetSignPath(SALTMD5)
	indPath := GetSignPath(INDMD5)
	douPath := GetSignPath(NORMD5)

	if CheckFileExists(priBlobPath) && CheckFileExists(pubBlobPath) && CheckFileExists(saltPath) && CheckFileExists(indPath) && CheckFileExists(douPath) {
		// Also check the legitimacy of the file
		// Destroy/Back up* the file too
		return false
	}
	return true
}
