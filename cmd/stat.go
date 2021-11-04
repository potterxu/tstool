/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/potterxu/tstool/io"
	"github.com/potterxu/tstool/mpts"
	"github.com/potterxu/tstool/util"
	"github.com/spf13/cobra"
)

// statCmd represents the stat command
var (
	interval int
	dynamic  bool

	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "statistics of ts stream",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("stat", args)
			runStat(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(statCmd)

	statCmd.Flags().IntVarP(&interval, "interval", "i", 3, "Statistics update interval in seconds")
	statCmd.Flags().BoolVarP(&dynamic, "dynamic", "d", false, "Dynamic updating statistics")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dynamicStatisticWorker(reader io.IOInterface) {
	cnt := make(map[int]int64)
	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = false
	table.TextAlignment = ui.AlignCenter

	start := time.Now()
	duration := time.Duration(interval) * time.Second
	stopped := false

	for {
		if time.Since(start) >= duration || stopped {
			table.Rows = [][]string{
				[]string{"pids", "ts rate/bps"},
			}

			keys := make([]int, len(cnt))
			i := 0
			for k := range cnt {
				keys[i] = k
				i++
			}
			sort.Ints(keys)
			for _, k := range keys {
				bitrate := cnt[k] / int64(interval)
				table.Rows = append(table.Rows, []string{fmt.Sprint(k), fmt.Sprint(bitrate)})
				delete(cnt, k)
			}

			table.SetRect(0, 0, 40, (len(keys) + 4))
			ui.Render(table)
			start = time.Now()
		}
		if stopped {
			break
		}
		buf := reader.Read()
		if buf != nil {
			tsPkt := mpts.Packet(buf)
			if tsPkt == nil {
				continue
			}
			cnt[tsPkt.PID] += int64(tsPkt.PayloadLength * 8)
		} else {
			stopped = true
		}
	}
}

func staticStatistics(reader io.IOInterface) {
	start := time.Now()
	cnt := make(map[int]int64)

	for {
		if buf := reader.Read(); buf != nil && time.Since(start) <= time.Duration(interval)*time.Second {
			tsPkt := mpts.Packet(buf)
			if tsPkt == nil {
				continue
			}
			cnt[tsPkt.PID] += int64(tsPkt.PayloadLength * 8)
		} else {
			break
		}
	}

	keys := make([]int, len(cnt))
	i := 0
	for k := range cnt {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	for _, k := range keys {
		bitrate := cnt[k] / int64(interval)
		fmt.Println(k, bitrate)
		delete(cnt, k)
	}

}

func validateArgs(args []string) bool {
	// multicast
	if len(args) < 3 {
		log.Fatal("need interface, multicast address and port")
		return false
	}
	return true
}

func runStat(args []string) {
	validateArgs(args)

	// multicast
	intf := args[0]
	ipAddr := args[1]
	port, err := strconv.Atoi(args[2])
	util.ExitOnErr(err)
	ioReader, err := io.Multicast(ipAddr, port, intf)
	util.ExitOnErr(err)
	ioReader.Open()
	defer ioReader.Close()

	if dynamic {
		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer ui.Close()

		go dynamicStatisticWorker(ioReader)

		uiEvents := ui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		}
	} else {
		staticStatistics(ioReader)
	}
}
