package main

// Result data received from front-end
type ResultRecord struct {
	ResultID string `json:"id"`
	Chart string `json:"chart"`
	TrialID string `json:"trial"`
	TimeTaken string `json:"timeTaken"`
	Answer string `json:"answer"`
}

type RequestData struct {
	Data []ResultRecord `json:"data"`
}

type chartDataset struct {
	Label string `json:"label"`
	Data []int `json:"data"`
	Fill bool `json:"fill"`
}

type chart struct {
	Labels []string `json:"labels"`
	Datasets []chartDataset `json:"datasets"`
}

type trial struct {
	Id int `json:"id"`
	Question string `json:"question"`
	Answers []string `json:"answers"`
	Chart chart `json:"chart"`
}

type trials struct {
	Trials []trial `json:"trials"`
}