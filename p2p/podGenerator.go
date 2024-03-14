package p2p

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/airchains-network/decentralized-sequencer/config"
	logs "github.com/airchains-network/decentralized-sequencer/log"
	"github.com/airchains-network/decentralized-sequencer/node/shared"
	"github.com/airchains-network/decentralized-sequencer/types"
	"github.com/airchains-network/decentralized-sequencer/utilis"
	v1 "github.com/airchains-network/decentralized-sequencer/zk/v1"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func BatchGeneration(wg *sync.WaitGroup) {
	defer wg.Done()
	GenerateUnverifiedPods()
}

func GenerateUnverifiedPods() {

	lds := shared.Node.NodeConnections.GetStaticDatabaseConnection()
	ldt := shared.Node.NodeConnections.GetTxnDatabaseConnection()

	ConfirmendTransactionIndex, err := lds.Get([]byte("batchStartIndex"), nil)
	if err != nil {
		err = lds.Put([]byte("batchStartIndex"), []byte("0"), nil)
		if err != nil {
			logs.Log.Error(fmt.Sprintf("Error in saving batchStartIndex in static db : %s", err.Error()))
			os.Exit(0)
		}
	}

	currentPodNumber, err := lds.Get([]byte("batchCount"), nil)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error in getting sssssss from static db : %s", err.Error()))
		os.Exit(0)
	}

	// @aakash please passs here the Last TrackAppHash of the Pods
	previousStateData, err := getPodStateFromDatabase()
	if err != nil {
		logs.Log.Error("Error in getting previous station data")
		os.Exit(0)
	}
	PreviousTrackAppHash := previousStateData.TracksAppHash

	SelectedMaster := MasterTracksSelection(Node, string(PreviousTrackAppHash))
	decodedMaster, err := peer.Decode(SelectedMaster)

	currentPodNumberInt, _ := strconv.Atoi(strings.TrimSpace(string(currentPodNumber)))
	batchNumber := currentPodNumberInt + 1

	//var batchInput *types.BatchStruct
	Witness, uZKP, MRH, batchInput, err := createPOD(ldt, ConfirmendTransactionIndex, currentPodNumber)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error in creating POD : %s", err.Error()))
		os.Exit(0)
	}

	TrackAppHash := generatePodHash(Witness, uZKP, MRH, currentPodNumber)
	podState := shared.GetPodState()
	tempMasterTrackAppHash := podState.MasterTrackAppHash

	// update pod state as per latest pod
	updateNewPodState(TrackAppHash, Witness, uZKP, MRH, uint64(batchNumber), batchInput)

	// Here the MasterTrack Will Broadcast the uZKP in the Network

	if decodedMaster == Node.ID() {

		// make master's vote true by default
		podState := shared.GetPodState()
		currentVotes := podState.Votes
		currentVotes[decodedMaster.String()] = shared.Votes{
			PeerID: decodedMaster.String(),
			Vote:   true,
		}
		podState.Votes = currentVotes
		shared.SetPodState(podState)
		Peers := peerList.GetPeers()
		peerCount := len(Peers)
		if peerCount == 0 {
			// if (no peers connected): update database and make next pod without voting process
			saveVerifiedPOD()        // save data to database
			GenerateUnverifiedPods() // generate next pod
		} else {
			// Preparing the Message that master track will gossip to the Network
			proofData := ProofData{
				PodNumber:    uint64(batchNumber),
				TrackAppHash: TrackAppHash,
			}

			// Marshal the proofData
			proofDataByte, err := json.Marshal(proofData)
			if err != nil {
				logs.Log.Error(fmt.Sprintf("Error in marshalling proof data : %s", err.Error()))
			}

			gossipMsg := types.GossipData{
				Type: "proof",
				Data: proofDataByte,
			}

			gossipMsgByte, err := json.Marshal(gossipMsg)
			if err != nil {
				logs.Log.Error("Error marshaling gossip message")
				return
			}

			logs.Log.Info("Sending proof result: %s")
			BroadcastMessage(context.Background(), Node, gossipMsgByte)
		}

	} else {
		if podState.MasterTrackAppHash != nil {
			fmt.Println(TrackAppHash)
			fmt.Println(tempMasterTrackAppHash)
			currentPodData := shared.GetPodState()
			if bytes.Equal(TrackAppHash, tempMasterTrackAppHash) {
				SendValidProof(CTX, currentPodData.LatestPodHeight, decodedMaster)
				return
			} else {
				SendInvalidProofError(CTX, currentPodData.LatestPodHeight, decodedMaster)
				return
			}
		} else {
			// pod state is nil, means master track has not yet broadcasted the proof
			// don't need to do anything..
		}
	}

}

