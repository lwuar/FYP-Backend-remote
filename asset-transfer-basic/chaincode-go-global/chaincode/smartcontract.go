package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	GlobalRootID         string `form:"globalRootID" json:"globalRootID" xml:"globalRootID"  binding:"required"`
	LocalChainID         string `form:"localChainID" json:"localChainID" xml:"localChainID"  binding:"required"`
	MerkleTreeRoot       string `form:"merkleTreeRoot" json:"merkleTreeRoot" xml:"merkleTreeRoot"  binding:"required"`
	GlobalChainTxHash    string `form:"globalChainTxHash" json:"globalChainTxHash" xml:"globalChainTxHash"  binding:"required"`
	GlobalChainBlockNum  int64  `form:"globalChainBlockNum" json:"globalChainBlockNum" xml:"globalChainBlockNum"  binding:"required"`
	GlobalChainTimeStamp int64  `form:"globalChainTimeStamp" json:"globalChainTimeStamp" xml:"globalChainTimeStamp"  binding:"required"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{
			GlobalRootID:         "testID",
			LocalChainID:         "testLocalChainID",
			MerkleTreeRoot:       "merkletreeroot test",
			GlobalChainTxHash:    "testGlobalChainTxHash",
			GlobalChainBlockNum:  1,
			GlobalChainTimeStamp: time.Now().Unix(),
		},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.GlobalRootID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface,
	globalRootID string,
	localChainID string,
	merkleTreeRoot string) error {

	exists, err := s.AssetExists(ctx, globalRootID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", globalRootID)
	}

	asset := Asset{
		GlobalRootID:         globalRootID,
		LocalChainID:         localChainID,
		MerkleTreeRoot:       merkleTreeRoot,
		GlobalChainTxHash:    "",
		GlobalChainBlockNum:  1,
		GlobalChainTimeStamp: time.Now().Unix()}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(globalRootID, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given certID.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, globalRootID string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(globalRootID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", globalRootID)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, globalRootID string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(globalRootID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}
