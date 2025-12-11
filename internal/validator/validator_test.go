package validator

import (
	"testing"
)

func TestValidate(t *testing.T) {
	// Create a minimal valid JPEG header + footer
	validJPEG := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00}

	padding := make([]byte, MinSize)
	validJPEG = append(validJPEG, padding...)
	validJPEG = append(validJPEG, 0xFF, 0xD9)

	tests := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name:    "Empty file",
			data:    []byte{},
			wantErr: ErrEmpty,
		},
		{
			name:    "Too small",
			data:    []byte{0xFF, 0xD8, 0xFF, 0xD9},
			wantErr: ErrInvalidSize,
		},
		{
			name:    "Not a JPEG (PNG signature)",
			data:    append([]byte{0x89, 'P', 'N', 'G'}, make([]byte, MinSize)...),
			wantErr: ErrNotJPEG,
		},
		{
			name:    "Valid JPEG structure",
			data:    validJPEG,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.data)
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

