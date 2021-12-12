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

	"github.com/potterxu/mpts"
	"github.com/potterxu/tstool/io"
	"github.com/potterxu/tstool/parser"
	"github.com/potterxu/tstool/util"
	"github.com/spf13/cobra"
)

var (
	psiCmd = &cobra.Command{
		Use:   "psi (filename)",
		Short: "parse psi",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			runParsePsi(args)
		},
	}
)

func init() {
	parseCmd.AddCommand(psiCmd)
}

func runParsePsi(args []string) {
	if !validateParseArgs(args) {
		return
	}

	inputFileName := args[0]

	patPayload := getFirstPayload(inputFileName, 0)
	if patPayload == nil {
		log.Fatalln("No PAT found")
	}
	pats := mpts.Psi(patPayload).TableData.PAT
	for _, pat := range pats {
		fmt.Println("Program:", pat.ProgramNum)

		pmtPID := pat.ProgramMapPID
		fmt.Println("\tPmtPID:", pmtPID)

		pmtPayload := getFirstPayload(inputFileName, pmtPID)
		pmt := mpts.Psi(pmtPayload).TableData.PMT
		fmt.Println("\tPCR:", pmt.PcrPID)
		fmt.Println("\tProgramDescriptors:")
		for _, d := range pmt.ProgramDescriptors {
			fmt.Println("\t\t Descriptor", d.Tag, d.Data)
		}
		for _, info := range pmt.ElementStreamInfos {
			fmt.Println("\tPID:", info.ElementaryPID)
			for _, d := range info.ESDescriptors {
				fmt.Println("\t\t Descriptor", d.Tag, d.Data)
			}
		}
	}
}

func getFirstPayload(fileName string, pid int) []byte {
	reader := io.FileReader(fileName)
	err := reader.Open()
	util.PanicOnErr(err)
	defer reader.Close()

	parser := parser.PayloadParser(pid)
	for {
		data, ok := reader.Read()
		if !ok {
			break
		}
		payload, ok := parser.Parse(mpts.Packet(data))
		if ok && payload != nil {
			return payload
		}
	}
	return nil
}
