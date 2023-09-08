package main

import (
	"strconv"
	"time"

	"github.com/BoltzExchange/boltz-lnd/lightning"
	"github.com/BoltzExchange/boltz-lnd/lnd"
	"github.com/BoltzExchange/boltz-lnd/logger"
)

const retryInterval = 15

var retryMessage = "Retrying in " + strconv.Itoa(retryInterval) + " seconds"

func connectToLnd(lnd *lnd.LND) *lightning.LightningInfo {
	lndInfo, err := lnd.GetInfo()

	if err != nil {
		logger.Warning("Could not connect to LND: " + err.Error())
		logger.Info(retryMessage)
		time.Sleep(retryInterval * time.Second)

		_ = lnd.Connect()
		return connectToLnd(lnd)
	} else {
		return lndInfo
	}
}

func waitForLightningSynced(lightning lightning.LightningNode) {
	info, err := lightning.GetInfo()

	if err == nil {
		if !info.Synced {
			logger.Warning("LND node not synced yet")
			logger.Info(retryMessage)
			time.Sleep(retryInterval * time.Second)

			waitForLightningSynced(lightning)
		}
	} else {
		logger.Error("Could not get LND info: " + err.Error())
		logger.Info(retryMessage)
		time.Sleep(retryInterval * time.Second)

		waitForLightningSynced(lightning)
	}
}
