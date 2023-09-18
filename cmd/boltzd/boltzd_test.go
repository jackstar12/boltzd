package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"testing"
	"time"

	boltz_lnd "github.com/BoltzExchange/boltz-lnd"
	"github.com/BoltzExchange/boltz-lnd/boltzrpc"
	"github.com/BoltzExchange/boltz-lnd/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func setup(t *testing.T) (*boltzrpc.Boltz, func()) {
	cfg := boltz_lnd.LoadConfig()
	cfg.RPC.NoTls = true
	cfg.RPC.NoMacaroons = true

	logger.InitLogger(cfg.LogFile, cfg.LogPrefix)

	Init(cfg)

	server := cfg.RPC.Grpc

	lis := bufconn.Listen(1024 * 1024)

	conn, err := grpc.DialContext(
		context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		logger.Fatal("error connecting to server: " + err.Error())
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			logger.Error("error connecting serving server: " + err.Error())
		}
	}()

	close := func() {
		err := lis.Close()
		if err != nil {
			logger.Error("error closing listener: " + err.Error())
		}
		server.Stop()
	}

	client := boltzrpc.Boltz{
		Client: boltzrpc.NewBoltzClient(conn),
		Ctx:    context.Background(),
	}

	return &client, close
}

func TestGetInfo(t *testing.T) {
	client, close := setup(t)
	defer close()

	info, err := client.GetInfo()

	assert.NoError(t, err)
	assert.Equal(t, "regtest", info.Network)
}

func TestDeposit(t *testing.T) {

	client, close := setup(t)
	defer close()

	tests := []struct {
		desc   string
		pairId string
		cli    func(string) string
	}{
		{"BTC Deposit", "BTC/BTC", btc_cli},
		{"Liquid Deposit", "L-BTC/BTC", liquid_cli},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			swap, err := client.Deposit(25, tc.pairId)
			assert.NoError(t, err)

			info, err := client.GetSwapInfo(swap.Id)
			assert.NoError(t, err)
			assert.Equal(t, boltzrpc.SwapState_PENDING, info.Swap.State)

			tc.cli("sendtoaddress " + swap.Address + " 0.0025")
			tc.cli("-generate 1")

			time.Sleep(500 * time.Millisecond)

			info, err = client.GetSwapInfo(swap.Id)
			assert.NoError(t, err)
			assert.Equal(t, boltzrpc.SwapState_SUCCESSFUL, info.Swap.State)

		})
	}

}

func TestReverseSwap(t *testing.T) {
	client, close := setup(t)
	defer close()

	tests := []struct {
		desc   string
		pairId string
		cli    func(string) string
	}{
		{"BTC", "BTC/BTC", btc_cli},
		{"Liquid", "L-BTC/BTC", liquid_cli},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			address := tc.cli("getnewaddress")
			fmt.Println("New address: " + address)
			swap, err := client.CreateReverseSwap(250000, address, true, tc.pairId)
			assert.NoError(t, err)

			time.Sleep(500 * time.Millisecond)

			info, err := client.GetSwapInfo(swap.Id)
			assert.NoError(t, err)
			assert.Equal(t, boltzrpc.SwapState_PENDING, info.ReverseSwap.State)

			btc_cli("-generate 1")
			liquid_cli("-generate 1")
			liquid_cli("rescanblockchain")
			time.Sleep(50000 * time.Millisecond)

			info, err = client.GetSwapInfo(swap.Id)
			assert.NoError(t, err)
			assert.Equal(t, boltzrpc.SwapState_SUCCESSFUL, info.ReverseSwap.State)
		})
	}
}

func TestReverseSwapZeroConf(t *testing.T) {
	client, close := setup(t)
	defer close()

	swap, err := client.CreateReverseSwap(250000, "", true, "BTC/BTC")
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	info, err := client.GetSwapInfo(swap.Id)
	assert.NoError(t, err)
	assert.Equal(t, boltzrpc.SwapState_SUCCESSFUL, info.ReverseSwap.State)

}

func sh(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()

	if err != nil {
		logger.Fatal("could not execute cmd: " + cmd + " err:" + err.Error())
	}

	return strings.TrimSuffix(string(out), "\n")
}

func btc_cli(cmd string) string {
	return sh("docker exec lnbits-legend-bitcoind-1 bitcoin-cli -rpcuser=lnbits -rpcpassword=lnbits -regtest " + cmd)
}

func liquid_cli(cmd string) string {
	return sh("docker exec lnbits-legend-elementsd-1 elements-cli " + cmd)
}
