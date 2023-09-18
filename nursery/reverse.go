package nursery

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/go-errors/errors"
	"github.com/vulpemventures/go-elements/address"
	"github.com/vulpemventures/go-elements/confidential"
	"github.com/vulpemventures/go-elements/payment"
	"github.com/vulpemventures/go-elements/transaction"
	"strconv"

	"github.com/BoltzExchange/boltz-lnd/boltz"
	"github.com/BoltzExchange/boltz-lnd/boltzrpc"
	"github.com/BoltzExchange/boltz-lnd/database"
	"github.com/BoltzExchange/boltz-lnd/logger"
	"github.com/BoltzExchange/boltz-lnd/utils"
	"github.com/btcsuite/btcd/btcutil"
)

func (nursery *Nursery) recoverReverseSwaps() error {
	logger.Info("Recovering pending Reverse Swaps")

	reverseSwaps, err := nursery.database.QueryPendingReverseSwaps()

	if err != nil {
		return err
	}

	for _, reverseSwap := range reverseSwaps {
		logger.Info("Recovering Reverse Swap " + reverseSwap.Id + " at state: " + reverseSwap.Status.String())

		// TODO: handle race condition when status is updated between the POST request and the time the streaming starts
		status, err := nursery.boltz.SwapStatus(reverseSwap.Id)

		if err != nil {
			logger.Warning("Boltz could not find Reverse Swap " + reverseSwap.Id + ": " + err.Error())
			continue
		}

		if status.Status != reverseSwap.Status.String() {
			logger.Info("Swap " + reverseSwap.Id + " status changed to: " + status.Status)
			nursery.handleReverseSwapStatus(&reverseSwap, *status, nil)

			if reverseSwap.State == boltzrpc.SwapState_PENDING {
				nursery.RegisterReverseSwap(reverseSwap, nil)
			}

			continue
		}

		logger.Info("Reverse Swap " + reverseSwap.Id + " status did not change")
		nursery.RegisterReverseSwap(reverseSwap, nil)
	}

	return nil
}

func (nursery *Nursery) RegisterReverseSwap(reverseSwap database.ReverseSwap, claimTransactionIdChan chan string) chan string {
	logger.Info("Listening to events of Reverse Swap " + reverseSwap.Id)

	go func() {
		stopListening := make(chan bool)
		stopHandler := make(chan bool)

		eventListenersLock.Lock()
		eventListeners[reverseSwap.Id] = stopListening
		eventListenersLock.Unlock()

		eventStream := make(chan *boltz.SwapStatusResponse)

		nursery.streamSwapStatus(reverseSwap.Id, "Reverse Swap", eventStream, stopListening, stopHandler)

		for {
			select {
			case event := <-eventStream:
				logger.Info("Reverse Swap " + reverseSwap.Id + " status update: " + event.Status)
				nursery.handleReverseSwapStatus(&reverseSwap, *event, claimTransactionIdChan)

				// The event listening can stop after the Reverse Swap has succeeded
				if reverseSwap.Status == boltz.InvoiceSettled {
					stopListening <- true
				}

				break

			case <-stopHandler:
				return
			}
		}
	}()

	return claimTransactionIdChan
}

func (nursery *Nursery) constructClaimTxBtc(reverseSwap *database.ReverseSwap, event boltz.SwapStatusResponse) (string, string, error) {

	lockupTransactionRaw, err := hex.DecodeString(event.Transaction.Hex)

	if err != nil {
		return "", "", errors.New("Could not decode lockup transaction: " + err.Error())
	}

	lockupTransaction, err := btcutil.NewTxFromBytes(lockupTransactionRaw)

	if err != nil {
		return "", "", errors.New("Could not parse lockup transaction: " + err.Error())
	}

	lockupAddress, err := boltz.WitnessScriptHashAddress(nursery.chainParams, reverseSwap.RedeemScript)

	if err != nil {
		return "", "", errors.New("Could not derive lockup address: " + err.Error())
	}

	lockupVout, err := nursery.findLockupVout(lockupAddress, lockupTransaction.MsgTx().TxOut)

	if err != nil {
		return "", "", errors.New("Could not find lockup vout of Reverse Swap " + reverseSwap.Id)
	}

	if lockupTransaction.MsgTx().TxOut[lockupVout].Value < int64(reverseSwap.OnchainAmount) {
		logger.Warning("Boltz locked up less onchain coins than expected. Abandoning Reverse Swap")
	}

	logger.Info("Constructing claim transaction for Reverse Swap " + reverseSwap.Id + " with output: " + lockupTransaction.Hash().String() + ":" + strconv.Itoa(int(lockupVout)))

	claimAddress, err := btcutil.DecodeAddress(reverseSwap.ClaimAddress, nursery.chainParams)

	if err != nil {
		return "", "", errors.New("Could not decode claim address of Reverse Swap: " + err.Error())
	}

	feeSatPerVbyte, err := nursery.mempool.GetFeeEstimation()

	if err != nil {
		return "", "", errors.New("Could not get LND fee estimation: " + err.Error())
	}

	logger.Info("Using fee of " + strconv.FormatInt(feeSatPerVbyte, 10) + " sat/vbyte for claim transaction")

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
}

