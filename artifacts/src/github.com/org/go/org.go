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

type Org struct {
	Name string `json:"name"`
	CompanyName string `json:"companyname"`
	PhoneNo string `json:"phoneno"`
	Email string `json:"email"`
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Org
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	orgs := []Org{}

	for i, org := range orgs {
		orgAsBytes, _ := json.Marshal(org)
		err := ctx.GetStub().PutState("ORG"+strconv.Itoa(i), orgAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) CreateOrg(ctx contractapi.TransactionContextInterface, name string, companyname string, phoneno string, email string) error {
	org := Org{
		Name: name,
		CompanyName: companyname,
		PhoneNo: phoneno,
		Email: email,
	}

	orgAsBytes, _ := json.Marshal(org)

	return ctx.GetStub().PutState(name, orgAsBytes)
}

func (s *SmartContract) QueryAllOrgs(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "ORG0"
	endKey := "ORG99"

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

		org := new(Org)
		_ = json.Unmarshal(queryResponse.Value, org)

		queryResult := QueryResult{Key: queryResponse.Key, Record: org}
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
