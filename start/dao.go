package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

func putString(stub *shim.ChaincodeStub, key string, value string) error {
	return stub.PutState(key, []byte(value))
}
func putFloat(stub *shim.ChaincodeStub, key string, value float32) error {
	return stub.PutState(key, Float32bytes(value))
}
func putInt(stub *shim.ChaincodeStub, key string, value int) error {
	return stub.PutState(key, []byte(strconv.Itoa(value)))
}

func getInt(stub *shim.ChaincodeStub, key string) (int, error) {
	someBytes, err := stub.GetState(key)
	if err != nil {
		l("error getting int " + key)
		return -1, err
	}
	number, err := strconv.Atoi(string(someBytes))
	if err != nil {
		l("error parsing int for " + key + ": " + string(someBytes))
		return -1, err
	}
	return number, nil
}

func getFloat(stub *shim.ChaincodeStub, key string) (float32, error) {
	b, err := stub.GetState(key)
	if err != nil {
		l("error getting float")
		return -1, err
	}
	return Float32frombytes(b), nil
}

func getString(stub *shim.ChaincodeStub, key string) (string, error) {
	b, err := stub.GetState(key)
	if err != nil {
		l("error getting string")
		return "", err
	}
	return string(b), nil
}