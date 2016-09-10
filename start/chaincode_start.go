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
const NAME = "Name"

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
	l("initing")

	var err error

	err = stub.PutState(NUMBER_OF_AGENTS, []byte("2"))
	if err != nil {
		l("error setting number of agents")
		return nil, err
	}
	err = writeAgent(stub, createAgentInternal("foo", 0, 50, 100, "bob"))
	if err != nil {
		l("error writing agent")
		return nil, err
	}
	err = writeAgent(stub, createAgentInternal("bar", 1, 98, 112, "jeff"))
	if err != nil {
		l("error writing agent")
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "updateAgent" {
		return updateAgent(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	if function == "getAgents" {
		jsonResp, err = getAgents(stub)
		if (err != nil) {
			return nil, err
		}
	}

	return []byte(jsonResp), nil
}

func updateAgent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	l("updating agent")
	rating, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		l("error parsing float")
		return nil, err
	}
	agentPost := createAgentPost(args[0], float32(rating), args[2])

	agentInternal, err := getAgentInternal(stub, args[0])
	if err != nil {
		l("error getting agent internal")
		return nil, err
	}

	agentInternal.Name = agentPost.Name
	agentInternal.NumberOfRatings = agentInternal.NumberOfRatings + 1
	agentInternal.TotalRating = agentInternal.TotalRating + agentPost.Rating

	writeAgent(stub, agentInternal)

	return nil, nil
}