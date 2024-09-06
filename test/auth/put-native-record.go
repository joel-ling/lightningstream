package testauth

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"

	"github.com/PowerDNS/lmdb-go/lmdb"
	"github.com/cucumber/godog"

	"github.com/PowerDNS/lightningstream/test/lmdb"
)

func AddStepPutNativeRecord(sc *godog.ScenarioContext) {
	sc.When(`^I put an LS-native record "([^"]*)" "([^"]*)" in DB "([^"]*)" `+
		`of LMDB env\. "([^"]*)"$`,
		putNativeRecord,
	)

	return
}

func putNativeRecord(ctx0 context.Context, key, val, dbName, envName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		buf *bytes.Buffer = bytes.NewBuffer(
			make([]byte, 0,
				lsNativeHeaderLen+len(val),
			),
		)

		put = func(txn *lmdb.Txn) (err error) {
			var (
				dbi lmdb.DBI
			)

			dbi, err = txn.OpenDBI(dbName, lmdb.Create)
			if err != nil {
				return
			}

			err = binary.Write(buf, binary.BigEndian,
				uint64(time.Now().UnixNano()),
			)
			if err != nil {
				return
			}

			err = binary.Write(buf, binary.BigEndian,
				uint64(txn.ID()),
			)
			if err != nil {
				return
			}

			err = binary.Write(buf, binary.BigEndian,
				uint64(0),
			)
			if err != nil {
				return
			}

			_, err = buf.Write(
				[]byte(val),
			)
			if err != nil {
				return
			}

			err = txn.Put(dbi,
				[]byte(key),
				buf.Bytes(),
				0,
			)
			if err != nil {
				return
			}

			return
		}
	)

	e = testlmdb.UpdateLMDBEnv(ctx, envName, put)
	if e != nil {
		return
	}

	return
}

const (
	lsNativeHeaderLen = 24
)
