package actions

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"

	"log"
	"net/http"
	"strconv"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

// CoreRequest is a bare-bones JSON request, essential for all other API calls
func CoreRequest(method string, URL string, data []byte, addHeaders map[string]string) (*http.Response, error) {
	// println(string(data[:len(data)]))
	client := &http.Client{}
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(data))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")
	if len(addHeaders) > 0 {
		for k, v := range addHeaders {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
		return resp, err
	}
	return resp, nil
}

// ParseMilisecondTime returns a human readable date from millisecond time string
func ParseMilisecondTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	t := time.Unix(0, msInt*int64(time.Millisecond))
	return t, nil
}

// PrintBanner prints a banner to stdout for testing purposes
func PrintBanner(message string) {
	fmt.Println("")
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(message)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("")
}

// GetDestination returns a http response fop a GET request to provided url
func GetDestination(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetDestinationStatus gets a status from a URL GET request
func GetDestinationStatus(url string) {
	resp, _ := GetDestination(url)
	if resp.StatusCode == 200 {
		fmt.Fprintf(color.Output, "%s", color.GreenString(resp.Status))
	} else {
		fmt.Fprintf(color.Output, "%s", color.RedString(resp.Status))
	}
	fmt.Println(" " + url)
}

// DownloadFromURL downloads file from specified URL
func DownloadFromURL(url string, path string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	if len(path) > 0 {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0700)
		}
	}

	// fmt.Println("Downloading", url, "to", path+"/"+fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(path + "/" + fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	if n < 1000 {
		os.Remove(path + "/" + fileName)
	} else {
		fmt.Println(n, "bytes downloaded from", path+"/"+fileName)
	}

}
