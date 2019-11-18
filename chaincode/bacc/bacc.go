package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	"time"
	"strconv"


	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


type ChainCode struct {
}

// 유저 구조체
type User struct {
	Phone   string `json:"phone"`
	Battery Battery `json:"battery"`
}

// 배터리 구조체
type Battery struct {
	BatteryStatusStart string `json:"bss"`
	BatteryStatusEnd string `json:"bse"`
	BatteryCount string `json:"bc"`
	Gps string `json:"gps"`
	Date string `json:"date"`
}


func (s *ChainCode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}



func (s *ChainCode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	if function == "addUser" {
		return s.addUser(APIstub, args)
	} else if function == "addBattery" {
		return s.addBattery(APIstub, args)
	} else if function == "getBattery" {
		return s.getBattery(APIstub, args)
	} else if function == "getAllBattery" {
		return s.getAllBattery(APIstub)
	} else if function == "getHistory" {
		return s.getHistory(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// 유저 등록
func (s *ChainCode) addUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("fail!")
	}
	var user = User{Phone: args[0]}
	userAsBytes, _ := json.Marshal(user)
	APIstub.PutState(args[0], userAsBytes)

	return shim.Success(nil)
}

// 데이터 입력
func (s *ChainCode) addBattery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	// 유저 정보 가져오기
	userAsBytes, err := APIstub.GetState(args[0])
	if err != nil{
		jsonResp := "\"Error\":\"Failed to get state for "+ args[0]+"\"}"
		return shim.Error(jsonResp)
	} else if userAsBytes == nil{ // no State! error
		jsonResp := "\"Error\":\"User does not exist: "+ args[0]+"\"}"
		return shim.Error(jsonResp)
	}
	// 상태 확인
	user := User{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 배터리 구조체 값 업데이트
	var data = Battery{BatteryStatusStart: args[1],BatteryStatusEnd: args[2], BatteryCount: args[3] , Gps: args[4], Date: time.Now().Format("20060102150405")}
	user.Battery=data

	// 월드스테이드 업데이트 
	userAsBytes, err = json.Marshal(user);
	APIstub.PutState(args[0], userAsBytes)

	return shim.Success([]byte("rating is updated"))

}

// 키값 데이터 조회
func (s *ChainCode) getBattery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	value, err := APIstub.GetState(args[0])
	if err != nil {
		return shim.Error("Failed to get Battery")
	}
	if value == nil {
		return shim.Error("value not found")
	}
	return shim.Success(value)
}

// 모든 데이터 조회
func (s *ChainCode) getAllBattery(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "00000000000"
	endKey := "999999999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllBatterys:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// 키 이력 조회
func (s *ChainCode) getHistory(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	batteryName := args[0]

	fmt.Printf("- start getHistoryForBattery: %s\n", batteryName)

	resultsIterator, err := stub.GetHistoryForKey(batteryName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()


	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForBattery returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


func main() {
	if err := shim.Start(new(ChainCode)); err != nil {
		fmt.Printf("Error starting ChainCode chaincode: %s", err)
	}
}