func createPOD(ldt *leveldb.DB, batchStartIndex []byte, limit []byte) (witness []byte, unverifiedProof []byte, MRH []byte, podData *types.BatchStruct, err error) {

	limitInt, _ := strconv.Atoi(strings.TrimSpace(string(limit)))

	batchStartIndexInt, _ := strconv.Atoi(strings.TrimSpace(string(batchStartIndex)))

	fmt.Println(limitInt)
	fmt.Println(batchStartIndexInt)

	var batch types.BatchStruct

	var From []string
	var To []string
	var Amounts []string
	var TransactionHash []string
	var SenderBalances []string
	var ReceiverBalances []string
	var Messages []string
	var TransactionNonces []string
	var AccountNonces []string

	for i := batchStartIndexInt; i < (config.PODSize * (limitInt + 1)); i++ {

		findKey := fmt.Sprintf("txns-%d", i+1)
		txData, err := ldt.Get([]byte(findKey), nil)
		if err != nil {
			i--
			time.Sleep(1 * time.Second)
			continue
		}
		var tx types.TransactionStruct
		err = json.Unmarshal(txData, &tx)
		if err != nil {
			logs.Log.Error(fmt.Sprintf("Error in unmarshalling tx data : %s", err.Error()))
			os.Exit(0)
		}

		senderBalancesCheck, err := utilis.GetBalance(tx.From, (tx.BlockNumber - 1))
		if err != nil {
			logs.Log.Error(fmt.Sprintf("Error in getting sender balance : %s", err.Error()))
			os.Exit(0)
		}

		receiverBalancesCheck, err := utilis.GetBalance(tx.To, (tx.BlockNumber - 1))
		if err != nil {
			logs.Log.Error(fmt.Sprintf("Error in getting reciver balance : %s", err.Error()))
			os.Exit(0)
		}

		accountNouceCheck, err := utilis.GetAccountNonce(context.Background(), tx.Hash, tx.BlockNumber)
		if err != nil {
			logs.Log.Error(fmt.Sprintf("Error in getting account nonce : %s", err.Error()))
			os.Exit(0)
		}

		From = append(From, tx.From)
		To = append(To, tx.To)
		Amounts = append(Amounts, tx.Value)
		TransactionHash = append(TransactionHash, tx.Hash)
		SenderBalances = append(SenderBalances, senderBalancesCheck)
		ReceiverBalances = append(ReceiverBalances, receiverBalancesCheck)
		Messages = append(Messages, tx.Input)
		TransactionNonces = append(TransactionNonces, tx.Nonce)
		AccountNonces = append(AccountNonces, accountNouceCheck)
	}

	batch.From = From
	batch.To = To
	batch.Amounts = Amounts
	batch.TransactionHash = TransactionHash
	batch.SenderBalances = SenderBalances
	batch.ReceiverBalances = ReceiverBalances
	batch.Messages = Messages
	batch.TransactionNonces = TransactionNonces
	batch.AccountNonces = AccountNonces

	witnessVector, currentStatusHash, proofByte, pkErr := v1.GenerateProof(batch, limitInt+1)
	if pkErr != nil {
		logs.Log.Error(fmt.Sprintf("Error in generating proof : %s", pkErr.Error()))
		return nil, nil, nil, nil, pkErr
	}
	logs.Log.Warn(fmt.Sprintf("Successfully generated  Unverified proof for Batch %s in the latest phase", strconv.Itoa(limitInt+1)))

	// marshal witnessVector
	witnessVectorByte, err := json.Marshal(witnessVector)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error in marshalling witness vector : %s", err.Error()))
	}

	// string to []byte currentStatusHash
	currentStatusHashByte, err := json.Marshal(currentStatusHash)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error in marshalling current status hash : %s", err.Error()))
		os.Exit(0)
	}

	return witnessVectorByte, proofByte, currentStatusHashByte, &batch, nil
}

