package cln

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BoltzExchange/boltz-lnd/cln/protos"
	"github.com/BoltzExchange/boltz-lnd/lightning"
	bolt11 "github.com/nbd-wtf/ln-decodepay"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Cln struct {
	Host       string `long:"cln.host" description:"gRPC host of the CLN daemon"`
	Port       int    `long:"cln.port" description:"gRPC port of the CLN daemon"`
	RootCert   string `long:"cln.rootcert" description:"Path to the root cert of the CLN gRPC"`
	PrivateKey string `long:"cln.privatekey" description:"Path to the client key of the CLN gRPC"`
	CertChain  string `long:"cln.certchain" description:"Path to the client cert of the CLN gRPC"`

	Client protos.NodeClient
}

const (
	serviceName = lightning.NodeTypeCln

	paymentFailure = "payment failed"

	msatFactor = 1000
)

var (
	ErrPaymentNotInitiated = errors.New("payment not initialized")

	paymentStatusFromGrpc = map[protos.ListpaysPays_ListpaysPaysStatus]lightning.PaymentState{
		protos.ListpaysPays_PENDING:  lightning.PaymentPending,
		protos.ListpaysPays_COMPLETE: lightning.PaymentSucceeded,
		protos.ListpaysPays_FAILED:   lightning.PaymentFailed,
	}
)

func (c *Cln) Connect() error {
	caFile, err := os.ReadFile(c.RootCert)
	if err != nil {
		return fmt.Errorf("could not read %s root certificate %s: %s", serviceName, c.RootCert, err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caFile) {
		return fmt.Errorf("could not parse %s root certificate", serviceName)
	}

	cert, err := tls.LoadX509KeyPair(c.CertChain, c.PrivateKey)
	if err != nil {
		return fmt.Errorf("could not read %s client certificate: %s", serviceName, err)
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   "cln",
		RootCAs:      caPool,
		Certificates: []tls.Certificate{cert},
	})

	con, err := grpc.Dial(c.Host+":"+strconv.Itoa(c.Port), grpc.WithTransportCredentials(creds))
	if err != nil {
		return fmt.Errorf("could not create %s gRPC client: %s", serviceName, err)
	}

	c.Client = protos.NewNodeClient(con)
	return nil
}

func (c *Cln) Name() string {
	return string(serviceName)
}

func (c *Cln) NodeType() lightning.LightningNodeType {
	return serviceName
}

func (c *Cln) GetInfo() (*lightning.LightningInfo, error) {
	info, err := c.Client.Getinfo(context.Background(), &protos.GetinfoRequest{})
	if err != nil {
		return nil, err
	}
	return &lightning.LightningInfo{
		Pubkey:      string(info.Id),
		BlockHeight: info.Blockheight,
		Version:     info.Version,
		Network:     info.Network,
		Synced:      *info.WarningBitcoindSync == "null" && *info.WarningLightningdSync == "null",
	}, nil
}

func (c *Cln) SanityCheck() (string, error) {
	info, err := c.Client.Getinfo(context.Background(), &protos.GetinfoRequest{})
	if err != nil {
		return "", err
	}

	return info.Version, nil
}

func (c *Cln) CreateInvoice(value int64, preimage []byte, expiry int64, memo string) (*lightning.AddInvoiceResponse, error) {
	parsed_expiry := uint64(expiry)
	invoice, err := c.Client.Invoice(context.Background(), &protos.InvoiceRequest{
		// wtf is this
		AmountMsat: &protos.AmountOrAny{
			Value: &protos.AmountOrAny_Amount{
				Amount: &protos.Amount{
					Msat: uint64(value),
				},
			},
		},
		Preimage:    preimage,
		Expiry:      &parsed_expiry,
		Description: memo,
	})

	if err != nil {
		return nil, err
	}

	return &lightning.AddInvoiceResponse{
		PaymentRequest: invoice.Bolt11,
		PaymentHash:    invoice.PaymentHash,
	}, nil
}

