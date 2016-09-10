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

func getName(stub *shim.ChaincodeStub, index int) (string, error) {
	uuidKey := strconv.Itoa(index) + NAME
	l("getting name from key: " + uuidKey)
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

func getTotalRating(stub *shim.ChaincodeStub, index int) (float32, error) {
	b, err := stub.GetState(strconv.Itoa(index) + TOTAL_RATING)
	if err != nil {
		l("error getting total rating")
		return -1, err
	}
	return Float32frombytes(b), nil
}

func getAverageRating(stub *shim.ChaincodeStub, index int) (float32, error) {
	l("getting average rating " + strconv.Itoa(index))
	var err error
	var totalRating float32
	var numberOfRatings int

	totalRating, err = getTotalRating(stub, index)
	if err != nil {
		l("error parsing total rating ")
		return -1, err
	}

	numberOfRatings, err = getNumberOfRatings(stub, index)
	if err != nil {
		l("error parsing number of ratings")
		return -1, err
	}

	return totalRating / float32(numberOfRatings), err
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

func incrementNumberOfAgents(stub *shim.ChaincodeStub) (error) {
	l("incrementing number of agents")
	index, err := getNumberOfAgents(stub)
	if err != nil {
		l("error getting number of agents ")
		return err
	}

	return stub.PutState(NUMBER_OF_AGENTS, []byte(strconv.Itoa(index + 1)))
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
		name, err := getName(stub, x)
		if err != nil {
			l("error getting agent uuid")
			return "", err
		}
		agents = append(agents, createAgentResponse(uuid, averageRating, numberOfRatings, name))
	}

	s, e := json.Marshal(agents)

	return string(s), e
}

func writeAgent(stub *shim.ChaincodeStub, agent AgentInternal) (error) {
	l("writing agent")
	if agent.Index == -1 {
		numberOfAgents, err := getNumberOfAgents(stub)
		if err != nil {
			l("error getting number of agents")
			return err
		}
		agent.Index = numberOfAgents
	}

	err := stub.PutState(strconv.Itoa(agent.Index) + UUID, []byte(agent.Uuid))
	l("writing uuid " + agent.Uuid + " to index " + strconv.Itoa(agent.Index))
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

func getAgentInternal(stub *shim.ChaincodeStub, uuid string) (AgentInternal, error) {
	l("getting agent internal")
	index, err := getAgentIndex(stub, uuid)
	if err != nil {
		l("error getting agent index")
		return AgentInternal{}, err
	}

	if index != -1 {
		l("found existing agent")
		rating, err := getTotalRating(stub, index)
		if err != nil {
			l("error getting total rating")
			return AgentInternal{}, err
		}
		numberOfRatings, err := getNumberOfRatings(stub, index)
		if err != nil {
			l("error getting number of ratings")
			return AgentInternal{}, err
		}
		name, err := getName(stub, index)
		if err != nil {
			l("error getting name")
			return AgentInternal{}, err
		}
		return createAgentInternal(uuid, index, rating, numberOfRatings, name), nil
	}

	l("no existing agent, creating new one")

	index, err = getNumberOfAgents(stub)
	if err != nil {
		l("error getting number of agents")
		return AgentInternal{}, err
	}
	incrementNumberOfAgents(stub)
	if err != nil {
		l("error incrementing number of agents")
		return AgentInternal{}, err
	}

	return createAgentInternal(uuid, index, 0, 0, ""), nil

}

func getAgentIndex(stub *shim.ChaincodeStub, uuid string) (int, error) {
	numberOfAgents, err := getNumberOfAgents(stub)
	if err != nil {
		l("error getting number of agents")
		return -1, err
	}
	for x := 0; x < numberOfAgents; x++ {
		myUuid, err := getUuid(stub, x)
		if err != nil {
			l("error getting agent uuid")
			return -1, err
		}
		if myUuid == uuid {
			return x, nil
		}
	}
	return -1, nil
}