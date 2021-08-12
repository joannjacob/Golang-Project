package cronjob

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"main_project/models"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type Property struct {
	csvPos          float64     `json:"csvPos"`
	rdsColumnName   string      `json:"rdsColumnName"`
	datatype        string      `json:"datatype"`
	hasDefaultValue bool        `json:"hasDefaultValue"`
	isDefault       interface{} `json:"default,omitempty"`
}
type Fields struct {
	SKUCode     Property `json:"SKUCode"`
	Name        Property `json:"Name"`
	Description Property `json:"Description"`
	Color       Property `json:"Color"`
	Size        Property `json:"Size"`
}

var wg sync.WaitGroup

func LoadConfiguration() map[string]int {
	var config Fields
	var result map[string]map[string]interface{}

	configFile, err := os.Open("./config/csvConfig.json")
	defer configFile.Close()
	if err != nil {
		log.WithFields(log.Fields{"Cron": "ImportData", "error": err}).Error("Error reading csv configuration file")
	}

	json.NewDecoder(configFile).Decode(&result)
	json.NewDecoder(configFile).Decode(&config)

	configMap := make(map[string]map[string]string)
	for k, v := range result["Fields"] {
		innerMap := make(map[string]string)
		for i, val := range v.(map[string]interface{}) {
			switch val.(type) {
			case int:
				innerMap[i] = strconv.Itoa(val.(int))
			case string:
				innerMap[i] = val.(string)
			case bool:
				innerMap[i] = strconv.FormatBool(val.(bool))
			case float64:
				innerMap[i] = fmt.Sprintf("%f", val.(float64)-1)
			}
		}
		configMap[k] = innerMap
	}
	modelMap := make(map[string]int)
	for _, value := range configMap {
		position, _ := strconv.ParseFloat(value["csvPos"], 64)
		modelMap[value["rdsColumnName"]] = int(position)

	}
	return modelMap
}

func saveProduct(line []string, modelMap map[string]int) {
	product := &models.Product{}
	productMap := make(map[string]interface{})
	for key, value := range modelMap {
		productMap[key] = line[value]

	}
	mapstructure.Decode(productMap, &product)
	product.Create()
	wg.Done()
}

var ImportData = func() {
	startTime := time.Now()
	modelMap := LoadConfiguration()

	csvFile, err := os.Open("./config/productData.csv")
	if err != nil {
		log.WithFields(log.Fields{"Cron": "ImportData", "error": err}).Error("Error opening product data input file")
	}
	defer csvFile.Close()
	row1, err := bufio.NewReader(csvFile).ReadSlice('\n')
	if err != nil {
		log.WithFields(log.Fields{"Cron": "ImportData", "error": err}).Error("Error in reading product data file")
	}
	_, err = csvFile.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		log.WithFields(log.Fields{"Cron": "ImportData", "error": err}).Error("Error in reading data file(seek)")
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			log.WithFields(log.Fields{"Cron": "ImportData"}).Info("Reached EOF")
			break
		} else if error != nil {
			log.WithFields(log.Fields{"Cron": "ImportData", "error": err}).Error("Error reading line in product data file")
		}
		wg.Add(1)
		go saveProduct(line, modelMap)
	}

	fmt.Println("Successfully imported data to database in ", time.Now().Sub(startTime))
	wg.Wait()
}
