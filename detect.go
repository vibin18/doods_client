package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)


type RequestData struct {
	DetectorName string `json:"detector_name"`
	Data         string `json:"data"`
	File         string `json:"file"`
	Detect       struct {
		Person     int64 `json:"person"`
		Cat        int64 `json:"cat"`
		Bicycle    int64 `json:"bicycle"`
		Car        int64 `json:"car"`
		Motorcycle int64 `json:"motorcycle"`
		Truck      int64 `json:"truck"`
		Bird       int64 `json:"bird"`
		Dog        int64 `json:"dog"`
		Horse      int64 `json:"horse"`
		cow        int64 `json:"cow"`
		elephant   int64 `json:"elephant"`
		bear       int64 `json:"bear"`
		umbrella   int64 `json:"umbrella"`
		handbag    int64 `json:"handbag"`
	} `json:"detect"`
}

type ResponseData struct {
	Detections []Detections `json:"detections"`
}

type Detections struct {
	Bottom     float64 `json:"bottom"`
	Confidence float64 `json:"confidence"`
	Label      string  `json:"label"`
	Left       float64 `json:"left"`
	Right      float64 `json:"right"`
	Top        float64 `json:"top"`
}

func NewRequestData() *RequestData {
	return &RequestData{}
}

func (c *RequestData) SetdetectOption(val int64) {
	c.Detect.Person = val
	c.Detect.Cat = val
	c.Detect.Bicycle = val
	c.Detect.Car = val
	c.Detect.Motorcycle = val
	c.Detect.Truck = val
	c.Detect.Bird = val
	c.Detect.Dog = val
	c.Detect.Horse = val
	c.Detect.cow = val
	c.Detect.elephant = val
	c.Detect.bear = val
	c.Detect.umbrella = val
	c.Detect.handbag = val
}

func (c *RequestData) Setdetector_name(val string) {
	c.DetectorName = val
}

func (c *RequestData) SetFiledata(val []byte)  {
	encoded := base64.StdEncoding.EncodeToString(val)
	c.Data = encoded
}


//func DetectImage(imageFile string, minProb int64, wg *sync.WaitGroup) (error, ResponseData ) {
func DetectImage(imageFile []byte, minProb int64) (error, ResponseData) {
	var ret ResponseData
	con := NewRequestData()
	con.SetFiledata(imageFile)
	con.Setdetector_name("default")
	con.SetdetectOption(minProb)

	prettyJSON, err := json.MarshalIndent(con, "", "    ")

	if err != nil {
		log.Panic(err)
	}

	responseBody := bytes.NewBuffer(prettyJSON)

	resp, err := http.Post("http://" + arg.DoodsServer + "/detect", "application/json", responseBody)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(body, &ret)
	if err != nil {
		log.Panic(err)
	}

	//wg.Done()
	return err, ret
}
