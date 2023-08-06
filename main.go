package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"simplepractice.com/wp-post-migrator/builder"
	wp "simplepractice.com/wp-post-migrator/wordpress"
)

func main() {
	resp, err := http.Get(fmt.Sprintf("%s?per_page=%d", wp.WP_POST_ENDPOINT, wp.WP_POSTS_PER_PAGE))
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
		var an wp.AuthorNameData
		json.Unmarshal(body, &an)
		v.AuthorName = an.Name
		img := v.YoastHeadJson.OgImage[0].Url
		builderImgUrl := builder.UploadImageToBuilder(img)
		v.YoastHeadJson.OgImage[0].Url = builderImgUrl
		builder.SendToBuilder(v)
	}

}
