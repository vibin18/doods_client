package main

import (
	"fmt"
	"github.com/andersfylling/snowflake"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"github.com/vibin18/doods_client/webhooks"
	"io/ioutil"
	"os"
)

type opts struct {
	File          string  `short:"f"  long:"file"      env:"FILE"  description:"Filename for detecting" default:"vibin3.jpg"`
	DoodsServer   string  `           long:"server"      env:"DOODS_SERVER"  description:"Server name or IP of doods server and port number" default:"192.168.4.1:8082"`
	DiscordToken  string  `           long:"token"      env:"DISCORD_TOKEN"  description:"Discord Webhook token"`
	WebhookId     uint64  `           long:"webhook"      env:"DISCORD_WEBHOOK_ID"  description:"Discord Webhook ID"`
	MinConfidence float64 `           long:"mincon"      env:"MINIMUM_CONFIDENCE"  description:"Minimum confidence level and Max is 100" default:"50"`
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
