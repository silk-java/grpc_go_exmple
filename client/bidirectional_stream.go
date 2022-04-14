package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
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

	stream, err := client.BidirectionalStream(context.Background())
	if err!=nil {
		log.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		err := stream.Send(&pb.Product{Id: int64(i)})
		if err!=nil {
			log.Fatal(err)
		}
		fmt.Println("发送商品：",i)
		recv, err := stream.Recv()
		if err==io.EOF {
			//表示接收此次流接收结束
			break
		}
		if err!=nil {
			log.Fatal(err)
			continue
		}
		fmt.Println("收到商品：",recv.GetId(),recv.GetName())
	}

}
