package casbin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyMatchDomain(t *testing.T) {
	for id, tt := range map[string]struct {
		RequestObject string
		PolicyObject  string
		Expected      bool
	}{
		"1": {"/api/v1/org/demo2/adapter-component-template", "/api/v1/org/:domain/*", true},
		"2": {"/api/v1/org/demo2/adapter-component-template", "/api/v1/org/:domain/adapter-component-template", true},
		"3": {"/api/v1/org/demo2/adapter-component-template", "/api/v1/org/:domain/adapters", false},
		"4": {"/api/v1/org/demo2/adapter-component-template", "/api/v1/org/*", true},
	} {
		require.Equal(t, tt.Expected, KeyMatch2(tt.RequestObject, tt.PolicyObject), id)
	}
}

func TestParseDomain(t *testing.T) {
	for id, tt := range map[string]struct {
		RequestObject string
		PolicyObject  string
		Expected      string
	}{
		"1": {"/api/v1/org/demo2/adapter-component-template", "/api/v1/org/:domain/*", "demo2"},
		"2": {"/api/v1/org/demo2/adapter-component-template", "/api/v1/org/:domain/adapter-component-template", "demo2"},
		"3": {"/api/v1/org/demo2/adapter-component-template", "/api/v1/org/:domain/adapter", "*"},
		"4": {"/api/v1/org/demo2/adapter-component-template", "/api/v1/org/adapters", "*"},
		"5": {"/api/bootstrap/adapter-component-template", "/api/:domain/*", "bootstrap"},
		"6": {"/bla/bootstrap/adapter-component-template", "/api/:domain/*", "*"},
	} {
		require.Equal(t, tt.Expected, ParseDomain(tt.RequestObject, tt.PolicyObject), id)
	}
}
