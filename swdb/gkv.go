package swdb

//gkv.go
//提供操作底层gkvdb的能力
import (
	"gitee.com/johng/gkvdb/gkvdb"
)

type triasdb struct {
	db   *gkvdb.DB
	path string
}

func (t *triasdb) initswdb(path string) (bool, error) {
	db, err := gkvdb.New(path)
	printErr(err,"initswdb")
	t.db = db
	t.path = path
	return true, err
}

func (t *triasdb) set(k []byte, v []byte) bool {
	t.db.Set(k, v)
	return true
}

func (t *triasdb) get(k []byte) []byte {
	return t.db.Get(k)
}

func (t *triasdb) del(k []byte) bool {
	t.db.Remove(k)
	return true
}

func (t *triasdb) close() bool {
	t.db.Close()
	return true
}
