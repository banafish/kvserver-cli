package main

import (
	"bufio"
	"fmt"
	"github.com/banafish/kvserver-cli/client"
	"os"
	"strconv"
	"strings"
)

var ck *client.Clerk

func main() {
	input := bufio.NewScanner(os.Stdin)
	var serverIDs []string

	fmt.Println("请输入要连接服务器的地址，以空格分隔")
	for {
		fmt.Print("---> ")
		input.Scan()
		str := strings.TrimSpace(input.Text())
		if str == "" {
			fmt.Println("输入有误，请重新输入")
			continue
		}
		serverIDs = strings.Split(str, " ")
		ck = client.MakeClerk(serverIDs)
		fmt.Println("ok")
		break
	}

	for {
		fmt.Print("---> ")
		input.Scan()
		arr := strings.Split(input.Text(), " ")
		switch arr[0] {
		case "":

		case "get":
			if len(arr) != 2 {
				fmt.Println("输入有误")
			} else {
				fmt.Println(ck.Get(arr[1]))
			}
		case "put":
			if len(arr) != 3 {
				fmt.Println("输入有误")
			} else {
				ck.Put(arr[1], arr[2])
				fmt.Println("ok")
			}
		case "append":
			if len(arr) != 3 {
				fmt.Println("输入有误")
			} else {
				ck.Append(arr[1], arr[2])
				fmt.Println("ok")
			}
		case "rf":
			if len(arr) != 2 {
				fmt.Println("输入有误")
			} else {
				fmt.Println(ck.GetRaftStat(false, arr[1]))
			}
		case "rflog":
			if len(arr) != 2 {
				fmt.Println("输入有误")
			} else {
				fmt.Println(ck.GetRaftStat(true, arr[1]))
			}
		case "svr":
			if len(arr) != 2 {
				fmt.Println("输入有误")
			} else {
				fmt.Println(ck.GetServerStat(arr[1]))
			}
		case "speed":
			if len(arr) != 4 {
				fmt.Println("输入有误")
			} else {
				num, _ := strconv.Atoi(arr[2])
				t, _ := strconv.Atoi(arr[3])
				client.SpeedTest(serverIDs, num, arr[1], t)
			}
		case "speedauto":
			if len(arr) != 4 {
				fmt.Println("输入有误")
			} else {
				num, _ := strconv.Atoi(arr[2])
				t, _ := strconv.Atoi(arr[3])
				for i := 1; i <= num; i++ {
					client.SpeedTest(serverIDs, i, arr[1], t)
				}
			}
		case "exit":
			return
		default:
			fmt.Println("非法输入")
		}
	}
}
