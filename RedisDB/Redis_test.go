package RedisDB

import (
	"fmt"
	"testing"
	"time"
)

func TestInitAndClose(t *testing.T) { // 用于测试驱动连接启动和关闭的测试函数
	go InitRedisAPI()
	go InitRedisCheck()
	time.Sleep(time.Second)
	CloseRedisAPI()
	time.Sleep(time.Second)
	CloseRedisCheck()
	time.Sleep(time.Second)
}

func TestAllIPNow(t *testing.T) {
	fmt.Printf("当前Redis数据库共有IP代理信息： %d 条\n", len(AllIPNow()))
}
