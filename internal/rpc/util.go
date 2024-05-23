package ziherpc

import (
	"encoding/binary"
	"io"
)

func packMessage(message []byte) []byte {
	length := uint32(len(message))
	buffer := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buffer[:4], length)
	copy(buffer[4:], message)
	return buffer
}

func unpackMessage(conn io.ReadWriteCloser) ([]byte, error) {
	// 1. read the length of request
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		return nil, err
	}
	// Convert the length prefix to an integer
	length := binary.BigEndian.Uint32(lengthBuf)
	// 2. read request from protobuf
	messageBuf := make([]byte, length)
	_, err = io.ReadFull(conn, messageBuf)

	if err != nil {
		return nil, err
	}

	return messageBuf, nil
}
