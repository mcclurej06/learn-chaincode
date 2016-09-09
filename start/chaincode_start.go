/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

const NUMBER_OF_AGENTS = "numberOfAgents"
const UUID = "UUID"
const TOTAL_RATING = "TotalRating"
const NUMBER_OF_RATINGS = "NumberOfRatings"

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var value string
	//var numberOfAgents int
	//var someBytes []byte
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	//someBytes, err = stub.GetState(NUMBER_OF_AGENTS)
	if (err != nil) {
		return nil, err
	}
	//numberOfAgents,err= strconv.Atoi(string(someBytes))
	if (err != nil) {
		return nil, err
	}

	//for x := 0; x < numberOfAgents - 1; x++ {
	//	if stub.GetState(x + UUID) {
	//
	//	}
	//}

	agentId := args[0]                            //rename for fun
	value = args[1]
	err = stub.PutState(agentId + TOTAL_RATING, []byte(value))
	err = stub.PutState(agentId + TOTAL_RATING, []byte(value))
	err = stub.PutState(agentId + NUMBER_OF_RATINGS, []byte(value))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	if function == "getAgents" {
		jsonResp, err = t.getAgents(stub)
		if (err != nil) {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return []byte(jsonResp), nil
}

type Agent struct {
	uuid            string        `json:"uuid"`
	averageRating   float32        `json:"averageRating"`
	numberOfRatings int        `json:"numberOfRatings"`
}
type Response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func (t *SimpleChaincode) getAgents(stub *shim.ChaincodeStub) (string, error) {
	var numberOfAgents int
	var someBytes []byte
	var err error

	someBytes, err = stub.GetState(NUMBER_OF_AGENTS)
	if err != nil {
		return "", err
	}
	numberOfAgents, err = strconv.Atoi(string(someBytes))
	if err != nil {
		return "", err
	}

	agents := []Agent{}

	for x := 0; x < numberOfAgents; x++ {

		var agent Agent

		agent.uuid, err = t.getUuid(stub, x)
		if err != nil {
			return "", err
		}
		agent.averageRating, err = t.getAverageRating(stub, x)
		if err != nil {
			return "", err
		}
		agent.numberOfRatings, err = t.getNumberOfRatings(stub, x)
		if err != nil {
			return "", err
		}

		agents = append(agents, agent)
	}

	s, e := json.Marshal(agents)

	return string(s), e
}

func (t *SimpleChaincode) getUuid(stub *shim.ChaincodeStub, index int) (string, error) {
	b, e := stub.GetState(strconv.Itoa(index) + UUID)
	return string(b), e
}

func (t *SimpleChaincode) getAverageRating(stub *shim.ChaincodeStub, index int) (float32, error) {
	var err error
	var totalRating int
	var numberOfRatings int

	b, err := stub.GetState(strconv.Itoa(index) + TOTAL_RATING)
	if err != nil {
		return -1, err
	}
	totalRating, err = strconv.Atoi(string(b))
	if err != nil {
		return -1, err
	}

	b, err =stub.GetState(strconv.Itoa(index) + NUMBER_OF_RATINGS)
	if err != nil {
		return -1, err
	}
	numberOfRatings, err = strconv.Atoi(string(b))
	if err != nil {
		return -1, err
	}

	return float32(totalRating) / float32(numberOfRatings), err
}

func (t *SimpleChaincode) getNumberOfRatings(stub *shim.ChaincodeStub, index int) (int, error) {
	b, err := stub.GetState(strconv.Itoa(index) + NUMBER_OF_RATINGS)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(b))
}