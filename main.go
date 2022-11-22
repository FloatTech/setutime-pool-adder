package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/FloatTech/AnimeAPI/pixiv"
	sql "github.com/FloatTech/sqlite"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// Pools 图片缓冲池
type imgpool struct {
	db   *sql.Sqlite
	dbmu sync.RWMutex
	path string
	max  int
	pool map[string][]*message.MessageSegment
}

var pool = &imgpool{
	db:   &sql.Sqlite{},
	path: pixiv.CacheDir,
	max:  10,
	pool: make(map[string][]*message.MessageSegment),
}

func (p *imgpool) add(imgtype string, id int64) error {
	p.dbmu.Lock()
	defer p.dbmu.Unlock()
	if err := p.db.Create(imgtype, &pixiv.Illust{}); err != nil {
		return err
	}
	// 查询P站插图信息
	illust, err := pixiv.Works(id)
	if err != nil {
		return err
	}
	// 添加插画到对应的数据库table
	if err := p.db.Insert(imgtype, illust); err != nil {
		return err
	}
	return nil
}

func main() {
	// 如果数据库不存在则下载
	pool.db.DBPath = "./SetuTime.db"
	err := pool.db.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	defer pool.db.Close()
	if err := pool.db.Create("缓存池", &pixiv.Illust{}); err != nil {
		panic(err)
	}
	f, err := os.Open("pid_in_pool.txt")
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(f)
	i := 1
	for s.Scan() {
		id, err := strconv.Atoi(s.Text())
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		pool.add("缓存池", int64(id))
		fmt.Print("\r[", i, "] add: ", id, "                   ")
		i++
	}
	println("\ncomplete")
}
