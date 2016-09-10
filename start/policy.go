package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

type PolicyEventResponse struct {
	AgentUuid        string        `json:"agentUuid"`
	Event            string        `json:"event"`
	PolicyHolderUuid string        `json:"policyHolderUuid"`
}

func createPolicyEventResponse(agentUuid string, event string, policyHolderUuid string) PolicyEventResponse {
	return PolicyEventResponse{AgentUuid: agentUuid, Event: event, PolicyHolderUuid: policyHolderUuid}
}

const NUMBER_OF_POLICY_EVENTS = "numberOfPolicyEvents"
const POLICY_EVENT_AGENT_UUID = "PolicyEventAgentUUID"
const POLICY_EVENT = "PolicyEvent"
const POLICY_EVENT_HOLDER_UUID = "PolicyHolderUuid"

func addPolicyEvent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	index, err := getNumberOfPolicyEvents(stub)
	if err != nil {
		l("error getting number of policy events")
		return nil, err
	}
	incrementNumberOfPolicyEvents(stub)

	policyEvent := createPolicyEventResponse(args[0], args[1], args[2]);

	err = putString(stub, strconv.Itoa(index) + POLICY_EVENT_AGENT_UUID, policyEvent.AgentUuid);
	if err != nil {
		l("error putting policy event agent uuid")
		return nil, err
	}
	err = putString(stub, strconv.Itoa(index) + POLICY_EVENT, policyEvent.Event);
	if err != nil {
		l("error putting policy event")
		return nil, err
	}
	err = putString(stub, strconv.Itoa(index) + POLICY_EVENT_HOLDER_UUID, policyEvent.PolicyHolderUuid);
	if err != nil {
		l("error putting policy event policy holder uuid")
		return nil, err
	}
	return nil, nil
}

func getPolicyEvents(stub *shim.ChaincodeStub) (string, error) {
	var err error

	numberOfPolicyEvents, err := getInt(stub, NUMBER_OF_POLICY_EVENTS)
	if err != nil {
		l("error getting number of policy events")
		return "", err
	}

	policyEvents := []PolicyEventResponse{}

	for x := 0; x < numberOfPolicyEvents; x++ {
		agentUuid, err := getPolicyEventAgentUuid(stub, x)
		if err != nil {
			l("error getting policy event agent uuid")
			return "", err
		}
		event, err := getPolicyEvent(stub, x)
		if err != nil {
			l("error getting policy event")
			return "", err
		}
		policyHolderUuid, err := getPolicyEventHolderUuid(stub, x)
		if err != nil {
			l("error getting policy event holder uuid")
			return "", err
		}
		policyEvents = append(policyEvents, createPolicyEventResponse(agentUuid, event, policyHolderUuid))
	}

	s, e := json.Marshal(policyEvents)

	return string(s), e
}

func getPolicyEventAgentUuid(stub *shim.ChaincodeStub, index int) (string, error) {
	policyEventAgentUuid, err := getString(stub, strconv.Itoa(index) + POLICY_EVENT_AGENT_UUID)
	if err != nil {
		l("error getting policy event agent uuid")
		return "", err
	}
	return policyEventAgentUuid, nil
}

func getPolicyEvent(stub *shim.ChaincodeStub, index int) (string, error) {
	policyEvent, err := getString(stub, strconv.Itoa(index) + POLICY_EVENT)
	if err != nil {
		l("error getting policy event")
		return "", err
	}
	return policyEvent, nil
}

func getPolicyEventHolderUuid(stub *shim.ChaincodeStub, index int) (string, error) {
	policyEventHolderUuid, err := getString(stub, strconv.Itoa(index) + POLICY_EVENT_HOLDER_UUID)
	if err != nil {
		l("error getting policy event holder uuid")
		return "", err
	}
	return policyEventHolderUuid, nil
}

func incrementNumberOfPolicyEvents(stub *shim.ChaincodeStub) (error) {
	l("incrementing number of policy events")
	numberOfPolicyEvents, err := getNumberOfPolicyEvents(stub)
	if err != nil {
		l("error getting number of agents ")
		return err
	}

	return stub.PutState(NUMBER_OF_POLICY_EVENTS, []byte(strconv.Itoa(numberOfPolicyEvents + 1)))
}

func getNumberOfPolicyEvents(stub *shim.ChaincodeStub) (int, error) {
	numberOfPolicyEvents, err := getInt(stub, NUMBER_OF_POLICY_EVENTS)
	if err != nil {
		l("error getting number of policy events")
		return -1, err
	}
	return numberOfPolicyEvents, nil
}