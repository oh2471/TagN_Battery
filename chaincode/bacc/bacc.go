package main

import (
	"encoding/json"
	"fmt"
	// "bytes"
	"time"
	// "strconv"


	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)


type ChainCode struct {
}

// // 유저 구조체
// type User struct {
// 	Phone   string `json:"phone"`
// 	Battery Battery `json:"battery"`
// }

// 배터리 구조체
type Battery struct {
	Phone string `json:"phone"`
	BatteryStatusStart string `json:"bss"`
	BatteryStatusEnd string `json:"bse"`
	StationName string `json:"sn"`
	StationGps string `json:"sgps"`
	StartDate string `json:"sdate"`
	EndDate string `json:"edate"`
}

var s_start = ""  // 배터리 충전 시작 상태
var d_start = ""  // 배터리 충전 시작 시간

func (s *ChainCode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	
	return shim.Success(nil)
}



func (s *ChainCode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()

	// if function == "addUser" {
	// 	return s.addUser(APIstub, args)
	// } else 
	  if function == "addBattery" {
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
// func (s *ChainCode) addUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 1 {
// 		return shim.Error("fail!")
// 	}
// 	var user = User{Phone: args[0]}
// 	userAsBytes, _ := json.Marshal(user)
// 	APIstub.PutState(args[0], userAsBytes)

// 	return shim.Success(nil)
// }

// 데이터 입력
func (s *ChainCode) addBattery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}
	// // 유저 정보 가져오기
	// userAsBytes, err := APIstub.GetState(args[0])
	// if err != nil{
	// 	jsonResp := "\"Error\":\"Failed to get state for "+ args[0]+"\"}"
	// 	return shim.Error(jsonResp)
	// } else if userAsBytes == nil{ // no State! error
	// 	jsonResp := "\"Error\":\"User does not exist: "+ args[0]+"\"}"
	// 	return shim.Error(jsonResp)
	// }
	// // 상태 확인
	// user := User{}
	// err = json.Unmarshal(userAsBytes, &user)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// 배터리 구조체 값 업데이트


	if args[4] == "1" {
		s_start = args[1] 
		d_start = time.Now().Format("20060102150405")
		_ = s_start
		_ = d_start
		return shim.Success([]byte("battery is update args true"))

	} else if args[4] == "0" {
		var data = Battery{Phone:args[0],BatteryStatusStart:s_start,BatteryStatusEnd:args[1],StationName:args[2],StationGps:args[3],StartDate:d_start ,EndDate: time.Now().Format("20060102150405")}
		userAsBytes,_:=json.Marshal(data)
		// 월드스테이드 업데이트 
		APIstub.PutState(args[0], userAsBytes)
		return shim.Success([]byte("battery is update args false"))
	}
	return shim.Error("Invalid Smart Contract Battery add.")
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
	var buffer string
	buffer ="["

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer += ","
		}

		buffer += string(response.Value)

		bArrayMemberAlreadyWritten = true
	}
	buffer += "]"
	return shim.Success([]byte(buffer))
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


	var buffer string
	buffer ="["

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer += ","
		}

			buffer += string(response.Value)

		bArrayMemberAlreadyWritten = true
	}
	buffer += "]"

	return shim.Success([]byte(buffer))
}


func main() {
	if err := shim.Start(new(ChainCode)); err != nil {
		fmt.Printf("Error starting ChainCode chaincode: %s", err)
	}
}