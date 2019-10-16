package repository_test

import (
	"repotags/repository"
	"testing"
)

func TestOpenConnection(t *testing.T) {
	_, err := repository.GetDBConnection()
	if err != nil {
		t.Fatal(err)
	}
}
