package uuidutil

import (
	"crypto/rand"
	"encoding/hex"
)

//GenUUID Generates a pseudo-random UUID string
//Consider using the go.uuid repository instead.
func GenUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	if n != len(uuid) || err != nil {
		return GenUUID()
	}
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7
	return hex.EncodeToString(uuid), nil
}