func saveVerifiedPOD() {

	podState := shared.GetPodState()
	batchInput := podState.Batch
	currentPodNumber := podState.LatestPodHeight
	currentPodNumberInt := int(currentPodNumber)
	batchJSON, err := json.Marshal(batchInput)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error in marshalling batch data : %s", err.Error()))
		os.Exit(0)
	}
	ldbatch := shared.Node.NodeConnections.GetDataAvailabilityDatabaseConnection()
	lds := shared.Node.NodeConnections.GetStaticDatabaseConnection()
	batchKey := fmt.Sprintf("batch-%d", currentPodNumberInt)
	err = ldbatch.Put([]byte(batchKey), batchJSON, nil)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error in writing batch data to file : %s", err.Error()))
		os.Exit(0)
	}
	err = lds.Put([]byte("batchStartIndex"), []byte(strconv.Itoa(config.PODSize*(currentPodNumberInt))), nil)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error in updating batchStartIndex in static db : %s", err.Error()))
		os.Exit(0)
	}

	err = lds.Put([]byte("batchCount"), []byte(strconv.Itoa(currentPodNumberInt)), nil)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error in updating batchCount in static db : %s", err.Error()))
		os.Exit(0)
	}

	err = os.WriteFile("data/batchCount.txt", []byte(strconv.Itoa(currentPodNumberInt)), 0666)
	if err != nil {
		panic("Failed to update batch number: " + err.Error())
	}

	podState.MasterTrackAppHash = nil
	shared.SetPodState(podState)
}

func generatePodHash(Witness, uZKP, MRH []byte, podNumber []byte) []byte {

	return MRH
}

func updateNewPodState(CombinedPodHash, Witness, uZKP, MRH []byte, podNumber uint64, batchInput *types.BatchStruct) {
	var podState *shared.PodState
	// empty votes
	votes := make(map[string]shared.Votes)
	fmt.Println("Hey Saving the States with VIOte")
	fmt.Println(podNumber)
	podState = &shared.PodState{
		LatestPodHeight:     podNumber,
		LatestPodHash:       MRH,
		LatestPodProof:      uZKP,
		LatestPublicWitness: Witness,
		Votes:               votes,
		TracksAppHash:       CombinedPodHash,
		Batch:               batchInput,
	}

	// save pod state to database
	shared.SetPodState(podState)

	// save pod data in local state
	updatePodStateInDatabase(podState)

}

func updatePodStateInDatabase(podState *shared.PodState) {
	stateConnection := shared.Node.NodeConnections.GetStateDatabaseConnection()

	podStateByte, err := json.Marshal(podState)
	if err != nil {
		logs.Log.Error(err.Error())
		os.Exit(0)
	}

	err = stateConnection.Put([]byte("podState"), podStateByte, nil)
	if err != nil {
		logs.Log.Error(err.Error())
		os.Exit(0)
	}
}

func getPodStateFromDatabase() (podStateData *shared.PodState, err error) {
	stateConnection := shared.Node.NodeConnections.GetStateDatabaseConnection()

	podStateDataByte, err := stateConnection.Get([]byte("podState"), nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(podStateDataByte, podStateData)
	if err != nil {
		return nil, err
	}

	return podStateData, nil

}