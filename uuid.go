package util

import (
	"github.com/erock530/go.util/uuidutil"
)

//GenUUID Generates a pseudo-random UUID string
func GenUUID() (string, error) {
	return uuidutil.GenUUID()
}
