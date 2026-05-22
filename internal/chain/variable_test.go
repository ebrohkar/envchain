package chain_test

import (
	"testing"

	"github.com/envchain/envchain/internal/chain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariable_Validate_Resolved(t *testing.T) {
	v := &chain.Variable{
		Name:   "PORT",
		Value:  "8080",
		Status: chain.StatusResolved,
	}
	err := v.Validate()
	assert.NoError(t, err)
}

func TestVariable_Validate_Missing(t *testing.T) {
	v := &chain.Variable{
		Name:   "PORT",
		Status: chain.StatusMissing,
	}
	err := v.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing")
	assert.Contains(t, err.Error(), "PORT")
}

func TestVariable_Validate_Empty(t *testing.T) {
	v := &chain.Variable{
		Name:   "TOKEN",
		Status: chain.StatusEmpty,
	}
	err := v.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestVarStatus_String(t *testing.T) {
	cases := []struct {
		status chain.VarStatus
		want   string
	}{
		{chain.StatusResolved, "resolved"},
		{chain.StatusMissing, "missing"},
		{chain.StatusEmpty, "empty"},
		{chain.StatusUnknown, "unknown"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, tc.status.String())
	}
}
