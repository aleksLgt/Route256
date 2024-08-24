package cmd

import (
	"testing"

	"github.com/stretchr/testify/suite"

	lomssuite "route256/loms/tests/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(lomssuite.ItemS))
}
