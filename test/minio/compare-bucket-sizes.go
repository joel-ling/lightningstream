package testminio

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

func AddStepCompareBucketSizes(sc *godog.ScenarioContext) {
	sc.Then(
		`^I should see that bucket size "([^"]*)" is greater than "([^"]*)"$`,
		compareBucketSizes,
	)

	return
}

func compareBucketSizes(ctx0 context.Context, sizeLabelL, sizeLabelS string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	assert.Greater(
		godog.T(ctx),
		ctx.Value(ctxKeyBucketSize{sizeLabelL}).(int),
		ctx.Value(ctxKeyBucketSize{sizeLabelS}).(int),
	)

	return
}
