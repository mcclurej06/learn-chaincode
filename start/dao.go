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

func getAverageRating(stub *shim.ChaincodeStub, index int) (float32, error) {
	l("getting average rating " + strconv.Itoa(index))
	var err error
	var totalRating int
	var numberOfRatings int

	b, err := stub.GetState(strconv.Itoa(index) + TOTAL_RATING)
	if err != nil {
		l("error getting total rating")
		return -1, err
	}
	totalRating= Float32frombytes(b)
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

func getNumberOfAgents(stub *shim.ChaincodeStub) (int, error) {
	someBytes, err := stub.GetState(NUMBER_OF_AGENTS)
	if err != nil {
		l("error getting number of agents")
		return -1, err
	}
	numberOfAgents, err := strconv.Atoi(string(someBytes))
	if err != nil {
		l("error parsing number of agents " + string(someBytes))
		return -1, err
	}
	return numberOfAgents, nil
}

func getAgents(stub *shim.ChaincodeStub) (string, error) {
	var numberOfAgents int
	var someBytes []byte
	var err error

	someBytes, err = stub.GetState(NUMBER_OF_AGENTS)
	if err != nil {
		l("error getting number of agents")
		return "", err
	}
	numberOfAgents, err = getNumberOfAgents(stub)
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
		averageRating, err := getAverageRating(stub, x)
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

func writeAgent(stub *shim.ChaincodeStub, agent AgentInternal) (error) {
	if agent.Index == -1 {
		numberOfAgents, err := getNumberOfAgents(stub)
		if err != nil {
			l("error getting number of agents")
			return err
		}
		agent.Index = numberOfAgents
	}

	err := stub.PutState(strconv.Itoa(agent.Index) + UUID, []byte(agent.Uuid))
	if err != nil {
		l("error putting uuid")
		return err
	}
	err = stub.PutState(strconv.Itoa(agent.Index) + TOTAL_RATING, Float32bytes(agent.TotalRating))
	if err != nil {
		l("error putting total rating")
		return err
	}
	err = stub.PutState(strconv.Itoa(agent.Index) + NUMBER_OF_RATINGS, []byte(strconv.Itoa(agent.NumberOfRatings)))
	if err != nil {
		l("error putting number of ratings")
		return err
	}
	err = stub.PutState(strconv.Itoa(agent.Index) + NAME, []byte(agent.Name))
	if err != nil {
		l("error putting name")
		return err
	}

	return nil
}