/*
Copyright © 2020 Logiq.ai <cli@logiq.ai>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	// "bitbucket/pkg/mod/github.com/go-delve/delve@v1.4.1/service/api"
	"errors"
	"fmt"
	"github.com/araddon/dateparse"
    "github.com/zenthangplus/goccm"
	"math/rand"
	"os"
	"regexp"

	// "runtime"
	"strings"
	"sync"
	"time"

	"github.com/logiqai/logiqctl/grpc_utils"
	"github.com/logiqai/logiqctl/api/v1/applications"

	"github.com/logiqai/logiqctl/api/v1/query"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
	"github.com/Knetic/govaluate"
	"github.com/logiqai/logiqctl/loglerpart"
)

const (
	OUTPUT_COLUMNS = "column"
	OUTPUT_RAW     = "raw"
	OUTPUT_JSON    = "json"
)

var st time.Time
var et time.Time
var pOnce sync.Once

var tq []time.Time // parallel time queue

var outch chan string =make(chan string, 100)

var outputMutex sync.Mutex

var ParMaxCopies int

// dateparse is a robust date/time parser
// here's is what it supports
// However, only yyyy-mm-dd hh:MM:ss.sss is used
//
/*
   "May 8, 2009 5:57:51 PM",
   "Mon Jan  2 15:04:05 2006",
   "Mon Jan  2 15:04:05 MST 2006",
   "Mon Jan 02 15:04:05 -0700 2006",
   "Monday, 02-Jan-06 15:04:05 MST",
   "Mon, 02 Jan 2006 15:04:05 MST",
   "Tue, 11 Jul 2017 16:28:13 +0200 (CEST)",
   "Mon, 02 Jan 2006 15:04:05 -0700",
   "Thu, 4 Jan 2018 17:53:36 +0000",
   "Mon Aug 10 15:44:11 UTC+0100 2015",
   "Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)",
   "12 Feb 2006, 19:17",
   "12 Feb 2006 19:17",
   "03 February 2013",
   "2013-Feb-03",
   //   mm/dd/yy
   "3/31/2014",
   "03/31/2014",
   "08/21/71",
   "8/1/71",
   "4/8/2014 22:05",
   "04/08/2014 22:05",
   "4/8/14 22:05",
   "04/2/2014 03:00:51",
   "8/8/1965 12:00:00 AM",
   "8/8/1965 01:00:01 PM",
   "8/8/1965 01:00 PM",
   "8/8/1965 1:00 PM",
   "8/8/1965 12:00 AM",
   "4/02/2014 03:00:51",
   "03/19/2012 10:11:59",
   "03/19/2012 10:11:59.3186369",
   // yyyy/mm/dd
   "2014/3/31",
   "2014/03/31",
   "2014/4/8 22:05",
   "2014/04/08 22:05",
   "2014/04/2 03:00:51",
   "2014/4/02 03:00:51",
   "2012/03/19 10:11:59",
   "2012/03/19 10:11:59.3186369",
   // Chinese
   "2014年04月08日",
   //   yyyy-mm-ddThh
   "2006-01-02T15:04:05+0000",
   "2009-08-12T22:15:09-07:00",
   "2009-08-12T22:15:09",
   "2009-08-12T22:15:09Z",
   //   yyyy-mm-dd hh:mm:ss
   "2014-04-26 17:24:37.3186369",
   "2012-08-03 18:31:59.257000000",
   "2014-04-26 17:24:37.123",
   "2013-04-01 22:43",
   "2013-04-01 22:43:22",
   "2014-12-16 06:20:00 UTC",
   "2014-12-16 06:20:00 GMT",
   "2014-04-26 05:24:37 PM",
   "2014-04-26 13:13:43 +0800",
   "2014-04-26 13:13:44 +09:00",
   "2012-08-03 18:31:59.257000000 +0000 UTC",
   "2015-09-30 18:48:56.35272715 +0000 UTC",
   "2015-02-18 00:12:00 +0000 GMT",
   "2015-02-18 00:12:00 +0000 UTC",
   "2017-07-19 03:21:51+00:00",
   "2014-04-26",
   "2014-04",
   "2014",
   "2014-05-11 08:20:13,787",
   // mm.dd.yy
   "3.31.2014",
   "03.31.2014",
   "08.21.71",
   //  yyyymmdd and similar
   "20140601",
   // unix seconds, ms
   "1332151919",
   "1384216367189",
*/

