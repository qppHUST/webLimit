package main

import (
	"fmt"
	"time"
	"webLimit/util"
)

func main() {
	//server := gin.Default()
	//server.Use(util.GetFixedWindowHandler())
	////server.Use(util.GetSlidingWindowHandler())
	//server.GET("/", func(context *gin.Context) {
	//	context.JSON(200, gin.H{
	//		"message": "root",
	//	})
	//})
	//server.Run(":8080")
	//testFixedWindow()
	testSlidingWindow()
}

func testSlidingWindow() {
	rate := util.NewSlideWindowLimitRate(3, time.Second*3, time.Second)
	fmt.Println("-----------------------------------------")
	print(rate.Acquire())
	fmt.Println("-----------------------------------------")
	print(rate.Acquire())
	fmt.Println("-----------------------------------------")
	fmt.Println("等待一秒")
	fmt.Println("-----------------------------------------")
	time.Sleep(time.Second)
	print(rate.Acquire())
	fmt.Println("-----------------------------------------")
	print(rate.Acquire())
	fmt.Println("-----------------------------------------")
	fmt.Println("等待两秒")
	fmt.Println("-----------------------------------------")
	time.Sleep(time.Second * 2)
	print(rate.Acquire())
	fmt.Println("-----------------------------------------")
	print(rate.Acquire())
	fmt.Println("-----------------------------------------")
	print(rate.Acquire())
}

func testFixedWindow() {
	rate := util.NewFixedWindow(2, time.Second*3)
	fmt.Println("-----------------------------------------")
	print(rate.TryAcquire())
	fmt.Println("-----------------------------------------")
	print(rate.TryAcquire())
	fmt.Println("-----------------------------------------")
	print(rate.TryAcquire())
	fmt.Println("-----------------------------------------")
	print(rate.TryAcquire())
	fmt.Println("-----------------------------------------")
	time.Sleep(time.Second * 3)
	print(rate.TryAcquire())
	fmt.Println("-----------------------------------------")
	print(rate.TryAcquire())
	fmt.Println("-----------------------------------------")
	print(rate.TryAcquire())
}

func print(flag bool) {
	if flag {
		fmt.Println("请求成功")
	} else {
		fmt.Println("请求失败")
	}
}
