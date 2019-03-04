package print

import (
  "bytes"
  // "encoding/json"
  "fmt"
  "gopkg.in/resty.v1"
  "io"
  "log"
  "math/rand"
  // "net/http"
  // "net/url"
  "os"
  "regexp"
)

type AuthSuccess struct {
	/* variables */
}

func Quicklatex (ipa string, localPath string) {

  /* send the request to Quicklatex */
  sendData := map[string]interface{}{
    "fcolor": "000000",
    "fsize": "99px",
    "formula": `$\textipa{` + ipa + `}$`,
    "mode": "0",
    "out": "1",
    "preamble": `\usepackage{tipa}`,
    "remhost": "github.com/uconn-ling/openHouseMap",
    "rnd": fmt.Sprintf("%f", rand.Float32()*100),
  }
  var sendDataString string = func (m map[string]interface{}) string {
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
  // log.Println(response.String())

  /* read out the url in the response */
  var url string
  regResp, _ := regexp.Compile(`\s([^ ]+)\s\d+ \d+ \d`)
  matches := regResp.FindStringSubmatch(response.String())
  if matches == nil {
    log.Fatal("Failed to get URL from Quicklatex")
  }
  url = matches[1]
  log.Println(url)

  /* download the file that the url points to */

  // Get the data
  response, err = resty.R().
    SetDoNotParseResponse(false).
    Get("https://quicklatex.com/cache3/9f/ql_d4da944d8baadd1ebf871eba46812e9f_l3.png")
    // Get(url)
  if err != nil {
    log.Fatal(err) // (when using the variable url) "first path segment in URL cannot contain colon"
  }

  // Create the file
  out, err := os.Create(localPath)
  if err != nil {
    log.Fatal(err)
  }

  // Write the body to the file
  _, err = io.Copy(out, response.RawBody()) // write from 2nd arg to 1st arg
  if err != nil {
    log.Fatal(err) // "http: read on closed response body"
  }
  response.RawBody().Close()
  out.Close()
}
