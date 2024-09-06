package testls

import (
	"context"
	"os"

	"github.com/cucumber/godog"
)

func AddStepCleanUp(sc *godog.ScenarioContext) {
	sc.After(cleanUp)

	return
}

func cleanUp(ctx0 context.Context, scenario *godog.Scenario, e0 error) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		process *os.Process
	)

	for _, process = range ctx.Value(ctxKeyLSProcs{}).([]*os.Process) {
		e = process.Kill()
		if e != nil {
			return
		}
	}

	e = os.RemoveAll(
		ctx.Value(ctxKeyTempDir{}).(string),
	)
	if e != nil {
		return
	}

	return
}
