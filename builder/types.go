package builder

import wp "simplepractice.com/wp-post-migrator/wordpress"

type convertPost func(wp.PostData) BuilderEntry[BlogPostData]
type convertAuthor func(wp.UserData) BuilderEntry[AuthorData]

type AuthorData struct {
	Name             string     `json:"name"`
	BiographicalInfo string     `json:"biographicalInfo"`
	Image            string     `json:"image"`
	AuthorPage       AuthorPage `json:"authorPage"`
}

type AuthorPage struct {
	Slug            string `json:"slug"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	Index           bool   `json:"index"`
}

type TopicDropdown struct {
	TopicDropdownSelection              string `json:"topics"`
	PracticeManagementDropdownSelection string `json:"practiceManagementSubTopics"`
	EHRSoftwareDropdownSelection        string `json:"ehrSoftwareSubTopics"`
	IndustryGuidanceDropdownSelection   string `json:"industryGuidanceSubTopics"`
	TreatmentsCareDropdownSelection     string `json:"treatmentsCareSubTopics"`
}

type TypeDropdown struct {
	TypeSelection string `json:"types"`
}

type PostCategory struct {
	Topic TopicDropdown `json:"topic"`
	Type  TypeDropdown  `json:"types"`
}

type BlogPostData struct {
	Title           string                `json:"title"`
	AuthorReference BuilderReference      `json:"author"`
	PublishDate     string                `json:"publishDate"`
	Content         string                `json:"articleContent"`
	FeaturedImage   BlogPostFeaturedImage `json:"image"`
	Download        BuilderReference      `json:"download"`
	TOC             TableOfContents       `json:"tableOfContents"`
	Category        PostCategory          `json:"category"`
	SeoSettings     SeoSettings           `json:"seoSettings"`
}

type BlogPostFeaturedImage struct {
	Upload  string `json:"upload"`
	AltText string `json:"altText"`
}

type SeoSettings struct {
	Slug            string `json:"slug"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	Index           bool   `json:"index"`
}

type TableOfContents struct {
	IncludeH2 bool `json:"includeHeading2"`
	IncludeH3 bool `json:"includeHeading3"`
}

type BuilderContentBlocks struct {
	Component BuilderContentComponent `json:"component"`
}

type BuilderContentComponent struct {
	Options BuilderContentOptions `json:"options"`
}

type BuilderContentData struct {
	Blocks      []BuilderContentBlocks `json:"blocks"`
	Image       string                 `json:"cardImage"`
	URL         string                 `json:"url"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
}

type BuilderContentOptions struct {
	Symbol BuilderContentSymbol `json:"symbol"`
}

type BuilderContentSymbol struct {
	Data BuilderContentSymbolData `json:"data"`
}

type BuilderContentSymbolData struct {
	Download BuilderReference `json:"download"`
	Title    string           `json:"articleTitle"`
	AltText  string           `json:"heroImageAltText"`
}

type BuilderDownloadData struct {
	FileName    string `json:"fileName"`
	FileUpload  string `json:"fileUpload"`
	PreservedId string `json:"preservedId"`
}

type BuilderEntry[T any] struct {
	Id              string `json:"id"`
	Data            T      `json:"data"`
	Name            string `json:"name"`
	PublishedStatus string `json:"published"`
	CreatedDate     int    `json:"createdDate"`
	LastUpdated     int64  `json:"lastUpdated"`
}

type BuilderImage struct {
	ImageUrl string `json:"url"`
}

type BuilderReference struct {
	Type  string `json:"@type"`
	Id    string `json:"id"`
	Model string `json:"model"`
}

type BuilderResourceData struct {
	Title           string               `json:"title"`
	PublishDate     int                  `json:"publishDate"`
	Content         string               `json:"content"`
	Image           BuilderResourceImage `json:"image"`
	Download        BuilderReference     `json:"download"`
	ResourcePostToc TableOfContents      `json:"tableOfContents"`
	Category        PostCategory         `json:"category"`
	SeoSettings     SeoSettings          `json:"seoSettings"`
}

type BuilderResourceImage struct {
	Upload  string `json:"upload"`
	AltText string `json:"altText"`
}

type BuilderResults[T any] struct {
	Results []BuilderEntry[T] `json:"results"`
}

type MigrationAsset struct {
	Data            MigrationAssetData `json:"data"`
	Name            string             `json:"name"`
	PublishedStatus string             `json:"published"`
}

type MigrationAssetData struct {
	FileName      string `json:"fileName"`
	BuilderImgUrl string `json:"builderImageUrl"`
}

type TestBlogPostData struct {
	AuthorReference BuilderReference `json:"author"`
}
