package main

import (
	"bytes"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
)

func main()  {
	var db152 *leveldb.DB
	var db_issue *leveldb.DB
    var err error
	var height152, height_issue int
	db152,height152, err = read_db("./check_hash.db")
	if err != nil {
		fmt.Println("leveldb.OpenFile:", err)
		return
	}

	db_issue,height_issue, err = read_db("./check_hash.db")
	if err != nil {
		fmt.Println("leveldb.OpenFile:", err)
		return
	}
    height := height152
	if height152 > height_issue {
		height = height_issue
	}
	for h:=0;h<=height;h++ {
		val152,err := db152.Get([]byte(strconv.Itoa(h)), nil)
		if err != nil {
			fmt.Println("db152.Get error:", err)
			return
		}
		val_issue,err := db_issue.Get([]byte(strconv.Itoa(h)), nil)
		if err != nil {
			fmt.Println("db_issue.Get error:", err)
			return
		}
		if bytes.Compare(val152, val_issue) !=0 {
			fmt.Println("bytes.Compare !=0")
			return
		}
	}
}

func read_db(path string) (*leveldb.DB, int, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		fmt.Println("leveldb.OpenFile:", err)
		return nil,0, err
	}
	blockheightbs,err := db.Get([]byte("blockheight"), nil)
	if err != nil {
		fmt.Println("blockheight:", err)
		return nil, 0, nil
	}
	height,err := strconv.Atoi(string(blockheightbs))
	if err != nil {
		fmt.Println("strconv.Atoi:", err)
		return nil, 0, nil
	}
	return db, height, nil
}
