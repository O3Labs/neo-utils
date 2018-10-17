package neoutils_test

import (
	"crypto/sha256"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func StringToNNSNameHash(v string) []byte {
	hash := sha256.Sum256([]byte(v))

	return neoutils.ReverseBytes(hash[:])
}

func NameHash(domain string, subdomain string) []byte {
	root := sha256.Sum256([]byte(domain))
	sub := sha256.Sum256([]byte(subdomain))
	full := append(sub[:], root[:]...)
	hash := sha256.Sum256(full)
	return hash[:]
}

func TestNNSHash(t *testing.T) {
	log.Printf("%x", StringToNNSNameHash("neo"))
}

// o3.neo = 231420791e6cfa9b6eee03f497cf0c24589435a92ab7b7b232a5951088afa3d6

func TestNameHash(t *testing.T) {
	// expected := "b8e55096b5871c43375881259c5664a366b67671e9ede571d369aee4a5188bc0"

	result := NameHash("neo", "o3")
	log.Printf("%x", neoutils.ReverseBytes(result))
}
func TestNNSResolver(t *testing.T) {
	//0020d86aa99be2f6689dc767c13be2f6a40098511b3c039d4266228656d79d1c64c6
	//046164647253c1077265736f6c766567c72871904920c0d977326620e4754a6c11878334

	//0020d86aa99be2f6689dc767c13be2f6a40098511b3c039d4266228656d79d1c64c6
	//046164647253c1077265736f6c766567c72871904920c0d977326620e4754a6c11878334
	//
	namehash := NameHash("neo", "o3")
	log.Printf("namehash = %x", namehash)
	scriptHash, err := smartcontract.NewScriptHash("348387116c4a75e420663277d9c02049907128c7")
	if err != nil {
		log.Printf("err = %v", err)
		t.Fail()
		return
	}
	// json := fmt.Sprintf("%x", result)
	args := []interface{}{[]byte("addr"), namehash, []byte("")}
	s := smartcontract.NewScriptBuilder()
	s.GenerateContractInvocationScript(scriptHash, "resolve", args)

	log.Printf("%x", s.ToBytes())
}