// parallel search related
func SearchWorker(i int, appname string) {
	var sst string
	if i==(len(tq)-2) {
		sst = timeFormat(tq[i+1])
	} else {
		ttmp:=tq[i+1]
		sst = timeFormat( ttmp.Add(time.Duration(-1*1000000000)))
	}
	et := timeFormat(tq[i])
	fmt.Println("SearchWorker ", i, "   app:", appname, "  st:", sst, "   et:", et)
    tt := rand.Intn(10)
    fmt.Printf("Job %d is running for %d second\n", i, tt)
    time.Sleep(time.Duration(int64(tt) * 1000000000))
	fmt.Printf("Job %d is done\n", i)

}


func ParallelExec(applicationV2s *[]*applications.ApplicationV2,
	appName string, args0 string, procId string, lastseen int64) {

	c := goccm.New(utils.FlagParCopies)


	if (applicationV2s!=nil) {
		ParMaxCopies = len(*applicationV2s)*(len(tq)-1)
		for _, app := range *applicationV2s {
			for i := 0; i < len(tq)-1; i++ {
				// fmt.Println("Name: ", app.Name, "  i=", i)

				c.Wait()
				go func(i int, appname string, procId string, lastseen int64) {
					DoQuery(i, appname, args0, procId, lastseen)
					// Sanity test
					// SearchWorker(i, appname)
					// This function have to when a goroutine has finished
					// Or you can use `defer c.Done()` at the top of goroutine.
					c.Done()
				}(i, app.Name, procId, lastseen)
			}
		}
		c.WaitAllDone()
	} else {
		ParMaxCopies = (len(tq)-1)
		for i := 0; i < len(tq)-1; i++ {
			c.Wait()
			go func(i int, appname string, procId string, lastseen int64) {
				DoQuery(i, appname, args0, procId, lastseen)
				// Sanity test
				// SearchWorker(i)
				// This function have to when a goroutine has finished
				// Or you can use `defer c.Done()` at the top of goroutine.
				c.Done()
			}(i, appName, procId, lastseen)
		}
		c.WaitAllDone()
	}
}


func ParallelSearch(cmdusage string, args0 string,  hasApp *bool, hasMultipleApps *bool, hasProc *bool,
	applicationV2s *[]*applications.ApplicationV2) {

	// var in *query.QueryProperties

	create_tq (1)
	*hasApp = true

	if *hasMultipleApps {

		ParallelExec(applicationV2s, "", args0, "", 0)
		/*for _, app := range *applicationV2s {
			ParallelExec(applicationV2s, app.Name, args0, "", 0)
		}*/

	} else {

		if len(*applicationV2s) > 0 {
			if utils.FlagProcId != "" {
				*hasProc = true
			}
			if *hasApp && *hasProc {
				proc, err := GetProcessByApplicationAndProc(utils.FlagAppName, utils.FlagProcId)
				utils.HandleError(err)

				ParallelExec(nil, (*applicationV2s)[0].Name, args0, proc.ProcID, 0)
				// DoQuery((*applicationV2s)[0].Name, args0, proc.ProcID, proc.LastSeen)
			} else if *hasApp {
				ParallelExec(nil, (*applicationV2s)[0].Name, args0, "", 0)
			} else {
				fmt.Println(cmdusage)
				// return
			}
		}
	}

}


