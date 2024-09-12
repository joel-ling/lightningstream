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

func countObjects(ctx0 context.Context, count int, bucket string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		n       int
		timeout context.Context
	)

	timeout, _ = context.WithTimeout(ctx, time.Minute)

	for n < count {
		select {
		case <-timeout.Done():
			e = timeout.Err()
			if e != nil {
				return
			}

		default:
			n, e = countObjectsOnce(timeout, bucket)
			if e != nil {
				return
			}
		}
	}

	assert.Equal(
		godog.T(ctx),
		count,
		n,
	)

	return
}

func countObjectsOnce(ctx context.Context, bucket string) (count int, e error) {
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
			e = object.Err
			if e != nil {
				return
			}

			if !recvOK {
				return
			}
		}

		count++
	}

	return
}
