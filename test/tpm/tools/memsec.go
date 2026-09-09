package tools

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/awnumar/memguard"
	"golang.org/x/sys/unix"
)

// Memory Hardening
// - RLIMIT_CORE = 0 (disable core dump)
// - prctl(PR_SET_DUMPABLE, 0) - (Prevent reading by other processes like ptrace-based debuggers, does not work for root user.)
// - prctl(PR_SET_SECCOMP) or seccomp-bpf — Without syscall filtering, any of your agents can be exploited to call `ptrace()` on sibling agents
// - Swap encryption — `mlock` prevents swapping of locked pages, but unlocked pages (your Go runtime overhead, stack frames) may still be swapped. Consider `zswap`/`zram` with encryption at the OS level as a complementary control
// - Signal Handling
// - GC Bypass (mmap memory)
// - Secure Wipe.

func WipeVault[C Credential](v *Vault[C]) {
	for i := range v.Credentials {
		WipeStruct(&v.Credentials[i])
	}
	v.Credentials = nil
}

func WipeStruct(v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		if f.Kind() == reflect.Slice && f.Type().Elem().Kind() == reflect.Uint8 && f.CanInterface() {
			memguard.WipeBytes(f.Interface().([]byte))
		}
	}
}

func WipeRawCredential(c *CredRest) {
	// wipePlainString(&c.Domain)
	// wipePlainString(&c.Username)
	// wipePlainString(&c.Password)
	// c.ID = 0
	memguard.WipeBytes(c.Password)
	c.Password = nil
}

// wipePlainString zeroes a plain Go string's backing bytes.
// Strings are immutable in Go, but their backing memory is just bytes.
// This is technically unsafe but necessary for sensitive data.
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

func HardenProcess() {
	// Belt-and-suspenders: memguard already sets RLIMIT_CORE=0 at init,
	// but we do it again explicitly to be certain.
	_ = unix.Setrlimit(unix.RLIMIT_CORE, &unix.Rlimit{Cur: 0, Max: 0})
 
	// Prevent ptrace-based inspection by other processes (non-root only).
	// Does not protect against root. Combine with seccomp-bpf for full coverage.
	if err := unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "[warn] prctl PR_SET_DUMPABLE: %v\n", err)
	}

	// seccomp-bpf

	// Swap Encryption

	// mlock()
 
	fmt.Println("[harden] Core dumps disabled. PR_SET_DUMPABLE=0.")
}

func GenerateRandomBytes(size int) []byte {
	b := make([]byte, 16) // 16 bytes
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	text := hex.EncodeToString(b)
	fmt.Println("[Rand] Generage a %d random bytes: %s", size, text)
	return b
}

func GetMD5Hash(text string) string {
   hash := md5.Sum([]byte(text))
   return hex.EncodeToString(hash[:])
}
