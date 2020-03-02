package chart

import (
	"analyze"
	"fmt"
	"sort"
	"strconv"
	"time"
)

// CollectData : collect data
type CollectData struct {
	Stream      string
	AnalyzeList []*analyze.StreamAnalyze
}

const latencyFormat string = `
<html>
<head>
	<title>%s</title>
	<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.8.0/Chart.bundle.min.js"></script>
	<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.8.0/Chart.min.js"></script>
</head>
<canvas id="line_chart" width="800" height="400"></canvas>
<script>
	var ctx = document.getElementById('line_chart');
	var myChart = new Chart(ctx, {
		type: 'line',
		data: {
			labels: [%s],
			datasets: [%s]
		},
		options: {
			responsive: false,
			scales: {
				yAxes: [{
					ticks: {
						beginAtZero: true
					}
				}]
			},
		}
	});
</script>
</html>
`

type sectionInfo struct {
	host     string
	latencys string
}

// MakeLatecyChart : make latency chart
// - chart.js
// - CollectData -> sectionInfo
// - current only one CollectData
func MakeLatecyChart(collectDataMap map[string]*CollectData) string {

	labels := ""
	stream := ""

	sectionInfos := make(map[int]*sectionInfo)
	for _, collectData := range collectDataMap {

		if len(stream) == 0 {
			stream = collectData.Stream
		}
		for index, analyzeData := range collectData.AnalyzeList {

			// label
			if index != 0 {
				labels += ", "
			}
			labels += "'" + analyzeData.CreateTime.Format(time.Stamp) + "'"

			// data set
			for _, latency := range analyzeData.LatencyList {

				if _, exist := sectionInfos[latency.Section]; exist == false {
					// crate + append
					sectionInfos[latency.Section] = &sectionInfo{latency.Host, strconv.FormatInt(latency.Latency, 10)}
				} else {
					// append
					sectionInfos[latency.Section].latencys += ", " + strconv.FormatInt(latency.Latency, 10)
				}
			}
		}
	}

	// map sort
	keys := make([]int, 0)
	for key := range sectionInfos {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	dataSets := ""
	for index, section := range keys {
		if index != 0 {
			dataSets += ", "
		}

		dataSets += fmt.Sprintf(`
		{
			label: '%d(%s)',
			data: [%s],
			backgroundColor:'rgba(0, 0, 0, 0)',
			borderColor: 'rgba(%d, %d, %d, 1)',
			borderWidth: 1
		}`, section, sectionInfos[section].host,
			sectionInfos[section].latencys,
			(section*50)%255, (section*100)%255, (section*150)%255)
	}

	return fmt.Sprintf(latencyFormat, stream, labels, dataSets)
}
