// SPDX-License-Identifier: Apache-2.0
package devmon_test

import (
	"testing"

	"github.com/red-hat-storage/hpessa-exporter/internal/devmon"
	"github.com/stretchr/testify/assert"
)

func TestSysUname(t *testing.T) {
	uname, err := devmon.Uname()
	assert.NoError(t, err)
	assert.NotNil(t, uname)
	assert.Greater(t, len(uname.Version), 0)
}