func SerialSearch(cmdusage string, args0 string, hasApp *bool, hasMultipleApps *bool, hasProc *bool,
	applicationV2s *[]*applications.ApplicationV2) {

	create_tq(0)

	*hasApp = true
	if *hasMultipleApps {
		ParMaxCopies = len(*applicationV2s)
		wg := sync.WaitGroup{}
		for _, app := range *applicationV2s {
			wg.Add(1)
			go func(app *applications.ApplicationV2, wg *sync.WaitGroup) {
				defer wg.Done()
				DoQuery(0, app.Name, args0, "", app.LastSeen)
			}(app, &wg)
		}
		wg.Wait()
	} else {
		ParMaxCopies = len(*applicationV2s)*(len(tq)-1)
		if len(*applicationV2s) > 0 {
			if utils.FlagProcId != "" {
				*hasProc = true
			}
			if *hasApp && *hasProc {
				proc, err := GetProcessByApplicationAndProc(utils.FlagAppName, utils.FlagProcId)
				utils.HandleError(err)
				DoQuery(0, (*applicationV2s)[0].Name, args0, proc.ProcID, proc.LastSeen)
			} else if *hasApp {
				fmt.Println("here 10")
				DoQuery(0, (*applicationV2s)[0].Name, args0, "", (*applicationV2s)[0].LastSeen)
			} else {
				fmt.Println(cmdusage)
			}
		}
	}
}

func FindApps(hasApp *bool, hasMultipleApps *bool, hasProc *bool,
	applicationV2s *[]*applications.ApplicationV2) {

	*hasApp = false
	*hasMultipleApps = false
	*hasProc = false
	if utils.FlagAppName == "" {
		a, err := RunSelectApplicationForNamespacePrompt(false)
		utils.HandleError(err)
		*applicationV2s = append(*applicationV2s, a)
	} else {
		if strings.Contains(utils.FlagAppName, ",") {
			apps := strings.Split(utils.FlagAppName, ",")
			for _, appI := range apps {
				*hasMultipleApps = true
				a, err := GetApplicationByName(appI)
				utils.HandleError(err)
				*applicationV2s = append(*applicationV2s, a)
			}
		} else {
			a, err := GetApplicationByName(utils.FlagAppName)
			utils.HandleError(err)
			*applicationV2s = append(*applicationV2s, a)
		}
	}
}

// var SearchPar int = 8

// in hours
// var ParPeriod int64 = 3


func create_period() []int {

	// supports hours, 1,2,3,4,6,8, 12, 24
	if (24%utils.FlagParPeriod)!=0	{
		err:= errors.New("unevenly divide day parPeriod")
		utils.HandleError(err)
	}

	var hr_period [] int
	for i:=0; i<24; i+=int(utils.FlagParPeriod) {
		hr_period = append(hr_period, i)
	}
	return hr_period
}

// 0 serial
// 1 parallel
func create_tq(mode int) {

	setTimeRange(0)

	if mode==0	 {
		tq = append(tq, st)
		tq = append(tq, et)
		return
	}

	if (24%utils.FlagParPeriod)!=0 {
		err := errors.New("unevenly divide parPeriod")
		utils.HandleError(err)
	}


	itime := st
	tq = append(tq, itime)
	itime = itime.Add(time.Duration(utils.FlagParPeriod * 60 * 60 * 1000000000))
	for {
		if (itime.Unix() > et.Unix()) {
			tq = append(tq, et)
			break
		} else {
			tq = append(tq, itime)
		}
			// milli-second accuracy
		itime = itime.Add(time.Duration(utils.FlagParPeriod * 60 * 60 * 1000000000))

	}

	// fmt.Println("debug tq=", tq)

}

// original

func timeFormat(t time.Time) string {
	if !utils.FlagSubSecond {
		ts := t.Format(time.RFC3339Nano)
		// fmt.Println("ts=", ts)
		return ts
	} else {
		return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%06dZ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		)
	}
}

