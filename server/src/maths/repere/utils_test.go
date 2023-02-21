package repere

import "testing"

func TestColor_ToARGB(t *testing.T) {
	tests := []struct {
		c     Color
		wantA uint8
		wantR uint8
		wantG uint8
		wantB uint8
	}{
		{"#FF00FF", 255, 255, 0, 255},
		{"#00Ff00Ff", 0, 255, 0, 255},
		{"#08Ff00Ff", 8, 255, 0, 255},
	}
	for _, tt := range tests {
		gotA, gotR, gotG, gotB := tt.c.ToARGB()
		if gotA != tt.wantA {
			t.Errorf("Color.ToARGB() gotA = %v, want %v", gotA, tt.wantA)
		}
		if gotR != tt.wantR {
			t.Errorf("Color.ToARGB() gotR = %v, want %v", gotR, tt.wantR)
		}
		if gotG != tt.wantG {
			t.Errorf("Color.ToARGB() gotG = %v, want %v", gotG, tt.wantG)
		}
		if gotB != tt.wantB {
			t.Errorf("Color.ToARGB() gotB = %v, want %v", gotB, tt.wantB)
		}
	}
}
