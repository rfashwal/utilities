package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUuidValidFailure(t *testing.T) {

	// Arrange.

	fakeUuid := "not a valid uuid"

	// Act.

	result := IsUuidValid(fakeUuid)

	// Assert.

	assert.False(t, result)
}

func TestIsUuidValidSuccessful(t *testing.T) {

	// Arrange.

	uuid := "76f44e68-e564-447d-90c9-af622bbe692a"

	// Act.

	result := IsUuidValid(uuid)

	// Assert.

	assert.True(t, result)
}

func TestIsUuidValid_section3_failure(t *testing.T) {

	// Arrange.

	uuid := "5a96689c-f25c-1d2c-8c96-72dc60a34799"

	// Act.

	result := IsUuidValid(uuid)

	// Assert.

	assert.False(t, result)
}

func TestIsUuidValid_section4_failure(t *testing.T) {

	// Arrange.

	uuid := "f43fa408-4acb-49a8-0c4d-dfb0a75ae559"

	// Act.

	result := IsUuidValid(uuid)

	// Assert.

	assert.False(t, result)
}
