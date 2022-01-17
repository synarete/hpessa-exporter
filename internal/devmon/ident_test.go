package devmon_test

import (
	"testing"

	"github.com/red-hat-storage/hpessa-exporter/internal/devmon"
	"github.com/stretchr/testify/assert"
)

func TestSelfIdent(t *testing.T) {
	ident := devmon.SelfIdent()
	assert.NotNil(t, ident)
	assert.NotNil(t, ident.User)
	assert.NotEqual(t, ident.Progname, "")
	assert.NotEqual(t, ident.User.Uid, "")
	assert.NotEqual(t, ident.User.Gid, "")
}
