package memory

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestManager_TranslateAddress(t *testing.T) {
	assert := assert.New(t)

	m := NewManager()
	file, err := os.Open("../rsrc/sample-init-dp.txt")
	assert.NoError(err)
	err = m.InitializeFromFile(file)
	assert.NoError(err)
	tests := []struct {
		cmd      string
		address  uint32
		expected int32
	}{
		{"RP", 2049, 0},
		{"RP", 19, -7},
		{"RP", 17, 3},
		{"RP", 18, 5000},
		{"RP", 1536, 10},
		{"RP", 1537, -20},

		{"TA", 2097162, 5130},
		{"TA", 2097674, 1034},
		{"TA", 2359306, 6666}, // page not resident
	}
	//log.Println(m.physicalMemory[16:20])

	for i, test := range tests {
		var res int32
		if test.cmd == "RP" {
			res = m.ReadPhysical(test.address)
		} else if test.cmd == "TA" {
			res, _ = m.TranslateAddress(test.address)
		}

		assert.Equal(test.expected, res)
		t.Logf("case %d\n", i)
	}
}
