package helpers_test

import (
	"github.com/2f4ek/lets-go-chat/internal/helpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSecureToken(t *testing.T) {
	token := helpers.GenerateSecureToken()

	assert.Equal(t, 32, len(token))
}
