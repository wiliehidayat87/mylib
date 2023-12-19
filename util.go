package mylib

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	cr "crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/mergermarket/go-pkcs7"
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

func (l *Utils) GetFormatTime(layout string) string {

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
	loc, _ := time.LoadLocation(l.TimeZone)

	t := time.Now()
	f := t.In(loc).Format(layout)

	return f
}

func (l *Utils) GetUniqId() string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

	t := time.Now()
	var formatId = t.In(loc).Format("20060102150405.000000")
	uniqId := strings.Replace(formatId, ".", "", -1)

	return uniqId
}

func (l *Utils) GetLogId() string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

	t := time.Now()
	var formatId = t.In(loc).Format("20060102150405")
	logId := strings.Replace(formatId, ".", "", -1)

	return logId
}

func (l *Utils) GetDate(dateFormat string) string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

	t := time.Now()
	var now = t.In(loc).Format(dateFormat)

	return now
}

func (l *Utils) GetDateTimeAdd(init string, add int, dateFormat string) string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

	t := time.Now()
	var now string

	if init == "hour" {
		now = t.Add(time.Hour * time.Duration(add)).In(loc).Format(dateFormat)
	} else if init == "minute" {
		now = t.Add(time.Minute * time.Duration(add)).In(loc).Format(dateFormat)
	} else if init == "second" {
		now = t.Add(time.Second * time.Duration(add)).In(loc).Format(dateFormat)
	} else if init == "day" {
		now = t.AddDate(0, 0, add).In(loc).Format(dateFormat)
	} else if init == "month" {
		now = t.AddDate(0, 1, 0).In(loc).Format(dateFormat)
	} else if init == "year" {
		now = t.AddDate(add, 0, 0).In(loc).Format(dateFormat)
	}

	return now
}

func (l *Utils) GetYesterday(day time.Duration) string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

	var format = "2006-01-02"

	now := time.Now()
	var curDate = now.In(loc).Format(format)

	t, _ := time.Parse(format, curDate)

	yesterday := 24 * day

	nano := t.Add(-yesterday * time.Hour).UnixNano()

	return time.Unix(0, nano).Format(format)
}

func (l *Utils) GetTomorrow(day time.Duration) string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

	var format = "2006-01-02"

	now := time.Now()
	var curDate = now.In(loc).Format(format)

	t, _ := time.Parse(format, curDate)

	tomorrow := 24 * day

	nano := t.In(loc).Add(tomorrow * time.Hour).UnixNano()

	return time.Unix(0, nano).Format(format)
}

func (l *Utils) GetYesterdayWithFormat(day time.Duration, formatDate string) string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

	//var format = "2006-01-02"

	now := time.Now()
	var curDate = now.In(loc).Format(formatDate)

	t, _ := time.Parse(formatDate, curDate)

	yesterday := 24 * day

	nano := t.Add(-yesterday * time.Hour).UnixNano()

	return time.Unix(0, nano).Format(formatDate)
}

func (l *Utils) GetTomorrowWithFormat(day time.Duration, formatDate string) string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

	//var format = "2006-01-02"

	now := time.Now()
	var curDate = now.In(loc).Format(formatDate)

	t, _ := time.Parse(formatDate, curDate)

	tomorrow := 24 * day

	nano := t.In(loc).Add(tomorrow * time.Hour).UnixNano()

	return time.Unix(0, nano).Format(formatDate)
}

func (l *Utils) GetDateAdd(format string, day int, month int, year int) string {

	//set timezone,
	loc, _ := time.LoadLocation(l.TimeZone)

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

	var (
		f *os.File
	)

	if append {
		f, _ = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, mode)
	} else {
		f, _ = os.OpenFile(file, os.O_CREATE|os.O_WRONLY, mode)
	}

	if _, err := f.WriteString(data); err != nil {
		fmt.Println(err)
	}
	defer f.Close()
}

func ReadOnFile(file string) string {

	content, _ := os.ReadFile(file)

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

	s := bufio.NewScanner(fz)
	return s.Bytes(), nil
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

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// Encrypt encrypts plain text string into cipher text string
func AES256CBCEnc(unencrypted string, secretKey string) (string, error) {
	key := []byte(secretKey)
	plainText := []byte(unencrypted)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf(`plainText: "%s" has error`, plainText)
	}
	if len(plainText)%aes.BlockSize != 0 {
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(cr.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return fmt.Sprintf("%x", cipherText), nil
}

// Decrypt decrypts cipher text string into plain text string
func AES256CBCDec(encrypted string, secretKey string) (string, error) {
	key := []byte(secretKey)
	cipherText, _ := hex.DecodeString(encrypted)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		panic("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return string(cipherText), nil
}
