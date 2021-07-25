package main

import (
	"bytes"
	"fmt"
	"github.com/andersfylling/snowflake"
	"github.com/nickname32/discordhook"
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

func NotifyDiscord(webhook snowflake.Snowflake, token string, imagefile []byte, imagename string, confidence_list []map[string]float64) string {

	desc := outFormat(confidence_list)
	imagefileIO :=bytes.NewReader(imagefile)
	hook := NewHookMatter()
	hook.SetHookMatterTitle(fmt.Sprintf("Objects with minimum %s %% probability found.", arg.MinConfidence))
	hook.SetHookMatterDescription(fmt.Sprintln(desc))
	hook.SetHookMatterImageFile(imagefileIO)
	hook.SetHookMatterImageName(imagename)

	if len(confidence_list) != 0 {
		wa, err := discordhook.NewWebhookAPI(webhook, token, true, nil)
		if err != nil {
			panic(err)
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
			panic(err)
		}
		return fmt.Sprintln(msg.ID)
	}
	return "No object with minimum probability found"
}
