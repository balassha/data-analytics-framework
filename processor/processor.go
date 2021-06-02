package processor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hellofresh/utils"
	"os"
	"sync"
	"time"

	"hellofresh/elastic"
)

var wg sync.WaitGroup

type order struct {
	Postcode string `json:"postcode"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
}

func (e *order) Unmarshal(b []byte) error {
	return json.Unmarshal(b, e)
}

func ProcessInputJson(fileName string, sleepTime int) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error to read [file=%v]: %v", fileName, err.Error())
	}

	r := bufio.NewReader(f)
	d := json.NewDecoder(r)

	i := 1
	wrapperArr := make([]elastic.IndexWrapper, 0)
	dataArr := make([]elastic.Data, 0)

	d.Token()
	for d.More() {
		elm := &order{}
		d.Decode(elm)
		//fmt.Printf("%v \n", elm)
		start, end, err := utils.GetStartAndEnd(elm.Delivery)
		if err != nil {
			return fmt.Errorf("rrror occured while getting start and end time: %v", err)
		}

		wrapper := elastic.IndexWrapper{
			Index: elastic.IndexName{Index: elastic.Index},
		}
		dat := elastic.Data{
			Postcode: elm.Postcode,
			Recipe:   elm.Recipe,
			Delivery: elm.Delivery,
			Start:    start,
			End:      end,
		}

		wrapperArr = append(wrapperArr, wrapper)
		dataArr = append(dataArr, dat)
		// Bulk request with 500000 entries
		if i%500000 == 0 {
			wg.Add(1)
			// Making bulk query to wait for 2s to avoid hitting Cordination execution from ES
			// This can be solved by scaling up ES
			time.Sleep(time.Duration(sleepTime) * time.Second)
			go ProcessWorker(wrapperArr, dataArr)
			wrapperArr = make([]elastic.IndexWrapper, 0)
			dataArr = make([]elastic.Data, 0)
			//break
		}
		i++
	}
	d.Token()
	wg.Wait()

	return nil
}

func ProcessWorker(indices []elastic.IndexWrapper, data []elastic.Data) {
	defer wg.Done()
	if err := elastic.BulkUpdate(indices, data); err != nil {
		fmt.Fprintf(os.Stderr, "Error while Bulk Update : %v", err)
	}
}
