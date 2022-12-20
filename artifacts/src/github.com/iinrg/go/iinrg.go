package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type IINRG struct {
	IINRGId string `json:"iinrgID"`
	Username string `json:"username"`
	Password string `json:"password"`
	Address string `json:"address"`
	City  string `json:"city"`
	State string `json:"state"`
	Pincode string `json:"pincode"`
	ContactNumber string `json:"contactNumber"`
	Email string `json:"email"`
	Website string `json:"website"`
	LicenseKey string `json:"licenseKey"`
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *IINRG
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	iinrgs := []IINRG{}

	for i, iinrg := range iinrgs {
		iinrgAsBytes, _ := json.Marshal(iinrg)
		err := ctx.GetStub().PutState("IINRG"+strconv.Itoa(i), iinrgAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) SignUpIINRG(ctx contractapi.TransactionContextInterface, iinrgID string, username string, password string, address string, city string, state string, pincode string, contactNumber string, email string, website string, licenseKey string) error {
	iinrg := IINRG{
		IINRGId: iinrgID,
		Username: username,
		Password: password,
		Address: address,
		City: city,
		State: state,
		Pincode: pincode,
		ContactNumber: contactNumber,
		Email: email,
		Website: website,
	}

	iinrgAsBytes, _ := json.Marshal(iinrg)

	return ctx.GetStub().PutState(iinrgID, iinrgAsBytes)
}

func (s *SmartContract) QueryAllIINRGs(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "IINRG0"
	endKey := "IINRG99"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		iinrg := new(IINRG)
		_ = json.Unmarshal(queryResponse.Value, iinrg)

		queryResult := QueryResult{Key: queryResponse.Key, Record: iinrg}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create org chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting org chaincode: %s", err.Error())
	}
}
