package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/BoltzExchange/boltz-lnd/boltz"
	"github.com/BoltzExchange/boltz-lnd/boltzrpc"
	"github.com/BoltzExchange/boltz-lnd/database"
	"github.com/BoltzExchange/boltz-lnd/lightning"
	"github.com/BoltzExchange/boltz-lnd/logger"
	"github.com/BoltzExchange/boltz-lnd/utils"
	"github.com/btcsuite/btcd/chaincfg"
)

type Wallet struct {
	MnemonicFile string
	Address      string
}

type SwapConfig struct {
	ChannelImbalanceThreshhold float64    `long:"swap.channel-imbalance-threshhold" description:"Threshhold to determine wheter or not a swap should be initiated"`
	AutoSwap                   bool       `long:"swap.auto-swap" description:"Automatically initiate swaps when a channel is inbalanced"`
	LiquidWallet               string     `long:"swap.liquid-wallet" description:"Seed phrase of liquid wallet to use for swaps"`
	BtcWallet                  string     `long:"swap.btc-wallet" description:"Seed phrase of bitcoin wallet to use for swaps"`
	AcceptZeroConf             bool       `long:"swap.accept-zero-conf" description:"Whether to accept zero conf on auto swaps"`
	Pair                       boltz.Pair `long:"swap.pair" description:"Which pair to use for autoswaps"`
}

func (cfg *SwapConfig) GetAddress(chainParams *chaincfg.Params, pair boltz.Pair) (string, error) {
	switch pair {
	case boltz.PairBtc:
		if cfg.BtcWallet != "" {
			privateKey, err := utils.LoadSeedPhrase(cfg.BtcWallet)
			if err != nil {
				return "", err
			}
			return boltz.PubKeyAddress(chainParams, privateKey.PubKey())
		} else {
			return "", errors.New("No btc wallet configured")
		}
		// TODO: case boltz.PairLiquid:
	}
	return "", errors.New("unknown pair")

}

type SwapRecommendation struct {
	Type    string
	Amount  uint
	Channel *lightning.LightningChannel
}

func getSwapRecommendations(channels []*lightning.LightningChannel, threshhold float64) []*SwapRecommendation {
	var recommendations []*SwapRecommendation

	for _, channel := range channels {
		balancedness := float64(channel.LocalMsat)/float64(channel.Capacity) - 0.5
		var swapType string
		if balancedness < -threshhold {
			swapType = "submarine"
		} else if balancedness > threshhold {
			swapType = "reversesubmarine"
		}
		if swapType != "" {
			recommendations = append(recommendations, &SwapRecommendation{
				Type:    swapType,
				Amount:  uint(math.Abs(float64(channel.LocalMsat) - float64(channel.Capacity)/2)),
				Channel: channel,
			})
		}
	}

	return recommendations
}

func (server *routedBoltzServer) StartChannelWatcher() error {
	cfg := server.swapConfig
	address, err := cfg.GetAddress(server.chainParams, cfg.Pair)
	if err != nil {
		return err
	}
	go func() {
		for {
			response, err := server.GetSwapRecommendations(context.Background(), &boltzrpc.GetSwapRecommendationsRequest{})

			if err != nil {
				logger.Warning("Could not fetch swap recommendations: " + err.Error())
			}
			swaps, err := server.database.QueryPendingSwaps()
			if err != nil {
				logger.Warning("Could not query pending swaps: " + err.Error())
			}

			reverseSwaps, err := server.database.QueryPendingReverseSwaps()
			if err != nil {
				logger.Warning("Could not query pending reverse swaps: " + err.Error())
			}

			for _, recommendation := range response.Swaps {
				logger.Info(fmt.Sprintf("Swap recommendation: %v", recommendation))
				if recommendation.Type == "reversesubmarine" {
					if slices.ContainsFunc(reverseSwaps, func(swap database.ReverseSwap) bool {
						return swap.ChanId == recommendation.Channel.Id
					}) {
						logger.Info("Already created swap for recommendation")
					} else {
						_, err := server.CreateReverseSwap(context.Background(), &boltzrpc.CreateReverseSwapRequest{
							Amount:         int64(recommendation.Amount),
							Address:        address,
							AcceptZeroConf: cfg.AcceptZeroConf,
							PairId:         string(cfg.Pair),
							ChanId:         recommendation.Channel.Id,
						})
						if err != nil {
							logger.Error("Could not act on swap recommendation : " + err.Error())
						}
					}
				} else if recommendation.Type == "submarine" {
					check := func(swap database.Swap) bool { return swap.ChanId == recommendation.Channel.Id }
					if !slices.ContainsFunc(swaps, check) {
						_, err := server.CreateSwap(context.Background(), &boltzrpc.CreateSwapRequest{
							Amount: int64(recommendation.Amount),
							PairId: string(cfg.Pair),
							ChanId: recommendation.Channel.Id,
						})
						if err != nil {
							logger.Error("Could not act on swap recommendation : " + err.Error())
						}
					}
				}
			}

			// TODO: make this configurable
			time.Sleep(10 * time.Second)
		}
	}()
	return nil
}
