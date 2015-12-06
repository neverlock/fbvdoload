package main

import (
	"fmt"
	//	"github.com/mitchellh/ioprogress"
	"github.com/neverlock/ioprogress"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}
var urlStr, fURL string
var DEBUG bool

func main() {
	DEBUG = false
	urlStr = "https://www.facebook.com/AnimeFreeWatch/videos/vb.801962233235589/888053707959774/"
	if len(os.Args) != 4 {
		fmt.Printf("Please use : %s %s [hd/sd] vdo_name.mp4\n", os.Args[0], urlStr)
		os.Exit(0)
	}
	urlStr = os.Args[1]

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil) // <-- URL-encoded payload
	if err != nil {
		defer func() {
			recover()
			log.Printf("Can'thttp.NewRequest\n")
		}()
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		defer func() {
			recover()
			log.Printf("Can't client.Do(req)\n")
		}()
		panic(err)
	}
	if DEBUG {
		fmt.Println(resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[2] {
	case "hd":
		fURL = getHDURL(string(body))
	case "sd":
		fURL = getSDURL(string(body))
	}
	//SD := getSDURL(string(body))
	//HD := getHDURL(string(body))
	out, err := os.Create(os.Args[3])
	defer out.Close()
	resp1, err := http.Get(fURL)
	defer resp1.Body.Close()
	if DEBUG {
		fmt.Println("Content-Lenght=", byteUnitStr(resp1.ContentLength))
	}
	progressR := &ioprogress.Reader{
		Reader:   resp1.Body,
		Size:     resp1.ContentLength,
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, ioprogress.DrawTextFormatBar(50)),
		//DrawFunc: ioprogress.DrawTerminalf(os.Stdout, ioprogress.DrawTextFormatBytes),
		//DrawFunc: ioprogress.DrawTerminal(os.Stdout),
	}
	io.Copy(out, progressR)

}

func getSDURL(b string) string {
	sd := regexp.MustCompile("\"sd_src\":\"(.*?)\"")
	sd_url := sd.FindString(b)
	if DEBUG {
		fmt.Printf("Cut from src:= %s\n", sd_url)
	}

	rep := regexp.MustCompile("\\\\")
	sd_url1 := rep.ReplaceAllString(sd_url, "")
	if DEBUG {
		fmt.Printf("Trim \\ := %s\n", sd_url1)
	}

	sd_url2 := strings.Split(sd_url1, "\"")
	if DEBUG {
		fmt.Printf("Get URL := %s\n", sd_url2[3])
	}
	return sd_url2[3]
}
func getHDURL(b string) string {
	hd := regexp.MustCompile("hd_src\":\"(.*?)\"")
	hd_url := hd.FindString(b)
	if DEBUG {
		fmt.Printf("Cut from src:= %s\n", hd_url)
	}

	rep := regexp.MustCompile("\\\\")
	hd_url1 := rep.ReplaceAllString(hd_url, "")
	if DEBUG {
		fmt.Printf("Trim \\ := %s\n", hd_url1)
	}

	hd_url2 := strings.Split(hd_url1, "\"")
	if DEBUG {
		fmt.Printf("Get URL := %s\n", hd_url2[2])
	}
	return hd_url2[2]
}
func byteUnitStr(n int64) string {
	var unit string
	size := float64(n)
	for i := 1; i < len(byteUnits); i++ {
		if size < 1000 {
			unit = byteUnits[i-1]
			break
		}

		size = size / 1000
	}

	return fmt.Sprintf("%.3g %s", size, unit)
}
