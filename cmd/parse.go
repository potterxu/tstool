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

	"github.com/spf13/cobra"
)

var (
	parsePid int

	parseCmd = &cobra.Command{
		Use:   "parse",
		Short: "parse ts",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			runParse(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.PersistentFlags().IntVar(&parsePid, "pid", -1, "Specify stream PID to be parsed")
}

func validateParseArgs(args []string) bool {
	// parse file
	if len(args) == 1 {
		log.Default().Println("Parse file", args[0])
		return true
	}
	log.Fatal("use -h check for usage")
	return false
}

func runParse(args []string) {
	log.Default().Println("Use -h to check usage of parse command")
}

// func capMulticast(args []string) {
// 	intf := args[0]
// 	ipAddr := args[1]
// 	port, err := strconv.Atoi(args[2])
// 	util.PanicOnErr(err)

// 	reader := io.UdpReader(ipAddr, port, intf)
// 	err = reader.Open()
// 	util.PanicOnErr(err)
// 	defer reader.Close()

// 	writer := io.FileWriter(outputFileName)
// 	err = writer.Open()
// 	util.PanicOnErr(err)
// 	defer writer.Close()

// 	log.Default().Println("Capture Start")
// 	go commonWorker(reader, writer)

// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt)

// 	capDuraion := time.Duration(captureTime) * time.Second
// 	timer := time.NewTimer(capDuraion)
// 	if captureTime <= 0 {
// 		timer.Stop()
// 	}

// 	refreshDuation := time.Duration(200) * time.Millisecond
// 	totalDuration := time.Duration(0)
// 	progressTimer := time.NewTimer(refreshDuation)

// 	running := true
// 	for running {
// 		select {
// 		case <-c:
// 			running = false
// 		case <-timer.C:
// 			running = false
// 		case <-progressTimer.C:
// 			totalDuration += refreshDuation
// 			if captureTime > 0 {
// 				fmt.Printf("\rCaptured for %v / %v", totalDuration.Round(time.Second), capDuraion)
// 			} else {
// 				fmt.Printf("\rCaptured for %v", totalDuration.Round(time.Second))
// 			}
// 			progressTimer.Reset(refreshDuation)
// 		}
// 	}
// 	progressTimer.Stop()
// 	fmt.Println()

// 	log.Default().Println("Capture done")
// }

// func validateCapArgs(args []string) bool {
// 	// multicast
// 	if len(args) < 3 {
// 		log.Fatal("need interface, multicast address and port")
// 		return false
// 	}
// 	return true
// }

// // parser
// parsers := make(map[int]*parser.ParserType)
// pktPos := int64(0)
// for {
// 	buf, _ := reader.Read()
// 	// if !ok || buf == nil (
// 	// 	continue
// 	// )
// 	for i := 0; i < len(buf); i += mpts.TS_PACKET_SIZE {
// 		pktBuf := buf[i : i+mpts.TS_PACKET_SIZE]
// 		pkt := mpts.Packet(pktBuf)

// 		pid := pkt.PID
// 		if pid == 8191 {
// 			continue
// 		}
// 		if parsers[pid] == nil {
// 			parsers[pid] = parser.Parser()
// 		}

// 		parsers[pid].Feed(pkt, pktPos)
// 		if pkt.AdaptationLength > 0 && pkt.Adaptation.PCRFlag {
// 			for p, parser := range parsers {
// 				parser.FeedPcr(pkt.Adaptation.PCR, pktPos)
// 				fmt.Println(p, parser.GetBitrate())
// 			}
// 		}
// 		pktPos++
// 	}
// }