func (nursery *Nursery) findVout(tx *transaction.Transaction, addr string) (uint32, error) {
	for vout, output := range tx.Outputs {
		//p, err := payment.FromScript(output.Script, &network.Regtest, nil)
		_, outputAddresses, _, err := txscript.ExtractPkScriptAddrs(output.Script, nursery.chainParams)

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

const LASSET = "5ac9f65c0efcc4775e0baec4ec03abdde22473cd3cf33c0419ca290e0751b225"

func (nursery *Nursery) constructClaimTxLiquid(reverseSwap *database.ReverseSwap, event boltz.SwapStatusResponse) (string, string, error) {
	net, _ := address.NetworkForAddress(reverseSwap.ClaimAddress)

	lockupTx, err := transaction.NewTxFromHex(event.Transaction.Hex)

	if err != nil {
		return "", "", errors.New("Could not parse lockup transaction: " + err.Error())
	}

	p, err := payment.FromScript(reverseSwap.RedeemScript, net, reverseSwap.BlindingKey.PubKey())
	lockupAddress, err := p.WitnessScriptHash()

	if err != nil {
		return "", "", errors.New("Could not derive lockup address: " + err.Error())
	}

	// TODO
	lockupVout, err := nursery.findVout(lockupTx, lockupAddress)
	lockupVout = 0
	err = nil

	txOut := lockupTx.Outputs[lockupVout]

	output, err := confidential.UnblindOutputWithKey(txOut, reverseSwap.BlindingKey.Serialize())
	if err != nil {
		return "", "", errors.New("Failed to unblind lockup tx: " + err.Error())
	}

	vOutAmount := output.Value

	if vOutAmount < reverseSwap.OnchainAmount {
		logger.Warning("Boltz locked up less onchain coins than expected. Abandoning Reverse Swap")
	}

	if err != nil {
		return "", "", errors.New("Could not find lockup lockupVout of Reverse Swap " + reverseSwap.Id)
	}

	//logger.Info("Constructing claim transaction for Reverse Swap " + reverseSwap.Id + " with output: " + lockupTransaction.Hash().String() + ":" + strconv.Itoa(int(lockupVout)))

	claimAddress, err := address.FromBlech32(reverseSwap.ClaimAddress)

	fmt.Println(claimAddress)

	if err != nil {
		return "", "", errors.New("Could not decode claim address of Reverse Swap: " + err.Error())
	}

	feeSatPerVbyte, err := nursery.mempool.GetFeeEstimation()

	if err != nil {
		return "", "", errors.New("Could not get LND fee estimation: " + err.Error())
	}

	logger.Info("Using fee of " + strconv.FormatInt(feeSatPerVbyte, 10) + " sat/vbyte for claim transaction")

	claimTx := &transaction.Transaction{}

	hash := lockupTx.TxHash()
	claimTx.AddInput(transaction.NewTxInput(hash[:], lockupVout))
	//asset, _ := hex.DecodeString(LASSET)
	//transaction.NewTxOutput(asset, confidential.ValueCommitment(), payment.From)

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
	return "", "", nil
}

// TODO: fail swap after "transaction.failed" event
func (nursery *Nursery) handleReverseSwapStatus(reverseSwap *database.ReverseSwap, event boltz.SwapStatusResponse, claimTransactionIdChan chan string) {
	parsedStatus := boltz.ParseEvent(event.Status)

	if parsedStatus == reverseSwap.Status {
		logger.Info("Status of Reverse Swap " + reverseSwap.Id + " is " + parsedStatus.String() + " already")
		return
	}

	switch parsedStatus {
	case boltz.TransactionMempool:
		fallthrough

	case boltz.TransactionConfirmed:
		if parsedStatus == boltz.TransactionMempool && !reverseSwap.AcceptZeroConf {
			break
		}

		err := nursery.database.SetReverseSwapLockupTransactionId(reverseSwap, event.Transaction.Id)

		if err != nil {
			logger.Error("Could not set lockup transaction id in database: " + err.Error())
			return
		}

		var claimTransactionId, serialized string
		if reverseSwap.PairId == boltz.PairLiquid {
			claimTransactionId, serialized, err = nursery.constructClaimTxLiquid(reverseSwap, event)
		} else {
			claimTransactionId, serialized, err = nursery.constructClaimTxBtc(reverseSwap, event)
		}

		if err != nil {
			logger.Error("Could not construct claim transaction: " + err.Error())
			return
		}

		err = nursery.broadcastTransaction(serialized, utils.CurrencyFromPair(reverseSwap.PairId))

		if err != nil {
			logger.Error("Could not finalize claim transaction: " + err.Error())
			return
		}

		err = nursery.database.SetReverseSwapClaimTransactionId(reverseSwap, claimTransactionId)

		if err != nil {
			logger.Error("Could not set claim transaction id in database: " + err.Error())
			return
		}

		if claimTransactionIdChan != nil {
			claimTransactionIdChan <- claimTransactionId
		}
	}

	err := nursery.database.UpdateReverseSwapStatus(reverseSwap, parsedStatus)

	if err != nil {
		logger.Error("Could not update status of Reverse Swap " + reverseSwap.Id + ": " + err.Error())
	}

	if parsedStatus.IsCompletedStatus() {
		err = nursery.database.UpdateReverseSwapState(reverseSwap, boltzrpc.SwapState_SUCCESSFUL, "")
	} else if parsedStatus.IsFailedStatus() {
		if reverseSwap.State == boltzrpc.SwapState_PENDING {
			err = nursery.database.UpdateReverseSwapState(reverseSwap, boltzrpc.SwapState_SERVER_ERROR, "")
		}
	}

	if err != nil {
		logger.Error("Could not update state of Reverse Swap " + reverseSwap.Id + ": " + err.Error())
	}
}
