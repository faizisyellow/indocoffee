package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestToken(t *testing.T) {

	t.Run("return uuid", func(t *testing.T) {
		tkn := UUID{
			Plaintoken: uuid.New().String(),
		}
		got := tkn.Generate()
		if len(got) == 0 {
			t.Errorf("expected have a string but got empty string")
		}
	})

}
