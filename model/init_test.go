package model

import (
	"testing"
)

func Test_makeTables(t *testing.T) {
	t.Run("Create tables.", func(t *testing.T) {
		if err := makeTables(); err != nil {
			t.Errorf("makeTables() error = %v", err)
		}
	})
}
