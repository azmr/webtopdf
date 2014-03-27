package main

import (
	"code.google.com/p/gofpdf"
	"fmt"
	"io/ioutil"
	//~ "log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	//~ "strings"
)

func main() {
	http.HandleFunc("/", pdfHandler)
	http.ListenAndServe(":8080", nil)
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	folder := "/home/andrew/pics/imgcollate" // Files are saved in a tmp subdir of this
	err := r.ParseForm()
	if err != nil {
		fmt.Errorf("%v", err)
	}
	fmt.Printf("%v - %v\n", r.Form, len(r.Form))
	if len(r.Form) == 0 {
		fmt.Fprint(w, `<html><title></title>
					<body><form>
					URL: <input type="text" name="url">
					Format: <select name="format"><option>A4</option><option>A3</option><option>A5</option><option>Letter</option><option>Legal</option><option>Scale</option></input>
					<input type="submit">
					</form></body>
					</html>`)
	} else {
		url, format := r.Form["url"][0], r.Form["format"][0]
		
		subfolder, err := ioutil.TempDir(folder, "ic")
		if err != nil {
			panic(err)
		}
		defer removeFolder(subfolder)
		images := collectImages(url, subfolder)
		convertToPDF(url, subfolder, images, format)
		http.ServeFile(w, r, subfolder+"/download.pdf")
	}
}

func removeFolder(folder string) {
	err := os.RemoveAll(folder)
		if err != nil {
			panic(err)
		}
}

// Downloads all the images (currently just jpgs) from the url to a subfolder of name title in containingFolder
func collectImages(sourceUrl, folder string) []imgSource {
	sourceHTML := getHTML(sourceUrl)
	
	//~ configureLogger(folder + "/" + "dl.txt", "", 0)
	
	var imgSources []imgSource
	
	// Finds all jpe?g files in HTML
	jpg := regexp.MustCompile("src=[\"']([^\"']*(\\.jpe?g))[\"']")
	imgURLs := jpg.FindAllStringSubmatch(sourceHTML, -1)
	imgURLsLen := len(imgURLs)
	for i := 1; i <= imgURLsLen; i++{
		imgURL := imgURLs[i-1][1]
		zeroedI := leadZeros(i, imgURLsLen)
		extension := imgURLs[i-1][2]
		filename := zeroedI + extension
		fileStr := folder + "/" + filename
		
		dlImg(imgURL, fileStr)
		fmt.Printf("%v of %d downloading\n", zeroedI, imgURLsLen)
		//~ log.Printf("%v\t-\t%v", fileStr, imgURL)
		imgSources = append(imgSources, imgSource{FileStr: fileStr, URL: imgURL})
	}
	return imgSources
}

// Converts selection of images in a folder to a pdf and outputs that in the same folder
func convertToPDF(pageUrl, outputFolder string, imgSources []imgSource, format string) {
	//~ outputFolder := containingFolder + "/" + filename // In case I later want to change output
	
	pdf := gofpdf.New("P", "mm", format, "/usr/share/fonts/")
	if format == "scale" {
		pdf = gofpdf.New("P", "mm", "A4", "/usr/share/fonts/")
	}
	
	for i, img := range imgSources {
		//~ fmt.Println(img.FileStr)
		makeImagePage(pdf, format, img, i)
		// makes the page the same dimensions as the image
		
	}
	err := pdf.OutputFileAndClose(outputFolder + "/download.pdf")
	if err != nil {
		fmt.Println(err)
	}
}

type imgSource struct {
	FileStr string
	URL     string
}



// Wraps GETting and 'saving' an image as a download
func dlImg(url, fileStr string) {
	img := getIMG(url)
	ioutil.WriteFile(fileStr, img, 0755)
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
// GETs an image from its url and returns as a []byte
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


func makeImagePage(pdf *gofpdf.Fpdf, format string, img imgSource, i int) {
	imgInfo := pdf.RegisterImage(img.FileStr, "")
	imgSize := gofpdf.SizeType{Wd: imgInfo.Width(), Ht: imgInfo.Height()}
	if format == "scale" {
			pdf.AddPageFormat("", imgSize)
			//func (f *Fpdf) Image(fileStr string, x, y, w, h float64, flow bool, tp string, link int, linkStr string)
			pdf.Image(img.FileStr, 0, 0, imgSize.Wd, imgSize.Ht, false, "", 0, img.URL)
		} else {
			pdf.AddPage()
			pgW, pgH, _ := pdf.PageSize(i)
			if imgSize.Wd/imgSize.Ht > pgW/pgH {
				if imgSize.Wd > pgW {
					imgSize.Wd = pgW
				}
				pdf.Image(img.FileStr, 0, 0, imgSize.Wd, 0, false, "", 0, img.URL)
			} else {
				if imgSize.Ht > pgH {
					imgSize.Ht = pgH
				}
				pdf.Image(img.FileStr, 0, 0, 0, imgSize.Ht, false, "", 0, img.URL)
			}
		}
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
// Convenience function that returns the number of difits in a %d representation of a number. e.g 84 in gives 2 out
func numDigits(i int) int {
	return len(strconv.Itoa(i))
}

//~ func configureLogger(fileStr string, prefix string, flags int) {
	//~ file, err := os.OpenFile(fileStr, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
		//~ if err != nil {
			//~ fmt.Print (err)
		//~ }
	//~ log.SetOutput(file)
	//~ log.SetPrefix(prefix)
	//~ log.SetFlags(flags)
//~ }

