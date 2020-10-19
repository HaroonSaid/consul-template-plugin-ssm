package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Checks if only 1 argument is provided
// Returns the first command line argument
func parseInput(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("No argument provided, exiting")
	} else if len(args) > 1 {
		return "", fmt.Errorf("Too many arguments provided only 1 argument is supported, exiting")
	}
	return args[0], nil
}

func retrieveEnv() (string, error) {
	filebyte, err := ioutil.ReadFile("/opt/environment")
	out := string(filebyte[:])
	out = strings.TrimSuffix(out, "\n")
	return out, err
}
func getTLSVersion(tr *http.Transport) string {
	switch tr.TLSClientConfig.MinVersion {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	}

	return "Unknown"
}

// Creates an AWS session
// Retrieves and decrypts a given parameter
func retrieveParam(paramName string, client http.Client, region string) (*ssm.GetParameterOutput, error) {
	ssmPath := paramName
	sess := session.Must(session.NewSession(&aws.Config{
		Region:     &region,
		HTTPClient: &client,
	}))
	svc := ssm.New(sess)
	decrypt := true
	out, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           &ssmPath,
		WithDecryption: &decrypt})
	return out, err
}

// return value from
func getParamValue(paramOutput *ssm.GetParameterOutput) string {
	return *paramOutput.Parameter.Value
}

// Return test param value
func getTestParamValue(param string) (string, error) {
	if param == "TEST_PARAM_VALUE" {
		return param, nil
	} else {
		return "", fmt.Errorf("Wrong value for test-mode, should be: TEST_PARAM_VALUE")
	}
}

func main() {
	// Bypass AWS calls in test mode
	testMode := flag.Bool("test-mode", false, "Enable test mode")
	flag.Parse()

	paramName, err := parseInput(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	// test mode
	if *testMode {
		out, err := getTestParamValue(paramName)
		if err != nil {
			log.Println("There was an error fetching/decrypting the parameter:", paramName)
			log.Fatal(err.Error())
		} else {
			fmt.Println(out)
			return
		}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	tr.ForceAttemptHTTP2 = true
	client := http.Client{Transport: tr}
	region := os.Getenv("AWS_REGION")
	out, err := retrieveParam(paramName, client, region)
	if err != nil {
		log.Println("There was an error fetching/decrypting the parameter:", paramName)
		log.Fatal(err.Error())
	} else {
		fmt.Println(getParamValue(out))
	}

}
