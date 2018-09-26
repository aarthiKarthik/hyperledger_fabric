package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	//"strings"
	//"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("pvt_data_cc0")

type SalesRecord struct {
	Customer  string `json:"customer"`
	Merchant  string `json:"merchant"`
	Product  string `json:"product"`
	PurchaseDate  string `json:"purchaseDate"`
}

type PvtSalesRecord struct {
	Customer  string `json:"customer"`
    Price  string `json:"price"`
}

type SmartContract struct {
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### pvt_data_cc0 Init ###########")
	logger.Info("Returning success")
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "initLedger" {
		return s.initLedger(stub)
	} else if fn == "addSalesAndPvtRecord" {
		return s.addSalesAndPvtRecord(stub, args)
	} else if fn == "getSalesRecords" {
		return s.getSalesRecords(stub, args)
	} else if fn == "getPvtSalesRecords" {
		return s.getPvtSalesRecords(stub, args)
	} //else if fn == "changeCustomer" {
//		return s.changeCustomer(stub, args)
//	}

	return shim.Error("Invalid function name.")
}

func (s *SmartContract) initLedger(stub shim.ChaincodeStubInterface) peer.Response {
	salesRecord := SalesRecord{Customer: "AK", Merchant: "Tokichoi", Product: "Dress", PurchaseDate: "24102018"}
	salesRecordAsBytes, _ := json.Marshal(salesRecord)
	stub.PutPrivateData("salesRecords", "PRE"+strconv.Itoa(0), salesRecordAsBytes)

	pvtRecord := PvtSalesRecord{Customer: "AK", Price: "100"}
	pvtRecordAsBytes, _ := json.Marshal(pvtRecord)
	stub.PutPrivateData("pvtRecords", "PRE"+strconv.Itoa(0), pvtRecordAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) addSalesAndPvtRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var SalesRecBytes []byte
	var PvtSalesRecBytes []byte
	var err error

	// Check arguments
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// Check if given record already exists
	key := "PRE" + args[0]
	data, err := stub.GetPrivateData("salesRecords", key)
	if err != nil {
		return shim.Error(err.Error())
	}
	if data != nil {
		return shim.Error("Record with such key already exists")
	}
	// Add sales records
	salesRec := SalesRecord{Customer: args[0], Merchant: args[1], Product: args[2], PurchaseDate: args[4]}
	SalesRecBytes, err = json.Marshal(salesRec)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("salesRecords", key, SalesRecBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Check if given private record already exists
	DataBytes, err := stub.GetPrivateData("pvtRecords", key)
	if err != nil {
		return shim.Error(err.Error())
	}
	if DataBytes != nil {
		return shim.Error("Private Record with such key already exists")
	}
	// Add private records
	pvtSalesRec := PvtSalesRecord{Customer: args[0], Price: args[3]}
	PvtSalesRecBytes, err = json.Marshal(pvtSalesRec)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("pvtRecords", key, PvtSalesRecBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) getSalesRecords(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//var err error

	// Check arguments
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := "PRE" + args[0]
	var salesRecordBytes []byte
	//var err error

	salesRecordBytes, _ = stub.GetPrivateData("salesRecords", key)
	if salesRecordBytes != nil {
			return shim.Success(salesRecordBytes)
	}

	return shim.Error("Sales Record not found")
}

func (s *SmartContract) getPvtSalesRecords(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// Check arguments
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := "PRE" + args[0]
	var pvtSalesRecordBytes []byte
	//var err error

	pvtSalesRecordBytes, _ = stub.GetPrivateData("pvtRecords", key)
	if pvtSalesRecordBytes != nil {
			return shim.Success(pvtSalesRecordBytes)
	}
	
	return shim.Error("Private Record not found")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting SmartContract chaincode: %s", err)
	}
}