func setTimeRange(lastSeen int64) {

	//var st time.Time
	//var et time.Time

	// dummy use
	in := &query.QueryProperties{
		Namespace: "dummyns",
		PageSize:  100,
		QType:     1,
	}

	var err error

	if utils.FlagBegTime != "" && utils.FlagEndTime != "" {

		st, err = dateparse.ParseLocal(utils.FlagBegTime)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("begin time (-b) error string <%s>", utils.FlagBegTime))
			//fmt.Println("ERR> beg time error string ", utils.FlagBegTime)
			//os.Exit(1)
		}
		et, err = dateparse.ParseLocal(utils.FlagEndTime)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("begin time (-b) error string <%s>", utils.FlagBegTime))
			//fmt.Println("ERR> end time error string ", utils.FlagEndTime)
			//os.Exit(1)
		}

	} else if utils.FlagLogsSince != "" && utils.FlagEndTime != "" {

		et, err = dateparse.ParseLocal(utils.FlagEndTime)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("end time (-e) error string <%s>", utils.FlagEndTime))
		}

		d, err := time.ParseDuration(utils.FlagLogsSince)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("Unable to parse duration (-s) <%s>", utils.FlagLogsSince))
		}
		st = et.Add(-1 * d)

	} else if utils.FlagLogsSince != "" && utils.FlagBegTime != "" {

		st, err = dateparse.ParseLocal(utils.FlagBegTime)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("begin time (-b) error string <%s>", utils.FlagBegTime))
		}

		d, err := time.ParseDuration(utils.FlagLogsSince)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("Unable to parse duration (-s) <%s>", utils.FlagLogsSince))
		}
		et = st.Add(d)

	} else if utils.FlagLogsSince != "" {

		if lastSeen > 0 {
			et = utils.GetStartTime(lastSeen)
		} else {
			et = time.Now().UTC()
		}

		d, err := time.ParseDuration(utils.FlagLogsSince)
		if err != nil {
			utils.HandleError2(err, fmt.Sprintf("Unable to parse duration (-s) <%s>", utils.FlagLogsSince))
		}
		st = et.Add(-1 * d)

	} else {
		// default search
		err := errors.New("Need to set search period or search time range (-s -b -e)")
		utils.HandleError(err)
		// utils.HandleError2(err, fmt.Sprintf("Need to set search period or search time range (-s -b -e)"))
	}

	//in.StartTime = timeFormat(st)
	//in.EndTime = timeFormat(et)

	if et.UnixNano()-st.UnixNano() < 0 {

		// doesn't work
		//in.StartTime = et.Format(time.RFC3339Nano)
		//in.EndTime = st.Format(time.RFC3339Nano)
		in.StartTime = timeFormat(st)
		in.EndTime = timeFormat(et)

		// fmt.Printf("ERR> EndTime %s is older than BegTime: %s \n", in.EndTime, in.StartTime)
		// fmt.Printf("ERR> difference is %d \n", et.UnixNano()-st.UnixNano())
		// os.Exit(1)
	} else {
		in.StartTime = timeFormat(et)
		in.EndTime = timeFormat(st)

	}
	pOnce.Do(func() {
		fmt.Println(" StartTime: ", in.EndTime)
		fmt.Println(" EndTime  : ", in.StartTime)
	})

	//fmt.Println("BegTime2:", in.StartTime)
	//fmt.Println("EndTime2:", in.EndTime)
}

