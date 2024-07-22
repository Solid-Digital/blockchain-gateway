package apperr

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestError_Copy(t *testing.T) {
	type fields struct {
		Error *Error
		Cause error
	}
	tests := []struct {
		name   string
		fields fields
		want   *Error
	}{
		{
			name: "",
			fields: fields{
				Error: &Error{
					NamedErrors: map[string][]string{"asdf": {"asd", "qwe"}, "qwerty": {"??", "!!"}},
					Message:     "one",
					RequestID:   "two",
					Status:      15,
					cause:       nil,
				},
			},
			want: &Error{
				NamedErrors: map[string][]string{"asdf": {"asd", "qwe"}, "qwerty": {"??", "!!"}},
				Message:     "one",
				RequestID:   "two",
				Status:      15,
				cause:       nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.fields.Error
			got := e.Copy()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}

			got.NamedErrors["asdf"][0] = "zzz"
			got.NamedErrors["gg"] = []string{"zzz"}

			if tt.fields.Error.NamedErrors["asdf"][0] == "zzz" {
				t.Errorf("Copy() = modifying the destination also modifies the source")
			}

			if len(tt.fields.Error.NamedErrors["gg"]) != 0 {
				t.Errorf("Copy() = modifying the destination also modifies the source")
			}

			spew.Dump(e, got)
		})
	}
}
