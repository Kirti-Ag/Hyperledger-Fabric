package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


type SmartContract struct {
	contractapi.Contract
}

type AccessLog struct {
	ID       string    `json:"ID"`
	DateTime string    `json:"datetime"`
	ReqType  string    `json:"reqtype"`
	Path     string    `json:"path"`
}

// WriteLog creates a new log entry to the world state with given details.
func (s *SmartContract) WriteLog(ctx contractapi.TransactionContextInterface, ip string, datetime string, reqtype string, path string) error {
	id := ip + "_" + datetime
	exists, err := s.LogExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the entry %s already exists", id)
	}

	accesslog := AccessLog{
		ID:       id,
		DateTime: datetime,
		ReqType:  reqtype,
		Path:     path,
	}
	accessJSON, err := json.Marshal(accesslog)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, accessJSON)
}

// ReadLog returns the asset stored in the world state with given id.
func (s *SmartContract) ReadLog(ctx contractapi.TransactionContextInterface, id string) (*AccessLog, error) {
	accessJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if accessJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var accesslog AccessLog
	err = json.Unmarshal(accessJSON, &accesslog)
	if err != nil {
		return nil, err
	}

	return &accesslog, nil
}

// LogExists returns true when asset with given ID exists in world state
func (s *SmartContract) LogExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	accessJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return accessJSON != nil, nil
}

func main() {
    assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
      log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
    }

    if err := assetChaincode.Start(); err != nil {
      log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
    }
}


