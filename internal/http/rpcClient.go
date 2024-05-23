package httpserver

import (
	ziherpc "entrytask/internal/rpc"
)

var (
	client *ziherpc.Client
)

func SetClient(c *ziherpc.Client) {
	client = c
}

func GetClient() *ziherpc.Client {
	return client
}
