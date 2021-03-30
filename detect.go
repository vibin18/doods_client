package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

//func DetectImage(imageFile string, minProb int64, wg *sync.WaitGroup) (error, ResponseData ) {
func DetectImage(imageFile io.Reader, minProb int64) (error, ResponseData) {
	var ret ResponseData
	con := NewRequestData()
	con.Setdata(imageFile)
	con.Setdetector_name("default")
	con.SetdetectOption(minProb)

	prettyJSON, err := json.MarshalIndent(con, "", "    ")

	if err != nil {
		println(err)
	}

	responseBody := bytes.NewBuffer(prettyJSON)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("http://192.168.4.1:8082/detect", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	err1 := json.Unmarshal(body, &ret)
	if err1 != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	//wg.Done()
	return err, ret
}
