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
	"os"
	"os/signal"
	"path"

	"github.com/potterxu/mpts"
	"github.com/potterxu/tstool/io"
	"github.com/potterxu/tstool/logger"
	"github.com/potterxu/tstool/parser"
	"github.com/potterxu/tstool/util"
	"github.com/spf13/cobra"
)

var (
	outputDir string
	parseCmd  = &cobra.Command{
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

	parseCmd.Flags().StringVarP(&outputDir, "outdir", "o", "out", "output directory")
}

func validateParseArgs(args []string) bool {
	// parse file
	if len(args) == 1 {
		logger.Logger.Println("Parse file", args[0])
		return true
	}
	logger.Logger.Fatal("use -h check for usage")
	return false
}

var (
	programNumToPcrParser = make(map[int]*parser.PcrParserType)
	pidToPayloadParser    = make(map[int]*parser.PayloadParserType)
	pidToOutfile          = make(map[int]string)
	pcrPidToOutfile       = make(map[int]string)
	pidToProgramNum       = make(map[int]int)

	pcrPIDs  = make([]int, 0)
	filename = ""
)

func runParse(args []string) {
	if !validateParseArgs(args) {
		return
	}

	filename = args[0]
	err := os.Mkdir(outputDir, 0755)
	util.PanicOnErr(err)

	dynamicParse()
	return

	patPayload := getFirstPayload(filename, 0)
	pats := mpts.Psi(patPayload).TableData.PAT

	for _, pat := range pats {
		programDir := fmt.Sprintf("Program_%v", pat.ProgramNum)
		programPath := path.Join(outputDir, programDir)
		err := os.Mkdir(programPath, 0755)
		util.PanicOnErr(err)
		// logger.Logger.Println("Parse Program", pat.ProgramNum)

		pmtPID := pat.ProgramMapPID
		pmtPayload := getFirstPayload(filename, pmtPID)
		pmt := mpts.Psi(pmtPayload).TableData.PMT

		pmtOutFile := fmt.Sprintf("pmt_%v.csv", pmtPID)
		pmtOutPath := path.Join(outputDir, programDir, pmtOutFile)
		pidToOutfile[pmtPID] = pmtOutPath

		pcrPID := pmt.PcrPID
		programNumToPcrParser[pat.ProgramNum] = parser.PcrParser(pcrPID)
		pcrOutFile := fmt.Sprintf("pcr_%v.csv", pcrPID)
		pcrOutPath := path.Join(outputDir, programDir, pcrOutFile)
		pcrPidToOutfile[pcrPID] = pcrOutPath
		pcrPIDs = append(pcrPIDs, pcrPID)

		for _, info := range pmt.ElementStreamInfos {
			PID := info.ElementaryPID
			outFile := fmt.Sprintf("%v.csv", PID)
			outPath := path.Join(outputDir, programDir, outFile)
			pidToOutfile[PID] = outPath
			pidToProgramNum[PID] = pat.ProgramNum
			pidToPayloadParser[PID] = parser.PayloadParser(PID)
		}
	}

	parsePcrPIDS()
	parsePesPIDs()
}

func parsePcrPIDS() {
	reader := io.FileReader(filename)
	err := reader.Open()
	util.PanicOnErr(err)
	defer reader.Close()

	head := fmt.Sprintln("Pos", "Pcr")

	pcrPIDtoWriter := make(map[int]io.IOWriterInterface)
	for _, PID := range pcrPIDs {
		w := io.FileWriter(pcrPidToOutfile[PID])
		err := w.Open()
		util.PanicOnErr(err)
		defer w.Close()
		w.Write([]byte(head))
		pcrPIDtoWriter[PID] = w
	}

	for {
		data, ok := reader.Read()
		if !ok {
			break
		}
		if data != nil {
			pkt := mpts.Packet(data)
			if ok {
				for _, parser := range programNumToPcrParser {
					pcr, pos, ok := parser.Parse(pkt)
					if ok {
						record := fmt.Sprintln(pos, pcr)
						pcrPIDtoWriter[pkt.PID].Write([]byte(record))
					}
				}
			}
		}
	}
}

func parsePesPIDs() {
	reader := io.FileReader(filename)
	err := reader.Open()
	util.PanicOnErr(err)
	defer reader.Close()

	head := fmt.Sprintln("Pos", "Size", "Pcr")

	pidToWriter := make(map[int]io.IOWriterInterface)
	for PID := range pidToPayloadParser {
		w := io.FileWriter(pidToOutfile[PID])
		err := w.Open()
		util.PanicOnErr(err)
		defer w.Close()
		w.Write([]byte(head))
		pidToWriter[PID] = w
	}

	for {
		data, ok := reader.Read()
		if !ok {
			break
		}
		if data != nil {
			pkt := mpts.Packet(data)
			if ok {
				for _, parser := range pidToPayloadParser {
					payload, pos, ok := parser.Parse(pkt)
					if ok {
						record := fmt.Sprintln(pos, len(payload), programNumToPcrParser[pidToProgramNum[pkt.PID]].GetPcr(pos))
						pidToWriter[pkt.PID].Write([]byte(record))
					}
				}
			}
		}
	}
}

var stopChan = make(chan os.Signal, 1)

// programNumToPcrParser = make(map[int]*parser.PcrParserType)
// 	pidToPayloadParser    = make(map[int]*parser.PayloadParserType)
// 	pidToOutfile          = make(map[int]string)
// 	pcrPidToOutfile       = make(map[int]string)
// 	pidToProgramNum       = make(map[int]int)

// 	pcrPIDs  = make([]int, 0)
// 	filename = ""

func pxu() {
	reader := io.FileReader(filename)
	err := reader.Open()
	util.PanicOnErr(err)
	defer reader.Close()

	pos := int64(0)
	pidToPayloadParser := make([]*parser.PayloadParserType, 8192)
	for pid, _ := range pidToPayloadParser {
		pidToPayloadParser[pid] = parser.PayloadParser(pid)
	}
	type pktRecord struct {
		Pid int
		Pos int64
		Pkt *mpts.PacketType
	}
	pktRecordChan := make(chan *pktRecord, 1000000)

	pidToPcrParser := make(map[int]*parser.PcrParserType)
	pcrPIDs := make([]bool, 8192)

	for {
		data, ok := reader.Read()
		if !ok {
			break
		}
		if data != nil {
			for i := 0; i < len(data); i += 188 {
				pktBuf := data[i : i+188]
				pkt := mpts.Packet(pktBuf)
				if pkt != nil {
					pid := pkt.PID
					pktRecordChan <- &pktRecord{pid, pos, pkt}
					pktRecordChan <- &pktRecord{-1, -1, nil} // a pause handling signal
					if pid == 0 {

					} else if pcrPIDs[pid] {
						// is pcrPid
						pidToPcrParser[pid].ParseWithPos(pkt, pos)
						for {
							record, _ := <-pktRecordChan
							if record.Pid == -1 {
								break
							}
							payload, payloadPos, ok := pidToPayloadParser[record.Pid].ParseWithPos(record.Pkt, record.Pos)
							if ok {
								// get pcr
							}
						}
					}
					pos++
				}
			}
		}
	}
	stopChan <- os.Interrupt
}

func dynamicParse() {
	go pxu()
	signal.Notify(stopChan, os.Interrupt)

	select {
	case <-stopChan:
	}
}
