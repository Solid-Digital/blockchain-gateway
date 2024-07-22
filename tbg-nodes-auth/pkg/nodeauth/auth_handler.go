package nodeauth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/unchain/tbg-nodes-auth/pkg/xorm"

	"github.com/unchainio/pkg/iferr"

	"github.com/unchainio/pkg/errors"
)

// Parse urls like
// http://host/v0/ws/{networkUUID}/*
// http://host/v0/{networkUUID}/*
func Parse(rawurl string) (map[string]string, error) {
	// TODO do this smarter and more configurable, currently the meanings of the path segments are hard-coded.
	params := make(map[string]string)
	params["protocol"] = "http"

	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, errors.Wrap(err, "invalid url")
	}

	segments := strings.Split(u.Path, "/")

	//apiVersion := segments[1]
	maybeProtocol := segments[2]
	nextSegment := 2

	switch maybeProtocol {
	case "ws":
		params["protocol"] = "websocket"
		nextSegment++
	default:
		params["protocol"] = "http"
	}

	networkUUID := segments[nextSegment]

	params["networkUUID"] = networkUUID

	return params, nil
}

func (s *Server) HandleAuth(w http.ResponseWriter, r *http.Request) {
	networkConfig := r.Header.Get("X-Unchain-Network-Config")
	protocol, networkUUID, err := parseURL(r.Header.Get("X-Original-Url"))

	if iferr.Respond(w, err, &iferr.ResponseOpts{Code: http.StatusBadRequest}) {
		return
	}

	user, pass, ok := r.BasicAuth()

	if !ok {
		http.Error(w, "No basic auth was provided", http.StatusBadRequest)
		return
	}

	network, _, creds, err := xorm.GetNetworkInterface(s.DB, networkUUID, protocol)
	if iferr.Respond(w, err, &iferr.ResponseOpts{Code: http.StatusInternalServerError}) {
		return
	}

	if network.NetworkConfiguration.String != networkConfig {
		http.Error(w, "Request not allowed", http.StatusForbidden)
		return
	}

	if !credsAreValid(user, pass, creds) {
		http.Error(w, "Request not allowed", http.StatusForbidden)
		return
	}

	fmt.Fprint(w, s.Meta.String())
	w.WriteHeader(200)
	fmt.Fprintf(w, "Request allowed\n")

	s.Log.Debugf("Request allowed")
}

func credsAreValid(user string, pass string, creds []*xorm.BasicAuthCreds) bool {
	for _, cred := range creds {
		if cred.Username == user && cred.Password == pass {
			return true
		}
	}

	return false
}

func parseURL(rawURL string) (protocol string, networkUUID string, err error) {
	params, err := Parse(rawURL)
	if err != nil {
		return "", "", err
	}

	return params["protocol"], params["networkUUID"], nil
}
