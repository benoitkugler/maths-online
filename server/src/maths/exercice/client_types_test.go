package exercice

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestClientExerciceJSON(t *testing.T) {
	for _, qu := range questions {
		cl := qu.toClient()

		b, err := json.Marshal(cl)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(b))
		fmt.Println()

		var cl2 ClientQuestion
		err = json.Unmarshal(b, &cl2)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(cl, cl2) {
			t.Fatal(cl, cl2)
		}
	}
}
