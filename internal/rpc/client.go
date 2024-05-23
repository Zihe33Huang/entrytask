package ziherpc

import (
	"entrytask/internal/rpc/pb"
	"errors"
	"fmt"
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
		var message []byte
		message, err = unpackMessage(client.conn)
		if err != nil {
			log.Println("rpc client: ", err.Error())
			client.terminateCalls(err)
			break
		}
		resp := pb.Response{}
		err = proto.Unmarshal(message, &resp)
		if err != nil {
			log.Printf("rpc: unmarshal response error: %s", err)
			client.terminateCalls(err)
			break
		}
		call := client.removeCall(resp.Header.Seq)
		switch {
		case call == nil:
			// it usually means that Write partially failed
			// and call was already removed.
			log.Printf("rpc: rpc: nil call")
		case resp.Header.Error != "":
			call.Error = fmt.Errorf(resp.Header.Error)
			call.done()
		default:
			////err = client.cc.ReadBody(call.Reply)
			//if err != nil {
			//	call.Error = errors.New("reading body " + err.Error())
			//}
			reply := call.Reply.(proto.Message)
			err = resp.Args.UnmarshalTo(reply)
			if err != nil {
				log.Printf("rpc: unmarshal response error: %s", err)
				break
			}

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
		conn:    conn,
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

	// 1. register this call.
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}
	// 2. prepare request
	h := pb.Header{}
	h.Seq = seq
	h.ServiceMethod = call.ServiceMethod

	pbReq := pb.Request{Hearder: &h}
	args := call.Args
	protoMessage, ok := args.(proto.Message)
	if !ok {
		log.Printf("rpc client: Args is not proto.Message", call.Args)
		return
	}
	newAny, err := anypb.New(protoMessage)
	if err != nil {
		call.Error = err
		call.done()
		log.Printf("rpc client: protobuf marshal new anypb error: %s", err)
	}
	pbReq.Args = newAny
	marshal, err := proto.Marshal(&pbReq)
	if err != nil {
		call.Error = err
		call.done()
		log.Printf("rpc client: marshal new request error: %s", err)
	}

	// 3. pack marshal with length
	message := packMessage(marshal)
	// Write the length prefix to the connection

	// 4. send request
	// Write the message payload to the connection
	_, err = client.conn.Write(message)
	if err != nil {
		log.Printf("rpc client: write request error: %s", err)
	}
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

// Go invokes the function asynchronously.
// It returns the Call structure representing the invocation.
func (client *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	client.send(call)
	return call
}

// Call invokes the named function, waits for it to complete,
// and returns its error status.
func (client *Client) Call(serviceMethod string, args, reply interface{}) error {
	call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}