func (c *Cln) NewAddress() (string, error) {
	res, err := c.Client.NewAddr(context.Background(), &protos.NewaddrRequest{
		//Addresstype: &protos.NewaddrRequest_BECH32,
	})
	if err != nil {
		return "", err
	}
	return *res.Bech32, nil
}

func (c *Cln) PaymentStatus(preimageHash string) (*lightning.PaymentStatus, error) {
	hash, err := hex.DecodeString(preimageHash)
	if err != nil {
		return nil, err
	}

	res, err := c.Client.ListPays(context.Background(), &protos.ListpaysRequest{
		PaymentHash: hash,
	})
	if err != nil {
		return nil, err
	}

	if len(res.Pays) == 0 {
		return nil, ErrPaymentNotInitiated
	}

	status := res.Pays[len(res.Pays)-1]

	// ListPays doesn't give a proper reason
	var failureReason *string
	if status.Status == protos.ListpaysPays_FAILED {
		reasonStr := paymentFailure
		failureReason = &reasonStr
	}

	return &lightning.PaymentStatus{
		State:         paymentStatusFromGrpc[status.Status],
		FailureReason: failureReason,
		Hash:          encodeOptionalBytes(status.PaymentHash),
		Preimage:      encodeOptionalBytes(status.Preimage),
		FeeMsat:       parseFeeMsat(status.AmountMsat, status.AmountSentMsat),
	}, nil
}

func (c *Cln) SendPayment(invoice string, feeLimit uint64, timeout int32) (<-chan *lightning.PaymentUpdate, error) {
	dec, err := bolt11.Decodepay(invoice)
	if err != nil {
		return nil, err
	}

	updateChan := make(chan *lightning.PaymentUpdate)

	go func() {
		updateChan <- &lightning.PaymentUpdate{
			IsLastUpdate: false,
			Update: lightning.PaymentStatus{
				State: lightning.PaymentPending,
				Hash:  dec.PaymentHash,
			},
		}

		defer close(updateChan)

		riskFactor := float64(0)
		timeoutUint := uint32(timeout)

		res, err := c.Client.Pay(context.Background(), &protos.PayRequest{
			Bolt11:     invoice,
			Riskfactor: &riskFactor,
			RetryFor:   &timeoutUint,
			Maxfee:     &protos.Amount{Msat: feeLimit * msatFactor},
		})
		if err != nil {
			errStr := parseClnError(err)
			updateChan <- &lightning.PaymentUpdate{
				IsLastUpdate: true,
				Update: lightning.PaymentStatus{
					State:         lightning.PaymentFailed,
					Hash:          dec.PaymentHash,
					FailureReason: &errStr,
				},
			}

			return
		}

		var state lightning.PaymentState
		var preimage string
		var failureReason *string

		switch res.Status {
		case protos.PayResponse_COMPLETE:
			state = lightning.PaymentSucceeded
			preimage = hex.EncodeToString(res.PaymentPreimage)

		case protos.PayResponse_FAILED:
			state = lightning.PaymentFailed

			// Pay doesn't give a proper reason
			reasonStr := paymentFailure
			failureReason = &reasonStr

		case protos.PayResponse_PENDING:
			state = lightning.PaymentPending
		}

		updateChan <- &lightning.PaymentUpdate{
			IsLastUpdate: true,
			Update: lightning.PaymentStatus{
				State:         state,
				FailureReason: failureReason,
				Hash:          dec.PaymentHash,
				Preimage:      preimage,
				FeeMsat:       parseFeeMsat(res.AmountMsat, res.AmountSentMsat),
			},
		}
	}()

	return updateChan, nil
}

func encodeOptionalBytes(data []byte) string {
	if data == nil {
		return ""
	}

	return hex.EncodeToString(data)
}

func parseFeeMsat(amountMsat *protos.Amount, amountSentMsat *protos.Amount) uint64 {
	if amountMsat == nil || amountSentMsat == nil {
		return 0
	}

	return amountSentMsat.Msat - amountMsat.Msat
}

func parseClnError(err error) string {
	spl := strings.Split(err.Error(), "\"")
	if len(spl) != 3 {
		return err.Error()
	}

	return spl[1]
}
