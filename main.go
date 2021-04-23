package main

import (
	"encoding/base64"
	"fmt"
	"github.com/andersfylling/snowflake"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type opts struct {
	File          string `short:"f"  long:"file"      env:"FILE"  description:"Filename for detecting" default:"vibin3.jpg"`
	DiscordToken  string `           long:"token"      env:"DISCORD_TOKEN"  description:"Discord Webhook token"`
	WebhookId     string `           long:"webhook"      env:"DISCORD_WEBHOOK_ID"  description:"Discord Webhook ID"`
	MinConfidence string `           long:"mincon"      env:"MINIMUM_CONFIDENCE"  description:"Minimum confidence level" default:"50"`
}

var (
	argparser *flags.Parser
	arg       opts
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

func (c *RequestData) Setdata(val io.Reader) error {
	content, _ := ioutil.ReadAll(val)
	encoded := base64.StdEncoding.EncodeToString(content)
	c.Data = encoded
	return nil
}

func initArgparser() {
	argparser = flags.NewParser(&arg, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

}

func main() {
	initArgparser()


	uint_webhook, _ := strconv.ParseUint(arg.WebhookId, 10, 64)
	webhook := snowflake.Snowflake(uint_webhook)
	int_confidence, _ := strconv.ParseFloat(arg.MinConfidence, 10)

	var ConfidenceMapList []map[string]float64
	if _, err := os.Stat(arg.File); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("File %q NOT found!", arg.File))
		os.Exit(1)
	}

	f, _ := os.Open(arg.File)
	reader := io.Reader(f)

	err, result := DetectImage(reader, 40)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		os.Exit(1)
	}
	f1, _ := os.Open(arg.File)
	reader1 := io.Reader(f1)

	for _, v := range result.Detections {
		if v.Confidence >= int_confidence {
			itemMap := map[string]float64{
				v.Label: v.Confidence,
			}
			ConfidenceMapList = append(ConfidenceMapList, itemMap)
		}
	}

	id := NotifyDiscord(webhook, arg.DiscordToken, reader1, "alert.jpg", ConfidenceMapList)
	fmt.Println(id)

}

