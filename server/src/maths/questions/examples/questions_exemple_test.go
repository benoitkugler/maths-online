package examples

// import (
// 	"encoding/json"
// 	"reflect"
// 	"testing"

// 	"github.com/benoitkugler/maths-online/maths/questions/client"
// )

// func TestClientExerciceJSON(t *testing.T) {
// 	for _, qu := range PredefinedQuestions {
// 		cl := qu.ToClient()

// 		b, err := json.Marshal(cl)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		var cl2 client.Question
// 		err = json.Unmarshal(b, &cl2)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if !reflect.DeepEqual(cl, cl2) {
// 			t.Fatal(cl, cl2)
// 		}
// 	}
// }

// func TestAllFields(t *testing.T) {
// 	_, err := loadAllFieldsQuestion()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }