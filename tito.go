package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const discountCodeType string = "PercentOffDiscountCode"
const _TEST_ = true

func main() {
	// configure logger
	log.SetFlags(0)

	// get API key from arguments and exit if none was supplied
	apiKey, error := getAPIKey()
	if error != nil {
		log.Fatal(error)
	}
	if _TEST_ {
		fmt.Println("API key: " + apiKey)
	}

	// set variables depending on event name
	// simply hard-coding event name for now
	eventName := "DODL"

	prefix := ""
	releaseId := ""
	filename := ""
	account := ""
	event := ""
	values := []string{}
	quantities := []string{}
	codes := []string{}

	// set variables depending on event name
	if eventName == "FFC" {
		prefix = "FFC24"
		releaseId = "1491301"
		filename = "./ffc.txt"
		account = "fast-flow-conf"
		event = "fast-flow-conf-2024"
		values = []string{"100.00", "20.00"}
		quantities = []string{"1", "100"}
		codes = []string{"100", "20"}
	} else if eventName == "DODL" {
		prefix = "DODL24"
		releaseId = "1456279"
		filename = "./dodl.txt"
		account = "devopsdays-london"
		event = "2024"
		values = []string{"100.00", "50.00", "20.00"}
		quantities = []string{"2", "5", "100"}
		codes = []string{"100", "50", "20"}
	}

	// create HTTP client and set timeout to 10 seconds
	httpClient := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://api.tito.io/v3/%s/%s/discount_codes", account, event)

	// get names of speakers
	names := getNames(filename)

	// loop through names and call ti.to to create discount codes
	for count := 0; count < len(names); count++ {

		for i := 0; i < len(values); i++ {
			code := fmt.Sprintf("%s_%s_%s", prefix, names[count], codes[i])

			// convert to bytes to be used in request body
			discountCode := fmt.Sprintf(`{"discount_code":{"code":"%s","type":"%s","value":"%s","quantity":"%s", "release_ids":[%s]}}`,
				code, discountCodeType, values[i], quantities[i], releaseId)
			requestBody := []byte(discountCode)

			// create a HTTP post request
			request, error := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
			checkError(error)

			// set headers
			token := fmt.Sprintf("Token token=%s", apiKey)
			request.Header.Add("Authorization", token)
			request.Header.Add("Accept", "application/json")
			request.Header.Add("Content-Type", "application/json")

			// create discount code
			if _TEST_ == true {
				fmt.Println("Creating code: " + discountCode)
			} else {
				response, error := httpClient.Do(request)
				checkError(error)

				msg := ""
				if response.StatusCode == 201 {
					msg = fmt.Sprintf("Created discount code: %s\n", code)
				} else {
					msg = fmt.Sprintf("Failed to create discount code: %s\n", code)
				}
				fmt.Println(msg)
			}
		}
	}
}

// get API key from arguments, throw an error if it hasn't been supplied
// rudimentary argument/error handling but will do for now
func getAPIKey() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("no API key specified - exiting")
	} else {
		apiKey := os.Args[1]
		return apiKey, nil
	}
}

// get names of speakers from file, convert to upper case and convert all spaces to underscores
func getNames(filename string) []string {
	var names []string

	file, error := os.Open(filename)
	checkError(error)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		name = strings.ToUpper(name)
		name = strings.ReplaceAll(name, " ", "_")
		names = append(names, name)
	}

	checkError(scanner.Err())
	return names
}

// error handling
func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
