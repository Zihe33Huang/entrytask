package ziherpc

import (
	"entrytask/internal/rpc/pb"
	"errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"io"
	"log"
	"net"
	"sync"
)

type Call struct {
	Seq           uint64
	ServiceMethod string      // format "<service>.<method>"
	Args          interface{} // arguments to the function
	Reply         interface{} // reply from the function
	Error         error       // if error occurs, it will be set
	Done          chan *Call  // Strobes when call is complete.
}

func (call *Call) done() {
	call.Done <- call
}

type Client struct {
	conn     io.ReadWriteCloser
	sending  sync.Mutex // protect following
	header   pb.Header
	mu       sync.Mutex // protect following
	seq      uint64
	pending  map[uint64]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("connection is shut down")

// Close the connection
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.conn.Close()
}

// IsAvailable return true if the client does work
func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown && !client.closing
}

func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

func (client *Client) receive() {
	var err error
	for err == nil {
		buf := make([]byte, 1024)
		client.conn.Read(buf)
		resp := pb.Response{}
		err := proto.Unmarshal(buf, &resp)
		if err != nil {
			log.Printf("rpc: unmarshal response error: %s", err)
			client.terminateCalls(err)
			return
		}
		call := client.removeCall(resp.Header.Seq)
		switch {
		case call == nil:
			// it usually means that Write partially failed
			// and call was already removed.
			log.Printf("rpc: rpc: nil call")
		//case h.Error != "":
		//	call.Error = fmt.Errorf(h.Error)
		//	err = client.cc.ReadBody(nil)
		//	call.done()
		default:
			////err = client.cc.ReadBody(call.Reply)
			//if err != nil {
			//	call.Error = errors.New("reading body " + err.Error())
			//}
			call.done()
		}
	}
	// error occurs, so terminateCalls pending calls
	client.terminateCalls(err)
}

func NewClient(conn io.ReadWriteCloser) (*Client, error) {
	client := &Client{
		seq:     1, // seq starts with 1, 0 means invalid call
		pending: make(map[uint64]*Call),
	}
	go client.receive()
	return client, nil
}

func Dial(network, address string) (client *Client, err error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	// close the connection if client is nil
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()
	return NewClient(conn)
}

func (client *Client) send(call *Call) {
	// make sure that the client will send a complete request
	client.sending.Lock()
	defer client.sending.Unlock()

	// register this call.
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	// prepare request header
	//client.header.ServiceMethod = call.ServiceMethod
	//client.header.Seq = seq
	pbReq := pb.Request{}
	pbReq.Hearder.Seq = seq
	pbReq.Hearder.ServiceMethod = call.ServiceMethod
	newAny, err := anypb.New(call.Args.(proto.Message))
	if err != nil {
		call.Error = err
		call.done()
		log.Printf("rpc: protobuf marshal new anypb error: %s", err)
	}
	pbReq.Args = newAny
	marshal, err := proto.Marshal(&pbReq)
	if err != nil {
		call.Error = err
		call.done()
		log.Printf("rpc: marshal new request error: %s", err)
	}
	client.conn.Write(marshal)
	// encode and send the request
	//if err := client.cc.Write(&client.header, call.Args); err != nil {
	//	call := client.removeCall(seq)
	//	// call may be nil, it usually means that Write partially failed,
	//	// client has received the response and handled
	//	if call != nil {
	//		call.Error = err
	//		call.done()
	//	}
	//}
}
