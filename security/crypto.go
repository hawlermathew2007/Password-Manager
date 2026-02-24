package security

import (
	"tools/data"
	"github.com/google/uuid"
)

type Security struct {
	DataManager 			*data.Manager
}

func (sec *Security) DecryptData() {
	// Master password -> Argon2 -> AES 256 -> decrypt eval.json -> Success? -> Decrypt overview.yml (no pass in here)
}

func (sec *Security) EncryptData() {

}

func (sec *Security) DecryptPass(accountID uuid.UUID) string {
	// Since the Password will be encrypted twice
	// decrypt cri.yml (has id and pass)

	// Decrypt here

	return sec.DataManager.LoadedCredsList[accountID]
}

