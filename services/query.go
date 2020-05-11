package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/logiqai/logiqctl/api/v1/applications"

	"github.com/logiqai/logiqctl/api/v1/namespace"

	"github.com/dustin/go-humanize"
	"github.com/manifoldco/promptui"

	"github.com/logiqai/logiqctl/api/v1/query"
	"github.com/logiqai/logiqctl/db"

	"github.com/logiqai/logiqctl/cfg"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

const (
	OUTPUT_COLUMNS = "column"
	OUTPUT_RAW     = "raw"
	OUTPUT_JSON    = "json"
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

func Query(c *cli.Context, config *cfg.Config, qType string) {
	startTime := c.String("st")
	endTime := c.String("et")
	debug := c.Bool("debug")
	tail := c.Bool("tail")
	filters := c.String("filter")
	output := c.String("output")

	filterValuesMap := getFilters(filters)

	isDebug := debug

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

	ns, app, lastSeen, _ := RunSelectApplicationForNamespacePrompt(config)

	in := &query.QueryProperties{
		Filters:   filterValuesMap,
		Namespace: ns,
		PageSize:  30,
		StartTime: time.Unix(lastSeen, 0).Add(-1 * st).Format(time.RFC3339),
	}

	if qType == "QUERY" {
		in.ApplicationNames = []string{app}
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

	errorCount := 0
	retry := true
	for {
		if !retry {
			break
		}
		dataResponse, err := client.GetDataNext(ctx, nextRequest)
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
		if dataResponse.Status == "NO_DATA" {
			fmt.Println("No data available for this query, \n\tMake sure that words are spelled correctly.\n\tTry changing the time range")
			return
		}
		data := dataResponse.GetData()
		if dataResponse.Status == "COMPLETE" && len(data) == 0 {
			fmt.Println("No more records available for this query.")
			return
		}

		if len(data) == 0 {
			if isDebug {
				fmt.Printf("Got No Data , Remaining : %d Status : %s\n", dataResponse.Remaining, dataResponse.Status)
			}
			time.Sleep(time.Second)
			continue
		}
		printData(dataResponse, tail, output)
		if !tail || (tail && dataResponse.Remaining <= 0) {
			retry = false
		} else {
			retry = true
		}
	}
}

func GetNext(c *cli.Context, config *cfg.Config) {
	output := c.String("output")

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

	dataResponse, err := client.GetDataNext(ctx, nextRequest)
	if err != nil {
		handleError(config, err)
		return
	}
	printData(dataResponse, false, output)
}

func printData(dataResponse *query.GetDataResponse, modeTail bool, output string) {
	data := dataResponse.GetData()
	for _, d := range data {
		printSyslogMessageForType(d, output)
	}
	if !modeTail {
		if dataResponse.Status != "COMPLETE" || dataResponse.Remaining != 0 {
			fmt.Printf("\nAtleast %d records remaining... Enter `n` or `next' to continue.\n", dataResponse.Remaining)
		} else {
			fmt.Println("No more records available for this query.")
			db.HandleGetDataStatus(dataResponse.Status)
		}
	}
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

func RunSelectApplicationForNamespacePrompt(config *cfg.Config) (string, string, int64, int64) {

	conn, err := grpc.Dial(config.Cluster, grpc.WithInsecure())
	if err != nil {
		handleError(config, err)
	}

	nsClient := namespace.NewNamespaceServiceClient(conn)

	nsResponse, err := nsClient.GetNamespaces(context.Background(), &namespace.NamespaceRequest{})

	if err != nil {
		handleError(config, err)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U000000BB {{ .Name | green }} ({{ .Details | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Details | red }})",
		Selected: "Selected namespace: {{ .Name | red }}",
	}
	var nameSpace string
	if len(nsResponse.Namespaces) > 0 {
		var nss []struct {
			Name    string
			Details string
		}
		for _, ns := range nsResponse.Namespaces {
			ls := time.Unix(ns.LastSeen, 0)
			nss = append(nss, struct {
				Name    string
				Details string
			}{
				Name:    ns.Namespace,
				Details: fmt.Sprintf("Last Seen %s", humanize.Time(ls)),
			})
		}

		nsPrompt := promptui.Select{
			Label:     fmt.Sprintf("Select a namespace"),
			Items:     nss,
			Templates: templates,
			Size:      6,
		}
		what, _, err := nsPrompt.Run()
		if err != nil {
			handleError(config, err)
		}
		nameSpace = nss[what].Name
	}
	appsClient := applications.NewApplicationsServiceClient(conn)
	appsResponse, err := appsClient.GetApplicationsV2(context.Background(), &applications.GetApplicationsRequest{
		Namespace: nameSpace,
	})
	if err != nil {
		handleError(config, err)
	}

	if len(appsResponse.Applications) > 0 {
		var apps []struct {
			Name      string
			Details   string
			LastSeen  int64
			FirstSeen int64
		}
		for _, app := range appsResponse.Applications {
			ls := time.Unix(app.LastSeen, 0)
			apps = append(apps, struct {
				Name      string
				Details   string
				LastSeen  int64
				FirstSeen int64
			}{
				Name:      app.Name,
				Details:   fmt.Sprintf("Last Seen %s", humanize.Time(ls)),
				LastSeen:  app.LastSeen,
				FirstSeen: app.LastSeen,
			})
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U000000BB {{ .Name | green }} ({{ .Details | red }})",
			Inactive: "  {{ .Name | cyan }} ({{ .Details | red }})",
			Selected: "Select application: {{ .Name | red }}",
		}

		whatPrompt := promptui.Select{
			Label:     fmt.Sprintf("Select an application to run query"),
			Items:     apps,
			Templates: templates,
			Size:      6,
		}

		what, _, err := whatPrompt.Run()
		if err != nil {
			handleError(config, err)
		}

		return nameSpace, apps[what].Name, apps[what].LastSeen, apps[what].FirstSeen
	}
	handleError(config, errors.New("Cannot find applications"))
	return "", "", 0, 0
}
