package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHashVal(t *testing.T) {
	getHash1 := GetHashVal("VATZ has Test 1")
	println("getHash1: ", getHash1)
	getHash2 := GetHashVal("VATZ has Test 2")
	println("getHash2: ", getHash2)
	getHash3 := GetHashVal("VATZ has Test 3")
	println("getHash3: ", getHash3)
	getHash4 := GetHashVal("VATZ has Test 2")
	println("getHash2: ", getHash2)
	assert.Equal(t, "459322276471", getHash1)
	assert.Equal(t, "513315553674", getHash2)
	assert.Equal(t, "344345442263", getHash3)
	assert.Equal(t, "513315553674", getHash4)
}
