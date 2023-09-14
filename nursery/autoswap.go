package nursery

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/BoltzExchange/boltz-lnd/lightning"
	"github.com/BoltzExchange/boltz-lnd/logger"
)

type SwapConfig struct {
	ChannelImbalanceThreshhold float64 `long:"swap.channel-imbalance-threshhold" description:"Threshhold to determine wheter or not a swap should be initiated"`
	AutoSwap                   bool    `long:"swap.auto-swap" description:"Automatically initiate swaps when a channel is inbalanced"`
}

func (nursery *Nursery) channelWatcher() {
	for {
		recommendations, err := nursery.GetSwapRecommendations()
		if err != nil {
			logger.Warning("Could not fetch swap recommendations: " + err.Error())
		}

		for _, channel := range recommendations {
			// TODO: initiate swap
			fmt.Println(channel)
		}

		time.Sleep(10 * time.Second)
	}
}

type SwapRecommendation struct {
	Type    string
	Amount  uint
	Channel *lightning.LightningChannel
}

func (nursery *Nursery) GetSwapRecommendations() ([]*SwapRecommendation, error) {
	if nursery.swapConfig == nil || nursery.swapConfig.ChannelImbalanceThreshhold == 0 {
		return nil, errors.New("Channel inbalance threshhold not set")
	}

	channels, err := nursery.lightning.ListChannels()

	if err != nil {
		return nil, err
	}

	return getSwapRecommendations(channels, nursery.swapConfig.ChannelImbalanceThreshhold), nil
}

func getSwapRecommendations(channels []*lightning.LightningChannel, threshhold float64) []*SwapRecommendation {
	var recommendations []*SwapRecommendation

	for _, channel := range channels {
		logger.Info(fmt.Sprint("Channel: ", *channel))
		var swapType string
		if channel.Balancedness() < -threshhold {
			swapType = "submarine"
		} else if channel.Balancedness() > threshhold {
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
