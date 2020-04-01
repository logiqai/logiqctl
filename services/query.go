package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/logiqai/logiq-box/api/v1/query"
	"github.com/logiqai/logiq-box/db"

	"github.com/logiqai/logiq-box/cfg"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

type QueryDebug struct {
	IdMap          map[string]int
	Count          int
	DuplicateCount int
}

func (q *QueryDebug) insert(id string) {
	q.Count = q.Count + 1
	if c, ok := q.IdMap[id]; ok {
		q.IdMap[id] = c + 1
		fmt.Printf("Found duplicate id: %s, %d\n", id, q.IdMap[id])
		q.DuplicateCount = q.DuplicateCount + 1
	} else {
		q.IdMap[id] = 1
	}
}

func (q *QueryDebug) Strings() string {
	return fmt.Sprintf("Count : %d \n Duplicate Count : %d\n", q.Count, q.DuplicateCount)
}

func Query(c *cli.Context, config *cfg.Config, applications []string, keyWords string, qType string) {

	startTime := c.String("st")
	endTime := c.String("et")
	debug := c.String("debug")
	filters := c.String("filter")
	filterValuesMap := getFilters(filters)

	isDebug := false
	if debug == "true" {
		isDebug = true
	}

	ctx := context.Background()
	conn, err := grpc.Dial(config.Cluster, grpc.WithInsecure())
	if err != nil {
		handleError(config, err)
		return
	}
	client := query.NewQueryServiceClient(conn)
	st := time.Hour
	if startTime != "" {
		st, err = parseTime(startTime)
		if err != nil {
			fmt.Printf("Invalid start time %s\n", st)
			return
		}
	} else {
		st = time.Hour
	}
	et := time.Second
	if endTime != "" {
		et, err = parseTime(endTime)
		if err != nil {
			fmt.Printf("Invalid start time %s\n", et)
			return
		}
	}

	in := &query.QueryProperties{
		Filters:   filterValuesMap,
		PageSize:  30,
		StartTime: time.Now().Add(-1 * st).Format(time.RFC3339),
		EndTime:   time.Now().Add(-1 * et).Format(time.RFC3339),
	}

	if qType == "QUERY" {
		in.ApplicationNames = applications
	} else if qType == "SEARCH" {
		in.KeyWord = keyWords
	}

	response, err := client.Query(ctx, in)
	if err != nil {
		handleError(config, err)
	}
	qId := response.GetQueryId()
	db.UpdateLatestQID(qId)
	if isDebug {
		fmt.Printf("Query id is : %s\n", qId)
	}
	nextRequest := &query.GetDataRequest{
		QueryId: qId,
	}

	time.Sleep(5 * time.Second)
	errorCount := 0
	retry := true
	for {
		if !retry {
			break
		}
		dataResponse, err := client.GetData(ctx, nextRequest)
		if err != nil {
			errorCount++
			if errorCount > 5 {
				fmt.Printf("Somthing Wrong with the request. ")
			}
			seconds5 := 0
			for {
				seconds5++
				fmt.Print("...")
				time.Sleep(time.Second)
				if seconds5 == 5 {
					break
				}
			}
			fmt.Print("\n")
			continue
		}
		data := dataResponse.GetData()
		if len(data) == 0 {
			if isDebug {
				fmt.Println("Got No Data ")
			}
			time.Sleep(5 * time.Second)
			continue
		}
		printData(dataResponse)
		retry = false
	}
}

func GetNext(c *cli.Context, config *cfg.Config) {
	qId, err := db.GetLastQuery()
	if err != nil {

	}
	if qId == "" {
		fmt.Println("No data available for the query. Please Query again...")
		return
	}

	ctx := context.Background()
	conn, err := grpc.Dial(config.Cluster, grpc.WithInsecure())
	if err != nil {
		handleError(config, err)
		return
	}

	client := query.NewQueryServiceClient(conn)
	nextRequest := &query.GetDataRequest{
		QueryId: qId,
	}

	dataResponse, err := client.GetData(ctx, nextRequest)
	if err != nil {
		handleError(config, err)
		return
	}
	printData(dataResponse)
}

func printData(dataResponse *query.GetDataResponse) {
	data := dataResponse.GetData()
	for _, d := range data {
		printSyslogMessageForType(d)
		time.Sleep(150 * time.Millisecond)
	}
	if "COMPLETE" != dataResponse.Status && dataResponse.Remaining == 0 {
		fmt.Println("Enter `n` or `next' to continue.")
	} else if "COMPLETE" != dataResponse.Status {
		fmt.Printf("\nAtleast %d records remaining... Enter `n` or `next' to continue.\n", dataResponse.Remaining)
	} else {
		fmt.Println("No more records available for this query.")
	}
	//TODO handle retry status MS6
	db.HandleGetDataStatus(dataResponse.Status)
}

func getFilters(filters string) map[string]*query.FilterValues {
	var filterValuesMap map[string]*query.FilterValues
	if len(filters) > 0 {
		filterValuesMap = make(map[string]*query.FilterValues)
		if strings.Contains(filters, ";") {
			filterArray := strings.Split(filters, ";")
			for _, filterItem := range filterArray {
				key, values := splitIndividualFilter(filterItem)
				filterValuesMap[key] = &query.FilterValues{Values: values}
			}
		} else {
			key, values := splitIndividualFilter(filters)
			filterValuesMap[key] = &query.FilterValues{Values: values}
		}
	}
	return filterValuesMap
}

func splitIndividualFilter(filter string) (string, []string) {
	if strings.Contains(filter, "=") {
		filterItemSplit := strings.Split(filter, "=")
		filterKey := filterItemSplit[0]
		var filterValues []string
		if strings.Contains(filterItemSplit[1], ",") {
			filterValues = strings.Split(filterItemSplit[1], ",")
		} else {
			filterValues = []string{filterItemSplit[1]}
		}
		return filterKey, filterValues
	}
	return "", nil
}

func parseTime(t string) (time.Duration, error) {
	d := strings.ToLower(t[len(t)-1:])
	i := strings.ToLower(t[:len(t)-1])
	iNum, err := strconv.ParseInt(i, 10, 64)
	if err != nil {
		return 0, err
	}
	switch d {
	case "d":
		return time.Hour * 24 * time.Duration(iNum), err
	case "h":
		return time.Hour * time.Duration(iNum), err
	case "m":
		return time.Minute * time.Duration(iNum), err
	default:
		return 0, errors.New("invalid Duration")
	}
}
