package bits

import "testing"

func TestCheckOddParity(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args byte
		want bool
	}{
		{"00000000", byte(0x00), false},
		{"00000001", byte(0x01), true},
		{"10100010", byte(0xA2), true},
		{"10100011", byte(0xA3), false},
		{"11010010", byte(0xD2), false},
		{"11010011", byte(0xD3), true},
		{"11111110", byte(0xFE), true},
		{"11111111", byte(0xFF), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckOddParity(tt.args); got != tt.want {
				t.Errorf("CheckOddParity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckEvenParity(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args byte
		want bool
	}{
		{"00000000", byte(0x00), true},
		{"00000001", byte(0x01), false},
		{"10100010", byte(0xA2), false},
		{"10100011", byte(0xA3), true},
		{"11010010", byte(0xD2), true},
		{"11010011", byte(0xD3), false},
		{"11111110", byte(0xFE), false},
		{"11111111", byte(0xFF), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckEvenParity(tt.args); got != tt.want {
				t.Errorf("CheckEvenParity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeOddParity(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name    string
		args    byte
		want    byte
		wantErr bool
	}{
		{"0000000", 0x00, 0x80, false},
		{"81-1010001", 0x51, 0x51, false},
		{"105-1101001", 0x69, 0xE9, false},
		{"127-1111111", 0x7F, 0x7F, false},
		{"128-10000000", 0x80, 0x00, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeOddParity(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeOddParity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeOddParity() = %v, want %v", got, tt.want)
			}
		})
	}
}
