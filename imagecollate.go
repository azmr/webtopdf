package main

import (
	//~ "code.google.com/p/gofpdf"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	//~ "strings"
)

func main() {
	url := "http://www.cutterandtailor.com/forum/index.php?showtopic=545"
	folder := "/home/andrew/pics/imgcollate" // Files are saved in a subfolder of this called title's string
	title := "MitchellSystemSleeves"
	firstFile := 1 // After slice of files collected, all files from firstFile to the end of the slice are downloaded. '1', not '0' is the first possible file.
	
	collectImages(url, folder, title, firstFile)
	//~ convertToPDF(url, title)
}

func collectImages(url, containingFolder, title string, startNum int) {
	sourceHTML := getHTML(url)
	//~ fmt.Println(sourceHTML)
	
	jpg := regexp.MustCompile("src=[\"']([^\"']*\\.jpe?g)[\"']")
	imgUrls := jpg.FindAllStringSubmatch(sourceHTML, -1)
	imgUrlsLen := len(imgUrls)
	dlLen := len(imgUrls[startNum-1:])
	
	folder := containingFolder + "/" + title
	os.MkdirAll(folder, 0755)
	configureLogger(folder + "/" + "dl.txt", "", 0)
	for i := startNum; i <= dlLen; i++{
		zeroedI := leadZeros(i, imgUrlsLen)
		filename := title + zeroedI
		dlImg(imgUrls[i-1][1], folder + "/" + filename)
		fmt.Printf("%v of %d downloading; %v of %d total\n", leadZeros(i+1-startNum, dlLen), dlLen, zeroedI, imgUrlsLen)
		log.Printf("%v\t-\t%v", filename, imgUrls[i-1][1])
	}
}

func dlImg(url, file string) {
	img := getIMG(url)
	ioutil.WriteFile(file, img, 0755)
}

// Returns a string of the HTML from the URL provided
func getHTML(url string) string{
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	HTML, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	return string(HTML)
}
func getIMG(url string) []byte{
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	img, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	return img
}

func printSliceSlice(sl [][]string) {
    fmt.Printf("Slice length = %d\r\n", len(sl))
    for i := 0; i < len(sl); i++ {
        fmt.Println(sl[i][1])
    }
}


// Returns string with prepended leading zeroes so that currentNum has the same numDigits as the highest number in the list
func leadZeros(currentNum, maxNum int) string {
	maxNumDigits := len(strconv.Itoa(maxNum))
	intAsStr := strconv.Itoa(currentNum)
	currentNumDigits := len(intAsStr)
	for currentNumDigits < maxNumDigits {
		intAsStr = "0"+intAsStr
		currentNumDigits++
	}
	return intAsStr
}
func numDigits(i int) int {
	return len(strconv.Itoa(i))
}

func configureLogger(filename string, prefix string, flags int) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			panic (err)
		}
	//defer file.Close() Doesn't work with this uncommented... why?
	log.SetOutput(file)
	log.SetPrefix(prefix)
	log.SetFlags(flags)
}
//~ func convertToPDF(url, name string) {
	//~ 
//~ }
