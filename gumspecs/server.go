package gumspecs

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type (
	HTTPServerSpecs struct {
		Host string `envconfig:"http_host"`
		Port int    `envconfig:"http_port"`
	}
)

func ReadHTTPServer() *HTTPServerSpecs {
	specs := &HTTPServerSpecs{}
	err := envconfig.Process(AppName, specs)
	if err != nil {
		log.Fatal(err)
	}

	if specs.Host == "" {
		return nil
	}

	return specs
}

func (s HTTPServerSpecs) String() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}
