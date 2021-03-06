// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	"github.com/uconn-ling/openHouseMap/src/print"
	"github.com/uconn-ling/openHouseMap/src/data"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "run \"make\" to build pdf",
	Long:  `put pictures and IPA into \"data\" folder as illustrated in \"test/data\"`,
	Run:   mainBuild,
}

const ()

var ()

func init() {
	rootCmd.AddCommand(makeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// makeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// makeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func mainBuild(cmd *cobra.Command, args []string) {
	// fmt.Println("hello world!")
	baseDataDir := "./inputData"
	baseCountryName := "country_"
	convertedPNGsDir := path.Join(baseDataDir, "pngs")
	countries := data.GetData(baseDataDir, baseCountryName)

	for endonym, c := range countries {
		// fmt.Println("### " + endonym)
		// fmt.Println(c)
		// fmt.Println(len(c.People))
		fmt.Println("Checking rendered png to tipa LATEX string # " + c.Endonym.Value + " # from country " + endonym)
		c.Endonym = print.Quicklatex(c.Endonym, convertedPNGsDir, true)

		for personName, p := range c.People {
			fmt.Println("Checking rendered png to tipa LATEX string # " + p.Name.Value + " # from person " + personName)
			p.Name = print.Quicklatex(p.Name, convertedPNGsDir, false)
			c.People[personName] = p
		}
		countries[endonym] = c
	}

	print.CreatePdf(countries)

	// print.Quicklatex(`"ro:man`, "./tmp/")
}
