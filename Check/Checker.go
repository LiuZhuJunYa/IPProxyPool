package Check

import (
	"ProxyPool/RedisDB"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"
)

type IPstatus struct {
	IP     string
	Status bool
}

var (
	wg    sync.WaitGroup // 设置等待队列
	mutex sync.Mutex     // 设置互斥锁
)

func CheckPool() { // 用于检测整个Redis数据库中的IP代理可用性
	var statusList []IPstatus
	waitForTest := RedisDB.AllIPForCheck()
	for index, value := range waitForTest {
		wg.Add(1) // 开启了一个协程
		if 0 == index%10000 && index != 0 {
			fmt.Println("已经完成", index, "个代理的检测......")
		}
		if 1000 < runtime.NumGoroutine() { // 最多可同时启动1000个协程进行检测
			time.Sleep(10 * time.Second)
		}
		go func(IPURL string) {
			status := TestIP(IPURL) // 放入检测函数检查IP代理的可用性
			mutex.Lock()            // 加锁
			defer mutex.Unlock()    // 确保在函数结束时解锁
			statusList = append(statusList, IPstatus{
				IP:     IPURL,
				Status: status,
			})
			wg.Done() // 完成了一个携程
		}(value)
	}
	wg.Wait() // 进程阻塞至所有协程工作结束

	for _, value := range statusList {
		RedisDB.ProxyStatus(value.IP, value.Status)
	}
}

//func TestIP(test string) bool { // 测试IP代理有效性的函数
//	proxyURL, err := url.Parse(test) // 解析IP代理信息
//	if err != nil {                  // 解析报错处理
//		//panic(err)
//		fmt.Println(test, "代理解析出错")
//	}
//	transport := &http.Transport{ // 设置代理通道
//		Proxy: http.ProxyURL(proxyURL),
//	}
//	client := &http.Client{ // 设置通讯客户端
//		Transport: transport,
//		Timeout:   10 * time.Second, // 设置超时时间，可根据需求调整
//	}
//	response, err := client.Get("http://httpbin.org/ip") // 向测试网站发送get请求
//	if err != nil {                                      // 代理不可用：没有回应报文
//		//fmt.Println("代理不可用：", err)
//		return false
//	}
//	defer response.Body.Close()
//	//body, err := ioutil.ReadAll(response.Body)
//	_, err = ioutil.ReadAll(response.Body)
//	if err != nil { // 代理不可用：中间人出了问题会丢包
//		//fmt.Println("读取响应失败：", err)
//		return false
//	}
//	if response.StatusCode == http.StatusOK { // 代理可用
//		//fmt.Println("代理可用，响应内容：", string(body))
//		return true
//	} else { // 代理不可用：与网站无法通讯
//		//fmt.Printf("代理不可用，响应状态码：%d\n", response.StatusCode)
//		return false
//	}
//}

func TestIP(test string) bool {
	proxyURL, err := url.Parse(test)
	if err != nil {
		log.SetOutput(ioutil.Discard) // 禁止日志输出
		log.Println(test, "代理解析出错")
		return false
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
	response, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		log.SetOutput(ioutil.Discard) // 禁止日志输出
		log.Println("代理不可用：", err)
		return false
	}
	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.SetOutput(ioutil.Discard) // 禁止日志输出
		log.Println("读取响应失败：", err)
		return false
	}
	if response.StatusCode == http.StatusOK {
		//fmt.Println("代理可用，响应内容：", string(body))
		return true
	} else {
		log.SetOutput(ioutil.Discard) // 禁止日志输出
		log.Printf("代理不可用，响应状态码：%d\n", response.StatusCode)
		return false
	}
}
