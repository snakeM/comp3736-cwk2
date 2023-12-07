package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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

const INSERT_NEW_RESULT string = `
	INSERT INTO
		Results
	(
		ExperimentID,
		ResultID,
		TrialID,
		Chart,
		TimeTaken,
		Answer
	)
	VALUES (
		@p1,
		@p2,
		@p3,
		@p4,
		@p5,
		@p6
	);
`

// ExperimentID is just the current timestamp
func generateExperimentID() string {
		// Get the current date and time
		currentTime := time.Now()
		// Format the date and time as a string
		return currentTime.Format("02012006-150405")
}

func handleResultData(c *gin.Context) {

	// Deserialize request body 
	var requestData RequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get database handle
	dbHandle, err := initDatabaseConnection()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	log.Println("Succesfully created connection to database")
	
	// Generate a unique ID for the experiment
	experimentID := generateExperimentID()

	// Perform insert query
	for _, item := range requestData.Data {
		_, err := dbHandle.Exec(
			INSERT_NEW_RESULT,
			experimentID,
			item.ResultID,
			item.TrialID,
			item.Chart,
			item.TimeTaken,
			item.Answer,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully added results to database"})
	
}

func getChartData(c *gin.Context) {
	d := []int{10, 20, 11, 25, 26, 30, 9}
	d2 := []int{20, 30, 11, 50, 32, 4, 23}
	dateLabels := []string{"1996", "2000", "2004", "2008", "2012", "2016", "2020", "2024"}
	lineData := LineData{
		Labels: dateLabels,
		Datasets: []LineChart{
			{
				Label: "USA",
				Data: d,
			},
			{
				Label: "UK",
				Data: d2,
			},
		},
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, lineData)
}


func main() {
	r := gin.Default()
	r.POST("/result/new", handleResultData)
	r.GET("/charts", getChartData)
	r.Run()
}
