package main

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
	return AgentResponse{
		Uuid:uuid,
		AverageRating:averageRating,
		NumberOfRatings:numberOfRatings,
		Name:name,
	}
}

func createAgentPost(uuid string, rating float32, name string) (AgentPost) {
	return AgentPost{
		Uuid:uuid,
		Rating:rating,
		Name:name,
	}
}

func createAgentInternal(uuid string, index int, totalRating float32, numberOfRatings int, name string) (AgentInternal) {
	return AgentInternal{
		Uuid:uuid,
		Index:index,
		TotalRating:totalRating,
		NumberOfRatings:numberOfRatings,
		Name:name,
	}
}