package s3

import (
	"io"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/google/wire"
	"github.com/minio/minio-go"
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/errors"
)

var FileStoreSet = wire.NewSet(NewFileStore, wire.Bind(new(ares.FileStore), new(FileStore)))

type FileStore struct {
	client *minio.Client
	log    logger.Logger
	cfg    *Config
}

func NewFileStore(log logger.Logger, cfg *Config) (*FileStore, error) {
	client, err := minio.New(cfg.URL, cfg.AccessKey, cfg.SecretAccessKey, cfg.SSL)

	if err != nil {
		return nil, errors.Wrap(err, "failed to establish minio connection")
	}

	return &FileStore{
		client: client,
		log:    log,
		cfg:    cfg,
	}, nil
}

func (ps FileStore) GetObject(id string) (io.ReadCloser, error) {
	r, err := ps.client.GetObject(ps.cfg.BucketName, id, minio.GetObjectOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return r, nil
}

func (ps FileStore) PutObject(id string, reader io.Reader) error {
	_, err := ps.client.PutObject(ps.cfg.BucketName, id, reader, -1, minio.PutObjectOptions{ContentType: "tar"})

	if err != nil {
		return errors.Wrap(err, "")
	}

	return err
}
