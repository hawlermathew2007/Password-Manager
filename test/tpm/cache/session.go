package cache

import (
	"fmt"
	"time"
)

// Session for TPM - Just TTL watcher ig
// Autolock
// Rules:
// - No LockedBuffer or Buffer with Cred exists and No Enclaves (Scramble, Destroy, Purge, WipeBytes, ScrambleBytes)
// - When on app is still opened, the idle Cred is still in TPM-based cache. If close, no TPM-based cache.

const (
	SESSIONTTL = 5 * time.Minute
	PWDTTL = 1 * time.Minute
)

// Waiting for Session to end 
func SessionWatcher(){
	select{
		case <- time.After(30 * time.Second):
			return
	}
}

// Capturing Signal for wiping the single enclave
func (ce *CredentialEnclave) SigWatcher() {
	select {
		case <- ce.WipedSig:
			fmt.Println("[Enclave] The enclave is manually wiped.")
			return // Del if want to continue receiving signal
	}
}

// RemainingTTL returns how long until the credential is auto-wiped.
func (ce *CredentialEnclave) RemainingTTL() time.Duration {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	if ce.enclave == nil {
		return 0
	}
	rem := ce.ttl - time.Since(ce.loadedAt)
	if rem < 0 {
		return 0
	}
	return rem
}
