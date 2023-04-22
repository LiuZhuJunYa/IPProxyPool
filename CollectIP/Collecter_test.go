package CollectIP

import (
	"ProxyPool/RedisDB"
	"fmt"
	"testing"
)

func TestGetProxyFrom89(t *testing.T) { // 测试网络爬虫的函数方法
	GetProxyFrom89()
}

func TestGetProxyFromLocal(t *testing.T) {
	URL := "C:\\Users\\ASUS\\Desktop\\2.txt"
	collectNUM, storeNUM := GetProxyFromLocal(URL)
	fmt.Printf("已经爬取了 %d 条代理信息\n", collectNUM)
	fmt.Printf("其中存储了 %d 条代理信息\n", storeNUM)
	for i, v := range RedisDB.AllIPNow() {
		fmt.Println(i, v)
	}
}
