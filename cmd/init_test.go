package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WriteLocalConfig(t *testing.T) {
	type test struct {
		name           string
		in             *Config
		expected       string
		shouldNotMatch bool
	}
	tests := []test{
		{name: "happy . for local", in: &Config{Repository: &Repository{RelativePath: "."}}, expected: "repository:\n    path: .\n", shouldNotMatch: false},
		{name: "sad . for local", in: &Config{Repository: &Repository{RelativePath: "."}}, expected: "repository:\n    path:.\n", shouldNotMatch: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			err := WriteLocalConfig(tt.in, b)
			assert.NoError(t, err, "No test in this table should generate errors")
			if tt.shouldNotMatch {
				assert.NotEqual(t, string(b.Bytes()), tt.expected)
			} else {
				assert.Equal(t, string(b.Bytes()), tt.expected)
			}
		})
	}
}
