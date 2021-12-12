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
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/potterxu/tstool/io"
	"github.com/potterxu/tstool/util"
	"github.com/spf13/cobra"
)

var (
	outputFileName string
	captureTime    int

	capCmd = &cobra.Command{
		Use:   "cap (interface ipAddress port)",
		Short: "capture ts multicast",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			runCap(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(capCmd)

	capCmd.Flags().StringVarP(&outputFileName, "output", "o", "./output.ts", "Output filename")
	capCmd.Flags().IntVarP(&captureTime, "time", "t", 0, "Captureing time")
}

func runCap(args []string) {
	if !validateCapArgs(args) {
		return
	}
	capMulticast(args)
}

func capMulticast(args []string) {
	intf := args[0]
	ipAddr := args[1]
	port, err := strconv.Atoi(args[2])
	util.PanicOnErr(err)

	reader := io.UdpReader(ipAddr, port, intf)
	err = reader.Open()
	util.PanicOnErr(err)
	defer reader.Close()

	writer := io.FileWriter(outputFileName)
	err = writer.Open()
	util.PanicOnErr(err)
	defer writer.Close()

	log.Default().Println("Capture Start")
	go commonWorker(reader, writer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	capDuraion := time.Duration(captureTime) * time.Second
	timer := time.NewTimer(capDuraion)
	if captureTime <= 0 {
		timer.Stop()
	}

	refreshDuation := time.Duration(200) * time.Millisecond
	totalDuration := time.Duration(0)
	progressTimer := time.NewTimer(refreshDuation)

	running := true
	for running {
		select {
		case <-c:
			running = false
		case <-timer.C:
			running = false
		case <-progressTimer.C:
			totalDuration += refreshDuation
			if captureTime > 0 {
				fmt.Printf("\rCaptured for %v / %v", totalDuration.Round(time.Second), capDuraion)
			} else {
				fmt.Printf("\rCaptured for %v", totalDuration.Round(time.Second))
			}
			progressTimer.Reset(refreshDuation)
		}
	}
	progressTimer.Stop()
	fmt.Println()

	log.Default().Println("Capture done")
}

func validateCapArgs(args []string) bool {
	// multicast
	if len(args) < 3 {
		log.Fatal("use -h check for usage")
		return false
	}
	return true
}
