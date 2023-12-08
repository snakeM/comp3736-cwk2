package main

// Result data received from front-end
type ResultRecord struct {
	ResultID string `json:"id"`
	Chart string `json:"chart"`
	TrialID string `json:"trial"`
	TimeTaken string `json:"timeTaken"`
	Answer string `json:"answer"`
}

type LineChart struct {
	Label string `json:"label"`
	Data []int `json:"data"`
}

type LineData struct {
	Labels []string `json:"labels"`
	Datasets []LineChart `json:"datasets"`
}

type RequestData struct {
	Data []ResultRecord `json:"data"`
}

type Charts struct {
	LineCharts []LineData `json:"lineCharts"`
}