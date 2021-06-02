package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"hellofresh/elastic"
	"hellofresh/models"
	"hellofresh/processor"
	"hellofresh/utils"
)

type Store interface {
	InitClient()
	IndexExists() (bool, error)
	CreateIndex() error
	CreateMapping() error
}

func main() {

	var store Store = new(models.PersistanceType)

	var (
		fileName    string
		postCode    string
		startTime   string
		endTime     string
		ingredients string
		sleepTime   int
		process     bool
	)

	flag.StringVar(&fileName, "f", "hf_test_calculation_fixtures.json", "JSON file path")
	flag.StringVar(&postCode, "p", "10120", "Postcode to find deliveries")
	flag.StringVar(&startTime, "s", "10AM", "Start time of delivery to search")
	flag.StringVar(&endTime, "e", "3PM", "End time of delivery to search")
	flag.StringVar(&ingredients, "i", "Potato,Veggie,Chops", "List of ingredients to search for Recipes. e.g.Potato,Veggie,Mushroom")
	flag.BoolVar(&process, "process", false, "Process the input json file and flush to Elasticsearch")
	flag.IntVar(&sleepTime, "t", 2, "Time to Sleep between Bulk Requests")
	flag.Parse()

	//Vaidate Input
	if !utils.IsNumeric(postCode) {
		fmt.Fprintf(os.Stderr, "Error : Invalid Postcode.")
		return
	}

	if !utils.IsTimeValid(startTime, "AM") || !utils.IsTimeValid(endTime, "PM") {
		fmt.Fprintf(os.Stderr, "Error : Start time or End time of delivery is not in expected format.")
		return
	}

	if len(ingredients) == 0 {
		fmt.Fprintf(os.Stderr, "Ingredients is not in the expected format.")
		return
	}

	//Initialize ES client
	store.InitClient()
	indexExists, err := store.IndexExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while checking if the Index exists : %v", err)
	}

	//Process the Json file
	if process {
		if !utils.FileExists(fileName) {
			fmt.Fprintf(os.Stderr, "File doesn't exist : %s", fileName)
			return
		}
		if !indexExists {
			if err := store.CreateIndex(); err != nil {
				fmt.Fprintf(os.Stderr, "CreateIndex failed : %v", err)
			} else {
				if err := store.CreateMapping(); err != nil {
					fmt.Fprintf(os.Stderr, "CreateMapping failed : %v", err)
				}
			}
		}

		//Process input JSON file and flush to Persistence layer
		err = processor.ProcessInputJson(fileName, sleepTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON file Processing failed : %v", err)
		}
		return
	}

	if !indexExists {
		fmt.Fprintf(os.Stderr, "Error : Index not found. Process data and then try again")
		return
	}

	// Generate the Final Response
	var response elastic.FinalResponse

	items, recipeList, err := elastic.GetUniqueRecipes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while getting unique recipes : %v", err)
	}

	response.Count = items
	array := make([]elastic.RecipeCountItem, 0)
	for _, v := range recipeList {
		item := elastic.RecipeCountItem{
			Recipe: v.Key,
			Count:  v.Count,
		}
		array = append(array, item)
	}
	response.CountPerRecipeList = array

	postcodeItem, err := elastic.GetMostDeliverdPostCode()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while getting most delivered postcode : %v", err)
	}

	value := elastic.BusiestPostcode{
		Postcode:      postcodeItem.Key,
		DeliveryCount: postcodeItem.Count,
	}
	response.BusiestPostcode = value

	startingTime, endingTime, _ := utils.GetStartAndEnd("Weekday " + startTime + " - " + endTime)
	var deliveryCount elastic.CountPerPostcodeAndTime
	deliveryCount.Postcode = postCode
	deliveryCount.From = startTime
	deliveryCount.To = endTime
	deliveryCount.DeliveryCount, err = elastic.GetDeliveriesToPostcodeWithinTimerange(postCode, startingTime, endingTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while getting Delivery count within time frame : %v", err)
	}
	response.CountPerPostcodeAndTime = deliveryCount

	response.MatchByName, err = elastic.GetRecipesUsingKeywords(strings.Split(ingredients, ","))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured while getting recipes using keywords : %v", err)
	}

	b, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(string(b), "\n")
}
