package toolkit

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)
	require.Len(t, s, 10, "Wrong length for RandomString")
}
