package test1

import (
	"fmt"
	"math"
)

// ChannelForPrime 协程判断质数
func ChannelForPrime() {
	// 定义工作协程数量、数据个数
	workerNum := 4
	count := 800
	num := make(chan int, 10)
	res := make(chan int, count)
	exit := make(chan bool, workerNum)

	// 产生2~count个数据
	go produce(num, count)

	// 创建n个协程，不断从num中取数据并进行判断
	for i := 0; i < workerNum; i++ {
		go worker(num, res, exit)
	}

	// 所有worker发送exit后，关闭所有通道并退出循环；必须关闭res才能正常遍历
	flag := 0
	for {
		if _, ok := <-exit; ok {
			flag++
			if flag == workerNum {
				close(exit)
				close(res)
				break
			}
		}
	}

	// 此处天坑；for时取了一个值，print时已经取到下一个值
	//for range res {
	//	fmt.Println(<-res)
	//}

	for val := range res {
		fmt.Println(val)
	}
}

func produce(num chan int, length int) {
	for i := 2; i <= length; i++ {
		num <- i
	}
	close(num)
}

func worker(num chan int, res chan int, exit chan bool) {
	for {
		if val, ok := <-num; !ok {
			exit <- true
			break
		} else {
			if val == 2 || val == 3 {
				res <- val
				continue
			}
			if checkPrime(val) {
				res <- val
			}
		}
	}
}

func checkPrime(val int) bool {
	for i := 2; i <= int(math.Sqrt(float64(val))); i++ {
		if val%i == 0 {
			return false
		}
	}
	return true
}
