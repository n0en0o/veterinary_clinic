package graphics

import (
	"io"
	"my-docker-app/backend/models"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func DrawPetHealthChart(records []models.HealthRecord, w io.Writer) {
	line := charts.NewLine()

	dates := make([]string, 0)
	weights := make([]opts.LineData, 0)
	temps := make([]opts.LineData, 0)

	for i := len(records) - 1; i >= 0; i-- {
		dates = append(dates, records[i].VisitDate)
		weights = append(weights, opts.LineData{Value: records[i].Weight})
		temps = append(temps, opts.LineData{Value: records[i].Temperature})
	}

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Показатели здоровья питомца"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true), Trigger: "axis"}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithXAxisOpts(opts.XAxis{Name: "Дата визита"}),
	)

	line.SetXAxis(dates).
		AddSeries("Вес (кг)", weights).
		AddSeries("Температура (°C)", temps)

	line.Render(w)

}