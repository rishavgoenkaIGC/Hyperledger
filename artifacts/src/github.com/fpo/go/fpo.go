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

type FPO struct {
	FPOId string `json:"fpoID"`
	Name string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Address string `json:"address"`
	City  string `json:"city"`
	State string `json:"state"`
	Pincode string `json:"pincode"`
	ContactNumber string `json:"contactNumber"`
	Email string `json:"email"`
	Website string `json:"website"`
	DateOfIncorporation string `json:"dateOfIncorporation"`
	PanNumber string `json:"panNumber"`
	RegistrationNumber string `json:"registrationNumber"`
	NumberOfShareholders string `json:"numberOfShareholders"`
	BankNumber string `json:"bankNumber"`
	AccountNumber string `json:"accountNumber"`
	IFSCCode string `json:"IFSCCode"`
	BankImageURL string `json:"bankImageURL"`
	DirectorName string `json:"directorName"`
	DirectorMobileNumber string `json:"directorMobileNumber"`
	DirectorGender string `json:"directorGender"`
	CopyOfBylaws string `json:"copyOfBylaws"`
	CertificateOfIncorporation string `json:"certificateOfIncorporation"`
	GSTCertificate string `json:"GSTCertificate"`
	LicenseKey string `json:"licenseKey"`
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *FPO
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fpos := []FPO{}

	for i, fpo := range fpos {
		fpoAsBytes, _ := json.Marshal(fpo)
		err := ctx.GetStub().PutState("FPO"+strconv.Itoa(i), fpoAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) SignUpFPO(ctx contractapi.TransactionContextInterface, fpoID string, name string, username string, password string, address string, city string, state string, pincode string, contactNumber string, email string, website string, dateOfIncorporation string, panNumber string, registrationNumber string, numberOfShareholders string, bankNumber string, accountNumber string, IFSCCode string, bankImageURL string, directorName string, directorMobileNumber string, directorGender string, copyOfBylaws string, certificateOfIncorporation string, GSTCertificate string, licenseKey string) error {
	fpo := FPO{
		FPOId: fpoID,
		Name: name,
		Username: username,
		Password: password,
		Address: address,
		City: city,
		State: state,
		Pincode: pincode,
		ContactNumber: contactNumber,
		Email: email,
		Website: website,
		DateOfIncorporation: dateOfIncorporation,
		PanNumber: panNumber,
		RegistrationNumber: registrationNumber,
		NumberOfShareholders: numberOfShareholders,
		BankNumber: bankNumber,
		AccountNumber: accountNumber,
		IFSCCode: IFSCCode,
		BankImageURL: bankImageURL,
		DirectorName: directorName,
		DirectorMobileNumber: directorMobileNumber,
		DirectorGender: directorGender,
		CopyOfBylaws: copyOfBylaws,
		CertificateOfIncorporation: certificateOfIncorporation,
		GSTCertificate: GSTCertificate,
		LicenseKey: licenseKey,
	}

	fpoAsBytes, _ := json.Marshal(fpo)

	return ctx.GetStub().PutState(fpoID, fpoAsBytes)
}

func (s *SmartContract) QueryAllFPOs(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "FPO0"
	endKey := "FPO99"

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

		fpo := new(FPO)
		_ = json.Unmarshal(queryResponse.Value, fpo)

		queryResult := QueryResult{Key: queryResponse.Key, Record: fpo}
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
