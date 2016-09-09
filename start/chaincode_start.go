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
	l("initing")

	var err error

	err = stub.PutState(NUMBER_OF_AGENTS, []byte("2"))
	if err != nil {
		return nil, err
	}
	err = stub.PutState("0" + UUID, []byte("foo"))
	if err != nil {
		return nil, err
	}
	err = stub.PutState("0" + TOTAL_RATING, []byte("50"))
	if err != nil {
		return nil, err
	}
	err = stub.PutState("0" + NUMBER_OF_RATINGS, []byte("100"))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("1" + UUID, []byte("bar"))
	if err != nil {
		return nil, err
	}
	err = stub.PutState("1" + TOTAL_RATING, []byte("98"))
	if err != nil {
		return nil, err
	}
	err = stub.PutState("1" + NUMBER_OF_RATINGS, []byte("100"))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	if function == "init" {
		return t.Init(stub, "init", args)
	}

	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}


func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var jsonResp string
	var err error

	if function == "getAgents" {
		jsonResp, err = t.getAgents(stub)
		if (err != nil) {
			return nil, err
		}
	}

	return []byte(jsonResp), nil
}


