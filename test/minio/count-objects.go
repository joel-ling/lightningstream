package testminio

import (
	"context"
	"time"

	"github.com/cucumber/godog"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
)

func AddStepCountObjects(sc *godog.ScenarioContext) {
	sc.Then(`^I should count a total of (\d+) object(?:s)? in "([^"]*)"$`,
		countObjects,
	)

	return
}

func AddStepCountObjectsNoteBucketSize(sc *godog.ScenarioContext) {
	sc.Then(`^I should count a total of (\d+) object(?:s)? in "([^"]*)" `+
		`\(total size "([^"]*)" bytes\)$`,
		countObjectsNoteBucketSize,
	)

	return
}

func countObjects(ctx context.Context, numExpected int, bucket string) (
	context.Context, error,
) {
	return countObjectsNoteBucketSize(ctx, numExpected, bucket, "")
}

func countObjectsNoteBucketSize(
	ctx0 context.Context, numExpected int, bucket, sizeLabel string,
) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		timeout context.Context

		numObjects int
		totalSizeB int
	)

	timeout, _ = context.WithTimeout(ctx, time.Minute)

	for numObjects < numExpected {
		select {
		case <-timeout.Done():
			e = timeout.Err()
			if e != nil {
				return
			}

		default:
			numObjects, totalSizeB, e = measureBucketSize(timeout, bucket)
			if e != nil {
				return
			}
		}
	}

	assert.Equal(
		godog.T(ctx),
		numExpected,
		numObjects,
	)

	if sizeLabel != "" {
		ctx = context.WithValue(ctx, ctxKeyBucketSize{sizeLabel}, totalSizeB)
	}

	return
}

type ctxKeyBucketSize struct {
	Label string
}

func measureBucketSize(ctx context.Context, bucket string) (
	numObjects, totalSizeB int, e error,
) {
	var (
		object  minio.ObjectInfo
		objects <-chan minio.ObjectInfo
		recvOK  bool
	)

	objects = minioClient.ListObjects(ctx, bucket,
		minio.ListObjectsOptions{},
	)

	for {
		select {
		case <-ctx.Done():
			e = ctx.Err()
			if e != nil {
				return
			}

		case object, recvOK = <-objects:
			if !recvOK {
				return
			}

			e = object.Err
			if e != nil {
				return
			}
		}

		numObjects++

		totalSizeB += int(object.Size)
	}

	return
}
