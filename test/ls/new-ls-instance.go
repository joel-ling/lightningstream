package testls

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cucumber/godog"

	"github.com/PowerDNS/lightningstream/test/lmdb"
	"github.com/PowerDNS/lightningstream/test/minio"
)

func AddStepNewLSInstance(sc *godog.ScenarioContext) {
	sc.Given(`^there is a new LS instance syncing "([^"]*)" to "([^"]*)"$`,
		newLSInstance,
	)

	sc.When(`^I start a new LS instance syncing "([^"]*)" to "([^"]*)"$`,
		newLSInstance,
	)

	return
}

func newLSInstance(ctx0 context.Context, envName, bucketName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		binPath string = filepath.Join(
			ctx.Value(ctxKeyTempDir{}).(string),
			"lightningstream",
		)

		builder *exec.Cmd = exec.Command("go", "build",
			"-o", binPath,
			"../cmd/lightningstream",
		)

		cfgPath string
		command *exec.Cmd
	)

	e = builder.Run()
	if e != nil {
		return
	}

	cfgPath, e = configureLS(ctx, envName, bucketName)
	if e != nil {
		return
	}

	command = exec.Command(binPath, "sync",
		"--config", cfgPath,
	)

	//command.Stderr = os.Stderr

	e = command.Start()
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyLSProcs{},
		append(
			ctx.Value(ctxKeyLSProcs{}).([]*os.Process),
			command.Process,
		),
	)

	return
}

func configureLS(ctx context.Context, envName, bucketName string) (
	cfgPath string, e error,
) {
	var (
		cfgFile *os.File
		cfgYaml string
		envPath string
	)

	envPath, e = testlmdb.GetPathToLMDBEnv(ctx, envName)
	if e != nil {
		return
	}

	cfgYaml = fmt.Sprintf(cfgTmpl,
		envName,
		envPath,
		testminio.MinioServerCred,
		testminio.MinioServerCred,
		bucketName,
		testminio.MinioServerAddr,
	)

	cfgFile, e = os.CreateTemp(
		ctx.Value(ctxKeyTempDir{}).(string),
		"lightningstream.yaml",
	)
	if e != nil {
		return
	}

	_, e = cfgFile.Write(
		[]byte(cfgYaml),
	)
	if e != nil {
		return
	}

	cfgPath = cfgFile.Name()

	return
}

//go:embed config.yaml
var cfgTmpl string
