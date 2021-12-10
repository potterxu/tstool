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
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/potterxu/tstool/io"
	"github.com/potterxu/tstool/util"
	"github.com/spf13/cobra"
)

// capCmd represents the stat command
var (
	outputFileName string
	captureTime    int

	capCmd = &cobra.Command{
		Use:   "cap",
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

	select {
	case <-c:
	case <-timer.C:
		log.Default().Println("Captures for", time.Duration(capDuraion))
	}

	log.Default().Println("Capture done")
}

func validateCapArgs(args []string) bool {
	// multicast
	if len(args) < 3 {
		log.Fatal("need interface, multicast address and port")
		return false
	}
	return true
}
