package RedisDB

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type KeyValuePairs struct { // 设置一个键值对结构体，方便记录
	Key   string // 键名
	Value int    // 值
}

var (
	RedisCheck *redis.Client          // 用于动态检测 IP 信息可用性的 Redis 驱动
	RedisAPI   *redis.Client          // 用于其他工作的 Redis 驱动
	ctxCheck   = context.Background() // v9 版本之后必须添加的一个上下文功能
	ctxAPI     = context.Background()

	InitValue = 10  // 初始获取IP的值为10
	MaxValue  = 100 // IP代理可用性最高为 100
)

func InitRedisCheck() { // 初始化、启动 RedisCheck 驱动的函数
	RedisCheck = redis.NewClient(&redis.Options{
		Addr:     "192.168.254.134:6379", // 设置 redis 的 IP 地址和端口
		Password: "root",                 // 设置密码
		DB:       0,                      // 使用默认数据库
	})
	_, redisErr := RedisCheck.Ping(ctxCheck).Result() // 测试 redis 是否已经连接成功
	if redisErr != nil {
		log.Println("RedisCheck 连接失败")
	} else {
		//log.Println("RedisCheck 初始化成功")
	}
}

func CloseRedisCheck() { // 关闭 RedisCheck 驱动
	err := RedisCheck.Close()
	if err != nil {
		fmt.Println("RedisCheck 连接关闭失败")
	}
	//fmt.Println("RedisCheck 驱动连接已关闭")
}

func InitRedisAPI() { // 初始化、启动 RedisAPI 驱动的函数
	RedisAPI = redis.NewClient(&redis.Options{
		Addr:     "192.168.254.134:6379", // 设置 redis 的 IP 地址和端口
		Password: "root",                 // 设置密码
		DB:       0,                      // 使用默认数据库
	})
	_, redisErr := RedisAPI.Ping(ctxAPI).Result() // 测试 redis 是否已经连接成功
	if redisErr != nil {
		log.Println("RedisAPI 连接失败")
	} else {
		//log.Println("RedisAPI 初始化成功")
	}
}

func CloseRedisAPI() { // 关闭 RedisAPI 驱动
	err := RedisAPI.Close()
	if err != nil {
		fmt.Println("RedisAPI 连接关闭失败")
	}
	//fmt.Println("RedisAPI 驱动连接已关闭")
}

func StoreAfterCollect(IPList []string) { // 专门用于收集IP后的存储功能
	InitRedisAPI()
	fmt.Println("Redis数据库连接初始化成功，正在准备存储数据......")

	for _, ipProxy := range IPList { // 将传进来的IP代理信息全部加入到Redis中
		err := RedisAPI.Set(ctxAPI, ipProxy, InitValue, 0).Err() // 往Redis中添加数据
		if nil != err {
			fmt.Println(ipProxy, "该代理信息添加失败")
		} else {
			//fmt.Println(ipProxy, "代理信息添加/更改成功")
		}
	}

	CloseRedisAPI()
	fmt.Println("已经完成数据存储，已关闭Redis数据库连接......")
}

func AllIPNow() (List []string) { // 用于取出当前池中所有的IP代理信息
	InitRedisAPI()
	List = RedisAPI.Keys(ctxAPI, "*").Val()
	CloseRedisAPI()
	return
}

func AllIPForCheck() (List []string) { // 取出所有IP代理信息，因为检测时独立运行的，要用另外的驱动程序
	InitRedisCheck()
	List = RedisCheck.Keys(ctxCheck, "*").Val()
	CloseRedisCheck()
	return
}

