package factory

import (
	"bitbucket.org/unchain/ares/pkg/hello"
	"bitbucket.org/unchain/ares/pkg/http"
)

func (f *Factory) HelloService() *hello.Service {
	return hello.NewService(f.Metadata())
}

func (f *Factory) HelloHandler() *http.HelloHandler {
	return http.NewHelloHandler(f.HelloService(), f.Logger())
}
