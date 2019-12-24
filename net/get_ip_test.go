package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExternalIP_WithIgnoreLoopbackSetToFalse(t *testing.T) {

	// Arrange.

	ignoreLoopback := false

	// Act.

	ip, err := GetIP(ignoreLoopback)

	// Assert.

	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", ip)
}

func TestExternalIP_WithIgnoreLoopbackSetToTrue(t *testing.T) {

	// Arrange.

	ignoreLoopback := true

	// Act.

	ip, err := GetIP(ignoreLoopback)

	// Assert.

	assert.Nil(t, err)
	assert.NotEqual(t, "127.0.0.1", ip)
}
