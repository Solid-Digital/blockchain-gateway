package internal

import (
	"encoding/json"

	"bitbucket.org/unchain/ares/gen/dto"

	"github.com/jmoiron/jsonq"

	"github.com/olivere/elastic"
	"github.com/unchainio/pkg/errors"
)

type HitParser struct {
	hit    *elastic.SearchHit
	source *jsonq.JsonQuery
	fields *jsonq.JsonQuery
	sort   *jsonq.JsonQuery

	timestamp   int64
	containerID string
	log         string
	logLevel    string

	message *LogMessage
}

func NewHitParser(hit *elastic.SearchHit) (*HitParser, error) {
	sourceMap := make(map[string]interface{})
	err := json.Unmarshal(*hit.Source, &sourceMap)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal source")
	}

	source := jsonq.NewQuery(sourceMap)
	fields := jsonq.NewQuery(hit.Fields)
	sort := jsonq.NewQuery(map[string]interface{}{"sort": hit.Sort})

	return &HitParser{
		hit:    hit,
		source: source,
		fields: fields,
		sort:   sort,
	}, nil
}

//func (p *HitParser) Timestamp() (int64, error) {
//	if p.timestamp != 0 {
//		return p.timestamp, nil
//	}
//
//	return p.parseTimestamp()
//}

func (p *HitParser) Timestamp() (int64, error) {
	var err error

	if p.timestamp == 0 {
		p.timestamp, err = p.parseTimestamp()
	}

	return p.timestamp, err
}

func (p *HitParser) parseTimestamp() (int64, error) {
	timestampFloat, err := p.sort.ArrayOfFloats("sort")

	if err != nil {
		return 0, errors.Wrap(err, "")
	}

	if len(timestampFloat) == 0 {
		return 0, errors.New("unable to parse timestamp")
	}

	return int64(timestampFloat[0]), nil
}

func (p *HitParser) ContainerID() (string, error) {
	var err error

	if p.containerID == "" {
		p.containerID, err = p.parseContainerID()
	}

	return p.containerID, err
}

func (p *HitParser) parseContainerID() (string, error) {
	containerID, err := p.source.String("docker", "container_id")

	if err != nil {
		return "", errors.Wrap(err, "")
	}

	return containerID, nil
}

func (p *HitParser) Log() (string, error) {
	var err error
	lm, err := p.LogMessage()

	if err != nil {
		return "", err
	}

	return lm.Text, nil
}

func (p *HitParser) LogMessage() (*LogMessage, error) {
	var err error

	if p.message == nil {
		p.message, err = p.parseLogMessage()
	}

	return p.message, err
}

func (p *HitParser) parseLogMessage() (*LogMessage, error) {
	log, err := p.source.String("log")

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	ilm := new(LogMessage)
	err = json.Unmarshal([]byte(log), ilm)

	if err != nil {
		ilm = &LogMessage{
			Text:     log,
			Level:    "unknown",
			Caller:   "",
			Time:     "",
			Function: "",
		}
	}

	return ilm, nil
}

func (p *HitParser) Parse() (*dto.LogLine, error) {
	timestamp, err := p.Timestamp()

	if err != nil {
		return nil, err
	}

	containerID, err := p.ContainerID()

	if err != nil {
		return nil, err
	}

	lm, err := p.LogMessage()

	if err != nil {
		return nil, err
	}

	return &dto.LogLine{
		Caller:     lm.Caller,
		Function:   lm.Function,
		InstanceID: containerID,
		Level:      lm.Level,
		Text:       lm.Text,
		Time:       lm.Time,
		Timestamp:  timestamp,
	}, nil
}
