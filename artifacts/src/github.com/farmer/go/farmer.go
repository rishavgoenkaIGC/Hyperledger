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

type Farmer struct {
	FarmerId string `json:"farmerID"`
	FPOId string `json:"fpoID"`
	Name string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
	Password string `json:"password"`
	DateOfBirth string `json:"dateOfBirth"`
	Gender string `json:"gender"`
	FatherName string `json:"fatherName"`
	Village string `json:"village"`
	Block string `json:"block"`
	District string `json:"district"`
	AadharNumber string `json:"aadharNumber"`
	BankName string `json:"bankName"`
	AccountNumber string `json:"accountNumber"`
	IFSCCode string `json:"IFSCCode"`
	PanNumber string `json:"panNumber"`
	ImageURL string `json:"imageURL"`
	FPOName string `json:"FPOName"`
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Farmer
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	farmers := []Farmer{}

	for i, farmer := range farmers {
		farmerAsBytes, _ := json.Marshal(farmer)
		err := ctx.GetStub().PutState("Farmer"+strconv.Itoa(i), farmerAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) SignUpFarmer(ctx contractapi.TransactionContextInterface, farmerID string, fpoID string, name string, mobileNumber string, password string, dateOfBirth string, gender string, fatherName string, village string, block string, district string, aadharNumber string, bankName string, accountNumber string, IFSCCode string, panNumber string, imageURL string, FPOName string) error {
	farmer := Farmer{
		FarmerId: farmerID,
		FPOId: fpoID,
		Name: name,
		MobileNumber: mobileNumber,
		Password: password,
		DateOfBirth: dateOfBirth,
		Gender: gender,
		FatherName: fatherName,
		Village: village,
		Block: block,
		District: district,
		AadharNumber: aadharNumber,
		BankName: bankName,
		AccountNumber: accountNumber,
		IFSCCode: IFSCCode,
		PanNumber: panNumber,
		ImageURL: imageURL,
		FPOName: FPOName,
	}

	farmerAsBytes, _ := json.Marshal(farmer)

	return ctx.GetStub().PutState(farmerID, farmerAsBytes)
}

func (s *SmartContract) QueryAllFarmers(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "Farmer0"
	endKey := "Farmer99"

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

		farmer := new(Farmer)
		_ = json.Unmarshal(queryResponse.Value, farmer)

		queryResult := QueryResult{Key: queryResponse.Key, Record: farmer}
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
