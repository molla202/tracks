func SubmitCurrentPod() (success bool) {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	jsonRpc, stationId, accountPath, accountName, addressPrefix, tracks, err := GetJunctionDetails()
	_ = tracks
	if err != nil {
		logs.Log.Error("can not get junctionDetails.json data: " + err.Error())
		return false
	}
	currentPodState := shared.GetPodState()

	podNumber := currentPodState.LatestPodHeight

	LatestPodStatusHash := currentPodState.LatestPodHash
	var LatestPodStatusHashStr string
	LatestPodStatusHashStr = string(LatestPodStatusHash)

	PreviousPodHash := currentPodState.PreviousPodHash
	var PreviousPodStatusHashStr string
	if PreviousPodHash == nil {
		PreviousPodStatusHashStr = ""
	} else {
		PreviousPodStatusHashStr = string(PreviousPodHash)
	}

	witnessByte := currentPodState.LatestPublicWitness

	registry, err := cosmosaccount.New(cosmosaccount.WithHome(accountPath))
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error creating account registry: %v", err))
		return false
	}

	newTempAccount, err := registry.GetByName(accountName)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error getting account: %v", err))
		return false
	}

	newTempAddr, err := newTempAccount.Address(addressPrefix)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error getting address: %v", err))
		return false
	}

	// Get the current account nonce (sequence)
	accountSequence, err := AccountNounceCheck(newTempAddr, jsonRpc)
	if err != nil {
		logs.Log.Error(fmt.Sprintf("Error getting account nonce: %v", err))
		return false
	}

	ctx := context.Background()
	gas := utilis.GenerateRandomWithFavour(100, 300, [2]int{120, 250}, 0.7)
	gasFees := fmt.Sprintf("%damf", gas)
	log.Info().Str("module", "junction").Str("Gas Fees Used to Validate VRF", gasFees)
	accountClient, err := cosmosclient.New(ctx, cosmosclient.WithAddressPrefix(addressPrefix), cosmosclient.WithNodeAddress(jsonRpc), cosmosclient.WithHome(accountPath), cosmosclient.WithGas("auto"), cosmosclient.WithFees(gasFees), cosmosclient.WithSequence(accountSequence))
	if err != nil {
		logs.Log.Error("Switchyard client connection error")
		logs.Log.Error(err.Error())

		return false
	}

	unixTime := time.Now().Unix()
	currentTime := fmt.Sprintf("%d", unixTime)

	msg := types.MsgSubmitPod{
		Creator:                newTempAddr,
		StationId:              stationId,
		PodNumber:              podNumber,
		MerkleRootHash:         LatestPodStatusHashStr,
		PreviousMerkleRootHash: PreviousPodStatusHashStr,
		PublicWitness:          witnessByte,
		Timestamp:              currentTime,
	}

	podDetails := QueryPod(podNumber)
	if podDetails != nil {
		log.Debug().Str("module", "junction").Msg("Pod already submitted")
		return true
	}

	for {
		txRes, errTxRes := accountClient.BroadcastTx(ctx, newTempAccount, &msg)
		if errTxRes != nil {
			errStr := errTxRes.Error()
			log.Error().Str("module", "junction").Str("Error", errStr).Msg("Error in SubmitPod Transaction")

			log.Debug().Str("module", "junction").Msg("Retrying SubmitPod transaction after 10 seconds..")
			time.Sleep(10 * time.Second)
		} else {
			currentPodState.InitPodTxHash = txRes.TxHash
			shared.SetPodState(currentPodState)
			log.Info().Str("module", "junction").Str("txHash", txRes.TxHash).Msg("Pod submitted successfully")
			return true
		}
	}
}
