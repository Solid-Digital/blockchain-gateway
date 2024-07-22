package internal_test

import (
	"encoding/json"
	"testing"

	"github.com/unchainio/pkg/xtest"

	"github.com/olivere/elastic"

	"bitbucket.org/unchain/ares/pkg/pipeline/internal"
	"github.com/unchainio/pkg/xconfig"

	"github.com/stretchr/testify/require"
)

func getHit(t *testing.T) *elastic.SearchHit {
	hit := new(elastic.SearchHit)
	err := json.Unmarshal(xtest.LoadBytes(t, "search_hit.json"), hit)

	require.NoError(t, err)
	require.NotNil(t, hit)

	return hit
}

func TestName(t *testing.T) {
	hit := getHit(t)

	p, err := internal.NewHitParser(hit)
	require.NoError(t, err)
	require.NotNil(t, p)

	timestamp, err := p.Timestamp()

	require.NoError(t, err)
	require.Equal(t, int64(1537798568000), timestamp)

	containerID, err := p.ContainerID()
	require.NoError(t, err)
	require.Equal(t, "f437f685e25be518fb91bebf75b3780d32f748b656cb1534a70a4876c174334e", containerID)

	log, err := p.Log()
	require.NoError(t, err)
	require.Equal(t, " [fabsdk/msp] 2018/09/24 14:16:08 UTC - msp.NewIdentityManager -> WARN Cryptopath not provided for organization [ordererorg], MSP stores not created\n", log)
}

func TestNewHitParser(t *testing.T) {
	var tcs map[string]struct {
		HitJSON string

		Timestamp  int64
		InstanceID string
		Text       string
		Level      string
		Caller     string
		Function   string
	}

	err := xconfig.Load(&tcs, xconfig.FromPaths("testdata/hitparser.toml"))
	require.NoError(t, err)

	for id, tc := range tcs {
		hit := new(elastic.SearchHit)
		err := json.Unmarshal([]byte(tc.HitJSON), hit)

		require.NoError(t, err, id)
		require.NotNil(t, hit, id)

		p, err := internal.NewHitParser(hit)
		require.NoError(t, err, id)
		require.NotNil(t, p, id)

		ll, err := p.Parse()
		require.NoError(t, err, id)

		require.Equal(t, tc.Timestamp, ll.Timestamp, id)
		require.Equal(t, tc.InstanceID, ll.InstanceID, id)
		require.Equal(t, tc.Text, ll.Text)
		require.Equal(t, tc.Level, ll.Level)
		require.Equal(t, tc.Caller, ll.Caller)
		require.Equal(t, tc.Function, ll.Function)
	}
}
