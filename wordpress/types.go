package wp

type MediaData struct {
	Guid    RenderedData `json:"guid"`
	AltText string       `json:"alt_text"`
}

type ProductAdIdData struct {
	ProductAds []int `json:"product_ad"`
}

type RenderedData struct {
	Rendered string `json:"rendered"`
}

type SeoImage struct {
	Url string `json:"url"`
}

type SeoRobots struct {
	Index string `json:"index"`
}

type UserData struct {
	Id            int           `json:"id"`
	Name          string        `json:"name"`
	YoastHeadJson YoastHeadJson `json:"yoast_head_json"`
	Slug          string        `json:"slug"`
	Description   string        `json:"description"`
}

type PostData struct {
	Title           RenderedData  `json:"title"`
	AuthorId        int           `json:"author"`
	Date            string        `json:"date"`
	Content         RenderedData  `json:"content"`
	Slug            string        `json:"slug"`
	FeaturedMediaId int           `json:"featured_media"`
	YoastHeadJson   YoastHeadJson `json:"yoast_head_json"`
}

type YoastHeadJson struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	OgImage     []SeoImage `json:"og_image"`
	Robots      SeoRobots  `json:"robots"`
}
