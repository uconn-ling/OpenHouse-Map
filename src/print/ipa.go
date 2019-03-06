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
	"strconv"

	"gopkg.in/resty.v1"
	"github.com/fogleman/gg"
	// "net/http"
	// "net/url"

	"os"
	"regexp"

	"github.com/uconn-ling/openHouseMap/src/utils"
	"github.com/uconn-ling/openHouseMap/src/data"
)

type AuthSuccess struct {
	/* variables */
}

func Quicklatex(input data.Transcription, localFolder string, isCountry bool) data.Transcription {

	/* generate hash */
	hash := utils.HashString(input.Value)
	newPath := path.Join(localFolder, hash+".png")
	// fmt.Println(newPath)
	input.Rendered = data.Picture{}
	input.Rendered.Path = newPath

	/* check if file is present */
	if _, err := os.Stat(newPath); err == nil { // file exists

		im, err := gg.LoadImage(newPath)
		if err != nil {
			panic(err)
		}
		input.Rendered.Width = float64(im.Bounds().Dx())
		input.Rendered.Height = float64(im.Bounds().Dy())
		// log.Printf("Image exists. input = %v", input)
		return input

	} else if os.IsNotExist(err) { // file does *not* exist

		/* send the request to Quicklatex */
		fcolor := "000000"
		if isCountry {
			fcolor = "FFFFFF"
		}
		sendData := map[string]interface{}{
			"fcolor":   fcolor,
			"fsize":    "99px",
			"formula":  `$\textipa{` + input.Value + `}$`,
			"mode":     "0",
			"out":      "1",
			"preamble": `\usepackage[tone]{tipa}\usepackage{lmodern}`,
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
		regResp, _ := regexp.Compile(`\s+([^ ]+)\s\d+ (\d+) (\d+)`)
		matches := regResp.FindStringSubmatch(response.String())
		if matches == nil {
			log.Fatalf("Failed to get URL from Quicklatex. The response was: %s", response.String())
		}
		url = matches[1]
		width, err := strconv.Atoi(matches[2])
		if err != nil {
			panic(err)
		}
		height, err := strconv.Atoi(matches[3])
		if err != nil {
			panic(err)
		}
		// log.Println("###" + url)
		// log.Printf("### width %v height %v", width, height)

		/* download the file that the url points to */
		if err := DownloadFile(newPath, url); err != nil {
			panic(err)
		}

		input.Rendered.Width = float64(width)
		input.Rendered.Height = float64(height)
		return input
	} else {
		// Schrodinger: file may or may not exist. See err for details.

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

	}
	return data.Transcription{} // this should be unreachable
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
