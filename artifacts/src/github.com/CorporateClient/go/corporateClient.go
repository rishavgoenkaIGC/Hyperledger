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

type CorporateClient struct {
	CorporateClientId string `json:"corporateClientID"`
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
	Record *CorporateClient
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	corporateClients := []CorporateClient{}

	for i, corporateClient := range corporateClients {
		corporateClientAsBytes, _ := json.Marshal(corporateClient)
		err := ctx.GetStub().PutState("CorporateClient"+strconv.Itoa(i), corporateClientAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) SignUpCorporateClient(ctx contractapi.TransactionContextInterface, corporateClientID string, username string, password string, address string, city string, state string, pincode string, contactNumber string, email string, website string, licenseKey string) error {
	corporateClient := CorporateClient{
		CorporateClientId: corporateClientID,
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

	corporateClientAsBytes, _ := json.Marshal(corporateClient)

	return ctx.GetStub().PutState(corporateClientID, corporateClientAsBytes)
}

func (s *SmartContract) QueryAllCorporateClients(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "CorporateClient0"
	endKey := "CorporateClient99"

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

		corporateClient := new(CorporateClient)
		_ = json.Unmarshal(queryResponse.Value, corporateClient)

		queryResult := QueryResult{Key: queryResponse.Key, Record: corporateClient}
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
