package exercice

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
)

func TestClientExerciceJSON(t *testing.T) {
	for _, qu := range questions {
		cl := qu.toClient()

		b, err := json.Marshal(cl)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("%q\n", string(b))
		fmt.Println()

		var cl2 client.Question
		err = json.Unmarshal(b, &cl2)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(cl, cl2) {
			t.Fatal(cl, cl2)
		}
	}
}
