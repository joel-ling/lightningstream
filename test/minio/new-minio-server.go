package testminio

import (
	"context"
	"net"
	"os/exec"
	"time"

	"github.com/cucumber/godog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	MinioServerAddr = "127.153.90.30:9030"
	MinioServerCred = "minioadmin"

	minioServerPath = "/tmp/minio"
)

func AddStepNewMinioServer(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new Minio server$`,
		newMinioServer,
	)

	return
}

func newMinioServer(ctx0 context.Context) (ctx context.Context, e error) {
	ctx = ctx0

	var (
		timeout context.Context

		server *exec.Cmd = exec.Command(minioServerPath,
			"server",
			"--address", MinioServerAddr,
			ctx.Value(ctxKeyTempDir{}).(string),
		)
	)

	e = server.Start()
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeySrvProc{}, server.Process)

	timeout, _ = context.WithTimeout(ctx, time.Second)

	for {
		select {
		case <-timeout.Done():
			return

		default:
			_, e = net.Dial("tcp", MinioServerAddr)
			if e == nil {
				return
			}
		}
	}

	return
}

type ctxKeySrvProc struct{}

func newMinioClient() (*minio.Client, error) {
	return minio.New(MinioServerAddr,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				MinioServerCred,
				MinioServerCred,
				"",
			),
		},
	)
}
