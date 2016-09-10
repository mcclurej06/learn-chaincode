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