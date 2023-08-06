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
	AcfData ProductAdIdData `json:"acf"`
	YoastHeadJson YoastHeadJson `json:"yoast_head_json"`
}

type AuthorNameData struct {
	Name string `json:"name"`
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
	WP_POSTS_PER_PAGE = 1
)