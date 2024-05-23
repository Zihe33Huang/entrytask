package main

import (
	ziherpc "entrytask/internal/rpc"
	"entrytask/internal/rpc/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"net"
	"time"
)

type UserService struct{}

func (us UserService) AppendUser(args User, reply *User) error {
	reply.Name = args.Name + "xxxxx"
	return nil
}

func (us UserService) Append(args Num, reply *Num) error {
	reply.Num = args.Num + 1
	return nil
}

func startServer(addr chan string) {
	var foo UserService
	if err := ziherpc.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	ziherpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple geerpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// send options
	// send request & receive response
	for i := 0; i < 5; i++ {
		//user := User{Name: "zihe"}
		num := Num{Num: 5}
		h := &pb.Header{
			ServiceMethod: "UserService.Append",
			Seq:           uint64(i),
		}
		a, err := anypb.New(&num)
		if err != nil {
			log.Println(err.Error())
		}

		request := pb.Request{
			Hearder: h,
			Args:    a,
		}

		//_ = c.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		//_ = cc.ReadHeader(h)
		//var reply string
		//_ = cc.ReadBody(&reply)
		marshal, err := proto.Marshal(&request)
		if err != nil {
			log.Println(err.Error())
		}
		conn.Write(marshal)

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err.Error())
		}
		resp := pb.Response{}
		proto.Unmarshal(buffer[:n], &resp)
		nn := Num{}
		resp.Args.UnmarshalTo(&nn)
		fmt.Println(nn.Num)
	}
	fmt.Scanln()
}
