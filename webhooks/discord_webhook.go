package webhooks

import (
	"bytes"
	"fmt"
	"github.com/andersfylling/snowflake"
	"github.com/nickname32/discordhook"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
)

type HookMatter struct {
	Embeditem discordhook.Embed
	ImageFile io.Reader
	ImageName string
}

func NewHookMatter() *HookMatter {
	return &HookMatter{}
}

func (h *HookMatter) SetHookMatterTitle(val string) {
	h.Embeditem.Title = val
}

func (h *HookMatter) SetHookMatterDescription(val string) {
	h.Embeditem.Description = val
}

func (h *HookMatter) SetHookMatterImageFile(val io.Reader) {
	h.ImageFile = val
}

func (h *HookMatter) SetHookMatterImageName(val string) {
	h.ImageName = val
}

func outFormat(clist []map[string]float64) string {
	var ret_string string
	for _, v := range clist {
		for d, r := range v {
			ret := fmt.Sprintf("Detected object %s with %f probablity.\n", d, r)
			ret_string = strings.Join([]string{ret_string, ret}, " ")
		}

	}
	return ret_string
}

func NotifyDiscord(webhookName snowflake.Snowflake, WebHookToken string, imageFile []byte, imagename string, minConfidence string, confidenceList []map[string]float64) {

	desc := outFormat(confidenceList)
	imageFileIO := bytes.NewReader(imageFile)
	hook := NewHookMatter()
	hook.SetHookMatterTitle(fmt.Sprintf("Objects with minimum %s %% probability found.", minConfidence))
	hook.SetHookMatterDescription(fmt.Sprintln(desc))
	hook.SetHookMatterImageFile(imageFileIO)
	hook.SetHookMatterImageName(imagename)

	if len(confidenceList) != 0 {
		wa, err := discordhook.NewWebhookAPI(webhookName, WebHookToken, true, nil)
		if err != nil {
			log.Panic(err)
		}

		msg, err := wa.Execute(nil, &discordhook.WebhookExecuteParams{
			Content: "A.I Detected a motion",

			Embeds: []*discordhook.Embed{
				{
					Title:       hook.Embeditem.Title,
					Description: hook.Embeditem.Description,
				},
			},
		}, hook.ImageFile, hook.ImageName)
		if err != nil {
			log.Panic(err)
		}
		imageId := fmt.Sprint(msg.ID)
		log.Infof("Image id: %s with minimum %s%% probability found", imageId, minConfidence)
	}else {
		log.Errorf("No object with minimum %s%% probability found", minConfidence)
	}
}
