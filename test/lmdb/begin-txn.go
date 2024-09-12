package testlmdb

import (
	"context"

	"github.com/PowerDNS/lmdb-go/lmdb"
	"github.com/cucumber/godog"
)

func AddStepBeginTxn(sc *godog.ScenarioContext) {
	sc.When(`^I begin a transaction in "([^"]*)"$`,
		beginTxn,
	)

	return
}

func beginTxn(ctx0 context.Context, envName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		env *lmdb.Env = ctx.Value(ctxKeyLMDBEnv{envName}).(*lmdb.Env)
		txn *lmdb.Txn
	)

	txn, e = env.BeginTxn(nil, 0)
	if e != nil {
		return
	}

	ctx = context.WithValue(ctx, ctxKeyLMDBTxn{}, txn)

	return
}

type ctxKeyLMDBTxn struct{}

func GetLMDBTxn(ctx context.Context) *lmdb.Txn {
	return ctx.Value(ctxKeyLMDBTxn{}).(*lmdb.Txn)
}
