package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
 * Complete the 'numDevices' function below.
 *
 * The function is expected to return an INTEGER.
 * The function accepts following parameters:
 *  1. STRING statusQuery
 *  2. INTEGER threshold
 *  3. STRING dateStr
 * https://jsonmock.hackerrank.com/api/iot_devices/search?status=<statusQuery>&page=<pageNumber>
 */

type Response struct {
	Page        int
	Per_page    int
	Total       int
	Total_pages int
	Data        []DeviceInfo
}

type DeviceInfo struct {
	Id              int
	Timestamp       int64
	Status          string
	Asset           Asset
	OperatingParams OperatingParams
	Parent          Parent
}

type Asset struct {
	Id    int
	Alias string
}

type OperatingParams struct {
	RotorSpeed    int
	Slack         float64
	RootThreshold float64
}

type Parent struct {
	Id    int
	Alias string
}

func numDevices(statusQuery string, threshold int32, dateStr string, currPage int, totalDevice int32) int32 {

	link := "https://jsonmock.hackerrank.com/api/iot_devices/search?status=" + statusQuery + "&page=" + strconv.Itoa(currPage)

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	bodyByrtes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	data := Response{}
	err = json.Unmarshal(bodyByrtes, &data)
	if err != nil {
		panic(err)
	}

	receivedTime, err := time.Parse("01-2006", dateStr)
	if err != nil {
		panic(err)
	}

	var counter int32 = 0
	for _, val := range data.Data {

		t := time.Unix(0, val.Timestamp*int64(time.Millisecond))

		if float64(threshold) < val.OperatingParams.RootThreshold {
			if t.Month().String() == receivedTime.Month().String() && t.Year() == receivedTime.Year() {

				counter += 1

			}

		}
	}

	totalDevice = totalDevice + counter

	if currPage < data.Total_pages {

		return numDevices(statusQuery, threshold, dateStr, currPage+1, totalDevice)
	} else {
		return totalDevice
	}

}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	statusQuery := readLine(reader)

	thresholdTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	threshold := int32(thresholdTemp)

	dateStr := readLine(reader)

	result := numDevices(statusQuery, threshold, dateStr, 1, 0)

	fmt.Fprintf(writer, "%d\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
