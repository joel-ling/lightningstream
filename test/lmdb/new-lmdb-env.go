package testlmdb

import (
	"context"
	"path/filepath"

	"github.com/PowerDNS/lmdb-go/lmdb"
	"github.com/cucumber/godog"
)

func AddStepNewLMDBEnv(sc *godog.ScenarioContext) {
	sc.Given(
		`^there is a new LMDB environment "([^"]*)" with (\d+) DBs at most$`,
		newLMDBEnv,
	)

	return
}

func newLMDBEnv(ctx0 context.Context, name string, maxDBs int) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		env *lmdb.Env

		path string = filepath.Join(
			ctx.Value(ctxKeyTempDir{}).(string),
			name,
		)
	)

	env, e = lmdb.NewEnv()
	if e != nil {
		return
	}

	e = env.SetMaxDBs(maxDBs)
	if e != nil {
		return
	}

	e = env.Open(path, lmdb.NoSubdir, 0644)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyLMDBEnv{name}, env)

	return
}

type ctxKeyLMDBEnv struct {
	Name string
}

func GetPathToLMDBEnv(ctx context.Context, name string) (string, error) {
	return ctx.Value(ctxKeyLMDBEnv{name}).(*lmdb.Env).Path()
}

func UpdateLMDBEnv(ctx context.Context, name string, txnOp lmdb.TxnOp) error {
	return ctx.Value(ctxKeyLMDBEnv{name}).(*lmdb.Env).Update(txnOp)
}

func ViewLMDBEnv(ctx context.Context, name string, txnOp lmdb.TxnOp) error {
	return ctx.Value(ctxKeyLMDBEnv{name}).(*lmdb.Env).View(txnOp)
}
