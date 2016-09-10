package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"encoding/json"
)

type AgentResponse struct {
	Uuid            string        `json:"uuid"`
	AverageRating   float32        `json:"averageRating"`
	NumberOfRatings int        `json:"numberOfRatings"`
	Name            string `json:"name"`
}

type AgentPost struct {
	Uuid   string
	Rating float32
	Name   string
}

type AgentInternal struct {
	Uuid            string
	Index           int
	TotalRating     float32
	NumberOfRatings int
	Name            string
}

func createAgentResponse(uuid string, averageRating float32, numberOfRatings int, name string) (AgentResponse) {
	l("creating agent response")
	return AgentResponse{
		Uuid:uuid,
		AverageRating:averageRating,
		NumberOfRatings:numberOfRatings,
		Name:name,
	}
}

func createAgentPost(uuid string, rating float32, name string) (AgentPost) {
	l("creating agent post")
	return AgentPost{
		Uuid:uuid,
		Rating:rating,
		Name:name,
	}
}

func createAgentInternal(uuid string, index int, totalRating float32, numberOfRatings int, name string) (AgentInternal) {
	l("creating agent internal")
	return AgentInternal{
		Uuid:uuid,
		Index:index,
		TotalRating:totalRating,
		NumberOfRatings:numberOfRatings,
		Name:name,
	}
}
func getAgentUuid(stub *shim.ChaincodeStub, index int) (string, error) {
	uuid, err :=getString(stub, strconv.Itoa(index) + AGENT_UUID)
	if err != nil {
		l("error getting agent uuid")
		return "", err
	}
	return uuid, nil
}

func getAgentName(stub *shim.ChaincodeStub, index int) (string, error) {
	name, err :=getString(stub, strconv.Itoa(index) + AGENT_NAME)
	if err != nil {
		l("error getting agent name")
		return "", err
	}
	return name, nil
}

func getAgentNumberOfRatings(stub *shim.ChaincodeStub, index int) (int, error) {
	numberOfAgents, err := getInt(stub, strconv.Itoa(index) + AGENT_NUMBER_OF_RATINGS)
	if err != nil {
		l("error getting number of ratings")
		return -1, err
	}
	return numberOfAgents, nil
}

func getAgentTotalRating(stub *shim.ChaincodeStub, index int) (float32, error) {
	totalRating, err := getFloat(stub, strconv.Itoa(index) + AGENT_TOTAL_RATING)
	if err != nil {
		l("error getting total rating")
		return -1, err
	}
	return totalRating, nil
}

func getAgentAverageRating(stub *shim.ChaincodeStub, index int) (float32, error) {
	l("getting average rating " + strconv.Itoa(index))
	var err error
	var totalRating float32
	var numberOfRatings int

	totalRating, err = getAgentTotalRating(stub, index)
	if err != nil {
		l("error parsing total rating ")
		return -1, err
	}

	numberOfRatings, err = getAgentNumberOfRatings(stub, index)
	if err != nil {
		l("error parsing number of ratings")
		return -1, err
	}

	return totalRating / float32(numberOfRatings), err
}

func getNumberOfAgents(stub *shim.ChaincodeStub) (int, error) {
	numberOfAgents, err := getInt(stub, NUMBER_OF_AGENTS)
	if err != nil {
		l("error getting number of agents")
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
		uuid, err := getAgentUuid(stub, x)
		if err != nil {
			l("error getting agent uuid")
			return "", err
		}
		averageRating, err := getAgentAverageRating(stub, x)
		if err != nil {
			l("error getting average rating")
			return "", err
		}
		numberOfRatings, err := getAgentNumberOfRatings(stub, x)
		if err != nil {
			l("error getting number of ratings")
			return "", err
		}
		name, err := getAgentName(stub, x)
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

	err := stub.PutState(strconv.Itoa(agent.Index) + AGENT_UUID, []byte(agent.Uuid))
	l("writing uuid " + agent.Uuid + " to index " + strconv.Itoa(agent.Index))
	if err != nil {
		l("error putting uuid")
		return err
	}
	err = stub.PutState(strconv.Itoa(agent.Index) + AGENT_TOTAL_RATING, Float32bytes(agent.TotalRating))
	if err != nil {
		l("error putting total rating")
		return err
	}
	err = stub.PutState(strconv.Itoa(agent.Index) + AGENT_NUMBER_OF_RATINGS, []byte(strconv.Itoa(agent.NumberOfRatings)))
	if err != nil {
		l("error putting number of ratings")
		return err
	}
	err = stub.PutState(strconv.Itoa(agent.Index) + AGENT_NAME, []byte(agent.Name))
	if err != nil {
		l("error putting name")
		return err
	}

	l("wrote agent, returning nil")
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
		rating, err := getAgentTotalRating(stub, index)
		if err != nil {
			l("error getting total rating")
			return AgentInternal{}, err
		}
		numberOfRatings, err := getAgentNumberOfRatings(stub, index)
		if err != nil {
			l("error getting number of ratings")
			return AgentInternal{}, err
		}
		name, err := getAgentName(stub, index)
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
		myUuid, err := getAgentUuid(stub, x)
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