package wp

type RenderedData struct {
	Rendered string `json:"rendered"`
}

type ProductAdIdData struct {
	ProductAds []int `json:"product_ad"`
}

type WPData struct {
	Content RenderedData `json:"content"`
	Title RenderedData `json:"title"`
	Excerpt RenderedData `json:"excerpt"`
	AuthorId int `json:"author"`
	AuthorName string
	AuthorImage string
	AuthorSlug string
	AcfData ProductAdIdData `json:"acf"`
	YoastHeadJson YoastHeadJson `json:"yoast_head_json"`
	Slug string `json:"slug"`
}

type AuthorData struct {
	Name string `json:"name"`
	YoastHeadJson YoastHeadJson `json:"yoast_head_json"`
	Slug string `json:"slug"`
	Id int `json:"id"`
}

type YoastHeadJson struct {
	OgImage []SeoImage `json:"og_image"`
}

type SeoImage struct {
	Url string `json:"url"`
}


type WPDataSlice []WPData

const (
	WP_USER_ENDPOINT = "https://spgatsbystg.wpengine.com/wp-json/wp/v2/users"
	WP_POST_ENDPOINT = "https://spgatsbystg.wpengine.com/wp-json/wp/v2/posts"
)