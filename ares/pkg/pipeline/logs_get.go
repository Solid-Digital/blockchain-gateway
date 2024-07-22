package pipeline

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"

	"bitbucket.org/unchain/ares/pkg/3p/apperr"

	"bitbucket.org/unchain/ares/pkg/ares"
	"bitbucket.org/unchain/ares/pkg/xorm"

	"bitbucket.org/unchain/ares/gen/dto"

	"github.com/unchainio/pkg/errors"

	"github.com/olivere/elastic"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/pipeline/internal"
)

func (s *Service) GetDeploymentLogs(orgName string, pipelineName string, envName string, from string, to string, limit string) ([]*dto.LogLine, *apperr.Error) {
	// limit is not used currently, just keep it here in case we need it later on
	fetchLogs, clearFn, err := fetchDeploymentLogsFn(s, orgName, pipelineName, envName, from, to)

	if err != nil {
		return nil, apperr.Internal.Wrap(err)
	}

	return parseAllLogs(fetchLogs, clearFn)
}

var (
	fetchDeploymentLogsFn = fetchDeploymentLogs
)

type DoFn func(ctx context.Context) (*elastic.SearchResult, error)
type ClearFn func(ctx context.Context) error

func fetchDeploymentLogs(s *Service, orgName, pipelineName string, envName string, from, to string) (DoFn, ClearFn, error) {
	var deployment *orm.Deployment

	appErr := ares.WrapTx(s.db, func(ctx context.Context, tx *sql.Tx) *apperr.Error {
		var appErr *apperr.Error

		org, pipeline, appErr := xorm.GetPpelineTx(ctx, tx, orgName, pipelineName)
		if appErr != nil {
			return appErr
		}

		env, appErr := xorm.GetOrgEnvironmentTx(ctx, tx, org, envName)
		if appErr != nil {
			return appErr
		}

		deployment, appErr = xorm.GetDeploymentTx(ctx, tx, org, pipeline, env)
		if appErr != nil {
			return appErr
		}
		return nil
	})

	if appErr != nil {
		return nil, nil, appErr
	}

	bq := elastic.NewBoolQuery()
	bq = bq.Must(elastic.NewTermQuery("kubernetes.labels.fullName.keyword", deployment.FullName))
	bq = bq.Must(elastic.NewRangeQuery("@timestamp").Gt(from).Lt(to))
	ss := elastic.NewSearchSource().DocvalueFields("@timestamp").Query(bq).Size(1000).Sort("@timestamp", true)

	scrollSvc := s.service.elastic.Scroll(s.cfg.Index).SearchSource(ss)
	//fmt.Printf(marshalSource(ss.Source()))
	return scrollSvc.Do, scrollSvc.Clear, nil
}

func marshalSource(i interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(i)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func parseAllLogs(getLogs DoFn, clearScrollFn ClearFn) ([]*dto.LogLine, *apperr.Error) {
	logs := make([]*dto.LogLine, 0)

	offset := 0

	for {
		var res *elastic.SearchResult
		var err error

		res, err = getLogs(context.Background())

		// no more logs, return
		if err == io.EOF {
			err2 := clearScrollFn(context.Background())
			if err2 != nil {
				return nil, apperr.Internal.Wrap(err2)
			}

			break
		} else if err != nil {
			return nil, apperr.Internal.Wrap(errors.Wrap(err, ""))
		}

		// Allocate the correct number of log lines after the first call to scroll.Do()
		if offset == 0 {
			logs = make([]*dto.LogLine, res.Hits.TotalHits)
		}

		offset, err = parseLogsBatch(offset, res, logs)

		if err != nil {
			return nil, apperr.Internal.Wrap(err)
		}
	}

	return logs, nil
}

func parseLogsBatch(offset int, res *elastic.SearchResult, logs []*dto.LogLine) (int, error) {
	for _, hit := range res.Hits.Hits {
		p, err := internal.NewHitParser(hit)

		if err != nil {
			return offset, err
		}

		ll, err := p.Parse()

		if err != nil {
			return offset, err
		}

		logs[offset] = ll
		offset++
	}

	return offset, nil
}

//
//func (s *Service) GetAdapterLogs2(adapter *models.Adapter, from, to, limit string) (interface{}, error) {
//	bq := elastic.NewBoolQuery()
//	bq = bq.Must(elastic.NewTermQuery("kubernetes.labels.adapterID.keyword", adapter.ID))
//	bq = bq.Must(elastic.NewRangeQuery("@timestamp").Gt(from).Lt(to))
//	ss := elastic.NewSearchSource().DocvalueFields("@timestamp").Query(bq).Size(1000).Sort("@timestamp", true)
//
//	return s.service.elastic.Scroll("adapters-*").SearchSource(ss).Do(context.Background())
//
//	return s.GetDeploymentLogs(adapter, from, to, limit)
//}
