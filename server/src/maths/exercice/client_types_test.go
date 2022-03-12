package exercice

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestClientExerciceJSON(t *testing.T) {
	cl := questions[0].toClient()

	b, err := json.Marshal(cl)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(b))

	var cl2 ClientQuestion
	err = json.Unmarshal(b, &cl2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(cl, cl2) {
		t.Fatal(cl, cl2)
	}
}
