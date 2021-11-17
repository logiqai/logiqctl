package services

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/logiqai/logiqctl/api/v1/forwards"
	"github.com/logiqai/logiqctl/grpc_utils"
	"github.com/logiqai/logiqctl/utils"
	"github.com/tatsushid/go-prettytable"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func getForwards() (*forwards.GetForwardsListResponse, error) {

	// in := &forwards.GetForwardsListResponse{}
	in := &empty.Empty{}

	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		//handleError(config, err)
		return nil, err
	}
	defer conn.Close()
	client := forwards.NewForwardServiceClient(conn)
	return client.GetForwardsList(grpc_utils.GetGrpcContext(), in)
	// return client.GetNamespaces(grpc_utils.GetGrpcContext(), &namespace.NamespaceRequest{})
}

func ListForwards() {
	response, err := getForwards()
	if err != nil {
		//handleError(config, err)
		return
	}
	if response != nil && len(response.Configs) > 0 {
		if !utils.PrintResponse(response) {
			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "fwd_id"},
				{Header: "name"},
				{Header: "schema"},
				{Header: "config"},
			}...)
			if err != nil {
				panic(err)
			}

			tbl.Separator = " | "

			for _, fw := range response.Configs {
				fwf := fw.Fields

				// s, _ := json.MarshalIndent(fwf, "", "    ")
				// fmt.Println("fwf=<", string(s), ">")

				id := fwf["id"]
				name := fwf["name"]
				schema := fwf["schema"]
				config := fwf["config"]

				// ff := map[string] interface{} {fwf["config"]}
				// widgets := dashboard["widgets"].([]interface{})

				/*
				data := fmt.Sprintf("%v", config)
				fmt.Println("data=<", data, ">")
				bdata := []byte(data)

				ff := map[string] interface{}{}
				err := protojson.Unmarshal(bdata, &ff)
				*/

				tt, err := protojson.Marshal(config)
				ff := map[string] interface{}{}
				err = json.Unmarshal(tt, &ff)
				if err!=nil {
					utils.HandleError(err)
				}

				cfgstr := ""
				for k , v := range ff {
					cfgstr = cfgstr + string(k) + ":" + fmt.Sprintf("%v", v) + ", "
				}

				tbl.AddRow(id.GetNumberValue(), name.GetStringValue(), schema.GetStringValue(), cfgstr)


			}
			tbl.Print()
		}
	}
}





///////////////////////////////

func getMappings() (*forwards.GetForwardsMappingListResponse, error) {

	// in := &forwards.GetForwardsListResponse{}
	in := &empty.Empty{}

	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		//handleError(config, err)
		return nil, err
	}
	defer conn.Close()
	client := forwards.NewForwardServiceClient(conn)
	return client.GetForwardsMappingList(grpc_utils.GetGrpcContext(), in)
	// return client.GetNamespaces(grpc_utils.GetGrpcContext(), &namespace.NamespaceRequest{})
}

func ListMappings() {

	for_response, err := getForwards()
	ffmap := make(map[string]string)
	if err != nil {
		//handleError(config, err)
		return
	}
	if for_response != nil && len(for_response.Configs) > 0 {
		for _, v := range for_response.Configs {
			vf := v.Fields
			id := fmt.Sprintf("%v", vf["id"].GetNumberValue())
			ffmap[id] = vf["name"].GetStringValue()
		}
	}

	response, err := getMappings()
	if err != nil {
		//handleError(config, err)
		return
	}
	if response != nil && len(response.Mappings) > 0 {
		if !utils.PrintResponse(response) {
			tbl, err := prettytable.NewTable([]prettytable.Column{
				{Header: "fwd_id"},
				{Header: "name"},
				{Header: "namespace"},
				{Header: "application"},
				{Header: "create"},
				{Header: "update"},
			}...)
			if err != nil {
				panic(err)
			}

			tbl.Separator = " | "

			for _, fw := range response.Mappings{
				fwf := fw.Fields
				id := fmt.Sprintf("%v", fwf["id"].GetNumberValue())
				name := ffmap[id]
				namespace := fwf["namespace"].GetStringValue()
				application:= fwf["application"].GetStringValue()
				create := fwf["createts"]
				update := fwf["updatets"]

				if namespace=="-1"	{
					namespace="*"
				}
				if application=="-1" {
					application="*"
				}

				tbl.AddRow(id, name, namespace, application,
					utils.GetTimeAsString(int64(create.GetNumberValue())),
					utils.GetTimeAsString(int64(update.GetNumberValue())))


			}
			tbl.Print()
		}
	}
}





