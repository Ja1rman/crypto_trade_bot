package analyzer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	
	"testing"
)

func TestCalculateMarketSize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CalculateMarketSize Suite")
}
