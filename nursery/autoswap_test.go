package nursery

import (
	"testing"

	"github.com/BoltzExchange/boltz-lnd/lightning"
	"github.com/stretchr/testify/assert"
)

func TestGetSwapRecommendations(t *testing.T) {

	channels := []*lightning.LightningChannel{
		{
			LocalMsat:  100,
			RemoteMsat: 100,
			Capacity:   200,
			Id:         "1",
		},
		{
			LocalMsat:  50,
			RemoteMsat: 100,
			Capacity:   150,
			Id:         "2",
		},
		{
			LocalMsat:  100,
			RemoteMsat: 50,
			Capacity:   150,
			Id:         "3",
		},
	}

	recommendations := getSwapRecommendations(
		channels,
		0.1,
	)

	assert.Equal(t, 2, len(recommendations))

	assert.Equal(t, &SwapRecommendation{
		Type:    "submarine",
		Amount:  25,
		Channel: channels[1],
	}, recommendations[0])

	assert.Equal(t, &SwapRecommendation{
		Type:    "reversesubmarine",
		Amount:  25,
		Channel: channels[2],
	}, recommendations[1])

}
