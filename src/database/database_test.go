package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRocksdb(t *testing.T) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	t.Logf("db path:%v", dir)
	db, err := NewDB(100, dir)
	if nil != err {
		t.Errorf("NewDB err:%v", err)
		return
	}

	srvIDMap := make(map[uint64]uint64)
	orderIDMap := make(map[uint64]uint64)

	srvIDs, orderIDs, err := db.GenSingleMessageID(1, 2, 100)
	if nil != err {
		t.Errorf("GenSingleMessageID err:%v", err)
		return
	}
	for i, v := range srvIDs {
		_, ok := srvIDMap[v]
		if true == ok {
			t.Errorf("GenSingleMessageID err id:%v repeated srv ids:%v", v, srvIDs)
			return
		}
		srvIDMap[v] = v

		_, ok = orderIDMap[orderIDs[i]]
		if true == ok {
			t.Errorf("GenSingleMessageID err id repeated order ids:%v", orderIDMap)
			return
		}
		orderIDMap[orderIDs[i]] = orderIDs[i]
	}
	t.Logf("  srvids 1:%v", srvIDs)
	t.Logf("orderids 1:%v", orderIDs)

	srvIDs, orderIDs, err = db.GenSingleMessageID(1, 2, 100)
	if nil != err {
		t.Errorf("GenSingleMessageID2 err:%v", err)
		return
	}
	t.Logf("  srvids 2:%v", srvIDs)
	t.Logf("orderids 2:%v", orderIDs)
	for i, v := range srvIDs {
		_, ok := srvIDMap[v]
		if true == ok {
			t.Errorf("GenSingleMessageID2 err id:%v repeated srv ids:%v", v, srvIDs)
			return
		}
		srvIDMap[v] = v

		_, ok = orderIDMap[orderIDs[i]]
		if true == ok {
			t.Errorf("GenSingleMessageID2 err id repeated order ids:%v", orderIDMap)
			return
		}
		orderIDMap[orderIDs[i]] = orderIDs[i]
	}
}
