package migrators

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"simplepractice.com/wp-post-migrator/builder"
	wp "simplepractice.com/wp-post-migrator/wordpress"
)

func MigratePosts(postsInfo ...int) {
	var pageNumber int

	if postsInfo[1] != 0 {
		pageNumber = postsInfo[1]
	} else {
		pageNumber = 1
	}
	if postsInfo[0] > 100 {
		fmt.Println("Too many posts requested")
		return
	}
	wp_url := fmt.Sprintf("%s?page=%d&per_page=%d", wp.WP_POST_ENDPOINT, pageNumber, postsInfo[0])
	resp, err := http.Get(wp_url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var jd wp.WPDataSlice
	json.Unmarshal(body, &jd)

	for _, v := range jd {
		resp, err = http.Get(fmt.Sprintf("%s/%d", wp.WP_USER_ENDPOINT, v.AuthorId))
		if err != nil {
			log.Fatalln(err)
		}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var an wp.AuthorData
		json.Unmarshal(body, &an)
		v.AuthorName = an.Name
		v.AuthorImage = an.YoastHeadJson.OgImage[0].Url
		v.AuthorSlug = an.Slug
		img := v.YoastHeadJson.OgImage[0].Url
		v.YoastHeadJson.OgImage[0].Url = builder.UploadImageToBuilder(img)
		v.AuthorImage = builder.UploadImageToBuilder(v.AuthorImage)
		fmt.Println(v)

		builder.SendPostToBuilder(v, builder.CreateBuilderPost)
	}
}

func MigrateAuthors(authorsInfo ...int) {

	var authorsPerPage int

	if authorsInfo[1] != 0 {
		authorsPerPage = authorsInfo[1]
	} else {
		authorsPerPage = 1
	}
	if authorsInfo[0] > 100 {
		fmt.Println("Too many posts requested")
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s?per_page=%d", wp.WP_USER_ENDPOINT, authorsPerPage))

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	var authorSlice []wp.AuthorData

	json.Unmarshal(body, &authorSlice)

	fmt.Print(authorSlice)
}