func ProxyStatus(IP string, status bool) {
	InitRedisCheck()
	defer CloseRedisCheck()

	if status { // 如果该代理可用，则直接变成 100
		err := RedisCheck.Set(ctxCheck, IP, MaxValue, 0).Err() // 往Redis中添加数据
		if nil != err {
			fmt.Println(IP, "该代理信息更新失败")
		} else {
			//fmt.Println(ipProxy, "代理信息添加/更改成功")
		}
	} else { // 如果代理不可用，有两种情况
		IPvalue, _ := RedisCheck.Get(ctxCheck, IP).Result()
		IPvalueInt, _ := strconv.Atoi(IPvalue) // 将值转换为 int 类型，方便后续运算
		if IPvalueInt != 0 {                   // 这个键当前的值不为 0，则值 -1
			IPvalueInt -= 1
			err := RedisCheck.Set(ctxCheck, IP, IPvalueInt, 0).Err() // 往Redis中添加数据
			if nil != err {
				fmt.Println(IP, "该代理信息更新失败")
			} else {
				//fmt.Println(ipProxy, "代理信息添加/更改成功")
			}
		} else { // 这个键当前的值为 0，则直接删掉
			RedisCheck.Del(ctxCheck, IP).Err()
		}
	}
}

func APICheckValue(IP string) (Value int) { // Api模块查询IP对应的可用值
	InitRedisAPI()
	defer CloseRedisAPI()

	IPvalue, _ := RedisAPI.Get(ctxAPI, IP).Result()
	Value, _ = strconv.Atoi(IPvalue) // 将值转换为 int 类型，方便后续运算

	return
}

func ReturnIP(num byte) bool { // 该函数的返回值暂时确定为 string 类型，后期可以在接口模块修改
	keyList := AllIPNow()               // 首先先取出所有的 Key，才能够取得值
	shuffledKeyList := shuffle(keyList) // 先将整个切片乱序
	var proxyList []KeyValuePairs
	InitRedisAPI()
	defer CloseRedisAPI()

	scoreStr0, _ := RedisAPI.Get(ctxAPI, shuffledKeyList[0]).Result()
	scoreInt0, _ := strconv.Atoi(scoreStr0)
	proxyList = append(proxyList, KeyValuePairs{
		Key:   shuffledKeyList[0],
		Value: scoreInt0,
	})

	for i := 1; i < len(keyList); i++ {
		scoreStr, _ := RedisAPI.Get(ctxAPI, shuffledKeyList[i]).Result()
		scoreInt, _ := strconv.Atoi(scoreStr)
		proxy := KeyValuePairs{
			Key:   shuffledKeyList[i],
			Value: scoreInt,
		}
		for index, value := range proxyList {
			if scoreInt >= value.Value {
				proxyList = append(proxyList[:index], append([]KeyValuePairs{proxy}, proxyList[index:]...)...)
				break
			}
		}
	}

	if len(proxyList) < int(num) {
		fmt.Println("当前代理池没有这么多代理信息，请获取一些之后再来吧")
		return false
	}

	for i := byte(0); i < num; i++ {
		fmt.Println("代理："+proxyList[i].Key, "该代理的可用性为：", proxyList[i].Value)
	}
	return true
}

func ReturnOneIP() string {
	keyList := AllIPNow()               // 首先先取出所有的 Key，才能够取得值
	shuffledKeyList := shuffle(keyList) // 先将整个切片乱序
	var proxyList []KeyValuePairs
	InitRedisAPI()
	defer CloseRedisAPI()

	scoreStr0, _ := RedisAPI.Get(ctxAPI, shuffledKeyList[0]).Result()
	scoreInt0, _ := strconv.Atoi(scoreStr0)
	proxyList = append(proxyList, KeyValuePairs{
		Key:   shuffledKeyList[0],
		Value: scoreInt0,
	})

	for i := 1; i < len(keyList); i++ {
		scoreStr, _ := RedisAPI.Get(ctxAPI, shuffledKeyList[i]).Result()
		scoreInt, _ := strconv.Atoi(scoreStr)
		proxy := KeyValuePairs{
			Key:   shuffledKeyList[i],
			Value: scoreInt,
		}
		for index, value := range proxyList {
			if scoreInt >= value.Value {
				proxyList = append(proxyList[:index], append([]KeyValuePairs{proxy}, proxyList[index:]...)...)
				break
			}
		}
	}
	return proxyList[0].Key
}

func shuffle(slice []string) []string { // shuffle 接受一个字符串切片，将其内部元素乱序后返回
	rand.Seed(time.Now().UnixNano())
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
