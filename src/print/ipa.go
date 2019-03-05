package print

import (
	"bytes"
	"net/http"
	"path"
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"

	"gopkg.in/resty.v1"
	// "net/http"
	// "net/url"

	"os"
	"regexp"

	"../utils"
)

type AuthSuccess struct {
	/* variables */
}

func Quicklatex(ipa string, localFolder string, cp string) {

	/* generate hash */
	// fmt.Println("## " + ipa + " ##")
	hash := utils.HashString(ipa)
	// fmt.Println(hash)

	newPath := path.Join(localFolder, hash+".png")
	// fmt.Println(newPath)
	fcolor := "000000"
	if cp == "c" {
		fcolor = "FFFFFF"
	}
	/* check if file is present */
	if _, err := os.Stat(newPath); err == nil {
		// file exists
		return
	} else if os.IsNotExist(err) {
		// file does *not* exist
		/* send the request to Quicklatex */
		sendData := map[string]interface{}{
			"fcolor":   fcolor,
			"fsize":    "99px",
			"formula":  `$\textipa{` + ipa + `}$`,
			"mode":     "0",
			"out":      "1",
			"preamble": `\usepackage{tipa}`,
			"remhost":  "github.com/uconn-ling/openHouseMap",
			"rnd":      fmt.Sprintf("%f", rand.Float32()*100),
		}
		var sendDataString string = func(m map[string]interface{}) string {
			b := new(bytes.Buffer)
			for key, value := range m {
				fmt.Fprintf(b, `%s=%s&`, key, value)
			}
			str := b.String()
			return str[:len(str)-1]
		}(sendData)
		// log.Printf(sendDataString)

		response, err := resty.R().
			SetBody(sendDataString).
			SetResult(&AuthSuccess{}).
			Post("https://www.quicklatex.com/latex3.f")
		if err != nil {
			log.Fatalln(err)
		}
		// log.Println("###" + response.String() + "----")

		/* read out the url in the response */
		var url string
		regResp, _ := regexp.Compile(`\s+([^ ]+)\s\d+ \d+ \d`)
		matches := regResp.FindStringSubmatch(response.String())
		if matches == nil {
			log.Fatal("Failed to get URL from Quicklatex")
		}
		url = matches[1]
		// log.Println("###" + url)

		/* download the file that the url points to */

		if err := DownloadFile(newPath, url); err != nil {
			panic(err)
		}
	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

	}
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
