package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"juhe.cn.weather_report/pb"
	"log"
	"time"
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


	/*客户端是流，服务是一次性返回 */
	//--------｜客户端｜----------｜服务端｜
	//-------｜流｜-----------  ｜接收完一次同步返回｜
	//客户端通过流的方式发送三个数据过去----start
	names, err := client.GetByStream(context.Background())
	city := []string{"苏州", "上海", "青岛"}
	for i := 0; i < 3; i++ {
		if err := names.Send(&pb.WeatherRequest{
			City: city[i],
			Key:  "",
		}); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 2)
	}
	//关闭发送，让服务端知道客户端已经发送完毕
	recv, err := names.CloseAndRecv()
	fmt.Println(recv, err)
	//----------------end-------

}
