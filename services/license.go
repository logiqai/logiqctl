package services

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/logiqai/logiqctl/api/v1/license"
	"github.com/logiqai/logiqctl/grpc_utils"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
)

var LicenseFile string

func SetLicense() error {

	if LicenseFile == "" {
		fmt.Println("Missing license file")
		return fmt.Errorf("Missing license file")
	} else {
		fmt.Println("license file:", LicenseFile)
		if fileBytes, err := ioutil.ReadFile(LicenseFile); err != nil {
			fmt.Println(err.Error())
			return err
		} else {
			println(string(fileBytes))
			conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
			if err != nil {
				return err
			}
			defer conn.Close()

			licensePayload := &license.LicensePayload{LicenseData: string(fileBytes)}

			client := license.NewLicenseServiceClient(conn)

			if license, err := client.UploadLicense(grpc_utils.GetGrpcContext(), licensePayload); err != nil {
				fmt.Println(err.Error())
				return err
			} else {
				printLicense(license)
				return nil

			}
		}
	}
}

func GetLicense() error {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer conn.Close()

	request := &license.LicenseRequest{}
	client := license.NewLicenseServiceClient(conn)
	if license, err := client.GetLicense(grpc_utils.GetGrpcContext(), request); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		printLicense(license)
	}
	return nil
}

func printLicense(l *license.LicenseResponse) {
	fmt.Println("License Type:", l.Type)
	fmt.Println("Status:", l.Status)
	fmt.Println("Message:", l.Message)
	fmt.Println("Expiry Date:", time.Unix(l.GetExpiryDate(), 0))
}
