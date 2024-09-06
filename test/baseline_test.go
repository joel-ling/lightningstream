package test

import (
	"testing"

	"github.com/cucumber/godog"

	"github.com/PowerDNS/lightningstream/test/auth"
	"github.com/PowerDNS/lightningstream/test/lmdb"
	"github.com/PowerDNS/lightningstream/test/ls"
	"github.com/PowerDNS/lightningstream/test/minio"
)

func TestBaselineFeatures(t *testing.T) {
	var (
		scenarioInitializer = func(sc *godog.ScenarioContext) {
			testlmdb.AddStepSetUp(sc)
			testlmdb.AddStepNewLMDBEnv(sc)
			testlmdb.AddStepCleanUp(sc)

			testminio.AddStepSetUp(sc)
			testminio.AddStepNewMinioServer(sc)
			testminio.AddStepCreateBucket(sc)
			testminio.AddStepCleanUp(sc)

			testls.AddStepSetUp(sc)
			testls.AddStepNewLSInstance(sc)
			testls.AddStepCleanUp(sc)

			testauth.AddStepPutNativeRecord(sc)
			testauth.AddStepGetNativeRecord(sc)

			return
		}

		options = &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/baseline.feature"},
			TestingT: t,
		}

		suite = godog.TestSuite{
			ScenarioInitializer: scenarioInitializer,
			Options:             options,
		}
	)

	if suite.Run() != 0 {
		t.Fatal()
	}

	return
}
