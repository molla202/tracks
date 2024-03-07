package blocksync

import (
	"encoding/json"
	"fmt"
	logs "github.com/airchains-network/decentralized-sequencer/log"
	"github.com/airchains-network/decentralized-sequencer/types"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
)

var txDbInstance *leveldb.DB
var blockDbInstance *leveldb.DB
var staticDbInstance *leveldb.DB
var batchesDbInstance *leveldb.DB
var proofDbInstance *leveldb.DB
var publicWitnessDbInstance *leveldb.DB
var daDbInstance *leveldb.DB
var mockDbInstance *leveldb.DB

// InitTxDb This function initializes a LevelDB database for transactions and returns a boolean indicating
// whether the initialization was successful.
func InitTxDb() bool {
	txDB, err := leveldb.OpenFile("data/leveldb/tx", nil)
	if err != nil {
		log.Fatal("Failed to open transaction LevelDB:", err)
		return false
	}
	txDbInstance = txDB
	return true
}

// InitBlockDb This function initializes a LevelDB database for storing blocks and returns a boolean indicating
// whether the initialization was successful.
func InitBlockDb() bool {
	blockDB, err := leveldb.OpenFile("data/leveldb/blocks", nil)
	if err != nil {
		log.Fatal("Failed to open block LevelDB:", err)
		return false
	}
	fmt.Println("blockDB", blockDB)
	blockDbInstance = blockDB

	// get

	// check is its assignes
	blockNumberByte, err := blockDB.Get([]byte("blockCount"), nil)

	if blockNumberByte == nil || err != nil {
		err = blockDB.Put([]byte("blockCount"), []byte("0"), nil)
		if err != nil {
			logs.Log.Error(fmt.Sprintf("Error in saving blockCount in blockDatabase : %s", err.Error()))
			//return false
			os.Exit(0)
		}
	}

	return true
}

// InitStaticDb This function initializes a static LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitStaticDb() bool {
	staticDB, err := leveldb.OpenFile("data/leveldb/static", nil)
	if err != nil {
		log.Fatal("Failed to open static LevelDB:", err)
		return false
	}
	staticDbInstance = staticDB
	return true
}

// InitBatchesDb This function initializes a batches LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitBatchesDb() bool {
	batchesDB, err := leveldb.OpenFile("data/leveldb/batches", nil)
	if err != nil {
		log.Fatal("Failed to open batches LevelDB:", err)
		return false
	}
	batchesDbInstance = batchesDB
	return true
}

// InitProofDb This function initializes a proof LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitProofDb() bool {
	proofDB, err := leveldb.OpenFile("data/leveldb/proof", nil)
	if err != nil {
		log.Fatal("Failed to open proof LevelDB:", err)
		return false
	}
	proofDbInstance = proofDB
	return true
}

func InitPublicWitnessDb() bool {
	publicWitnessDB, err := leveldb.OpenFile("data/leveldb/publicWitness", nil)
	if err != nil {
		log.Fatal("Failed to open public witness LevelDB:", err)
		return false
	}
	publicWitnessDbInstance = publicWitnessDB
	return true
}

func InitDaDb() bool {
	daDB, err := leveldb.OpenFile("data/leveldb/da", nil)
	if err != nil {
		log.Fatal("Failed to open da LevelDB:", err)
		return false
	}
	da := types.DAStruct{
		DAKey:             "0",
		DAClientName:      "0",
		BatchNumber:       "0",
		PreviousStateHash: "0",
		CurrentStateHash:  "0",
	}

	daBytes, err := json.Marshal(da)

	daDbInstance = daDB
	daBytes, err = daDbInstance.Get([]byte("batch_0"), nil)
	if daBytes == nil || err != nil {
		err = daDbInstance.Put([]byte("batch_0"), daBytes, nil)
		if err != nil {
			logs.Log.Error(fmt.Sprintf("Error in saving daBytes in da Database : %s", err.Error()))
			return false
		}
	}

	return true
}
func InitMockDb() bool {
	mockDb, err := leveldb.OpenFile("data/leveldb/mockda", nil)
	if err != nil {
		log.Fatal("Failed to open da LevelDB:", err)
		return false
	}
	mockDbInstance = mockDb
	return true
}

// InitDb This function  initializes three different databases and returns true if all of them are
// successfully initialized, otherwise it returns false.
func InitDb() bool {
	txStatus := InitTxDb()
	blockStatus := InitBlockDb()
	staticStatus := InitStaticDb()
	batchesStatus := InitBatchesDb()
	proofStatus := InitProofDb()
	publicWitnessStatus := InitPublicWitnessDb()
	daDbInstanceStatus := InitDaDb()
	mockDbInstanceStatus := InitMockDb()

	if txStatus && blockStatus && staticStatus && batchesStatus && proofStatus && publicWitnessStatus && daDbInstanceStatus && mockDbInstanceStatus {
		return true
	} else {
		return false
	}
}

// GetTxDbInstance This function returns the instance of the air-leveldb database.
func GetTxDbInstance() *leveldb.DB {
	return txDbInstance
}

// GetBlockDbInstance This function returns the instance of the block database.
func GetBlockDbInstance() *leveldb.DB {
	return blockDbInstance
}

// GetStaticDbInstance This function  is returning the instance of the LevelDB database that was
// initialized in the InitStaticDb function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetStaticDbInstance() *leveldb.DB {
	return staticDbInstance
}

// GetBatchesDbInstance This function  is returning the instance of the LevelDB database that was
// initialized in the InitBatchesDb function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetBatchesDbInstance() *leveldb.DB {
	return batchesDbInstance
}

// GetProofDbInstance This function  is returning the instance of the LevelDB database that was
// initialized in the InitProofDb function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetProofDbInstance() *leveldb.DB {
	return proofDbInstance
}

func GetPublicWitnessDbInstance() *leveldb.DB {
	return publicWitnessDbInstance
}

func GetDaDbInstance() *leveldb.DB {
	return daDbInstance
}

func GetMockDbInstance() *leveldb.DB {
	return mockDbInstance
}
