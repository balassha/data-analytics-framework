package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	esURI           = "http://localhost:9200/"
	Index           = "orders1"
	ContentType     = "Content-Type"
	ApplicationJson = "application/json"
	Mapping         = "/_mapping"
	Bulk            = "/_bulk"
	Search          = "/_search?size=10000&pretty=true"
)

var (
	HttpClient *http.Client
)

func InitClient() {
	HttpClient = &http.Client{}
}

func IndexExists() (bool, error) {
	req, err := http.NewRequest(http.MethodHead, esURI+Index, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set(ContentType, ApplicationJson)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	} else {
		return true, nil
	}
}

func CreateIndex() error {

	req, err := http.NewRequest(http.MethodPut, esURI+Index, nil)
	if err != nil {
		return err
	}
	req.Header.Set(ContentType, ApplicationJson)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with Status : %v,%v", resp.StatusCode, string(responseData))
	}
	return nil
}

func DeleteIndex() error {
	postBody, _ := json.Marshal("")
	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodDelete, esURI+Index, responseBody)
	if err != nil {
		return err
	}

	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with Status : %v", resp.StatusCode)
	}
	return nil
}

func CreateMapping() error {
	mapping := &mapping{
		Properties: properties{
			Postcode: keyword{Type: "keyword"},
			Recipe:   keyword{Type: "keyword"},
			Delivery: keyword{Type: "keyword"},
			Start:    keyword{Type: "integer"},
			End:      keyword{Type: "integer"},
		},
	}
	postBody, _ := json.Marshal(mapping)
	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPut, esURI+Index+Mapping, responseBody)
	if err != nil {
		return err
	}
	req.Header.Set(ContentType, ApplicationJson)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with Status : %v,%v", resp.StatusCode, string(responseData))
	}
	return nil
}

func BulkUpdate(indices []IndexWrapper, data []Data) error {
	var reqBody bytes.Buffer
	enc := json.NewEncoder(&reqBody)
	for i, item := range indices {
		enc.Encode(item)
		enc.Encode(data[i])
	}
	req, err := http.NewRequest(http.MethodPost, esURI+Index+Bulk, &reqBody)
	if err != nil {
		return err
	}
	req.Header.Set(ContentType, ApplicationJson)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with Status : %v,%v", resp.StatusCode, string(responseData))
	} else {
		return nil
	}
}

func GetUniqueRecipes() (int, []RecipeCount, error) {
	reqBody := &uniqueRecipes{
		Aggregator: aggregator{
			Recipes: terms{
				Terms: field{
					Field: "recipe",
				},
			},
		},
	}

	postBody, _ := json.Marshal(reqBody)
	body := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodGet, esURI+Index+Search, body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set(ContentType, ApplicationJson)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	var data aggregatorResponse
	if err = json.Unmarshal(responseData, &data); err != nil {
		return 0, nil, err
	}

	if data.Aggregations.Recipe.Buckets != nil {
		return len(data.Aggregations.Recipe.Buckets), data.Aggregations.Recipe.Buckets, nil
	} else {
		return 0, nil, err
	}
}

func GetMostDeliverdPostCode() (RecipeCount, error) {
	size := 1
	reqBody := &uniqueRecipes{
		Aggregator: aggregator{
			Recipes: terms{
				Terms: field{
					Field: "postcode",
					Size:  &size,
				},
			},
		},
	}

	postBody, _ := json.Marshal(reqBody)
	body := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodGet, esURI+Index+Search, body)
	if err != nil {
		return RecipeCount{}, err
	}
	req.Header.Set(ContentType, ApplicationJson)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return RecipeCount{}, err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RecipeCount{}, err
	}

	var data aggregatorResponse
	if err = json.Unmarshal(responseData, &data); err != nil {
		return RecipeCount{}, err
	}
	if data.Aggregations.Recipe.Buckets != nil {
		return data.Aggregations.Recipe.Buckets[0], nil
	} else {
		return RecipeCount{}, err
	}
}

func GetDeliveriesToPostcodeWithinTimerange(postCode string, gte int, lte int) (int, error) {
	request := &query{
		Query: boolType{
			Bool: must{
				Must: []matchRange{
					{
						Match: &match{
							Postcode: postCode,
						},
					},
					{
						Range: &rangeType{
							Start: &start{
								Gte: &gte,
							},
						},
					},
					{
						Range: &rangeType{
							End: &end{
								Lte: &lte,
							},
						},
					},
				},
			},
		},
	}

	postBody, _ := json.Marshal(request)
	body := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodGet, esURI+Index+Search, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set(ContentType, ApplicationJson)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data aggregatorResponse
	if err = json.Unmarshal(responseData, &data); err != nil {
		return 0, err
	}
	response := 0
	if data.Hits != nil && data.Hits.Hits != nil {
		response = len(data.Hits.Hits)
	}

	return response, nil
}

func GetRecipesUsingKeywords(keywords []string) ([]string, error) {
	response := make([]string, 0)
	query := ""
	for i, v := range keywords {
		query = query + "(*" + v + "*)"
		if i != len(keywords)-1 {
			query = query + " OR "
		}
	}
	request := &regexQuery{
		Query: queryString{
			QueryString: queryDefaultField{
				Query:        query,
				DefaultField: "recipe",
			},
		},
		Collapse: collapse{
			Field: "recipe",
		},
	}
	postBody, _ := json.Marshal(request)
	body := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodGet, esURI+Index+Search, body)
	if err != nil {
		return response, err
	}
	req.Header.Set(ContentType, ApplicationJson)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	var data aggregatorResponse
	if err = json.Unmarshal(responseData, &data); err != nil {
		return response, err
	}

	if data.Hits != nil && data.Hits.Hits != nil {
		for _, v := range data.Hits.Hits {
			response = append(response, v.Array.Recipe...)
		}
	}

	return response, nil
}
