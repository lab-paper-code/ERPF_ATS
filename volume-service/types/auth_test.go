package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	t.Run("test Auth", testAuth)
}

func testAuth(t *testing.T) {
	auth := CheckAuthKey("id21", "pw33", "aWQyMTpwdzMz")
	assert.True(t, auth)
}
