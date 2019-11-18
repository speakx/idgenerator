package database

import (
	"environment/cfgargs"
	"environment/logger"
	"environment/rocksdbimp"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// DB 数据库+缓存
type DB struct {
	db         *rocksdbimp.RocksdbImp
	instanceid int
	timeStamp  uint64
	incrID     uint32
	sync.Mutex
}

// NewDB create
func NewDB(cfg *cfgargs.SrvConfig) (*DB, error) {
	db := &DB{
		instanceid: cfg.Info.ID,
		db:         rocksdbimp.NewRocksdbImp(),
	}
	return db, db.init(cfg.DB.Path)
}

func (d *DB) init(dbPath string) error {
	return d.db.OpenDB(dbPath)
}

func (d *DB) genIncrID(ts uint64) uint32 {
	if d.timeStamp != ts {
		d.incrID = 0
		d.timeStamp = ts
		return d.incrID
	}

	d.incrID++
	if d.incrID > 0xFFFFF {
		logger.Error("db.gen.incrid overflow ts:", d.timeStamp, " incrid:", d.incrID)
	}
	return d.incrID
}

// GenSingleMessageID srvid、orderid
func (d *DB) GenSingleMessageID(fromUID, toUID uint64, num uint32) ([]uint64, []uint64, error) {
	d.Lock()
	defer d.Unlock()

	// make conv prefix
	midID := uint64(0)
	if fromUID > toUID {
		midID = fromUID / 2
	} else {
		midID = toUID / 2
	}

	// load db
	dbOrderIDKey := fmt.Sprintf("single_%016X_orderid", midID)
	dbOrderIDVal, err := d.db.Get([]byte(dbOrderIDKey))
	if nil != err {
		logger.Error("single.message.id load key:", dbOrderIDKey, " failed, err:", err)
		return nil, nil, err
	}
	curOrderID := uint64(0)
	if len(dbOrderIDVal.Data()) > 0 {
		curOrderID, err = strconv.ParseUint(string(dbOrderIDVal.Data()), 10, 64)
		if nil != err {
			logger.Error("single.message.id load key:", dbOrderIDKey,
				" val:", string(dbOrderIDVal.Data()), " atoi failed, err:", err)
			return nil, nil, err
		}
	}

	// make id
	srvIDs := make([]uint64, num)
	orderIDs := make([]uint64, num)

	timeStamp := uint64(time.Now().Unix()) << 32 & 0xFFFFFFFF00000000
	srvID := uint32(d.instanceid<<20) & 0xFFF00000
	for i := 0; i < int(num); i++ {
		curOrderID++
		seqID := d.genIncrID(timeStamp)

		srvIDs[i] = timeStamp | uint64(srvID) | uint64(seqID)
		orderIDs[i] = curOrderID
	}

	// save
	newDBOrderIDVal := strconv.FormatUint(uint64(curOrderID), 10)
	err = d.db.Put([]byte(dbOrderIDKey), []byte(newDBOrderIDVal))
	if nil != err {
		logger.Error("single.message.id save key:", dbOrderIDKey,
			" val", newDBOrderIDVal, " failed, err:", err)
		return nil, nil, err
	}

	return srvIDs, orderIDs, nil
}
