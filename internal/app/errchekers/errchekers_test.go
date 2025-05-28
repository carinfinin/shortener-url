package errchekers

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestErrCheckAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ErrCheckAnalyzer, "./...")
}
