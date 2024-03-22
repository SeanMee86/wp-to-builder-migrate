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
	resp, err := http.Get(fmt.Sprintf("%s?page=%d&per_page=%d", wp.WP_POST_ENDPOINT, pageNumber, postsInfo[0]))
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
		var an wp.UserData
		json.Unmarshal(body, &an)
		v.AuthorName = an.Name
		v.AuthorImage = an.YoastHeadJson.OgImage[0].Url
		v.AuthorSlug = an.Slug
		img := v.YoastHeadJson.OgImage[0].Url
		v.YoastHeadJson.OgImage[0].Url = builder.UploadImageToBuilder(img)
		v.AuthorImage = builder.UploadImageToBuilder(v.AuthorImage)

		builder.SendPostToBuilder(v, builder.CreateBuilderPost)
	}
}

func MigrateAuthors() {
	
	userSlice := wp.GetWordpressUsers(100, 2)

	fmt.Println(len(userSlice))

	for _, v := range userSlice {
		
		isAuthor := wp.CheckIfIsAuthor(v.Id)

		if isAuthor {
			if len(v.YoastHeadJson.OgImage) > 0 {
				img := v.YoastHeadJson.OgImage[0].Url
				v.YoastHeadJson.OgImage[0].Url = builder.UploadImageToBuilder(img)
			} else {
				yoastData := append(v.YoastHeadJson.OgImage, wp.SeoImage{
					Url: "https://cdn.builder.io/api/v1/image/assets%2F6a01583f2cca435aa4d0f2561c553e0c%2F0c25f9072b2c482982d69195fa10f19b?quality=60&width=200&height=200",
				})
				v.YoastHeadJson.OgImage = yoastData
			}
			builder.SendAuthorToBuilder(v, builder.CreateAuthorPost)
		}
	}
}
