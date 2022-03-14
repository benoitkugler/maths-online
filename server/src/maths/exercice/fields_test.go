package exercice

import (
	"testing"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
)

func TestFieldInstance_validateAnswerSyntax(t *testing.T) {
	tests := []struct {
		field   fieldInstance
		args    client.Answer
		wantErr bool
	}{
		{NumberFieldInstance{}, client.NumberAnswer{}, false},
		{NumberFieldInstance{}, client.RadioAnswer{}, true},
	}
	for _, tt := range tests {
		err := tt.field.validateAnswerSyntax(tt.args)
		if (err != nil) != tt.wantErr {
			t.Fatalf("NumberFieldInstance.validateAnswerSyntax() error = %v, wantErr %v", err, tt.wantErr)
		}
		if err != nil {
			_ = err.(InvalidFieldAnswer).Error()
		}
	}
}

func TestNumberFieldInstance_evaluateAnswer(t *testing.T) {
	tests := []struct {
		field         fieldInstance
		args          client.Answer
		wantIsCorrect bool
		wantErr       bool
	}{
		{NumberFieldInstance{Answer: 1}, client.NumberAnswer{Value: 1}, true, false},
		{NumberFieldInstance{Answer: 1}, client.NumberAnswer{Value: 1.1}, false, false},
	}
	for _, tt := range tests {
		f := tt.field
		gotIsCorrect, err := f.evaluateAnswer(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("NumberFieldInstance.evaluateAnswer() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if gotIsCorrect != tt.wantIsCorrect {
			t.Errorf("NumberFieldInstance.evaluateAnswer() = %v, want %v", gotIsCorrect, tt.wantIsCorrect)
		}
	}
}
