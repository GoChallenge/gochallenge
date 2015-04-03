package model_test

import (
	"testing"

	"github.com/gochallenge/gochallenge/model"
	"github.com/stretchr/testify/require"
)

func TestErrorOutput(t *testing.T) {
	err := model.ErrNotFound
	require.Equal(t, err.Error(), "Not found")
}
