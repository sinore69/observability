package main

import (
	"log"
	"math"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	cpuPlot := widgets.NewPlot()
	cpuPlot.Title = "CPU Usage (%)"
	cpuPlot.Data = make([][]float64, 1)
	cpuPlot.Data[0] = make([]float64, 70)
	cpuPlot.SetRect(0, 0, 70, 21)
	cpuPlot.AxesColor = ui.ColorWhite
	cpuPlot.LineColors[0] = ui.ColorCyan
	cpuPlot.Marker = widgets.MarkerBraille

	memPlot := widgets.NewPlot()
	memPlot.Title = "RAM Usage (%)"
	memPlot.Data = make([][]float64, 1)
	memPlot.Data[0] = make([]float64, 70)
	memPlot.SetRect(71, 0, 140, 21)
	memPlot.AxesColor = ui.ColorWhite
	memPlot.LineColors[0] = ui.ColorGreen
	memPlot.Marker = widgets.MarkerBraille

	netPlot := widgets.NewPlot()
	netPlot.Title = "Network I/O (KB/s)"
	netPlot.Data = [][]float64{
		make([]float64, 70),
		make([]float64, 70),
	}
	netPlot.SetRect(0, 21, 140, 45)
	netPlot.AxesColor = ui.ColorWhite
	netPlot.LineColors[0] = ui.ColorYellow
	netPlot.LineColors[1] = ui.ColorMagenta
	netPlot.Marker = widgets.MarkerBraille

	// Network baseline
	prevNetIO, _ := net.IOCounters(false)
	ticker := time.NewTicker(1 * time.Second).C

	ui.Render(cpuPlot, memPlot, netPlot)

	for {
		select {
		case <-ticker:
			// CPU usage
			cpuPercent, _ := cpu.Percent(0, false)
			cpuPlot.Data[0] = append(cpuPlot.Data[0][1:], cpuPercent[0])

			// Memory usage
			vm, _ := mem.VirtualMemory()
			memPlot.Data[0] = append(memPlot.Data[0][1:], vm.UsedPercent)

			// Network I/O
			netIO, _ := net.IOCounters(false)
			sentKB := float64(netIO[0].BytesSent-prevNetIO[0].BytesSent) / 1024
			recvKB := float64(netIO[0].BytesRecv-prevNetIO[0].BytesRecv) / 1024
			prevNetIO = netIO

			netPlot.Data[0] = append(netPlot.Data[0][1:], math.Round(sentKB*10)/10)
			netPlot.Data[1] = append(netPlot.Data[1][1:], math.Round(recvKB*10)/10)

			ui.Render(cpuPlot, memPlot, netPlot)

		case e := <-ui.PollEvents():
			if e.Type == ui.KeyboardEvent && (e.ID == "q" || e.ID == "<C-c>") {
				return
			}
		}
	}
}

