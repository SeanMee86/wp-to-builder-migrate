package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	wp "simplepractice.com/wp-post-migrator/wordpress"
)

type BuilderData struct {
	Title         string `json:"title"`
	Content       string `json:"content"`
	Author        string `json:"author"`
	Excerpt       string `json:"excerpt"`
	ProductAdId   int    `json:"productAdId"`
	FeaturedImage string `json:"featuredImage"`
	Slug          string `json:"slug"`
	AuthorImage   string `json:"authorImage"`
	AuthorSlug    string `json:"authorSlug"`
}

type BuilderImage struct {
	ImageUrl string `json:"url"`
}

type Builder struct {
	Data            BuilderData `json:"data"`
	Name            string      `json:"name"`
	PublishedStatus string      `json:"published"`
}

type convert func(wp.WPData) Builder

const (
	BUILDER_WRITE_API_ENDPOINT = "https://builder.io/api/v1/write"
)

func getToken() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("BUILDER_AUTH_TOKEN")
}

func UploadImageToBuilder(img string) string {
	imgUrlSlice := strings.Split(img, ".")

	imgType := imgUrlSlice[len(imgUrlSlice)-1]

	resp, err := http.Get(img)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://builder.io/api/v1/upload", bytes.NewBuffer(body))
	token := getToken()

	req.Header.Add("Content-Type", fmt.Sprintf("image/%s", imgType))
	req.Header.Add("Authorization", token)

	resp, err = client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, err = io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var builderImgUrl BuilderImage

	json.Unmarshal(body, &builderImgUrl)

	return builderImgUrl.ImageUrl
}

func CreateBuilderPost(wpData wp.WPData) Builder {
	var prod_ad wp.ProductAdIdData
	return Builder{
		Name: wpData.Title.Rendered,
		Data: BuilderData{
			Author:        wpData.AuthorName,
			Title:         wpData.Title.Rendered,
			Content:       wpData.Content.Rendered,
			Excerpt:       wpData.Excerpt.Rendered,
			ProductAdId:   wpData.AcfData.ProductAds[0],
			FeaturedImage: wpData.YoastHeadJson.OgImage[0].Url,
			Slug:          wpData.Slug,
			AuthorImage:   wpData.AuthorImage,
			AuthorSlug:    wpData.AuthorSlug,
		},
		PublishedStatus: "published",
	}
}

func SendPostToBuilder(wpData wp.WPData, generateBuilderData convert) {
	bd := generateBuilderData(wpData)

	bs, err := json.Marshal(bd)

	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/legacy-blog-article", BUILDER_WRITE_API_ENDPOINT), bytes.NewBuffer(bs))

	if err != nil {
		log.Fatalln(err)
	}

	token := getToken()
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
}
