package nursery

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/BoltzExchange/boltz-lnd/boltzrpc"
	"github.com/BoltzExchange/boltz-lnd/database"
	"github.com/BoltzExchange/boltz-lnd/logger"
	"github.com/stretchr/testify/assert"
	"github.com/vulpemventures/go-elements/network"
)

func GetSwapAndTx() (*database.ReverseSwap, string) {
	bytesTx, _ := os.ReadFile("data/test.json")
	var test map[string]interface{}
	json.Unmarshal(bytesTx, &test)
	lockupTx := test["txHex"].(string)

	swap := database.ReverseSwap{
		Id:                  "v7UI3A",
		PairId:              "L-BTC/BTC",
		ChanId:              "",
		State:               boltzrpc.SwapState_PENDING,
		Error:               "",
		AcceptZeroConf:      true,
		Invoice:             "lnbcrt2500u1pjsnd2xpp5fvp4ugmzckef52sp6zcarvfqpxu4t78vznhalmvvcapsksyyrdxqdpz2djkuepqw3hjqnpdgf2yxgrpv3j8qz95xqrrsssp5vzqlz5jngjg25xw6pc7awxn72kllavkfdmcdlxjy9k4k8zhctrtq9qyyssqqsqcvz75ykvhzy608fpwq64e788a4t4td6098vwjapxysnc06hj9ynl6cu5uev6dhxd3p6yczuhlval5tjupa76y5hdjpzvv9rzv9wspfr29z8",
		ClaimAddress:        "el1qqg8dx39yxme2u5sus57h904ewwc66d0qa8prgjtlrq3jflfhkqxvsp0wuzm4e602ptzywta04akj90tp7drquu0gdax8w7njc",
		OnchainAmount:       248474,
		TimeoutBlockHeight:  1596,
		LockupTransactionId: "",
		ClaimTransactionId:  "",
	}
	swap.Preimage, _ = hex.DecodeString("43ead68753ae7a74050520cd8fcc1ad7428d38c849f9637e14ff5ca37035250b")
	swap.PrivateKey, _ = database.ParsePrivateKey("b9de9557297f6c8ad3922c15d8fc3d53996ffcaff4f1d5759dfb7aca02cb1c12")
	swap.RedeemScript, _ = hex.DecodeString("8201208763a914ce1535e921dc39282543b2a49c892163db134c3f882102d691618b230bbf33e971ff482b8e7dbd51a3335379fabc2300221839645a67a06775023c06b1752102acdac1115c5d83983069d08742b612126ed6fa181c00e95d46dd0991c75697e468ac")
	swap.BlindingKey, _ = database.ParsePrivateKey("263d898560071f5a50977c3017d644d85b97cce19ff2b98e70e7753278eb32ed")
	return &swap, lockupTx
}

func TestLiquidClaimTx(t *testing.T) {
	logger.InitLogger("test.log", "[TEST] ")

	fee := int64(0)
	swap, lockupTx := GetSwapAndTx()

	_, _, err := constructClaimTxLiquid(swap, lockupTx, fee)

	assert.NoError(t, err)
}

func TestLockupAddr(t *testing.T) {
	logger.InitLogger("test.log", "[TEST] ")

	swap, _ := GetSwapAndTx()

	lockupAddress, err := getLockupAddress(swap, &network.Regtest)

	assert.NoError(t, err)
	assert.Equal(t, "el1qqdx8rxlpqfhhyfxn355t2sl8s565jkndyy7u3e6glkhuh43crxryrr2peh0cutzqnhtjtt9aekexe89ukqll2nw2q6kuz0shu4n6hf6r4xp", lockupAddress)
}
