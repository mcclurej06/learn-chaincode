package main

type AgentResponse struct {
	Uuid            string        `json:"uuid"`
	AverageRating   float32        `json:"averageRating"`
	NumberOfRatings int        `json:"numberOfRatings"`
}

func createAgentResponse(uuid string, averageRating float32, numberOfRatings int) (AgentResponse) {
	return AgentResponse{Uuid:uuid, AverageRating:averageRating, NumberOfRatings:numberOfRatings}
}