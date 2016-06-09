package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
)

var (
	cryptKey []byte // a key used in crypto/cipher algorithm
)

// SetKey changes a key that package using for crypto/cipher algorithm
func SetKey(key []byte) error {
	if diff := aes.BlockSize - len(key)%aes.BlockSize; diff > 0 {
		for diff > 0 {
			key = append(key, 0)
			diff--
		}
	}
	cryptKey = key
	return nil
}

// GetKey returns a key used in crypto/cipher algorithm
func GetKey() []byte {
	if cryptKey == nil {
		SetKey([]byte("hard-coded key:)"))
	}
	return cryptKey
}

// EncryptString encrypts string with crypto/cipher, salting it and makes base64.StdEncoding, returns blank string if encoding fails
func EncryptString(data string) string {
	// cypher encryption
	result, err := EncryptData([]byte(data))
	if err != nil {
		return ""
	}

	// salting
	salt := []byte{':'}
	if cryptKey := GetKey(); len(cryptKey) > 0 {
		salt = append(salt, cryptKey[0])
	}
	result = append(result, salt...)

	//base64 encoding
	return base64.StdEncoding.EncodeToString(result)
}

// DecryptString decodes base64.StdEncoding string un-salting it and then decrypts it with crypto/cipher, returns original value or error
func DecryptString(data string) string {

	// base64 decoding
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return data
	}

	// checking and removing salt
	salt := ":"
	if cryptKey := GetKey(); len(cryptKey) > 0 {
		salt += string(cryptKey[0])
	}
	saltIdx := len(decodedData) - len(salt)
	if saltIdx < 0 || salt != string(decodedData[saltIdx:]) {
		return data
	}
	decodedData = decodedData[0:saltIdx]

	// making cypher decryption
	result, err := DecryptData(decodedData)
	if err != nil {
		return data
	}

	return string(result)
}

// EncryptData encrypts given data with crypto/cipher algorithm
func EncryptData(data []byte) ([]byte, error) {
	var buffer bytes.Buffer

	writer, err := EncryptWriter(&buffer)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// DecryptData decrypts given data with crypto/cipher algorithm
func DecryptData(encodedData []byte) ([]byte, error) {
	result := make([]byte, len(encodedData))

	reader, err := EncryptReader(bytes.NewReader(encodedData))
	if err != nil {
		return nil, err
	}

	_, err = reader.Read(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// EncryptReader decrypts given stream with crypto/cipher algorithm
func EncryptReader(rawReader io.Reader) (io.Reader, error) {

	cryptKey := GetKey()
	cipherBlock, err := aes.NewCipher(cryptKey)
	if err != nil {
		return nil, err
	}

	iv := cryptKey[:aes.BlockSize]
	stream := cipher.NewOFB(cipherBlock, iv)

	return &cipher.StreamReader{S: stream, R: rawReader}, nil
}

// EncryptWriter encrypts given stream with crypto/cipher algorithm
func EncryptWriter(rawWriter io.Writer) (io.Writer, error) {

	cryptKey := GetKey()
	cipherBlock, err := aes.NewCipher(cryptKey)
	if err != nil {
		return nil, err
	}

	iv := cryptKey[:aes.BlockSize]
	stream := cipher.NewOFB(cipherBlock, iv)

	return &cipher.StreamWriter{S: stream, W: rawWriter}, nil
}

// CryptToURLString encrypts given string with base64 and hex encoding
func CryptToURLString(raw []byte) string {
	result := hex.EncodeToString([]byte(base64.StdEncoding.EncodeToString(raw)))

	return result
}

// CryptAsURLString encrypts given string with base64 and hex encoding
func CryptAsURLString(rawString string) string {
	result := hex.EncodeToString([]byte(base64.StdEncoding.EncodeToString([]byte(rawString))))

	return result
}

// DecryptURLString decode given encoded string and returns decoded value
func DecryptURLString(encodedString string) (string, error) {

	var result string

	data, err := hex.DecodeString(encodedString)
	if err != nil {
		return result, err
	}

	data, err = base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return result, err
	}

	if data != nil {
		result = string(data)
	}

	return result, nil
}

// PasswordEncode encode inputed password with using salt, if no salt it will use default one
func PasswordEncode(password string, salt string) string {

	hasher := md5.New()
	if salt == "" {
		salt := ":"
		if len(password) > 2 {
			salt += password[0:1]
		}
		hasher.Write([]byte(password + salt))
	} else {
		hasher.Write([]byte(salt + password))
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

// PasswordCheck compare inputed password with stored one
func PasswordCheck(password string, input string) bool {

	password = strings.TrimSpace(password)
	input = strings.TrimSpace(input)

	salt := ""

	tmp := strings.Split(password, ":")
	if len(tmp) == 2 {
		password = tmp[0]
		salt = tmp[1]
	}

	return PasswordEncode(input, salt) == password
}
