package lnd

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/chainrpc"
	"github.com/lightningnetwork/lnd/lnrpc/walletrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"io/ioutil"
	"strconv"
)

type LightningClient interface {
	GetInfo() (*lnrpc.GetInfoResponse, error)
	GetNodeInfo(pubkey string) (*lnrpc.NodeInfo, error)
	ListChannels() (*lnrpc.ListChannelsResponse, error)
	ClosedChannels() (*lnrpc.ClosedChannelsResponse, error)
	GetChannelInfo(chanId uint64) (*lnrpc.ChannelEdge, error)
	ListInactiveChannels() (*lnrpc.ListChannelsResponse, error)
	ForceCloseChannel(channelPoint string) (lnrpc.Lightning_CloseChannelClient, error)
}

type LND struct {
	Host        string `long:"lnd.host" description:"gRPC host of the LND node"`
	Port        int    `long:"lnd.port" description:"gRPC port of the LND node"`
	Macaroon    string `long:"lnd.macaroon" description:"Path to a macaroon file of the LND node"`
	Certificate string `long:"lnd.certificate" description:"Path to a certificate file of the LND node"`

	ctx context.Context

	client        lnrpc.LightningClient
	walletKit walletrpc.WalletKitClient
	chainNotifier chainrpc.ChainNotifierClient
}

func (lnd *LND) Connect() error {
	creds, err := credentials.NewClientTLSFromFile(lnd.Certificate, "")

	if err != nil {
		return errors.New(fmt.Sprint("could not read LND certificate: ", err))
	}

	con, err := grpc.Dial(lnd.Host+":"+strconv.Itoa(lnd.Port), grpc.WithTransportCredentials(creds))

	if err != nil {
		return errors.New(fmt.Sprint("could not create gRPC client: ", err))
	}

	lnd.client = lnrpc.NewLightningClient(con)
	lnd.walletKit = walletrpc.NewWalletKitClient(con)
	lnd.chainNotifier = chainrpc.NewChainNotifierClient(con)

	if lnd.ctx == nil {
		macaroonFile, err := ioutil.ReadFile(lnd.Macaroon)

		if err != nil {
			return errors.New(fmt.Sprint("could not read LND macaroon: ", err))
		}

		macaroon := metadata.Pairs("macaroon", hex.EncodeToString(macaroonFile))
		lnd.ctx = metadata.NewOutgoingContext(context.Background(), macaroon)
	}

	return nil
}

func (lnd *LND) GetInfo() (*lnrpc.GetInfoResponse, error) {
	return lnd.client.GetInfo(lnd.ctx, &lnrpc.GetInfoRequest{})
}

func (lnd *LND) ListChannels() (*lnrpc.ListChannelsResponse, error) {
	return lnd.client.ListChannels(lnd.ctx, &lnrpc.ListChannelsRequest{})
}

func (lnd *LND) AddInvoice(value int64, memo string) (*lnrpc.AddInvoiceResponse, error) {
	return lnd.client.AddInvoice(lnd.ctx, &lnrpc.Invoice{
		Value: value,
		Memo:  memo,
	})
}

func (lnd *LND) NewAddress() (string, error) {
	response, err := lnd.client.NewAddress(lnd.ctx, &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_WITNESS_PUBKEY_HASH,
	})

	if err != nil {
		return "", err
	}

	return response.Address, err
}

func (lnd *LND) LookupInvoice(preimageHash []byte) (*lnrpc.Invoice, error) {
	return lnd.client.LookupInvoice(lnd.ctx, &lnrpc.PaymentHash{
		RHash: preimageHash,
	})
}

func (lnd *LND) EstimateFee(confTarget int32) (*walletrpc.EstimateFeeResponse, error) {
	return lnd.walletKit.EstimateFee(lnd.ctx, &walletrpc.EstimateFeeRequest{
		ConfTarget: confTarget,
	})
}

func (lnd *LND) RegisterBlockListener(channel chan *chainrpc.BlockEpoch) error {
	client, err := lnd.chainNotifier.RegisterBlockEpochNtfn(lnd.ctx, &chainrpc.BlockEpoch{})

	go func() {
		for {
			// TODO: reconnection logic
			test, _ := client.Recv()
			channel <- test
		}
	}()

	return err
}
