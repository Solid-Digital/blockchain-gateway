package domain_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

const (
	DefaultAccount       = "0xfdfa8d41f986c80904bf4825402e788f3121e7af"
	NonRegisteredAccount = "0xfac399e49f5b6867af186390270af252e683b154"

	NonExistingContractAddress = "0xd0a6e6c54dbc68db5db3a091b171a77407ff7ccf"
)
