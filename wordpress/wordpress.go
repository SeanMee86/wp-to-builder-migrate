package wp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	WP_USER_ENDPOINT  = "https://spgatsby.wpengine.com/wp-json/wp/v2/users"
	WP_POST_ENDPOINT  = "https://spgatsby.wpengine.com/wp-json/wp/v2/posts"
	WP_MEDIA_ENDPOINT = "https://spgatsby.wpengine.com/wp-json/wp/v2/media"
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

	var jd []PostData
	json.Unmarshal(body, &jd)

	if len(jd) > 0 {
		return true
	} else {
		return false
	}
}

func GetAuthorSlug(id int) string {
	resp, err := http.Get(fmt.Sprintf("%s/%d", WP_USER_ENDPOINT, id))

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var authorData UserData

	json.Unmarshal(body, &authorData)

	return authorData.Slug

}

func GetWordpressPosts(postAmt int, pageNum int) []PostData {
	resp, err := http.Get(fmt.Sprintf("%s?page=%d&per_page=%d&filter[orderby]=date&order=asc", WP_POST_ENDPOINT, pageNum, postAmt))
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var postSlice []PostData
	json.Unmarshal(body, &postSlice)
	return postSlice
}

func GetWordpressMedia(id int) MediaData {
	resp, err := http.Get(fmt.Sprintf("%s/%d", WP_MEDIA_ENDPOINT, id))

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var imgData MediaData

	json.Unmarshal(body, &imgData)

	return imgData
}

func GetWordpressUsers(userAmt int, pageNum int) []UserData {
	resp, err := http.Get(fmt.Sprintf("%s?per_page=%d&page=%d", WP_USER_ENDPOINT, userAmt, pageNum))

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
