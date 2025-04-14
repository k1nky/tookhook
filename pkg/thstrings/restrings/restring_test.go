package restrings

import (
	"testing"

	"github.com/k1nky/tookhook/pkg/thstrings"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestString_Compile(t *testing.T) {
	tests := []struct {
		name      string
		s         *String
		wantError error
	}{
		{
			name: "Empty",
			s: &String{
				RegExp: "",
			},
		},
		{
			name: "Invalid RegExp value",
			s: &String{
				RegExp: ")",
			},
			wantError: thstrings.ErrCompile,
		},
		{
			name: "Valid RegExp value",
			s: &String{
				RegExp: ".*",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Compile()
			assert.ErrorIs(t, got, tt.wantError)
		})
	}

}

func TestString_MarshalYAML(t *testing.T) {
	s := String{RegExp: "abc"}
	got, err := yaml.Marshal(s)
	assert.NoError(t, err)
	assert.Equal(t, []byte("abc\n"), got)
}

func TestString_UnmarshalYAML(t *testing.T) {
	s := String{}
	err := yaml.Unmarshal([]byte("abc"), &s)
	assert.NoError(t, err)
	assert.Equal(t, String{RegExp: "abc"}, s)
}