func postQuery(ti int,
	applicationName,
	searchTerm,
	procId string,
	lastSeen int64) (string, query.QueryServiceClient, error) {
	//fmt.Println("Enter postQuery2")
	//fmt.Println("searchTerm ", searchTerm)
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return "", nil, err
	}
	client := query.NewQueryServiceClient(conn)

	in := &query.QueryProperties{
		Namespace: utils.GetDefaultNamespace(),
		PageSize:  utils.GetPageSize(),
		QType:     query.QueryType_Fetch,
	}

	if applicationName != "" {
		in.ApplicationNames = []string{applicationName}
		//  process filter is only available for query
		if procId != "" {
			filterValuesMap := make(map[string]*query.FilterValues)
			filterValuesMap["ProcId"] = &query.FilterValues{
				Values: []string{procId},
			}
			in.Filters = filterValuesMap
		}

	}

	if searchTerm != "" {
		in.KeyWord = searchTerm
		if utils.FlagRegex {
			fmt.Println("Advanced Search Enabled")
			in.QType = query.QueryType_AdvanceSearch
			in.IsAdvanceQuery = true
			_, err := govaluate.NewEvaluableExpression(searchTerm)
			if err != nil {
				utils.HandleError(err)
			}

		} else {
			in.QType = query.QueryType_Search
			in.IsAdvanceQuery = false
		}

	}
	if (ti != -1) {
		var sst string
		if ti == (len(tq) - 2) {
			sst = timeFormat(tq[ti+1])
		} else {
			ttmp := tq[ti+1]
			sst = timeFormat(ttmp.Add(time.Duration(-1 * 1000000000)))
		}

		in.StartTime = sst
		in.EndTime = timeFormat(tq[ti])
		// fmt.Println("here 33 ti:", ti, "  sst=", sst, "   in.EndTime", in.EndTime)
	}

	//in.StartTime = utils.GetTimeAsString(lastSeen)
	//in.EndTime= utils.GetTimeAsString(lastSeen-60*60)
	tm := time.Unix(lastSeen, 0)
	in.StartTime = tm.Format(time.RFC3339)
	//in.StartTime = fmt.Sprint(tm.UTC())
	ttm := tm.Add(-1*1000000000*60*60*24)
	in.EndTime = ttm.Format(time.RFC3339)
	//in.EndTime = fmt.Sprint(ttm.UTC())
	//in.StartTime= timeFormat(tm)
	//in.EndTime= timeFormat(tm.Add(-1*1000000000*60*60))
	fmt.Println("   BegTime: ", in.StartTime)
	fmt.Println("   EndTime: ", in.EndTime)
	// fmt.Println("here 34 st:", in.StartTime, " et:", in.EndTime, "  appName=", applicationName, "  ti=", ti, "searh=", in.KeyWord)

	queryResponse, err := client.Query(grpc_utils.GetGrpcContext(), in)
	if err != nil {
		matched, _ := regexp.MatchString(`^\s*rpc error: code = InvalidArgument desc =\s*$`, err.Error())
		if (matched) {
			// hide out-of-bound lastseen argument
			// fmt.Println("capture lastseen error here <", err.Error(),">")
			return "", nil, nil
		} else {
			fmt.Println("ERR> ", err.Error())
			return "", nil, err
		}
	}

	return queryResponse.QueryId, client, nil
}

func WriteFile() {

	var f *os.File

	if utils.FlagFile != "" {
			/*
				fn, _ := os.Stat(utils.FlagFile)
				if fn != nil {
					fmt.Printf("Outfile file %s already exists, please remove it before proceed\n", utils.FlagFile)
					os.Exit(-1)
					//utils.HandleError2(err, fmt.Sprintf("Outfile file %s already exists, cannot override", utils.FlagFile))
				}
			*/

		if fTmp, err := os.Create(utils.FlagFile); err != nil {
			// if fTmp, err := os.OpenFile(utils.FlagFile, os.O_CREATE|os.O_WRONLY, 0600); err != nil {}
			// if fTmp, err := os.OpenFile(utils.FlagFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {}
			// outputMutex.Unlock()
			utils.HandleError2(err, fmt.Sprintf("Err> Unable to write to file: %s \n", err.Error()))
		} else {
			f = fTmp
			// fmt.Printf("Info> Writing output to %s\n", utils.FlagFile)
		}

		defer f.Close()

		// magic number
		linecnt:=0
		searchcnt:=0

		var tmpstr string
		for {

			select {
			case tmpstr = <-outch:
				{
					if (tmpstr == "FINISHINGSEARCH") {
						searchcnt += 1
					} else {
						linecnt += 1
					}
				}
			case <-time.After(time.Minute * 60):
				{
					fmt.Println("ERR> time out: 60 minutes")
					break
				}
			}

			if tmpstr!="FINISHINGSEARCH" {
				if _, err := f.WriteString(tmpstr); err != nil {
					fmt.Printf("ERR> Cannot write file: %s\n", err.Error())
					break
				}
			}
			/*
				if stat, err := os.Stat(utils.FlagFile); err == nil {
					if stat.Size() > int64(utils.FlagMaxFileSize*1048576) {
						fmt.Printf("Info> Max file size reached. Control file size using -m\n")
						return
						// os.Exit(1)
					}
				}
			*/
			if linecnt > utils.FlagMaxLogLines {
				fmt.Printf("Info> Max log line %d reached. Control file size using -m\n", utils.FlagMaxLogLines)
				break
			}

			if (searchcnt >= ParMaxCopies) {
				fmt.Printf("Info> Successful writing to file: %s, %d streams accounted for\n", utils.FlagFile, searchcnt)
				break
			}
		}
			// outputMutex.Unlock
	}
	// fmt.Println("Exit WriteFile()")
}


