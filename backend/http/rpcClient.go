package httpserver

import (
	ziherpc "entrytask/backend/rpc"
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
