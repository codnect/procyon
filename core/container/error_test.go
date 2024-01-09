package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotFoundError_Error(t *testing.T) {
	err := &notFoundError{ErrorString: "anyError"}
	assert.Equal(t, "anyError", err.Error())
}
