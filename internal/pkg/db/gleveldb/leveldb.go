package gleveldb

import (
	"cloud-batch/configs"
	"cloud-batch/internal/pkg/logging"
	"fmt"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var db *leveldb.DB

//打开数据库
func init() {
	//数据存储路径和一些初始文件
	if db == nil {
		var err error
		db, err = leveldb.OpenFile(fmt.Sprintf("%s/%s", configs.Server.RuntimeRootPath, configs.Server.LevelDBPath), nil)
		if err != nil {
			logging.Logger.Fatalf("%+v", err)
		}
	}
}

func Save(key string, value string) error {

	return db.Put([]byte(key), []byte(value), nil)
}

func Get(key string) ([]byte, error) {

	value, err := db.Get([]byte(key), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return value, nil
}

func Delete(key string) error {
	err := db.Delete([]byte(key), nil)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func GetIterator() iterator.Iterator {
	return db.NewIterator(nil, nil)
}
