package main

import (
	"fmt"
	"github.com/andersfylling/snowflake"
	"github.com/jessevdk/go-flags"
	"github.com/vibin18/doods_client/webhooks"
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

	byteImage, err := ioutil.ReadFile(arg.File) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}


	err, result := DetectImage(byteImage, 40)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		os.Exit(1)
	}


	for _, v := range result.Detections {
		if v.Confidence >= int_confidence {
			itemMap := map[string]float64{
				v.Label: v.Confidence,
			}
			ConfidenceMapList = append(ConfidenceMapList, itemMap)
		}
	}

	id := webhooks.NotifyDiscord(webhook, arg.DiscordToken, byteImage, "alert.jpg", arg.MinConfidence, ConfidenceMapList)
	fmt.Println(id)

}

