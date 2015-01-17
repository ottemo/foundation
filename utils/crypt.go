package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"
)

// crypt.go providing an centralized way for bi-directional crypt of secure data,
//   - note that SetKey() makes change for entire application,
// 	   so, if you want local effect you should restore it after usage
//   - normally application should take care about SetKey() on init and you should not touch it
//   - if SetKey() was not called during application init then default hard-coded key will be used
//
//   Example 1:
//     source := "just test"
//     encoded := utils.EncryptStringBase64(source)
//     decoded := utils.DecryptStringBase64(encoded)
//     println( "'" + source + "' --encode--> '" + encoded + "' --decode--> '" +  decoded + "'")
//
//   Output:
//     'just test' --encode--> 'Ddryse1yNL5z' --decode--> 'just test'
//
//
//   Example 2:
//     sampleData := []byte("It is just a sample.")
//
//     outFile, _ := os.OpenFile("sample.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
//     defer outFile.Close()
//     writer, _ := utils.EncryptWriter(outFile)
//     writer.Write(sampleData)
//
//     inFile, _ := os.OpenFile("sample.txt", os.O_RDONLY, 0600)
//     defer inFile.Close()
//     reader, _ := utils.EncryptReader(inFile)
//     readBuffer := make([]byte, 10)
//
//     reader.Read(readBuffer)
//     println(string(readBuffer))
//     reader.Read(readBuffer)
//     println(string(readBuffer))
//
//   Output:
//     It is just
//      a sample.
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
