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

const INSERT_DATASET string = `
		INSERT INTO
			Datasets
		(
			Data
		)
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

func generateRandData(n int) []int {
	max := 50
	min := 0

	var randData = make([]int, n)
	for i := 0; i < n; i++ {
		randData[i] = generateRandNum(min, max)
	}


	return randData
}

func generateRandLine() []LineChart {
	countries := []string{"USA", "Belgium", "Great Britain", "Spain", "Italy", "France", "Greece", "Japan"}
	datasets := make([]LineChart, len(countries))
	for i, country := range countries {
		datasets[i] = LineChart{Label: country, Data: generateRandData(12)}
	}
	return datasets
}

func generateDataset(c *gin.Context) {
	var allLineData Charts

	dateLabels := []string{"1972", "1976", "1980", "1984", "1988", "1992", "1996", "2000", "2004", "2008", "2012", "2016"}
	allLineData.LineCharts = make([]LineData, 10)
	for i := 0; i < 10; i++ {
		allLineData.LineCharts[i] = LineData{
			Labels: dateLabels,
			Datasets: generateRandLine(),
		}
	}
	c.JSON(http.StatusOK, allLineData)
}

func getChartData(c *gin.Context) {
	lineData, err := ReadJSONFile("cmd/data.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening data file."})
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, lineData)
}

func ReadJSONFile(filePath string) (Charts, error) {
	
	file, err := os.Open(filePath)
	if err != nil {
		return Charts{nil}, err
	}
	defer file.Close()

	var jsonData Charts
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonData)
	if err != nil {
		return Charts{nil}, err
	}

	return jsonData, nil
}