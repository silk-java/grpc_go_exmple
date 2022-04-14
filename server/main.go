package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"juhe.cn.weather_report/pb"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

const key = "2d1b16a202************"
func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterQueryServer(s, &weatherServer{})
	log.Printf("server listening at %v", listen.Addr())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to server:%v", err)
	}

}
/**
 * @author micro.cloud.fly
 * @date 2022/3/25 3:08 下午
 * @desc 天气预报服务端
 */
type weatherServer struct {
	pb.UnimplementedQueryServer
}
type WeaResp struct {
	Reason string `json:"reason"`
	Result struct {
		City     string `json:"city"`
		Realtime struct {
			Temperature string `json:"temperature"`
			Humidity    string `json:"humidity"`
			Info        string `json:"info"`
			Wid         string `json:"wid"`
			Direct      string `json:"direct"`
			Power       string `json:"power"`
			Aqi         string `json:"aqi"`
		} `json:"realtime"`
		Future []struct {
			Date        string `json:"date"`
			Temperature string `json:"temperature"`
			Weather     string `json:"weather"`
			Wid         struct {
				Day   string `json:"day"`
				Night string `json:"night"`
			} `json:"wid"`
			Direct string `json:"direct"`
		} `json:"future"`
	} `json:"result"`
	ErrorCode int64 `json:"error_code"`
}

//客户端通过流的方式发送数据过来，那么此时服务端，需要不断的读取客户端发送来
//的数据，等全部收到之后，一次行把数据返回给客户端
func (ws *weatherServer) GetByStream(qgs pb.Query_GetByStreamServer) error {
	var respArr []*pb.WeatherResponse
	for {
		recv, err := qgs.Recv()
		if err == io.EOF {
			return qgs.SendAndClose(&pb.StreamResp{Results: respArr})
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		//不断取得从客户端发来的数据,不断调用接口
		log.Println("收到：", recv.GetCity())
		resp := httpGet(recv.GetCity())
		log.Println("聚合返回：", resp)
		var weaResp WeaResp
		_ = json.Unmarshal([]byte(resp), &weaResp)
		respArr = append(respArr, &pb.WeatherResponse{
			Reason: weaResp.Reason,
			Result: &pb.Result{
				City:     weaResp.Result.City,
				Realtime: &pb.Realtime{Aqi: weaResp.Result.Realtime.Aqi},
			},
			ErrorCode: weaResp.ErrorCode,
		})
	}
}
func (ws *weatherServer) GetByName(ctx context.Context, weaRequest *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	log.Println("收到：", weaRequest.City)
	resp := httpGet(weaRequest.GetCity())
	var weaResp WeaResp
	err := json.Unmarshal([]byte(resp), &weaResp)
	fu := &pb.Future{
		Date:        weaResp.Result.Future[0].Date,
		Temperature: weaResp.Result.Future[0].Temperature,
		Weather:     weaResp.Result.Future[0].Weather,
		Direct:      weaResp.Result.Future[0].Direct,
	}
	fuArr := []*pb.Future{fu}
	return &pb.WeatherResponse{
		Reason: weaResp.Reason,
		Result: &pb.Result{
			City:     weaResp.Result.City,
			Realtime: &pb.Realtime{Aqi: weaResp.Result.Realtime.Aqi},
			Future:   fuArr,
		},
		ErrorCode: weaResp.ErrorCode,
	}, err
}

//客户端一次发送一个省份名字，服务端通过流的方式，每次查到这个省份的一个城市后，就写入流中
//返回给客户端
func (ws *weatherServer) ReturnByStream(request *pb.CityRequest, qrs pb.Query_ReturnByStreamServer) error {
	//取出客户端发送来省份名
	jiangsu_city := []string{"徐州", "苏州", "南京", "镇江"}
	zhejiang_city := []string{"宁波", "舟山", "杭州", "温州"}
	if request.GetProvince() == "江苏" {
		for _, s := range jiangsu_city {
			err := qrs.Send(&pb.CityResp{Cityname: s})
			if err != nil {
				return err
			}
			time.Sleep(time.Second * 2)
		}
	} else {
		for _, s := range zhejiang_city {
			err := qrs.Send(&pb.CityResp{Cityname: s})
			if err != nil {
				return err
			}
			time.Sleep(time.Second * 2)
		}
	}
	return nil
}

//客户端和服务端都是使用的流，客户端发送一个，服务端就返回一个，直到结束
func (ws *weatherServer) BidirectionalStream(stream pb.Query_BidirectionalStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//此时表示客户端发送结束了
			return nil
		}
		if err != nil {
			//此时表示真正有错误了
			return err
		}
		//此时接收到一个请求，服务端就返回一个请求
		fmt.Println("收到请求：",req.Id)
		err = stream.Send(&pb.Product{
			Id:   req.Id,
			Name: strconv.Itoa(int(req.Id)) + "name",
		})
		fmt.Println("回复客户端：",req.Id)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 5)
	}
	return nil
}
func (ws *weatherServer) GetById(ctx context.Context, weaRequest *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	log.Println("收到：", weaRequest.GetCity())
	resp := httpGet(weaRequest.GetCity())
	log.Println("聚合返回：", resp)
	var weaResp WeaResp
	err := json.Unmarshal([]byte(resp), &weaResp)
	fu := &pb.Future{
		Date:        weaResp.Result.Future[0].Date,
		Temperature: weaResp.Result.Future[0].Temperature,
		Weather:     weaResp.Result.Future[0].Weather,
		Direct:      weaResp.Result.Future[0].Direct,
	}
	fuArr := []*pb.Future{fu}
	return &pb.WeatherResponse{
		Reason: weaResp.Reason,
		Result: &pb.Result{
			City:     weaResp.Result.City,
			Realtime: &pb.Realtime{Aqi: weaResp.Result.Realtime.Aqi},
			Future:   fuArr,
		},
		ErrorCode: weaResp.ErrorCode,
	}, err
}

func httpGet(cityName string) string {
	url := "http://apis.juhe.cn/simpleWeather/query?key=" + key + "&city=" + cityName
	log.Println(url)
	res, _ := http.Get(url)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
