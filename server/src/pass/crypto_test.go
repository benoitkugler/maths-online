package pass

import (
	"testing"
)

func TestEncryption(t *testing.T) {
	enc := newEncrypter("5s64qsd897e4q87mùlds54")
	for i := range [200]int{} {
		v1 := 456 + 100*int64(i)
		s, err := newEncryptedID(v1, enc)
		if err != nil {
			t.Fatal(err)
		}
		v2, err := s.Decrypt(enc)
		if err != nil {
			t.Fatal(err)
		}
		if v1 != v2 {
			t.Fatal(v1, v2)
		}
	}
}
