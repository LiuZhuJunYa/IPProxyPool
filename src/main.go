package main

import (
	"ProxyPool/API"
	"ProxyPool/Check"
	"ProxyPool/CollectIP"
	"fmt"
	"os"
	"os/exec"
)

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func main() {
	for {
		var choice int
		var LocalAddress string
		var num byte
		fmt.Println("--------------------------------------")
		fmt.Println("---        欢迎使用IP代理池         ---")
		fmt.Println("--------------------------------------")
		fmt.Println("1. 网上爬取代理 IP 或者在本地导入 IP 代理")
		fmt.Println("2. 单独进行一次 10 轮的 IP 信息更新迭代（不会影响正常的程序运行）")
		fmt.Println("3. 统计当前 IP 代理池中 IP 代理的数量")
		fmt.Println("4. 提取 http 代理信息")
		fmt.Println("5. 对自己创建的靶机模拟 SQL 注入攻击")
		fmt.Println("6. 退出 IP 代理池")
		fmt.Print("请输入 1~3：")
		fmt.Scanln(&choice)
		runCmd("cmd", "/c", "cls")
		switch choice {
		case 1:
			fmt.Println("--------------------------------------")
			fmt.Println("---     欢迎使用IP代理池 功能1     ---")
			fmt.Println("--------------------------------------")
			fmt.Println("1. 从 89 IP 免费代理网站爬取 IP 代理信息（40 页）")
			fmt.Println("2. 从本地导入 IP 代理信息")
			fmt.Println("3. 返回上一级")
			fmt.Print("请输入1~3：")
			fmt.Scanln(&choice)
			switch choice {
			case 1:
				collectNUM, storeNUM := CollectIP.GetProxyFrom89()
				fmt.Printf("已经爬取了 %d 条代理信息\n", collectNUM)
				fmt.Printf("其中存储了 %d 条代理信息\n", storeNUM)
				runCmd("cmd", "/c", "pause")
				runCmd("cmd", "/c", "cls")
				break
			case 2:
				fmt.Print("请输入本地地址：")
				fmt.Scanln(&LocalAddress)
				collectNUM, storeNUM := CollectIP.GetProxyFromLocal(LocalAddress)
				fmt.Println("已完成从 " + LocalAddress + " 导入 IP 代理信息数据")
				fmt.Printf("已经爬取了 %d 条代理信息\n", collectNUM)
				fmt.Printf("其中存储了 %d 条代理信息\n", storeNUM)
				runCmd("cmd", "/c", "pause")
				runCmd("cmd", "/c", "cls")
				break
			default:
				runCmd("cmd", "/c", "cls")
				break
			}
			break
		case 2:
			fmt.Println("--------------------------------------")
			fmt.Println("---     欢迎使用IP代理池 功能2     ---")
			fmt.Println("--------------------------------------")
			go func() {
				for i := 0; ; i++ {
					Check.CheckPool()
					fmt.Printf("已经完成 %d 轮IP代理池的检测\n", i+1)
				}
			}()
			fmt.Println("已经开始 IP 代理池的动态更新迭代")
			runCmd("cmd", "/c", "pause")
			runCmd("cmd", "/c", "cls")
			break
		case 3:
			fmt.Println("--------------------------------------")
			fmt.Println("---     欢迎使用IP代理池 功能3     ---")
			fmt.Println("--------------------------------------")
			num1, num2 := API.Statistics()
			fmt.Println("当前动态代理池情况：")
			fmt.Println("共有 IP 代理", num1, "条")
			fmt.Println("值为 100 的代理", num2, "条")
			runCmd("cmd", "/c", "pause")
			runCmd("cmd", "/c", "cls")
			break
		case 4:
			fmt.Println("--------------------------------------")
			fmt.Println("---     欢迎使用IP代理池 功能4     ---")
			fmt.Println("--------------------------------------")
			fmt.Print("请输入要提取的 IP 代理信息数量：(请输入小于 100 的数值)")
			fmt.Scanln(&num)
			if 100 < num {
				fmt.Printf("请输入小于 100 的数值\n")
			} else {
				flag := API.GetIPproxy(num)
				if flag {
					fmt.Printf("已经成功提取了 %d 条代理信息\n", num)
				}
			}
			runCmd("cmd", "/c", "pause")
			runCmd("cmd", "/c", "cls")
			break
		case 5:
			fmt.Println("--------------------------------------")
			fmt.Println("---     欢迎使用IP代理池 功能5     ---")
			fmt.Println("--------------------------------------")
			fmt.Println("SQL 注入攻击的效果请具体在靶机中看")
			fmt.Println("正在进行普通的 SQL 注入攻击......")
			API.SQLinjectionASCII()
			fmt.Println("已经完成普通版的SQL注入攻击")
			fmt.Println("--------------------------------------")
			fmt.Println("--------------------------------------")
			fmt.Println("--------------------------------------")
			fmt.Println("正在进行通过 IP 代理的 SQL 注入攻击......")
			API.SQLinjectionIP()
			fmt.Println("已经完成IP代理版的SQL注入攻击")
			runCmd("cmd", "/c", "pause")
			runCmd("cmd", "/c", "cls")
			break
		case 6:
			return
		default:
			runCmd("cmd", "/c", "cls")
			continue
		}
	}
}
