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

type Sammunati struct {
	SammunatiId string `json:"sammunatiID"`
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
	Record *Sammunati
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	sammunatis := []Sammunati{}

	for i, sammunati := range sammunatis {
		sammunatiAsBytes, _ := json.Marshal(sammunati)
		err := ctx.GetStub().PutState("Sammunati"+strconv.Itoa(i), sammunatiAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) SignUpSammunati(ctx contractapi.TransactionContextInterface, sammunatiID string, username string, password string, address string, city string, state string, pincode string, contactNumber string, email string, website string, licenseKey string) error {
	sammunati := Sammunati{
		SammunatiId: sammunatiID,
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

	sammunatiAsBytes, _ := json.Marshal(sammunati)

	return ctx.GetStub().PutState(sammunatiID, sammunatiAsBytes)
}

func (s *SmartContract) QueryAllSammunatis(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "Sammunati0"
	endKey := "Sammunati99"

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

		sammunati := new(Sammunati)
		_ = json.Unmarshal(queryResponse.Value, sammunati)

		queryResult := QueryResult{Key: queryResponse.Key, Record: sammunati}
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
