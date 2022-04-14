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


	/*客户端一次性发送过去数据，而服务端通过流的方式返回*/
	//--------｜客户端｜----------｜服务端｜
	//--------｜一次发送｜-----------｜流返回｜
	//客户端端一次发送一个省份，服务端通过流的方式返回这个省份的所有城市
	stream, err := client.ReturnByStream(context.Background(), &pb.CityRequest{Province: "江苏"})
	//由于服务端是流的方式返回数据，所以此时要从流中循环读取返回数据
	var city []string
	for {
		recv, err := stream.Recv()
		if err == nil {
			fmt.Println("服务端返回：", recv.Cityname)
			city = append(city, recv.Cityname)
		}
		if err == io.EOF {
			fmt.Println("服务端已经全部返回，客户端接受完毕！")
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("江苏所有的城市是：", city)


}
