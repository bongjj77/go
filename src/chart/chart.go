package chart

import (
	"fmt"
	"strconv"
	"time"
)

// Section :
type Section struct {
	Number      int
	Host        string
	LatencyList []int64
}

// CollectedData :
type CollectedData struct {
	Stream   string
	Times    []time.Time
	Sections []*Section
}

// NewCollectedData :
func NewCollectedData(sectionCount int) *CollectedData {
	collectedData := &CollectedData{Times: make([]time.Time, 0), Sections: make([]*Section, sectionCount)}
	for index := range collectedData.Sections {
		collectedData.Sections[index] = &Section{Number: index, LatencyList: make([]int64, 0)}
	}

	return collectedData
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
			title: {
				display: true,
				text: 'Media streaming latecy chart'
			},
			scales: {
				xAxes: [{
					display: true,
					scaleLabel: {
						display: true,
						labelString: 'Time'
					}
				}],
				yAxes: [{
					display: true,
					scaleLabel: {
						display: true,
						labelString: 'Latency(ms)'
					},
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

// MakeLatecyChart : make latency chart
// - chart.js
func MakeLatecyChart(collectedData *CollectedData) string {

	labels := ""
	latencyLineList := make([]string, len(collectedData.Sections))
	for index, collectedTime := range collectedData.Times {
		if index != 0 {
			labels += ", "

			for sectionIndex := range latencyLineList {
				latencyLineList[sectionIndex] += ", "
			}

		}
		labels += fmt.Sprintf("'%d:%d:%d'", collectedTime.Hour(), collectedTime.Minute(), collectedTime.Second())

		for sectionIndex := range latencyLineList {
			latencyLineList[sectionIndex] += strconv.FormatInt(collectedData.Sections[sectionIndex].LatencyList[index], 10)
		}
	}

	dataSets := ""
	for index, section := range collectedData.Sections {
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
		}`, index, section.Host,
			latencyLineList[index],
			(index*50)%255, (index*100)%255, (index*150)%255)

	}

	return fmt.Sprintf(latencyFormat, collectedData.Stream, labels, dataSets)
}
