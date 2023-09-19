package liquid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiquidClient(t *testing.T) {
	cookie := "/home/jacksn/repos/lnbits-legend/docker/data/elements/liquidregtest/.cookie"
	client := Connect("http://localhost:18884", cookie)

	info, err := client.GetBlockchainInfo()

	assert.NoError(t, err)
	assert.Equal(t, "liquidregtest", info.Chain)
}
