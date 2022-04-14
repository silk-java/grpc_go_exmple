package main

import (
	"context"
	"google.golang.org/grpc"
	"juhe.cn.weather_report/pb"
	"log"
)

/**
 * @author micro.cloud.fly
 * @date 2022/3/25 3:54 下午
 * @desc
 */

func main() {


	dial, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer dial.Close()
	client := pb.NewQueryClient(dial)
	/*  客户端和服务端是同步的*/
	//--------｜客户端｜----------｜服务端｜
	//-------｜一次发送｜---------｜同步返回｜
	pinyin, err := client.GetById(context.Background(), &pb.WeatherRequest{
		City: "北京",
		Key:  "",
	})
	log.Println("返回：", pinyin)
	log.Println("错误：", err)


}
