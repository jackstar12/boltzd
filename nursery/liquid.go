package nursery

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/BoltzExchange/boltz-lnd/database"
	"github.com/BoltzExchange/boltz-lnd/logger"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/vulpemventures/go-elements/address"
	"github.com/vulpemventures/go-elements/confidential"
	"github.com/vulpemventures/go-elements/network"
	"github.com/vulpemventures/go-elements/payment"
	"github.com/vulpemventures/go-elements/transaction"
)

const LASSET = "5ac9f65c0efcc4775e0baec4ec03abdde22473cd3cf33c0419ca290e0751b225"

func findVout(tx *transaction.Transaction, addr string) (uint32, error) {
	for vout, output := range tx.Outputs {
		//p, err := payment.FromScript(output.Script, &network.Regtest, nil)
		_, outputAddresses, _, err := txscript.ExtractPkScriptAddrs(output.Script, &chaincfg.RegressionNetParams)

		// Just ignore outputs we can't decode
		if err != nil {
			continue
		}

		for _, outputAddress := range outputAddresses {
			if outputAddress.EncodeAddress() == addr {
				return uint32(vout), nil
			}
		}
	}

	return 0, errors.New("could not find lockup vout")
}

func getLockupAddress(reverseSwap *database.ReverseSwap, net *network.Network) (string, error) {
	p, err := payment.FromScript(reverseSwap.RedeemScript, net, reverseSwap.BlindingKey.PubKey())
	if err != nil {
		return "", err
	}

	lockupAddress, err := p.ConfidentialWitnessScriptHash()
	if err != nil {
		return "", err
	}

	return lockupAddress, nil
}

func constructClaimTxLiquid(reverseSwap *database.ReverseSwap, lockupTxHex string, feeSatPerVbyte int64) (string, string, error) {
	blindingKey := reverseSwap.BlindingKey
	net, _ := address.NetworkForAddress(reverseSwap.ClaimAddress)

	lockupTx, err := transaction.NewTxFromHex(lockupTxHex)

	if err != nil {
		return "", "", errors.New("Could not parse lockup transaction: " + err.Error())
	}

	lockupAddress, err := getLockupAddress(reverseSwap, net)

	/*
		lockupAddress, err := p.ConfidentialWitnessScriptHash()
		logger.Info("Lockup address: " + lockupAddress)
		lockupAddress, err = p.ConfidentialWitnessPubKeyHash()
		logger.Info("Lockup address: " + lockupAddress)
		lockupAddress, err = p.ConfidentialPubKeyHash()
		logger.Info("Lockup address: " + lockupAddress)
		lockupAddress, err = p.ConfidentialScriptHash()
		logger.Info("Lockup address: " + lockupAddress)

		lockupAddress, err = p.ConfidentialWitnessScriptHash()
	*/
	logger.Info("Lockup address: " + lockupAddress)

	//err = nil

	if err != nil {
		return "", "", errors.New("Could not derive lockup address: " + err.Error())
	}

	// TODO
	lockupVout, err := findVout(lockupTx, lockupAddress)
	lockupVout = 0
	err = nil

	txOut := lockupTx.Outputs[lockupVout]

	p, err := payment.FromScript(txOut.Script, net, blindingKey.PubKey())
	lockupAddress, err = p.ConfidentialWitnessScriptHash()
	logger.Info("Lockup address: " + lockupAddress)

	output, err := confidential.UnblindOutputWithKey(txOut, reverseSwap.BlindingKey.Serialize())
	if err != nil {
		return "", "", errors.New("Failed to unblind lockup tx: " + err.Error())
	}

	vOutAmount := output.Value

	logger.Info(fmt.Sprintf("Unblineded output value %v", vOutAmount))

	if vOutAmount < reverseSwap.OnchainAmount {
		logger.Warning("Boltz locked up less onchain coins than expected. Abandoning Reverse Swap")
	}

	if err != nil {
		return "", "", errors.New("Could not find lockup lockupVout of Reverse Swap " + reverseSwap.Id)
	}

	//logger.Info("Constructing claim transaction for Reverse Swap " + reverseSwap.Id + " with output: " + lockupTransaction.Hash().String() + ":" + strconv.Itoa(int(lockupVout)))

	_, err = address.FromBlech32(reverseSwap.ClaimAddress)

	if err != nil {
		return "", "", errors.New("Could not decode claim address of Reverse Swap: " + err.Error())
	}

	claimTx := &transaction.Transaction{}

	hash := lockupTx.TxHash()
	claimTx.AddInput(transaction.NewTxInput(hash[:], lockupVout))

	// asset, _ := hex.DecodeString(LASSET)
	// generator := confidential.NewZKPGeneratorFromBlindingKeys()
	// transaction.NewTxOutput(asset, generator., payment.From)

	/*
		claimTransaction, err := boltz.ConstructTransaction(
			[]boltz.OutputDetails{
				{
					LockupTransaction: lockupTransaction,
					Vout:              lockupVout,
					OutputType:        boltz.SegWit,
					RedeemScript:      reverseSwap.RedeemScript,
					PrivateKey:        reverseSwap.PrivateKey,
					Preimage:          reverseSwap.Preimage,
				},
			},
			claimAddress,
			feeSatPerVbyte,
		)

		if err != nil {
			return "", "", err
		}

		serialized, err := boltz.SerializeTransaction(claimTransaction)

		if err != nil {
			return "", "", errors.New("Could not serialize tx: " + err.Error())
		}

		return claimTransaction.TxHash().String(), serialized, nil
	*/

	claimTxBytes, err := claimTx.Serialize()

	if err != nil {
		return "", "", errors.New("Could not serialize tx: " + err.Error())
	}

	return claimTx.TxHash().String(), hex.EncodeToString(claimTxBytes), nil
}
