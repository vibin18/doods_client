package main

import (
	"bytes"
	"fmt"
	"github.com/andersfylling/snowflake"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"github.com/vibin18/doods_client/webhooks"
	"io/ioutil"
	"net/http"
	"os"
)

type opts struct {
	File          string  `short:"f"  long:"file"      env:"FILE"  description:"Filename for detecting" default:"vibin3.jpg"`
	DoodsServer   string  `           long:"server"      env:"DOODS_SERVER"  description:"Server name or IP of doods server and port number" default:"192.168.4.1:8082"`
	DiscordToken  string  `           long:"token"      env:"DISCORD_TOKEN"  description:"Discord Webhook token"`
	WebhookId     uint64  `           long:"webhook"      env:"DISCORD_WEBHOOK_ID"  description:"Discord Webhook ID"`
	MinConfidence float64 `           long:"mincon"      env:"MINIMUM_CONFIDENCE"  description:"Minimum confidence level and Max is 100" default:"50"`
	CameraId 	  string `            long:"camera"      env:"CAMERA_NAME"  description:"Name of the camera"`
	ShinobiExporter   string  `       long:"exporter"      env:"SHINOBI_EXPORTER"  description:"Server name or IP of shinobi_exporter and port number" default:"192.168.4.5:8880"`
}

var (
	argparser *flags.Parser
	arg       opts
)

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
	webhook := snowflake.Snowflake(arg.WebhookId)
	minConfidence := arg.MinConfidence
	if !(minConfidence >= 0 && minConfidence <= 100) {
		log.Panicf("Minimum confidence should between 0-100, got %f", minConfidence)
	}
	shinobiExporterUrl := fmt.Sprintf("http://%s/hit",arg.ShinobiExporter)
	var jsonStr = []byte(fmt.Sprintf(`{"title": %s}`, arg.CameraId))
	req, err := http.NewRequest("GET", shinobiExporterUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	log.Infof("Sending hit to exporter %s", arg.ShinobiExporter)
	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("Error sending exporter request %s ",err)
	}
	if resp.StatusCode == 200 {
		log.Infof("Hit succesfully sent")
	}
	defer resp.Body.Close()

	var ConfidenceMapList []map[string]float64

	if _, err := os.Stat(arg.File); os.IsNotExist(err) {
		log.Panicf("File %s NOT found!", arg.File)
	}

	byteImage, err := ioutil.ReadFile(arg.File)
	if err != nil {
		log.Panicf("Byte conversion failed!")
	}

	err, result := DetectImage(byteImage, minConfidence)
	if err != nil {
		log.Panicf("Failed to detect image")
	}

	for _, v := range result.Detections {
		if v.Confidence >= minConfidence {
			itemMap := map[string]float64{
				v.Label: v.Confidence,
			}
			ConfidenceMapList = append(ConfidenceMapList, itemMap)
		}
	}


	webhooks.NotifyDiscord(webhook, arg.DiscordToken, byteImage, "alert.jpg", minConfidence, ConfidenceMapList)
}
