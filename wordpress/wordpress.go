package wp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RenderedData struct {
	Rendered string `json:"rendered"`
}

type ProductAdIdData struct {
	ProductAds []int `json:"product_ad"`
}

type WPData struct {
	Content       RenderedData `json:"content"`
	Title         RenderedData `json:"title"`
	Excerpt       RenderedData `json:"excerpt"`
	AuthorId      int          `json:"author"`
	AuthorName    string
	AuthorImage   string
	AuthorSlug    string
	AcfData       ProductAdIdData `json:"acf"`
	YoastHeadJson YoastHeadJson   `json:"yoast_head_json"`
	Slug          string          `json:"slug"`
}

type UserData struct {
	Id            int           `json:"id"`
	Name          string        `json:"name"`
	YoastHeadJson YoastHeadJson `json:"yoast_head_json"`
	Slug          string        `json:"slug"`
	Description   string        `json:"description"`
}

type YoastHeadJson struct {
	OgImage []SeoImage `json:"og_image"`
}

type SeoImage struct {
	Url string `json:"url"`
}

type UserDataSlice []UserData

type WPDataSlice []WPData

const (
	WP_USER_ENDPOINT = "https://spgatsbystg.wpengine.com/wp-json/wp/v2/users"
	WP_POST_ENDPOINT = "https://spgatsbystg.wpengine.com/wp-json/wp/v2/posts"
)

func CheckIfIsAuthor(id int) bool {
	resp, err := http.Get(fmt.Sprintf("%s?author=%d", WP_POST_ENDPOINT, id))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var jd WPDataSlice
	json.Unmarshal(body, &jd)

	if len(jd) > 0 {
		return true
	} else {
		return false
	}
}

func GetWordpressUsers() []UserData {
	resp, err := http.Get(fmt.Sprintf("%s?per_page=3", WP_USER_ENDPOINT))

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	var authorSlice []UserData

	json.Unmarshal(body, &authorSlice)

	return authorSlice
}
