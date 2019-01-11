package alienvault

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestJobPluginValidation(t *testing.T) {
	if len(plugins) == 0 {
		t.Error("No valid plugins declared")
	}

	warnings, errors := validateJobPlugin(plugins[len(plugins)-1], "plugin")
	assert.Equal(t, 0, len(warnings))
	assert.Equal(t, 0, len(errors))

	warnings, errors = validateJobPlugin("This plugin does not exist", "plugin")
	assert.Equal(t, 0, len(warnings))
	require.Equal(t, 1, len(errors))
}

func TestJobSourceValidation(t *testing.T) {

	var flagtests = []struct {
		in    string
		valid bool
	}{
		{"raw", true},
		{"syslog", true},
		{"raw2", false},
		{"", false},
		{"invalid", false},
	}

	for _, tt := range flagtests {
		t.Run(tt.in, func(t *testing.T) {
			_, errors := validateJobSourceFormat(tt.in, "source")
			assert.Equal(t, tt.valid, len(errors) == 0)
		})
	}

}

func TestIPValidation(t *testing.T) {

	var flagtests = []struct {
		in    string
		valid bool
	}{
		{"1.2.3.4", true},
		{"199.199.199.199", true},
		{"", false},
		{"something", false},
		{"1.1.1", false},
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
	}

	for _, tt := range flagtests {
		t.Run(tt.in, func(t *testing.T) {
			_, errors := validateIP(tt.in, "ip")
			assert.Equal(t, tt.valid, len(errors) == 0)
		})
	}
}
