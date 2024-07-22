package ethereum_client

type Config struct {
	Host            string
	SynchronousMode bool
	GasPrice        *int64
	GasLimit        uint64
}
