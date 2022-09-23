package mylib

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {
		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

func HttpDial(url string, log Logging) error {
	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {

		log.Write("error",
			true,
			fmt.Sprintf("Site unreachable : %s, error: %#v", url, err),
		)

	}

	defer conn.Close()

	return err
}

func HttpClient(timeout time.Duration) *http.Client {
	//ref: Copy and modify defaults from https://golang.org/src/net/http/transport.go
	//Note: Clients and Transports should only be created once and reused
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			// Modify the time to wait for a connection to establish
			Timeout:   1 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Transport: &transport,
		Timeout:   timeout * time.Second,
	}

	return &client
}

func Post(client *http.Client, log Logging, headers map[string]string, url string, bodyRequest []byte) (string, string, error) {

	startTime, _ := strconv.Atoi(GetLogId())

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyRequest))

	if len(headers) != 0 {
		for k, v := range headers {

			req.Header.Set(k, v)
		}
	}

	req.Close = false

	if err != nil {
		log.Write("error",
			false,
			fmt.Sprintf("Error Occured : %#v", err),
		)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Write("error",
			false,
			fmt.Sprintf("Error sending request to API endpoint : %#v", err),
		)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Write("error",
			false,
			fmt.Sprintf("Couldn't parse response body : %#v", err),
		)
	}

	endTime, _ := strconv.Atoi(GetLogId())
	elapse := endTime - startTime

	return string(responseBody), strconv.Itoa(elapse), err
}

func Get(client *http.Client, log Logging, contentType string, url string) (string, string, error) {

	startTime, _ := strconv.Atoi(GetLogId())

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", contentType)
	req.Close = false

	if err != nil {
		log.Write("error",
			false,
			fmt.Sprintf("Error Occured : %#v", err),
		)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Write("error",
			false,
			fmt.Sprintf("Error sending request to API endpoint : %#v", err),
		)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Write("error",
			false,
			fmt.Sprintf("Couldn't parse response body : %#v", err),
		)
	}

	endTime, _ := strconv.Atoi(GetLogId())
	elapse := endTime - startTime

	return string(responseBody), strconv.Itoa(elapse), err
}

func GetFormatTime(layout string) string {

	// Standard GO Constant Format :

	// ANSIC       = "Mon Jan _2 15:04:05 2006"
	// UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	// RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	// RFC822      = "02 Jan 06 15:04 MST"
	// RFC822Z     = "02 Jan 06 15:04 -0700"
	// RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	// RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	// RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"
	// RFC3339     = "2006-01-02T15:04:05Z07:00"
	// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	// Kitchen     = "3:04PM"
	// // Handy time stamps.
	// Stamp      = "Jan _2 15:04:05"
	// StampMilli = "Jan _2 15:04:05.000"
	// StampMicro = "Jan _2 15:04:05.000000"
	// StampNano  = "Jan _2 15:04:05.000000000"

	// Using Manual Format :
	// 1. date yyyy-mm-dd = 2006-01-02
	// 2. time hhhh:ii:ss = 15:04:05

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	t := time.Now()
	f := t.In(loc).Format(layout)

	return f
}

func GetUniqId() string {

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	t := time.Now()
	var formatId = t.In(loc).Format("20060102150405.000000")
	uniqId := strings.Replace(formatId, ".", "", -1)

	return uniqId
}

func GetLogId() string {

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	t := time.Now()
	var formatId = t.In(loc).Format("20060102150405")
	logId := strings.Replace(formatId, ".", "", -1)

	return logId
}

func GetDate(dateFormat string) string {

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	t := time.Now()
	var now = t.In(loc).Format(dateFormat)

	return now
}

func GetYesterday(day time.Duration) string {

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	var format = "2006-01-02"

	now := time.Now()
	var curDate = now.In(loc).Format(format)

	t, _ := time.Parse(format, curDate)

	yesterday := 24 * day

	nano := t.Add(-yesterday * time.Hour).UnixNano()

	return time.Unix(0, nano).Format(format)
}

func GetTomorrow(day time.Duration) string {

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	var format = "2006-01-02"

	now := time.Now()
	var curDate = now.In(loc).Format(format)

	t, _ := time.Parse(format, curDate)

	tomorrow := 24 * day

	nano := t.In(loc).Add(tomorrow * time.Hour).UnixNano()

	return time.Unix(0, nano).Format(format)
}

