package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"log"
)

func DecryptSecrets(secrets []string, passphrase string) map[string]string {
	return parseSecrets(secrets, passphrase)
}

func GetMD5Hash(text string) []byte {
	hash := md5.Sum([]byte(text))
	return hash[:]
}

func bytes_to_key(data []byte, salt []byte, output int32) []byte {
	merged := string(data) + string(salt)
	output = 48

	key := []byte(GetMD5Hash(merged))
	final_key := key
	for len(final_key) < int(output) {
		hash := GetMD5Hash(string(key) + string(merged))
		key = []byte(hash)
		final_key = []byte(string(final_key) + string(key))
	}
	return final_key[0:output]
}

type SecretAddedBy struct {
	name string
	id   string
}

type Secret struct {
	id        string
	key       string
	value     string
	comment   string
	addedDate string
	addedBy   SecretAddedBy
}
type Secrets struct {
	data map[string]*Secret
}

func parseSecrets(secrets []string, passcode string) map[string]string {
	parsed_secrets := make(map[string]string)

	for _, ciphertext := range secrets {
		decodedText, err := base64.StdEncoding.DecodeString(ciphertext)
		// fmt.Println(decodedText)
		if err != nil {
			log.Fatal("error:", err)
		}
		salted := decodedText[0:8]
		if string(salted) != "Salted__" {
			log.Fatal("error:", "Invalid encrypted data")
		}
		salt := decodedText[8:16]
		key_iv := bytes_to_key([]byte(passcode), salt, 48)
		key := key_iv[:32]
		iv := key_iv[32:]
		decrypted := decrypt(key, decodedText[16:], iv)
		var secretData map[string]interface{}
		json.Unmarshal(decrypted, &secretData)
		secretKey := secretData["key"]

		// addedBy := secretData["addedBy"].(map[string]interface{})

		parsed_secrets[secretKey.(string)] = secretData["value"].(string)
	}
	return parsed_secrets
}

func decrypt(key []byte, ciphertext []byte, iv []byte) []byte {

	// fmt.Println(key)
	newCipher, _ := aes.NewCipher([]byte(key))
	cfbdec := cipher.NewCBCDecrypter(newCipher, iv)
	decipher := make([]byte, len(ciphertext))
	cfbdec.CryptBlocks(decipher, ciphertext)
	decipher = removeBadPadding(decipher)
	return decipher
}

func removeBadPadding(b64 []byte) []byte {
	last := b64[len(b64)-1]
	if last > 16 {
		return b64
	}
	return b64[:len(b64)-int(last)]
}
