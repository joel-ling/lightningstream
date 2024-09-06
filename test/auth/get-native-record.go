package testauth

import (
	"context"
	"time"

	"github.com/PowerDNS/lmdb-go/lmdb"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"

	"github.com/PowerDNS/lightningstream/test/lmdb"
)

func AddStepGetNativeRecord(sc *godog.ScenarioContext) {
	sc.Then(`^I should get LS-native record "([^"]*)" "([^"]*)" `+
		`in DB "([^"]*)" of env\. "([^"]*)"$`,
		getNativeRecord,
	)

	return
}

func getNativeRecord(
	ctx0 context.Context, key, valExpect, dbName, envName string,
) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		timeout   context.Context
		valActual []byte

		get = func(txn *lmdb.Txn) (err error) {
			var (
				dbi lmdb.DBI
			)

			dbi, err = txn.OpenDBI(dbName, 0)
			if err != nil {
				return
			}

			valActual, err = txn.Get(dbi,
				[]byte(key),
			)
			if err != nil {
				return
			}

			return
		}
	)

	timeout, _ = context.WithTimeout(ctx, time.Minute)

loop:
	for {
		select {
		case <-timeout.Done():
			return

		default:
			e = testlmdb.ViewLMDBEnv(ctx, envName, get)
			if e == nil {
				break loop
			}
		}
	}

	assert.Equal(
		godog.T(ctx),
		[]byte(valExpect),
		valActual[lsNativeHeaderLen:],
	)

	return
}
