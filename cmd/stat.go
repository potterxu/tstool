/*
Copyright © 2021 Potter Xu <xujingchuan1995@gmail.com>

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
	"io/ioutil"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/potterxu/mpts"
	"github.com/potterxu/tstool/algo"
	"github.com/potterxu/tstool/io"
	"github.com/potterxu/tstool/logger"
	"github.com/potterxu/tstool/util"
	"github.com/spf13/cobra"
)

// statCmd represents the stat command
var (
	interval int

	statCmd = &cobra.Command{
		Use:   "stat (interface ipAddress port)",
		Short: "statistics of ts stream",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			runStat(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(statCmd)

	statCmd.Flags().IntVarP(&interval, "interval", "i", 400, "Statistics update interval in milliseconds")
}

func multicastStatWorker(reader io.IOReaderInterface) {
	rows := make([][]string, 0)
	pktTimeQueue := make(map[int][]int64)
	const MAX_TIME_INTERVAL int64 = 2000

	timer := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer timer.Stop()

	logger.Logger.SetOutput(ioutil.Discard)
	defer logger.Logger.SetOutput(os.Stderr)

	for {
		select {
		case <-timer.C:
			rows = [][]string{
				{"pids", "ts rate/bps"},
			}

			pids := make([]int, len(pktTimeQueue))
			i := 0
			for k := range pktTimeQueue {
				pids[i] = k
				i++
			}
			sort.Ints(pids)
			for _, pid := range pids {
				if len(pktTimeQueue[pid]) < 2 {
					continue
				}
				lastTime := time.Now().UnixMilli()
				offset := algo.BinarySearchNoSmallerThan(pktTimeQueue[pid], lastTime-MAX_TIME_INTERVAL)
				pktTimeQueue[pid] = pktTimeQueue[pid][offset:]

				if len(pktTimeQueue[pid]) > 0 {
					dur := lastTime - pktTimeQueue[pid][0]
					bitrate := int64(len(pktTimeQueue[pid])) * int64(mpts.TS_PACKET_SIZE) * 8 * 1000 / int64(dur)
					rows = append(rows, []string{fmt.Sprint(pid), fmt.Sprint(bitrate)})
				} else {
					delete(pktTimeQueue, pid)
				}
			}
			displayInTableUI(rows)
		default:
		}
		buf, ok := reader.Read()
		if !ok {
			break
		}
		if buf != nil {
			for i := 0; i < len(buf); i += 188 {
				pktBuf := buf[i : i+188]
				tsPkt := mpts.Packet(pktBuf)
				if tsPkt == nil {
					logger.Logger.Fatal("invalid ts packet buffer ", pktBuf)
					continue
				}

				curTime := time.Now().UnixMilli()
				if pktTimeQueue[tsPkt.PID] == nil {
					pktTimeQueue[tsPkt.PID] = make([]int64, 0)
				}
				pktTimeQueue[tsPkt.PID] = append(pktTimeQueue[tsPkt.PID], curTime)
			}
		}
	}
}

func validateStatArgs(args []string) bool {
	// multicast
	if len(args) < 3 {
		logger.Logger.Fatal("use -h check for usage")
		return false
	}
	return true
}

func multicastStat(args []string) {
	intf := args[0]
	ipAddr := args[1]
	port, err := strconv.Atoi(args[2])
	util.PanicOnErr(err)

	var reader io.IOReaderInterface

	if true {
		// multicast
		reader = io.UdpReader(ipAddr, port, intf)
	}

	err = reader.Open()
	util.PanicOnErr(err)
	defer reader.Close()

	startUI()
	defer ui.Close()

	go multicastStatWorker(reader)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	uiEvents := ui.PollEvents()

	running := true
	for running {
		select {
		case <-c:
			running = false
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				running = false
			}
		}
	}
}

func runStat(args []string) {
	if !validateStatArgs(args) {
		return
	}
	multicastStat(args)
}

func startUI() {
	err := ui.Init()
	util.PanicOnErr(err)
}

func displayInTableUI(rows [][]string) {
	table := widgets.NewTable()
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = false
	table.TextAlignment = ui.AlignCenter
	table.Rows = rows
	table.SetRect(0, 0, 40, len(rows)+2)
	ui.Clear()
	ui.Render(table)
}
