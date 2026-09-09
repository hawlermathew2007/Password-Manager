package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"test/tools"
	// "runtime"
	"github.com/awnumar/memguard"
)

// Refined PwdCaVe caching method with Session Key
// Memguard
// - Read & Write Memory Protectiom
// - Kernel-level Immutability
// - Secure Termination
// Memory Hardening in memsec.go

// Hack my own sys
// Process Isolation via Privilege Separation

// func CheckingHardware -> Decision 0 or 1
// 0 : Hardware not compatible => Software Enclaves
// 1 : Hard compatible => Hardware Enclaves (TEE/EPK & ORAM) (Maybe never... Bcuz im poor and have none of those hardware)

// Rules for sensi buf
// - Mutual Exclustion needed
// - Must be wiped after finish using (WipeBytes, to nil, Scamble & Destroy)

type CredentialEnclave struct {
	typeCred  tools.CredInd
	mu        sync.Mutex // Make it only accessible with TCROSS
	enclave   *memguard.Enclave // encrypted ciphertext at rest
	loadedAt  time.Time
	ttl       time.Duration
	timer 	  *time.Timer
	WipedSig  chan struct{} // closed when the enclave has been wiped
}

// IsAlive reports whether the credential is still in memory.
func (ce *CredentialEnclave) IsAlive() bool {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	return ce.enclave != nil
}

// Passing Plain JSON to Software Enclaves and Start TTL for the Enclaves.
func NewCredentialEnclave(jsonBytes []byte, isPwd bool) (*CredentialEnclave, error) {
	var check tools.CredRest
	if err := json.Unmarshal(jsonBytes, &check); err != nil {
		memguard.WipeBytes(jsonBytes)
		return nil, fmt.Errorf("invalid credential JSON: %w", err)
	}
	tools.WipeRawCredential(&check)

	// NewBufferFromBytes copies jsonBytes into guarded memory AND wipes
	// jsonByteskatomically. Single operation, no window between copy and wipe.
	buf := memguard.NewBufferFromBytes(jsonBytes)
	// jsonBytes is now zeroed by the call above — do not use again
	jsonBytes = nil

	if !buf.IsAlive() {
		return nil, fmt.Errorf("failed to allocate guarded buffer")
	}

	enc := buf.Seal() // Help produce an Enclave
	buf = nil
	if enc == nil {
		return nil, fmt.Errorf("seal failed")
	}

	ttl := SESSIONTTL
	if isPwd {
		ttl = PWDTTL
	}

	ce := &CredentialEnclave{
		enclave:  enc,
		loadedAt: time.Now(),
		ttl:      ttl,
		WipedSig: make(chan struct{}),
	}

	ce.timer = time.AfterFunc(ttl, func() {
		fmt.Println("\n[Enclave] TTL expired — wiping credential from memory")
		ce.Wipe()
		if ce.IsAlive() {
			fmt.Println("The system lied")
		} else {
			fmt.Println("It is dead for real")
		}
		ce.WipedSig <- struct{}{}
	})

	return ce, nil
}

// Use enclave to output the Locked buffer
func (ce *CredentialEnclave) GetLockedBuffer() (*memguard.LockedBuffer, error) {
	ce.mu.Lock()
	if ce.enclave == nil {
		ce.mu.Unlock()
		return nil, fmt.Errorf("[Cache] No enclave bruh.")
	}
	enc := ce.enclave
	ce.mu.Unlock()

	// Decrypt into a guarded LockedBuffer (immutable by default after Open)
	lb, err := enc.Open()
	if err != nil {
		return nil, fmt.Errorf("[Cache] Unable to unlock enclave.")
	}
	return lb, nil
}

func DestroyLockedBuffer(lb *memguard.LockedBuffer) {
	lb.Melt()
	lb.Scramble() // random overwrite before deallocate
	lb.Destroy()
}

// Wipe destroys the enclave immediately.
func (ce *CredentialEnclave) Wipe() {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if ce.enclave == nil {
		return // already wiped
	}

	// Open => LockedBuffer so we can Scramble before final destroy.
	// This ensures even the encrypted backing bytes are overwritten.
	lb, err := ce.enclave.Open()
	if err == nil {
		lb.Melt()
		lb.Scramble()
		lb.Destroy()
	}

	ce.enclave = nil
	fmt.Println("[enclave] Credential wiped. All LockedBuffers destroyed.")
}

