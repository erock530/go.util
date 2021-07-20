//Package bits is for manipulation/error checking at the bit level
package bits

import "errors"

//CheckOddParity returns true if byte passed odd parity bit check
func CheckOddParity(b byte) bool {
	return (((b & 0x80) > 0) != ((b & 0x40) > 0) != ((b & 0x20) > 0) != ((b & 0x10) > 0) != ((b & 0x08) > 0) != ((b & 0x04) > 0) != ((b & 0x02) > 0) != ((b & 0x01) > 0))
}

//CheckEvenParity returns true if byte passed even parity bit check
func CheckEvenParity(b byte) bool {
	return (((b & 0x80) > 0) == ((b & 0x40) > 0) == ((b & 0x20) > 0) == ((b & 0x10) > 0) == ((b & 0x08) > 0) == ((b & 0x04) > 0) == ((b & 0x02) > 0) == ((b & 0x01) > 0))
}

//EncodeOddParity encodes the parity bit as the MSB
func EncodeOddParity(b byte) (byte, error) {
	if b > 0x7f {
		return 0x00, errors.New("Byte can be a max value of 0x7f")
	}

	if CheckOddParity(b + 0x80) {
		b += 0x80
	}
	return b, nil
}
