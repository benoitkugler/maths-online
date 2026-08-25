package automatismes

import (
	"testing"

	"github.com/benoitkugler/maths-online/server/src/sql/editor"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestGetTrivialAutomatismesIn_match(t *testing.T) {
	tests := []struct {
		args      GetAutomatismesIn
		questions editor.TagGroup
		want      bool
	}{
		{
			GetAutomatismesIn{}, editor.TagGroup{}, false,
		},
		{
			GetAutomatismesIn{editor.Premiere, "TECHNO"}, editor.TagGroup{
				TagIndex: editor.TagIndex{Level: editor.Premiere},
			}, false,
		},
		{
			GetAutomatismesIn{editor.Premiere, ""}, editor.TagGroup{
				TagIndex: editor.TagIndex{Level: editor.Premiere, Chapter: "AUTOMATISMES"},
			}, true,
		},
		{
			GetAutomatismesIn{editor.Premiere, "TECHNO"}, editor.TagGroup{
				TagIndex:  editor.TagIndex{Level: editor.Premiere, Chapter: "AUTOMATISMES"},
				SubLevels: []string{"TECHNO"},
			}, true,
		},
		{
			GetAutomatismesIn{editor.Troisieme, ""}, editor.TagGroup{
				TagIndex: editor.TagIndex{Level: editor.Troisieme},
			}, true,
		},
		{
			GetAutomatismesIn{editor.Troisieme, "TECHNO"}, editor.TagGroup{
				TagIndex: editor.TagIndex{Level: editor.Troisieme},
			}, false,
		},
	}
	for _, tt := range tests {
		got := tt.args.match(tt.questions)
		tu.Assert(t, got == tt.want)
	}
}
