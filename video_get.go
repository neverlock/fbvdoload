package main

import (
	"fmt"
	"github.com/mitchellh/ioprogress"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func main() {
	urlStr := "https://www.facebook.com/AnimeFreeWatch/videos/vb.801962233235589/888053707959774/"

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
	fmt.Println(resp.Status)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	SD := getSDURL(string(body))
	HD := getHDURL(string(body))
	fmt.Println("Download", SD, "and", HD)
	out, err := os.Create("HD.mp4")
	defer out.Close()
	resp1, err := http.Get(HD)
	defer resp1.Body.Close()
	fmt.Println("Content-Lenght=", byteUnitStr(resp1.ContentLength))
	//	n, err := io.Copy(out, resp1.Body)
	//	fmt.Printf("Write %v bytes\n", n)
	progressR := &ioprogress.Reader{
		Reader:   resp1.Body,
		Size:     resp1.ContentLength,
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, ioprogress.DrawTextFormatBar(50)),
	}
	io.Copy(out, progressR)

}

func getSDURL(b string) string {
	sd := regexp.MustCompile("\"sd_src\":\"(.*?)\"")
	sd_url := sd.FindString(b)
	//	fmt.Printf("Cut from src:= %s\n", sd_url)

	rep := regexp.MustCompile("\\\\")
	sd_url1 := rep.ReplaceAllString(sd_url, "")
	//	fmt.Printf("Trim \\ := %s\n", sd_url1)

	sd_url2 := strings.Split(sd_url1, "\"")
	//	fmt.Printf("Get URL := %s\n", sd_url2[3])
	return sd_url2[3]
}
func getHDURL(b string) string {
	hd := regexp.MustCompile("hd_src\":\"(.*?)\"")
	hd_url := hd.FindString(b)
	//	fmt.Printf("Cut from src:= %s\n", hd_url)

	rep := regexp.MustCompile("\\\\")
	hd_url1 := rep.ReplaceAllString(hd_url, "")
	//	fmt.Printf("Trim \\ := %s\n", hd_url1)

	hd_url2 := strings.Split(hd_url1, "\"")
	//	fmt.Printf("Get URL := %s\n", hd_url2[2])
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
