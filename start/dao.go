package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"encoding/json"
)

func getUuid(stub *shim.ChaincodeStub, index int) (string, error) {
	uuidKey := strconv.Itoa(index) + UUID
	l("getting uuid from key: " + uuidKey)
	b, e := stub.GetState(uuidKey)
	return string(b), e
}

func getNumberOfRatings(stub *shim.ChaincodeStub, index int) (int, error) {
	b, err := stub.GetState(strconv.Itoa(index) + NUMBER_OF_RATINGS)
	if err != nil {
		l("error getting number of ratings")
		return -1, err
	}
	return strconv.Atoi(string(b))
}

func (t *SimpleChaincode) getAverageRating(stub *shim.ChaincodeStub, index int) (float32, error) {
	l("getting average rating " + strconv.Itoa(index))
	var err error
	var totalRating int
	var numberOfRatings int

	b, err := stub.GetState(strconv.Itoa(index) + TOTAL_RATING)
	if err != nil {
		l("error getting total rating")
		return -1, err
	}
	totalRating, err = strconv.Atoi(string(b))
	if err != nil {
		l("error parsing total rating " + string(b))
		return -1, err
	}

	b, err = stub.GetState(strconv.Itoa(index) + NUMBER_OF_RATINGS)
	if err != nil {
		l("error getting number of ratings")
		return -1, err
	}
	numberOfRatings, err = strconv.Atoi(string(b))
	if err != nil {
		l("error parsing number of ratings" + string(b))
		return -1, err
	}

	return float32(totalRating) / float32(numberOfRatings), err
}

func (t *SimpleChaincode) getAgents(stub *shim.ChaincodeStub) (string, error) {
	var numberOfAgents int
	var someBytes []byte
	var err error

	someBytes, err = stub.GetState(NUMBER_OF_AGENTS)
	if err != nil {
		l("error getting number of agents")
		return "", err
	}
	numberOfAgents, err = strconv.Atoi(string(someBytes))
	if err != nil {
		l("error parsing number of agents " + string(someBytes))
		return "", err
	}

	agents := []AgentResponse{}

	for x := 0; x < numberOfAgents; x++ {
		uuid, err := getUuid(stub, x)
		if err != nil {
			l("error getting agent uuid")
			return "", err
		}
		averageRating, err := t.getAverageRating(stub, x)
		if err != nil {
			l("error getting average rating")
			return "", err
		}
		numberOfRatings, err := getNumberOfRatings(stub, x)
		if err != nil {
			l("error getting number of ratings")
			return "", err
		}

		agents = append(agents, createAgentResponse(uuid, averageRating, numberOfRatings))
	}

	s, e := json.Marshal(agents)

	return string(s), e
}