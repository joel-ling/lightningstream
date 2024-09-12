package testminio

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/minio/minio-go/v7"
)

func AddStepNewBucket(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new bucket "([^"]*)"$`,
		newBucket,
	)

	return
}

func newBucket(ctx0 context.Context, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	e = minioClient.MakeBucket(ctx, name,
		minio.MakeBucketOptions{},
	)
	if e != nil {
		return
	}

	return
}
