package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestError(t *testing.T) {
	reqError := NewRequestError(nil, InternalServerError, "test", nil)

	assert.Equal(t, "test", reqError.Msg)
	assert.Equal(t, InternalServerError, reqError.Status)
}
