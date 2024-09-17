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

func AddStepDelNativeRecord(sc *godog.ScenarioContext) {
	sc.When(`^in the transaction I mark LS-native record "([^"]*)" `+
		`in DB "([^"]*)" as deleted$`,
		delNativeRecord,
	)

	return
}

func delNativeRecord(ctx0 context.Context, key, dbName string) (
	ctx context.Context, e error,
) {
	ctx = ctx0

	var (
		buf *bytes.Buffer = bytes.NewBuffer(
			make([]byte, 0, lsNativeHeaderLen),
		)

		dbi lmdb.DBI

		txn *lmdb.Txn = testlmdb.GetLMDBTxn(ctx)
	)

	dbi, e = txn.OpenDBI(dbName, 0)
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

	_, e = buf.Write(
		[]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
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
