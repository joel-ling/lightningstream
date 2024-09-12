package testminio

import (
	"context"
	"time"

	"github.com/cucumber/godog"
	"github.com/minio/minio-go/v7"
)

func AddStepCopyIthObject(sc *godog.ScenarioContext) {
	sc.When(`^I copy object i=(\d+) from "([^"]*)" to "([^"]*)"$`,
		copyIthObject,
	)

	return
}

func copyIthObject(ctx0 context.Context, i int, source, target string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		j       int
		object  minio.ObjectInfo
		objects <-chan minio.ObjectInfo
		timeout context.Context
	)

	objects = minioClient.ListObjects(ctx, source,
		minio.ListObjectsOptions{},
	)

	timeout, _ = context.WithTimeout(ctx, time.Second)

	for j = 0; j <= i; j++ {
		select {
		case <-timeout.Done():
			e = timeout.Err()
			if e != nil {
				return
			}

		case object = <-objects:
			e = object.Err
			if e != nil {
				return
			}
		}
	}

	_, e = minioClient.CopyObject(ctx,
		minio.CopyDestOptions{
			Bucket: target,
			Object: object.Key,
		},
		minio.CopySrcOptions{
			Bucket: source,
			Object: object.Key,
		},
	)
	if e != nil {
		return
	}

	return
}
