package ziherpc

import (
	"entrytask/internal/rpc/pb"
	"errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

type Server struct {
	serviceMap sync.Map
}

type request struct {
	header       *pb.Header
	argv, replyv reflect.Value // argv and replyv of request
	mtype        *methodType
	svc          *service
}

func NewServer() *Server {
	return &Server{}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = NewServer()

var invalidRequest = struct{}{}

func (server *Server) serveConn(conn io.ReadWriteCloser) {
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)  // wait until all request are handled
	for {
		req, err := server.readRequest(conn)
		if err != nil {
			if req == nil {
				break // it's not possible to recover, so close the connection
			}
			req.header.Error = err.Error()
			server.sendResponse(conn, req, sending)
			continue
		}

		wg.Add(1)
		go server.handleRequest(conn, req, sending, wg)
	}
	wg.Wait()
	_ = conn.Close()
}

func (server *Server) readRequest(conn io.ReadWriteCloser) (*request, error) {
	// 1. unpack message
	messageBuf, err := unpackMessage(conn)
	if err != nil {
		log.Println("rpc server: unpackMessage error:", err)
		return nil, err
	}

	// 2.
	pbReq := pb.Request{}
	err = proto.Unmarshal(messageBuf, &pbReq)
	if err != nil {
		log.Printf("rpc server: read request unmarshal error: %v", err)
	}

	// 3. prepare request
	req := &request{header: pbReq.Hearder}
	req.svc, req.mtype, err = server.findService(req.header.ServiceMethod)
	if err != nil {
		return req, err
	}
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()
	argvi := req.argv.Interface()

	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	message, ok := argvi.(proto.Message)
	if !ok {
		log.Panicln("rpc server: the message is not a proto.Message")
	}
	pbReq.Args.UnmarshalTo(message)
	return req, nil

}

func (server *Server) handleRequest(conn io.ReadWriteCloser, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	err := req.svc.call(req.mtype, req.argv, req.replyv)
	if err != nil {
		req.header.Error = err.Error()
		server.sendResponse(conn, req, sending)
		return
	}
	//rv := req.replyv.Interface()
	//s := rv.(proto.Message)
	//newAny, err := anypb.New(s)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//pbReq := pb.Request{
	//	Hearder: req.header,
	//	Args:    newAny,
	//}
	server.sendResponse(conn, req, sending)
}

func (server *Server) sendResponse(conn io.ReadWriteCloser, req *request, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()

	pbResp := pb.Response{Header: req.header}
	// if the type of value is proto.Message, i.e. the type is defined in .proto file
	if message, ok := req.replyv.Interface().(proto.Message); ok {
		newAny, err := anypb.New(message)
		if err != nil {
			log.Println("protobuf marshal response error:", err)
			req.header.Error = err.Error()
		} else {
			pbResp.Args = newAny
		}
	} else {
		log.Println("the message is not a proto.Message")
		req.header.Error = "you must transfer proto.Message type"
	}

	marshal, err := proto.Marshal(&pbResp)
	if err != nil {
		log.Println("marshal request error: %v", err)
		req.header.Error = "rpc server marshal request error: " + err.Error()
	}

	message := packMessage(marshal)

	_, err = conn.Write(message)
	if err != nil {
		log.Printf("write request error: %v", err)
	}

}

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.serveConn(conn)
	}
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection.
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
//   - exported method of exported type
//   - two arguments, both of exported type
//   - the second argument is a pointer
//   - one return value, of type error
func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}

// Register publishes the receiver's methods in the DefaultServer.
func Register(rcvr interface{}) error { return DefaultServer.Register(rcvr) }