func DoQuery(ti int, appName, searchTerm, procId string, lastSeen int64) {

	//fmt.Println("Enter DoQuery2")
	//fmt.Println("BegTime:", utils.FlagBegTime)
	//fmt.Println("EndTime:", utils.FlagEndTime)

	search := (searchTerm != "")

	queryId, client, err := postQuery(ti, appName, searchTerm, procId, lastSeen)

	utils.HandleError(err)

	// outputMutex.Lock()

	if queryId != "" {

    /*
		if utils.FlagFile != "" {
			writeToFile = true
			/ *
			fn, _ := os.Stat(utils.FlagFile)
			if fn != nil {
				fmt.Printf("Outfile file %s already exists, please remove it before proceed\n", utils.FlagFile)
				os.Exit(-1)
				//utils.HandleError2(err, fmt.Sprintf("Outfile file %s already exists, cannot override", utils.FlagFile))
			}
			* /

			if fTmp, err := os.OpenFile(fmt.Sprint(appName, "-", utils.FlagFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
				// outputMutex.Unlock()
				utils.HandleError2(err, fmt.Sprintf("Err> Unable to write to file: %s \n", err.Error()))
			} else {
				// fmt.Printf("Info> file opened %s\n", utils.FlagFile)
				f = fTmp
				fmt.Printf("Info> Writing output to %s\n", utils.FlagFile)
			}

			defer f.Close()
			// outputMutex.Unlock()
		}
       */

		for {
			var response *query.GetDataResponse
			var err error
			if search {
				response, err = client.GetDataPrevious(grpc_utils.GetGrpcContext(), &query.GetDataRequest{
					QueryId: queryId,
				})
				if err != nil {
					utils.HandleError(err)
				}
			} else {
				response, err = client.GetDataNext(grpc_utils.GetGrpcContext(), &query.GetDataRequest{
					QueryId: queryId,
				})
				if err != nil {
					utils.HandleError(err)
				}
			}

			if len(response.Data) > 0 {
				for _, entry := range response.Data {

					pp := "NonePat"
					for kk := range entry.StructuredData {
						if entry.StructuredData[kk].Key == "PatternId" {
							pp = entry.StructuredData[kk].Values[0]
							break
						}
					}

					if utils.FlagEnablePsmod {
						//if pp=="NoPat"

						loglerpart.IncLogLineCount()

						if pp == "NonePat" {
							msg := fmt.Sprintf("%s", entry.Message)
							// fmt.Println("msg=", msg)
							PS := loglerpart.ProcessLogCmd(msg)
							pp = loglerpart.PsCheckAndReturnTag(PS, msg)
						}
					}

					if utils.FlagFile != "" {
						line := fmt.Sprintf("%s %s %s %s %s :: %s",
							entry.Timestamp,
							pp,
							entry.SeverityString,
							entry.AppName,
							entry.ProcID,
							entry.Message,
						)
						if strings.HasSuffix(line, "\n") {
							line = strings.ReplaceAll(line, "\n", "")
						}
						line = fmt.Sprintf("%s\n", line)

						// output to channel
						outch <- line
						/*
						if _, err := f.WriteString(line); err != nil {
							fmt.Printf("Info> Cannot write file: %s\n", err.Error())
							return
							// os.Exit(1)
						}
						if stat, err := os.Stat(utils.FlagFile); err == nil {
							if stat.Size() > int64(utils.FlagMaxFileSize*1048576) {
								fmt.Printf("Info> Max file size reached. Control file size using -m\n")
								return
								// os.Exit(1)
							}
						}
						*/
					} else {
						PrintSyslogMessageForType(entry, "raw")
						time.Sleep(20 * time.Millisecond)
					}
				}
			} else {
				if response.Remaining <= 0 && response.Status == "COMPLETE" {

					if utils.FlagFile != "" {
						outch <- "FINISHINGSEARCH"
					}

					return
					//os.Exit(0)
				}
				// fmt.Println("response.Status=", response.Status)
				time.Sleep(2 * time.Second)
			}
		}
	}
}
