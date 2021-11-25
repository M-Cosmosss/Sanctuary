package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteStore(t *testing.T) {
	err := Init()
	assert.Nil(t, err)

}
