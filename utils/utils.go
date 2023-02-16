package utils

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"strconv"
)

func MakeUniqueValue(pName, pAddr string, pPort int) string {
	return pName + pAddr + strconv.Itoa(pPort)
}

func GetHashVal(hashString string) string {
	h := md5.Sum([]byte(hashString)) //nolint:gosec
	// Truncate the hash to 5 bytes and convert to uint64
	n := (uint64(h[0]) << 32) | (uint64(h[1]) << 24) | (uint64(h[2]) << 16) | (uint64(h[3]) << 8) | uint64(h[4])
	// Take the remainder after dividing by 10^12 and format as a zero-padded string
	return fmt.Sprintf("%012d", n%(1e12))
}
