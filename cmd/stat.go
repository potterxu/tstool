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
	"github.com/potterxu/mpts"
	"github.com/potterxu/tstool/io"
	"github.com/potterxu/tstool/util"
	"github.com/spf13/cobra"
)

// statCmd represents the stat command
var (
	interval int

	statCmd = &cobra.Command{
		Use:   "stat",
		Short: "statistics of ts stream",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			runStat(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(statCmd)

	statCmd.Flags().IntVarP(&interval, "interval", "i", 3, "Statistics update interval in seconds")
}

func multicastStatWorker(reader io.IOReaderInterface) {
	pktCnt := make(map[int]int64)
	rows := make([][]string, 0)

	start := time.Now()
	duration := time.Duration(interval) * time.Second

	for {
		if time.Since(start) >= duration {
			rows = [][]string{
				{"pids", "ts rate/bps"},
			}

			keys := make([]int, len(pktCnt))
			i := 0
			for k := range pktCnt {
				keys[i] = k
				i++
			}
			sort.Ints(keys)
			for _, k := range keys {
				bitrate := pktCnt[k] / int64(interval)
				rows = append(rows, []string{fmt.Sprint(k), fmt.Sprint(bitrate)})
				delete(pktCnt, k)
			}
			displayInTableUI(rows)
			start = time.Now()
		}
		buf, ok := reader.Read()
		if ok && buf != nil {
			for i := 0; i < len(buf); i += 188 {
				pktBuf := buf[i : i+188]
				tsPkt := mpts.Packet(pktBuf)
				if tsPkt == nil {
					log.Fatal("invalid ts packet buffer ", pktBuf)
					continue
				}
				pktCnt[tsPkt.PID] += int64(tsPkt.PayloadLength * 8)
			}
		}
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

func multicastStat(args []string) {
	intf := args[0]
	ipAddr := args[1]
	port, err := strconv.Atoi(args[2])
	util.ExitOnErr(err)

	var reader io.IOReaderInterface

	if true {
		// multicast
		reader, err = io.UdpReader(ipAddr, port, intf)
		util.ExitOnErr(err)
	}

	reader.Open()
	defer reader.Close()

	startUI()
	defer listenUIEvent()

	go multicastStatWorker(reader)
}

func runStat(args []string) {
	validateArgs(args)
	multicastStat(args)
}

func startUI() {
	err := ui.Init()
	util.ExitOnErr(err)
}

func listenUIEvent() {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			ui.Close()
			return
		}
	}
}

func displayInTableUI(rows [][]string) {
	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = false
	table.TextAlignment = ui.AlignCenter
	table.Rows = rows
	table.SetRect(0, 0, 40, len(rows)+4)
	ui.Clear()
	ui.Render(table)
}
