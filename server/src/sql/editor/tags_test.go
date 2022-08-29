package editor

import (
	"reflect"
	"sort"
	"testing"
)

func TestCommonTags(t *testing.T) {
	tests := []struct {
		args [][]string
		want []string
	}{
		{
			[][]string{
				{"1", "2"},
				{"1", "3"},
				{"3"},
			},
			nil,
		},
		{
			[][]string{
				{"1", "2"},
				{"1", "3"},
				{"1", "3"},
			},
			[]string{"1"},
		},
		{
			[][]string{
				{"1", "2", "3"},
				{"1", "3"},
				{"1", "3"},
			},
			[]string{"1", "3"},
		},
	}
	for _, tt := range tests {
		got := CommonTags(tt.args)
		sort.Strings(got)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("CommonTags() = %v, want %v", got, tt.want)
		}
	}
}