func GetYesterdayWithFormat(day time.Duration, formatDate string) string {

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	//var format = "2006-01-02"

	now := time.Now()
	var curDate = now.In(loc).Format(formatDate)

	t, _ := time.Parse(formatDate, curDate)

	yesterday := 24 * day

	nano := t.Add(-yesterday * time.Hour).UnixNano()

	return time.Unix(0, nano).Format(formatDate)
}

func GetTomorrowWithFormat(day time.Duration, formatDate string) string {

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	//var format = "2006-01-02"

	now := time.Now()
	var curDate = now.In(loc).Format(formatDate)

	t, _ := time.Parse(formatDate, curDate)

	tomorrow := 24 * day

	nano := t.In(loc).Add(tomorrow * time.Hour).UnixNano()

	return time.Unix(0, nano).Format(formatDate)
}

func GetDateAdd(format string, day int, month int, year int) string {

	//set timezone,
	loc, _ := time.LoadLocation(timezone)

	t := time.Now()

	now := t.In(loc).AddDate(year, month, day).Format(format)

	return now
}

func BytesToString(data []byte) string {
	return string(data[:])
}

func InlinePrintingXML(xmlString string) string {
	var unformatXMLRegEx = regexp.MustCompile(`>\s+<`)
	unformatBetweenTags := unformatXMLRegEx.ReplaceAllString(xmlString, "><") // remove whitespace between XML tags
	return strings.TrimSpace(unformatBetweenTags)                             // remove whitespace before and after XML
}

func Concat(args ...string) string {

	var b bytes.Buffer

	for _, arg := range args {
		b.WriteString(arg)
	}

	return b.String()
}

// WriteOnFile function
// Args:
// 1. data = @string
// 2. file = @string
// 3. append = @boolean
// 4. mod = @string
func WriteOnFile(data string, file string, append bool, mode os.FileMode) {

	var f *os.File
	var err error

	if append {
		f, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, mode)
	} else {
		f, err = os.OpenFile(file, os.O_CREATE|os.O_WRONLY, mode)
	}

	if err != nil {
		//log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data); err != nil {
		//log.Println(err)
	}

}

func ReadOnFile(file string) string {

	content, _ := ioutil.ReadFile(file)

	return string(content)
}

func ReduceWords(words string, start int, length int) string {

	runes := []rune(words)
	inputFmt := string(runes[start:length])

	return inputFmt
}

func Base64EncStd(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

func Base64DecStd(data string) string {

	sDec, _ := b64.StdEncoding.DecodeString(data)
	return string(sDec)
}

func Base64EncUrl(data string) string {
	return b64.URLEncoding.EncodeToString([]byte(data))
}

func Base64DecUrl(data string) string {

	sDec, _ := b64.URLEncoding.DecodeString(data)
	return string(sDec)
}

func RNG(min int, max int) int {

	rand.Seed(time.Now().UnixNano())

	return rand.Intn(max-min+1) + min
}

func CounterZeroNumber(length int) string {

	var wordNumbers string

	for w := 0; w < length; w++ {
		wordNumbers += "0"
	}

	return wordNumbers
}

func RemoveTabAndEnter(str string) string {
	space := regexp.MustCompile(`\s+`)
	r := space.ReplaceAllString(str, " ")

	return r
}

func GetMD5(s string) string {

	// Secret key
	data := md5.Sum([]byte(s))
	secretKey := hex.EncodeToString(data[:])

	return secretKey
}

func ReadASingleValueInFile(filename string, keyword string) string {

	var path []string

	if _, err := os.Stat(filename); err == nil {

		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {

			line := scanner.Text()

			if strings.Contains(line, keyword) {

				path = strings.Split(line, "=")
				break
			}

		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}

	return path[1]

}

func ContainsInt(i []int, e int) bool {
	for _, a := range i {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ReadGzFile(filename string) ([]byte, error) {
	fi, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func Copy(srcFolder string, destFolder string) bool {

	cpCmd := exec.Command("cp", "-pa", srcFolder, destFolder)
	err := cpCmd.Run()

	if err == nil {
		return true
	} else {
		return false
	}
}
