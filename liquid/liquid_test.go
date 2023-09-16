package liquid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLiquidClient(t *testing.T) {
	cookie := "/home/jacksn/repos/lnbits-legend/docker/data/elements/liquidregtest/.cookie"
	client := Connect("http://localhost:18884", cookie)

	info, err := client.GetBlockchainInfo()

	assert.NoError(t, err)
	assert.Equal(t, "liquidregtest", info.Chain)
}
