package testminio

import (
	"context"

	"github.com/cucumber/godog"
	"github.com/minio/minio-go/v7"
)

func AddStepCreateBucket(sc *godog.ScenarioContext) {
	sc.Given(`^there is a bucket "([^"]*)"$`,
		createBucket,
	)

	return
}

func createBucket(ctx0 context.Context, name string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		client *minio.Client
	)

	client, e = newMinioClient()
	if e != nil {
		return
	}

	e = client.MakeBucket(ctx, name,
		minio.MakeBucketOptions{},
	)
	if e != nil {
		return
	}

	return
}
