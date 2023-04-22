package CollectIP

import (
	"ProxyPool/RedisDB"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetProxyFrom89() (CollectNum, StoreNum int) { // 本程序暂时只支持从89免费IP代理网站爬取IP代理信息
	var IPList []string // 用来存储爬取的IP代理信息

	for page := 0; page < 40; page++ { // 89IP代理大概每天大概有40页左右可用
		time.Sleep(time.Second) // 防止爬虫时被ban掉IP
		fmt.Println("正在爬取第", strconv.Itoa(page+1), "页的数据......")
		urlAddress := "https://www.89ip.cn/index_" + strconv.Itoa(page+1) + ".html" // 可以更改url的数字
		res, err := http.Get(urlAddress)                                            // 向网站发送get请求，获取未经整理的HTML原文
		if err != nil {                                                             // 处理请求出错的情况
			log.Fatal(err)
		}
		defer func(Body io.ReadCloser) { // 要记得关闭 IO 通道，避免资源消耗
			Body.Close()
		}(res.Body)
		if res.StatusCode != 200 { // 检测是否被反爬，一般来说免费网站应该不至于
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		dom, err := goquery.NewDocumentFromReader(res.Body) // 初始化goquery阅读器
		if err != nil {
			log.Fatalln(err)
		}
		dom.Find(".layui-table").Find("tr").Each(func(i int, selection *goquery.Selection) { // HTML中的表格类名称class="layui-table"
			ip := selection.Find("td")                // 去具体的列中读取数据
			ipSplit := strings.Split(ip.Text(), "\n") // 通过换行截取数据
			if 2 < len(ipSplit) {
				IPList = append(IPList, "http://"+InformationConversion(ipSplit[1])+":"+InformationConversion(ipSplit[2])) // 第一段数据是IP，第二段数据是端口
			}
		})
	}

	//for index, value := range IPList {
	//	fmt.Println(index, value)
	//}

	CollectNum = len(IPList)               // 记录收集IP代理条目
	lastPoolNum := len(RedisDB.AllIPNow()) // 记录存储有效IP代理条目
	RedisDB.StoreAfterCollect(IPList)
	StoreNum = len(RedisDB.AllIPNow()) - lastPoolNum
	return
}

func InformationConversion(ord string) (after string) { // 通过抓包发现爬取的数据为：\t\t\txxx.xxx.xxx.xxx\t\t，要做一步数据处理
	old := strings.TrimRight(ord, "\t") // 去除后端所有的制表符
	after = strings.TrimLeft(old, "\t") // 去除前端所有的制表符
	return
}

func GetProxyFromLocal(urlAddress string) (CollectNum, StoreNum int) { // 用于读取并且存储本地txt文本的IP代理信息
	content, err := ioutil.ReadFile(urlAddress) //由于文件的打开和关闭操作已经封装在ReadFile中了，所以不需要额外添加open\close
	if err != nil {                             //读取有误
		fmt.Println("读取出错，错误为:", err)
	}

	var listWithDuel []string
	IPList := strings.Split(string(content), "\n") // 按行分割数据
	for _, v := range IPList {                     // 处理一下尾部的"\r"
		IPWithDuel := strings.TrimSuffix(v, "\r")
		IPWithDuel = "http://" + IPWithDuel
		listWithDuel = append(listWithDuel, IPWithDuel)
	}

	CollectNum = len(listWithDuel)         // 记录收集IP代理条目
	lastPoolNum := len(RedisDB.AllIPNow()) // 记录存储有效IP代理条目
	RedisDB.StoreAfterCollect(listWithDuel)
	StoreNum = len(RedisDB.AllIPNow()) - lastPoolNum
	return
}
