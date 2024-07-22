package ares_test

import (
	"testing"

	"github.com/go-openapi/strfmt"

	"bitbucket.org/unchain/ares/gen/client"
	httptransport "github.com/go-openapi/runtime/client"
)

func TestTransport(t *testing.T) {
	cli := client.New(httptransport.New("", "", nil), strfmt.Default)

	_, _ = cli.Component.CreateActionVersion(nil, httptransport.BearerToken("sdasd"))
}
