package parser

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const URL = "https://vin.drom.ru"

type Report struct {
	Vin      string `json:"vin"`
	Volume   int    `json:"volume"`
	Power    int    `json:"power"`
	CarPlate string `json:"carplate"`
	Color    string `json:"color"`
	Type     string `json:"type"`
}

type CarData struct {
	Image  string
	Report Report
}

func Parse(carPlate string) (*CarData, error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var html string
	var carItems, carImages []*cdp.Node
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("%s/report/%s/", URL, carPlate)),
		chromedp.WaitReady("[data-ftid=\"car-info-item\"], [data-app-root=\"auto-story\"] > .b-media-cont", chromedp.ByQuery),
		chromedp.Text("[data-app-root=\"auto-story\"]", &html, chromedp.ByQuery),
		chromedp.Nodes("[data-ftid=\"car-info-photo\"]", &carImages, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		chromedp.Nodes("[data-ftid=\"car-info-item\"]", &carItems, chromedp.ByQueryAll, chromedp.AtLeast(0)),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse drom by %s", carPlate)
	}

	if strings.Contains(html, "Мы не смогли найти автомобиль с указанным номером кузова") {
		return nil, errors.New(fmt.Sprintf("failed to find car by %s", carPlate))
	}

	return parse(ctx, carItems, carImages)
}

func parse(ctx context.Context, carItems, carImages []*cdp.Node) (*CarData, error) {
	var image string
	for _, item := range carImages {
		if v, ok := item.Attribute("src"); ok {
			image = v
			break
		}
	}

	var report Report
	for _, item := range carItems {
		var innerHTML string
		err := chromedp.Run(ctx,
			chromedp.InnerHTML(item.FullXPath(), &innerHTML),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse innerHTML")
		}

		if strings.Contains(innerHTML, "VIN") {
			report.Vin = getValue(innerHTML)
		}
		if strings.Contains(innerHTML, "Госномер") {
			report.CarPlate = getValue(innerHTML)
		}
		if strings.Contains(innerHTML, "Цвет") {
			report.Color = getValue(innerHTML)
		}
		if strings.Contains(innerHTML, "Тип ТС") {
			report.Type = getValue(innerHTML)
		}
		if strings.Contains(innerHTML, "Объем") {
			report.Volume = getIntValue(innerHTML, " ")
		}
		if strings.Contains(innerHTML, "Мощность") {
			report.Power = getIntValue(innerHTML, " ")
		}
	}

	return &CarData{
		Image:  image,
		Report: report,
	}, nil
}

func getValue(s string) string {
	items := strings.Split(s, "span>")
	if len(items) > 1 {
		return strings.TrimSpace(items[1])
	}

	return ""
}

func getIntValue(s, sep string) int {
	v := getValue(s)
	items := strings.Split(v, sep)
	if len(items) > 0 {
		i, _ := strconv.Atoi(strings.TrimSpace(items[0]))
		return i
	}

	return 0
}
