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
	sc.When(
		`^in the transaction I put an LS-native record "([^"]*)" "([^"]*)" `+
			`in DB "([^"]*)"$`,
		putNativeRecord,
	)

	return
}

func putNativeRecord(ctx0 context.Context, key, val, dbName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		buf *bytes.Buffer = bytes.NewBuffer(
			make([]byte, 0,
				lsNativeHeaderLen+len(val),
			),
		)

		dbi lmdb.DBI

		txn *lmdb.Txn = testlmdb.GetLMDBTxn(ctx)
	)

	dbi, e = txn.OpenDBI(dbName, lmdb.Create)
	if e != nil {
		return
	}

	e = binary.Write(buf, binary.BigEndian,
		uint64(time.Now().UnixNano()),
	)
	if e != nil {
		return
	}

	e = binary.Write(buf, binary.BigEndian,
		uint64(txn.ID()),
	)
	if e != nil {
		return
	}

	e = binary.Write(buf, binary.BigEndian,
		uint64(0),
	)
	if e != nil {
		return
	}

	_, e = buf.Write(
		[]byte(val),
	)
	if e != nil {
		return
	}

	e = txn.Put(dbi,
		[]byte(key),
		buf.Bytes(),
		0,
	)
	if e != nil {
		return
	}

	return
}

const (
	lsNativeHeaderLen = 24
)
