package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

func StringToByte(text string) []byte {
	var response []byte
	for _, item := range text {
		number := int(item)
		response = append(response, byte(number))
	}
	return response
}

func SliceIntToSliceHex(slice []int) []string {
	var response []string
	for _, value := range slice {
		hex := strconv.FormatInt(int64(value), 16)
		hex = strings.ToUpper(hex)
		response = append(response, hex)
	}
	return response
}

func Checksum(command []byte) byte {
	var checkSum = 0
	for _, number := range command {
		checkSum ^= int(number)
	}
	checkSum ^= 0xff
	return byte(checkSum)
}

func EncryptCommand(command []byte, secret byte) []byte {
	var encrypted []byte
	for _, number := range command {
		encrypted = append(encrypted, number^secret)
	}
	return encrypted
}

func GetMD5Hash(text string, length int) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	hash := hex.EncodeToString(hasher.Sum(nil))
	if length > 0 {
		return hash[:length]
	}
	return hash
}

func HexToByte(hex string) []byte {
	stringSlice := SplitBy(hex, 2)
	var byteSlice []byte
	for _, hex := range stringSlice {
		newHex := "0x" + hex
		newByte, _ := strconv.ParseInt(newHex, 0, 16)
		byteSlice = append(byteSlice, byte(newByte))
	}
	return byteSlice
}

func SplitBy(s string, n int) []string {
	var ss []string
	for i := 1; i < len(s); i++ {
		if i%n == 0 {
			ss = append(ss, s[:i])
			s = s[i:]
			i = 1
		}
	}
	ss = append(ss, s)
	return ss
}

func ExtractCrcHigh(crc int) int {
	number := 0xFF00
	result := (crc & number) >> 8
	return result
}

func ExtractCrcLow(crc int) int {
	number := 0x00FF
	result := crc & number
	return result
}

func ByteToBooleanSlice(bytes []byte) []bool {
	var slice = make([]bool, 0)
	for _, byte := range bytes {
		slice = append(slice, byte&1 != 0)
		slice = append(slice, byte&2 != 0)
		slice = append(slice, byte&4 != 0)
		slice = append(slice, byte&8 != 0)
		slice = append(slice, byte&16 != 0)
		slice = append(slice, byte&32 != 0)
		slice = append(slice, byte&64 != 0)
		slice = append(slice, byte&128 != 0)
	}
	return slice
}

func Contains(obj []any, key any) bool {
	for _, item := range obj {
		if item == key {
			return true
		}
	}
	return false
}

func Every(obj []any) bool {
	for i := 1; i < len(obj); i++ {
		if obj[i] != obj[0] {
			return false
		}
	}
	return true
}
