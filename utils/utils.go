package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type CtaVariation struct {
	OldHTML   string
	Shortcode string
}

var ctaVar1 = CtaVariation{
	OldHTML:   "<a href=\"https://www.simplepractice.com/sign-up/\" target=\"_blank\" rel=\"noopener\"><img loading=\"lazy\" decoding=\"async\" class=\"alignnone wp-image-3021 size-full\" src=\"https://spgatsby.wpengine.com/wp-content/uploads/2023/05/Pollen-Banner-Variation-2.png\" alt=\"Sign up for a free 30 day trial of SimplePractice\" width=\"690\" height=\"108\" /></a>",
	Shortcode: "<shortcode symbol=\"ContentBanner-General30daytrial\"></shortcode>",
}

var ctaVar2 = CtaVariation{
	OldHTML:   "<a href=\"https://www.simplepractice.com/sign-up/\" target=\"_blank\" rel=\"noopener\"><img loading=\"lazy\" decoding=\"async\" class=\"alignnone wp-image-3021 size-full\" src=\"https://spgatsby.wpengine.com/wp-content/uploads/2023/05/Pollen-Button-Variation-1-1.png\" alt=\"Sign up for a free 30 day trial of SimplePractice\" width=\"690\" height=\"108\" /></a>",
	Shortcode: "<shortcode symbol=\"ContentBanner-GeneralThrive\"></shortcode>",
}

var ctaVar3 = CtaVariation{
	OldHTML:   "<a href=\"https://www.simplepractice.com/sign-up/\" target=\"_blank\" rel=\"noopener noreferrer\"><img loading=\"lazy\" decoding=\"async\" class=\"alignnone wp-image-3021 size-full\" src=\"https://spgatsby.wpengine.com/wp-content/uploads/2023/05/Pollen-Banner-Variation-3-v2.png\" alt=\"Sign up for a free 30 day trial of SimplePractice\" width=\"690\" height=\"108\" /></a>",
	Shortcode: "<shortcode symbol=\"ContentBanner-General30daytrial\"></shortcode>",
}

var ctaVar4 = CtaVariation{
	OldHTML:   "<a href=\"https://www.simplepractice.com/sign-up/\"><img loading=\"lazy\" decoding=\"async\" class=\"alignnone  wp-image-25026\" src=\"https://spgatsby.wpengine.com/wp-content/uploads/2023/05/Pollen-Banner-Variation-1-300x47.png\" alt=\"A banner with a photo of a computer showing a HIPAA-compliant electronic scheduling system and a clickable button to try SimplePractice free for 30 days\" width=\"811\" height=\"127\" /></a>",
	Shortcode: "<shortcode symbol=\"ContentBanner-General30daytrial\"></shortcode>",
}

var ctaVar5 = CtaVariation{
	OldHTML:   "<a href=\"https://www.simplepractice.com/sign-up/\" target=\"_blank\" rel=\"noopener\"><img loading=\"lazy\" decoding=\"async\" class=\"alignnone wp-image-3021 size-full\" src=\"https://spgatsby.wpengine.com/wp-content/uploads/2023/05/Pollen-Button-Variation-2-1.png\" alt=\"Sign up for a free 30 day trial of SimplePractice\" width=\"690\" height=\"108\" /></a>",
	Shortcode: "<shortcode symbol=\"ContentBanner-GeneralThrive\"></shortcode>",
}

var practitionerNumber = CtaVariation{
	OldHTML:   "[acf field=\"number_of_customers\" post_id=\"options\"]",
	Shortcode: "<shortcode options=\"globalPractitionerCount.globalPractitionerCountShort\"></shortcode>",
}

func ReadCsv(filePath string, path string) (string, string, string, string) {
	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV:", err)
	}

	// Process the records
	for i, record := range records {
		if record[0] == "Blog" {
			if record[2] == path {
				return record[3], record[4], record[5], ""
			}
		}
		if record[0] == "Resource" {
			if strings.TrimSpace(strings.ReplaceAll(record[2], "https://www.simplepractice.com", "")) == path {
				return record[3], record[4], record[5], record[6]
			}
		}
		if i != 0 && record[0] == "" {
			break
		}
	}

	return "", "", "", ""
}

func DeleteFile(filePath string) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist, no need to delete.")
			return
		}
	}

	err = os.Remove(filePath)

	if err != nil {
		log.Fatal(err)
	}
}

func EncodeParam(s string) string {
	return url.QueryEscape(s)
}

func HtmlImgToShortcode(htmlString string, endingShortcode string) string {
	htmlString = strings.ReplaceAll(htmlString, ctaVar1.OldHTML, ctaVar1.Shortcode)
	htmlString = strings.ReplaceAll(htmlString, ctaVar2.OldHTML, ctaVar2.Shortcode)
	htmlString = strings.ReplaceAll(htmlString, ctaVar3.OldHTML, ctaVar3.Shortcode)
	htmlString = strings.ReplaceAll(htmlString, ctaVar4.OldHTML, ctaVar4.Shortcode)
	htmlString = strings.ReplaceAll(htmlString, ctaVar5.OldHTML, ctaVar5.Shortcode)
	htmlString = strings.ReplaceAll(htmlString, practitionerNumber.OldHTML, practitionerNumber.Shortcode)
	htmlString = strings.ReplaceAll(htmlString, "Monarch", "TherapyFinder")
	htmlString = strings.ReplaceAll(htmlString, "meetmonarch.com", "therapyfinder.com")
	// htmlString = strings.ReplaceAll(htmlString, "</p>", "</p><p><br></p>")
	shortcode := fmt.Sprintf("<shortcode symbol=\"%s\"></shortcode>", endingShortcode)
	htmlString = fmt.Sprintf("%s %s", htmlString, shortcode)
	return htmlString
}

func WriteToErrorLogs(msg string) {
	logFile, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Create a new logger
	errorLog := log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	errorLog.Print(msg)
}

func ScrapeResourceHtml(path string) string {

	endpoint := fmt.Sprintf("https://www.simplepractice.com%s", path)

	resp, err := http.Get(endpoint)

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".cta-banner").SetHtml("<shortcode symbol=\"ContentBanner-General30daytrial\"></shortcode>")

	doc.Find(".EmailCapture-module--emailCaptureContainer--2ffc8").SetHtml("")

	selection := doc.Find("#mainContent")

	html, err := selection.Html()

	if err != nil {
		log.Fatal(err)
	}

	start := "<style"
	stop := "style>"
	for {
		startIndex := strings.Index(html, start)
		if startIndex == -1 {
			break
		}
		stopIndex := strings.Index(html, stop) + len(stop)
		res := html[:startIndex] + html[stopIndex:]
		html = strings.ReplaceAll(res, "\n\n", "\n")
	}

	
	html = strings.ReplaceAll(html, "<br/>", "")
	return html
}
