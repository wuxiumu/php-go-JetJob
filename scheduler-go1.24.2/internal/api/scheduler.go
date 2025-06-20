package api

import (
	"log"
	"time"
)

func StartScheduler() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			log.Println("定时扫描任务并分发...")
			// 1. 查询所有激活任务
			// 2. 判断哪些节点空闲
			// 3. 分发任务
		}
	}
}
