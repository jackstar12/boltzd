package liquid

import (
	"context"
	"encoding/base64"
	"github.com/BoltzExchange/boltz-lnd/logger"
	"github.com/ybbus/jsonrpc/v3"
	"os"
)

type Client struct {
	Rpc jsonrpc.RPCClient
	Ctx context.Context
}

type BlockchainInfo struct {
	Chain string `json:"chain"`
}

func Connect(url string, cookieFile string) *Client {
	cookie, err := os.ReadFile(cookieFile)
	if err != nil {
		logger.Fatal("Could not read file")
	}
	return &Client{
		Rpc: jsonrpc.NewClientWithOpts(url, &jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString(cookie),
			},
		}),
		Ctx: context.Background(),
	}
}

func (client Client) Call(method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	response, err := client.Rpc.Call(client.Ctx, method, params)
	if response.Error != nil {
		return nil, response.Error
	}
	return response, err
}

func (client Client) GetBlockchainInfo() (*BlockchainInfo, error) {
	var info *BlockchainInfo
	err := client.Rpc.CallFor(client.Ctx, &info, "getblockchaininfo")
	return info, err
}
