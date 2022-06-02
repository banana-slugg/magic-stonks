package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/chromedp/chromedp"
)

type webhookData struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	keyPath := filepath.Dir(exePath)
	keysPath := filepath.Join(keyPath, "keys.json")

	hookData := webhookData{}

	keys, err := ioutil.ReadFile(keysPath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(keys), &hookData)
	if err != nil {
		panic(err)
	}

	// dummy session just for the webhook, not sure if actually needed
	dg, err := discordgo.New(":3")
	if err != nil {
		panic(err)
	}

	data, err := scrape()
	if err != nil {
		panic(err)
	}

	webhook := discordgo.WebhookParams{Embeds: []*discordgo.MessageEmbed{
		{
			Title: "Stonks :chart_with_upwards_trend:",
			URL:   "https://www.mtgstocks.com/interests",
			Color: 0x85bb65,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   data[0],
					Value:  fmt.Sprintf("- %v\n- %v (%v)", data[1], data[2], data[4]),
					Inline: false,
				},
				{
					Name:   data[5],
					Value:  fmt.Sprintf("- %v\n- %v (%v)", data[6], data[7], data[9]),
					Inline: false,
				},
				{
					Name:   data[10],
					Value:  fmt.Sprintf("- %v\n- %v (%v)", data[11], data[12], data[14]),
					Inline: false,
				},
			},
		},
	}}

	_, err = dg.WebhookExecute(hookData.ID, hookData.Token, false, &webhook)

	if err != nil {
		panic(err)
	}

}
func scrape() ([]string, error) {
	URL := "https://www.mtgstocks.com/interests"
	data := []string{}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(URL),
		chromedp.Text(`<tbody`, &res, chromedp.NodeVisible),
	)
	if err != nil {
		return nil, err
	}
	rows := strings.Split(res, "\n")

	for _, row := range rows {
		cells := strings.Split(row, "\t")
		for _, cell := range cells {
			if len(cell) > 0 {
				data = append(data, cell)
			}
		}

		if len(data) >= 15 {
			break
		}

	}

	return data, nil

}
