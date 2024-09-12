package testlmdb

import (
	"context"

	"github.com/PowerDNS/lmdb-go/lmdb"
	"github.com/cucumber/godog"
)

func AddStepCommitTxn(sc *godog.ScenarioContext) {
	sc.When(`^I commit the transaction$`,
		commitTxn,
	)

	return
}

func commitTxn(ctx0 context.Context) (ctx context.Context, e error) {
	ctx = ctx0

	var (
		txn *lmdb.Txn = GetLMDBTxn(ctx)
	)

	e = txn.Commit()
	if e != nil {
		return
	}

	return
}
