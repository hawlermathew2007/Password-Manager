package security

import (

)

type Security struct {

}

func (sec *Security) DecryptData() {
	// Master password -> Argon2 -> AES 256 -> decrypt eval.json -> Success? -> Decrypt overview.yml (no pass in here)
}

func (sec *Security) EncryptData() {

}

func (sec *Security) DecryptPass() {
	// Since the Password will be encrypted twice
	// decrypt cri.yml (has id and pass)
}

