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
	"unsafe"

	"github.com/awnumar/memguard"
)

func wipePlainString(s *string) {
	if len(*s) == 0 {
		return
	}
	// Obtain a mutable view of the string's backing array
	// And also delete the original back string cuz go usually makes copies
	b := unsafe.Slice(unsafe.StringData(*s), len(*s))
	memguard.WipeBytes(b)
	*s = ""
}

func printPwd(password *memguard.LockedBuffer) {
	// Use credentials here (e.g. authenticate to git.com).
	// These local vars live on Go heap — wipePlainString is called
	// after this function returns.
	fmt.Printf("=> password: %s\n", password.Bytes()) // shown for demo only
}

func testAccess(ce *cache.CredentialEnclave){
	var err error
	select{
		case <- time.After(30 * time.Second):
			err = ce.DoSthWithPwd(printPwd)
			fmt.Println(ce.RemainingTTL())
			// Start ur day here
			if err != nil {
				fmt.Println(err)
				memguard.SafePanic(err)
			}
			return
	}
}

func WipeRawCredential(c *tools.CredRest) {
	// wipePlainString(&c.Domain)
	// wipePlainString(&c.Username)
	// wipePlainString(&c.Password)
	// c.ID = 0
	memguard.WipeBytes([]byte(c.Password))
	c.Password = nil
}

func TCROSSCache() {
	// ── a) Signal handling: wipe on SIGINT / SIGTERM ──────────
	// CatchInterrupt covers SIGINT. We also catch SIGTERM manually.
	memguard.CatchInterrupt()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGHUP)
	signal.Notify(sigCh, syscall.SIGINT)

	// ── b) OS hardening ───────────────────────────────────────
	tools.HardenProcess()

	// ── c) Build the credential JSON (simulating vault read) ──
	// In production this would come from an encrypted vault file,
	// not a literal — but the ingestion pattern is identical.
	// rawSurf := tools.CredSurf{
	// 	ID: 			0,
	// 	Domain:   []byte("git.com"),
	// 	Username: []byte("mamacoco"),
	// }
	rawRest := tools.CredRest{
		Token: tools.GenerateRandomBytes(16),
		Password: []byte("M@m@L0v$"),
	}
	jsonBytes, err := json.Marshal(rawRest)
	if err != nil {
		memguard.SafePanic(err)
	}
	tools.WipeRawCredential(&rawRest)

	// ── d) Load into enclave with 5-minute TTL ────────────────
	const ttl = 1 * time.Minute

	ce, err := cache.NewCredentialEnclave(jsonBytes, true)
	// jsonBytes already wiped by NewCredentialEnclave → Move()
	if err != nil {
		memguard.SafePanic(err)
	}
	fmt.Printf("[enclave] Credential sealed. TTL = %v. Will auto-wipe at %v\n",
		ttl, time.Now().Add(ttl).Format(time.Kitchen))
	
	go testAccess(ce)

	// ── e) Demo: read the credential back safely ──────────────
	fmt.Println("\n[demo] Reading credential from enclave...")
	err = ce.DoSthWithPwd(printPwd)
	if err != nil {
		memguard.SafePanic(err)
	}

	// ── g) TTL / signal / manual-wipe race ───────────────────
	fmt.Println("\n[wait] Waiting for TTL expiry, SIGINT, or SIGTERM...")
	fmt.Printf("       RemainingTTL = %v\n", ce.RemainingTTL().Round(time.Second))
	fmt.Println("       (press Ctrl-C to trigger early wipe via SIGINT)")

	// In a real agent you'd keep serving requests here.
	// We wait for the enclave to be wiped (either by TTL or signal).
	select {
		case sig := <-sigCh:
			fmt.Printf("\n[signal] Received %v — wiping immediately\n", sig)
			ce.Wipe()
		case <- ce.WipedSig:
			fmt.Println("\n[TTL] Credential auto-wiped by timer")
			close(ce.WipedSig)
			ce = nil
	}

	fmt.Println("\n[exit] Secure termination. Calling memguard.Purge() + SafeExit(0).")
	memguard.Purge()
	memguard.SafeExit(0)
}
