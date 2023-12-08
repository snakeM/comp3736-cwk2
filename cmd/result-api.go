package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

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
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"message": "Successfully added results to database"})
	
}

func generateRandNum(min int, max int) int {
	return rand.Intn(max-min) + min
}

func generateRandData(n int, seed int) []int {
	var min int = 10 * (2*seed + 1)
	var max int = 20 * (2*seed + 1)

	var randData = make([]int, n)
	for i := 0; i < n; i++ {
		randData[i] = generateRandNum(min, max)
	}

	return randData
}

func getCountrySet() []string {
	var countries []string = []string{"USA", "Belgium", "Great Britain", "Spain", "Italy", "France", "Greece", "Japan"}

	rand.Shuffle(len(countries), func(i, j int) { countries[i], countries[j] = countries[j], countries[i] })
	return countries[0:4]
}

func generateDatasets() ([]chartDataset, []chartDataset) {
	countries := getCountrySet()
	lineDatasets := make([]chartDataset, len(countries))
	areaDatasets := make([]chartDataset, len(countries))
	for i, country := range countries {
		data := generateRandData(12, i)
		lineDatasets[i] = chartDataset{Label: country, Data: data, Fill: false}
		areaDatasets[i] = chartDataset{Label: country, Data: data, Fill: true}
	}
	return lineDatasets, areaDatasets
}

func generateDataset(c *gin.Context) {
	
	var allTrials trials
	allTrials.Trials = make([]trial, 20)
	for i := 0; i < 10; i++ {
		lineDataset, areaDataset := generateDatasets()

		lineTrial := trial{
			Id: i,
			Question: QUESTIONS[0],
			Answers: []string{"USA", "Great Britain", "Spain", "Greece"},
			Chart: chart{
				Labels: OLYMPIC_YEARS,
				Datasets: lineDataset,
			},
		}
		areaTrial := trial{
			Id: i+10,
			Question: QUESTIONS[0],
			Answers: []string{"Spain", "Italy", "Greece", "Belgium"},
			Chart: chart{
				Labels: OLYMPIC_YEARS,
				Datasets: areaDataset,
			},
		}
		allTrials.Trials[i] = lineTrial
		allTrials.Trials[10+i] = areaTrial

	}
	c.JSON(http.StatusOK, allTrials)

}

func getChartData(c *gin.Context) {
	lineData, err := ReadJSONFile("cmd/data.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening data file."})
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, lineData)
}

func ReadJSONFile(filePath string) (trials, error) {
	
	file, err := os.Open(filePath)
	if err != nil {
		return trials{nil}, err
	}
	defer file.Close()

	var jsonData trials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonData)
	if err != nil {
		return trials{nil}, err
	}

	return jsonData, nil
}

