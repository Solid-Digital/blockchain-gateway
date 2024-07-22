package pipeline_test

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"bitbucket.org/unchain/ares/pkg/testhelper/xrequire"

	"bitbucket.org/unchain/ares/pkg/pipeline"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/require"
	"github.com/unchainio/pkg/xtest"

	"bitbucket.org/unchain/ares/gen/orm"
)

var update = flag.Bool("update", false, "update golden files")

func TestGetAdapterLogs(t *testing.T) {
	var tcs = map[string]struct {
		InputPath  string
		OutputPath string
	}{
		"test1": {"logs/test1.input.json", "logs/test1.output.json"},
		"test2": {"logs/test1.input.json", "logs/test1.output.json"},
	}

	s := pipeline.NewService(
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	for id, tc := range tcs {
		first := true

		inputBytes := xtest.LoadBytes(t, tc.InputPath)

		res := new(elastic.SearchResult)
		err := json.Unmarshal(inputBytes, res)
		require.NoError(t, err, id)

		// Set the test hook for fetching logs from elastic
		pipeline.SetFetchAdapterLogsFn(func(s *pipeline.Service, orgName, pipelineName string, envName string, from, to string) (pipeline.DoFn, pipeline.ClearFn, error) {
			return func(ctx context.Context) (*elastic.SearchResult, error) {
					if !first {
						return nil, io.EOF
					}
					first = false

					return res, nil
				}, func(ctx context.Context) error {
					return nil
				}, nil
		})

		// We are testing s.GetDeploymentLogs
		actual, err := s.GetDeploymentLogs("", "", "", "", "", "")
		xrequire.NoError(t, err, id)

		actualBytes, err := json.MarshalIndent(actual, "", "  ")
		xrequire.NoError(t, err, id)

		outputPath := filepath.Join("testdata", tc.OutputPath)

		if *update {
			err := ioutil.WriteFile(outputPath, actualBytes, 0644)
			xrequire.NoError(t, err, id)
		}

		expectedBytes, err := ioutil.ReadFile(outputPath)
		xrequire.NoError(t, err, id)

		require.Equal(t, string(expectedBytes), string(actualBytes), id)
	}

	pipeline.ResetFetchAdapterLogsFn()
}

func (s *TestSuite) TestService_GetDeploymentLogs() {
	if testing.Short() {
		s.T().Skip()
	}

	org1, pipeline1, env1, deployment := s.factory.DeploymentFromService()

	cases := map[string]struct {
		Organization *orm.Organization
		Pipeline     *orm.Pipeline
		Environment  *orm.Environment
		From         string
		To           string
		Limit        string
		Success      bool
	}{
		"get deployment logs": {
			org1,
			pipeline1,
			env1,
			"now-2d",
			"now",
			"",
			true,
		},
		// this test reuses the deployment from the previous test case
		// to speed up the result
		"pipeline not from organization": {
			s.factory.Organization(true),
			pipeline1,
			env1,
			"now-2d",
			"now",
			"",
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			s.ares.Log.Printf("namespace: %s", org1.Name)
			s.ares.Log.Printf("pod name: %s", deployment.FullName)

			pass := false
			iterations := 1

			if tc.Success {
				iterations = 30
			}

			for i := 0; i < iterations; i++ {
				s.ares.Log.Printf("fetching logs, attempt: %d", i+1)
				response, err := s.service.GetDeploymentLogs(tc.Organization.Name, tc.Pipeline.Name, tc.Environment.Name, tc.From, tc.To, tc.Limit)

				if tc.Success {
					xrequire.NoError(t, err)
					if len(response) > 0 {
						pass = true
						break
					}

					time.Sleep(1 * time.Second)
				} else {
					xrequire.Error(t, err)
				}
			}

			if tc.Success {
				require.True(t, pass)
			}
		})
	}
}
