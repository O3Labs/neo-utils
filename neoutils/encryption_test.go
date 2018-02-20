package neoutils

import (
	"log"
	"testing"
)

func TestDecrypt(t *testing.T) {
	encryptedString := "nkGfKuzVOBCxMpKZJ3ejjSYlTTP8-XnSKfoXMndPU5w="
	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	wallet, _ := GenerateFromWIF(wif)

	decrypted := Decrypt(wallet.PrivateKey, encryptedString)
	log.Printf("%v", decrypted)
}

func TestEncryptWithAddress(t *testing.T) {
	alice, err := NewWallet()
	if err != nil {
		log.Printf("%+v", err)
		t.Fail()
		return
	}

	message := "Hello bob"

	encryptedText := Encrypt([]byte(alice.PublicKey), message)
	log.Printf("encrypted = %v", encryptedText)
}

func TestEncryptMessage(t *testing.T) {
	alice, err := NewWallet()
	if err != nil {
		log.Printf("%+v", err)
		t.Fail()
		return
	}

	bob, err := NewWallet()
	if err != nil {
		log.Printf("%+v", err)
		t.Fail()
		return
	}

	randomPerson, err := NewWallet()
	if err != nil {
		log.Printf("%+v", err)
		t.Fail()
		return
	}

	log.Printf(`
		alice's adress = %+v\n
		bob's adress = %+v\n
		alice's public key = %+v\n
		bob's public key = %+v\n
		`,
		alice.Address,
		bob.Address,
		bytesToHex(alice.PublicKey),
		bytesToHex(bob.PublicKey),
	)

	alicesSharedSecret := alice.ComputeSharedSecret(bob.PublicKey)
	bobsSharedSecret := bob.ComputeSharedSecret(alice.PublicKey)

	randomPersonSharedSecretWithAlice := randomPerson.ComputeSharedSecret(alice.PublicKey)
	randomPersonSharedSecretWithBob := randomPerson.ComputeSharedSecret(bob.PublicKey)

	hexSlice := bytesToHex(alicesSharedSecret)
	hexBob := bytesToHex(bobsSharedSecret)
	hexRandomPerson := bytesToHex(randomPersonSharedSecretWithAlice)
	hexRandomPersonWithBob := bytesToHex(randomPersonSharedSecretWithBob)
	log.Printf("alice's shared secret (bob's public key) = %+v\n", hexSlice)
	log.Printf("bob's shared secret (alice's public key) = %+v\n", hexBob)
	log.Printf("random person's shared secret (alice's public key) = %+v\n", hexRandomPerson)
	log.Printf("random person's shared secret (bob's public key) = %+v\n", hexRandomPersonWithBob)

	message := "Hello bob"

	encryptedText := Encrypt(alicesSharedSecret, message)
	log.Printf("alice sent %v to bob\n", encryptedText)
	log.Printf("bob sent %v to alice\n", Encrypt(bobsSharedSecret, "Hey there Alice"))

	decryptedText := Decrypt(bobsSharedSecret, encryptedText)

	log.Printf(`bob uses his shared secret to read alice's message. "%v" \n`, decryptedText)
	log.Printf(`random person tries to use his shared secret between he and alice to read bob message sent by alice "%v" \n`, Decrypt(randomPersonSharedSecretWithAlice, encryptedText))

	log.Printf(`random person tries to use his shared secret between he and bob to read bob message sent by alice "%v" \n`, Decrypt(randomPersonSharedSecretWithBob, encryptedText))
}
