package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"simplepractice.com/wp-post-migrator/utils"
	wp "simplepractice.com/wp-post-migrator/wordpress"
)

const (
	BUILDER_WRITE_API_ENDPOINT   = "https://builder.io/api/v1/write"
	BUILDER_CONTENT_API_ENDPOINT = "https://cdn.builder.io/api/v3/content"
	BUILDER_UPLOAD_API_ENDPOINT  = "https://builder.io/api/v1/upload"
	BLOG_AUTHOR_IMAGE_FOLDER     = "&folder=a20c2bde530e473aaf1db08ed82b1d5a"
	BLOG_RESOURCE_ASSET_FOLDER   = "&folder=3b3cef64dcf449818707f63643453833"
	MODEL_BLOG_POST              = "blog-post"
	MODEL_DOWNLOADS              = "downloads"
	MODEL_MIGRATION_ASSETS       = "migration-assets"
	MODEL_BLOG_AUTHORS           = "blog-authors"
	MODEL_CONTENT_TOPICS         = "content-topics"
	MODEL_CONTENT_TYPE           = "content-type"
	MODEL_RESOURCE_POSTS         = "content-resource-post"
)

func addBuilderToken(req *http.Request) {
	token := getToken()
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
}

func buildTopicsDropdown(topic string, subTopic string) TopicDropdown {
	var topicDropdownSelection TopicDropdown
	switch topic {
	case "Practice management":
		topicDropdownSelection = TopicDropdown{
			TopicDropdownSelection:              topic,
			PracticeManagementDropdownSelection: subTopic,
		}
	case "Industry guidance":
		topicDropdownSelection = TopicDropdown{
			TopicDropdownSelection:            topic,
			IndustryGuidanceDropdownSelection: subTopic,
		}
	case "Treatments & care":
		topicDropdownSelection = TopicDropdown{
			TopicDropdownSelection:          topic,
			TreatmentsCareDropdownSelection: subTopic,
		}
	case "EHR software":
		topicDropdownSelection = TopicDropdown{
			TopicDropdownSelection:       topic,
			EHRSoftwareDropdownSelection: subTopic,
		}
	}

	return topicDropdownSelection
}

func doesAssetExist(fileName string) (string, bool) {
	apiKey := getMainApiKey()
	assetUrl := fmt.Sprintf(
		"%s/%s?apiKey=%s&cachebust=true&query.data.fileName=%s",
		BUILDER_CONTENT_API_ENDPOINT,
		MODEL_MIGRATION_ASSETS,
		apiKey,
		fileName,
	)
	resp, err := http.Get(assetUrl)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	var entries BuilderResults[MigrationAssetData]

	json.Unmarshal(body, &entries)

	if len(entries.Results) > 0 {
		return entries.Results[0].Data.BuilderImgUrl, true
	}
	return "", false
}

func getAllEntries(modelName string) BuilderResults[BuilderContentData] {
	resp, err := http.Get(fmt.Sprintf("%s/%s?apiKey=%s&cachebust=true&limit=100", BUILDER_CONTENT_API_ENDPOINT, modelName, getApiKey()))

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	var entries BuilderResults[BuilderContentData]

	json.Unmarshal(body, &entries)

	return entries
}

func getApiKey() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("BUILDER_API_KEY")
}

func getAuthorId(authorId int) string {

	apiKey := getApiKey()

	authorSlug := wp.GetAuthorSlug(authorId)

	authorQuery := fmt.Sprintf(
		"%s/%s?apiKey=%s&cachebust=true&query.data.authorPage.slug=/%s/",
		BUILDER_CONTENT_API_ENDPOINT,
		MODEL_BLOG_AUTHORS,
		apiKey,
		authorSlug,
	)

	resp, err := http.Get(authorQuery)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var content BuilderResults[BuilderContentData]

	json.Unmarshal(body, &content)

	return content.Results[0].Id
}

func getDownloadByID(id string, pathname string) (BuilderReference, error) {
	apiKey := getApiKey()
	downloadQuery := fmt.Sprintf("%s/%s?apiKey=%s&cachebust=true&query.data.preservedId=%s", BUILDER_CONTENT_API_ENDPOINT, MODEL_DOWNLOADS, apiKey, id)
	resp, err := http.Get(downloadQuery)

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var download BuilderResults[BuilderDownloadData]

	json.Unmarshal(body, &download)

	if len(download.Results) == 0 {
		err := fmt.Errorf("no download found at: %s", pathname)
		return BuilderReference{}, err
	}

	return BuilderReference{
		Type:  "@builder.io/core:Reference",
		Id:    download.Results[0].Id,
		Model: MODEL_DOWNLOADS,
	}, nil
}

func getMainApiKey() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("BUILDER_API_KEY_MAIN")
}

func getMainToken() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("BUILDER_AUTH_TOKEN_MAIN")
}

func getToken() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("BUILDER_AUTH_TOKEN")
}

func CreateAuthorEntry(authorData wp.UserData) BuilderEntry[AuthorData] {
	var index bool
	slug := fmt.Sprintf("/%s/", authorData.Slug)
	if authorData.YoastHeadJson.Robots.Index == "index" {
		index = true
	} else {
		index = false
	}
	return BuilderEntry[AuthorData]{
		Name: authorData.Name,
		Data: AuthorData{
			Name:             authorData.Name,
			BiographicalInfo: authorData.Description,
			Image:            authorData.YoastHeadJson.OgImage[0].Url,
			AuthorPage: AuthorPage{
				Slug:            slug,
				MetaDescription: authorData.YoastHeadJson.Description,
				Index:           index,
			},
		},
		PublishedStatus: "published",
	}
}

func CreatePostEntry(wpData wp.PostData) BuilderEntry[BlogPostData] {

	var index bool

	if wpData.YoastHeadJson.Robots.Index == "index" {
		index = true
	} else {
		index = false
	}

	var builderImage string

	authorId := getAuthorId(wpData.AuthorId)
	slug := fmt.Sprintf("/%s/", wpData.Slug)
	media := wp.GetWordpressMedia(wpData.FeaturedMediaId)
	if media.Guid.Rendered != "" {
		builderImage = UploadAssetToBuilder(media.Guid.Rendered, BLOG_RESOURCE_ASSET_FOLDER, "")
	} else {
		utils.WriteToErrorLogs(fmt.Sprintf("No Feature Image at: %s", slug))
	}

	topic, subtopic, endSymbolName, _ := utils.ReadCsv("blog-post-topics.csv", wpData.Slug)

	if topic == "" || subtopic == "" {
		err := fmt.Errorf("no entry found at: %s", slug)
		utils.WriteToErrorLogs(err.Error())
	}

	var altText string

	if media.AltText == "" {
		altText = "No Alt Text Found"
	} else {
		altText = media.AltText
	}

	htmlContent := utils.HtmlImgToShortcode(wpData.Content.Rendered, endSymbolName)

	topicDropdownSelection := buildTopicsDropdown(topic, subtopic)

	return BuilderEntry[BlogPostData]{
		Name:        html.UnescapeString(wpData.Title.Rendered),
		LastUpdated: time.Now().UnixMilli(),
		Data: BlogPostData{
			Title: html.UnescapeString(wpData.Title.Rendered),
			AuthorReference: BuilderReference{
				Type:  "@builder.io/core:Reference",
				Id:    authorId,
				Model: MODEL_BLOG_AUTHORS,
			},
			PublishDate: wpData.Date,
			Content:     htmlContent,
			FeaturedImage: BlogPostFeaturedImage{
				Upload:  builderImage,
				AltText: altText,
			},
			TOC: TableOfContents{
				IncludeH2: true,
				IncludeH3: false,
			},
			Category: PostCategory{
				Topic: topicDropdownSelection,
				Type: TypeDropdown{
					TypeSelection: "Article",
				},
			},
			SeoSettings: SeoSettings{
				Slug:            slug,
				MetaTitle:       wpData.YoastHeadJson.Title,
				MetaDescription: wpData.YoastHeadJson.Description,
				Index:           index,
			},
		},
		PublishedStatus: "published",
	}
}

func CreateDownloadEntry(downloadData BuilderEntry[BuilderDownloadData], assetUrl string) BuilderEntry[BuilderDownloadData] {
	return BuilderEntry[BuilderDownloadData]{
		Data: BuilderDownloadData{
			FileName:    downloadData.Data.FileName,
			FileUpload:  assetUrl,
			PreservedId: downloadData.Id,
		},
		Name:            downloadData.Name,
		PublishedStatus: downloadData.PublishedStatus,
	}
}

func CreateMigrationAssetEntry(fileName string, builderAssetUrl string) MigrationAsset {
	return MigrationAsset{
		Data: MigrationAssetData{
			FileName:      fileName,
			BuilderImgUrl: builderAssetUrl,
		},
		Name:            fileName,
		PublishedStatus: "published",
	}
}

func CreateResourceEntry(resource BuilderEntry[BuilderContentData]) BuilderEntry[BuilderResourceData] {
	ref, err := getDownloadByID(resource.Data.Blocks[0].Component.Options.Symbol.Data.Download.Id, resource.Data.URL)

	if err != nil {
		utils.WriteToErrorLogs(err.Error())
	}

	typeName, topic, subtopic, endSymbolName := utils.ReadCsv("resource-post-topics.csv", resource.Data.URL)

	if topic == "" || subtopic == "" || typeName == "" {
		err := fmt.Errorf("no entry found at: %s", resource.Data.URL)
		utils.WriteToErrorLogs(err.Error())
	}

	topicDropdown := buildTopicsDropdown(topic, subtopic)

	content := utils.ScrapeResourceHtml(resource.Data.URL)

	htmlContent := utils.HtmlImgToShortcode(content, endSymbolName)

	imageURL := UploadAssetToBuilder(resource.Data.Image, BLOG_RESOURCE_ASSET_FOLDER, fmt.Sprintf("%s-hero", resource.Name))

	return BuilderEntry[BuilderResourceData]{
		Data: BuilderResourceData{
			Title:       resource.Name,
			PublishDate: resource.CreatedDate,
			Content:     htmlContent,
			Image: BuilderResourceImage{
				Upload:  imageURL,
				AltText: resource.Data.Blocks[0].Component.Options.Symbol.Data.AltText,
			},
			Download: ref,
			ResourcePostToc: TableOfContents{
				IncludeH2: true,
				IncludeH3: false,
			},
			Category: PostCategory{
				Topic: topicDropdown,
				Type: TypeDropdown{
					TypeSelection: typeName,
				},
			},
			SeoSettings: SeoSettings{
				Slug:            strings.ReplaceAll(resource.Data.URL, "/resource", ""),
				MetaTitle:       resource.Data.Title,
				MetaDescription: resource.Data.Description,
				Index:           true,
			},
		},
		Name:            resource.Name,
		LastUpdated: time.Now().UnixMilli(),
		PublishedStatus: "published",
	}
}

func DeleteAllEntries(modelName string) {
	entries := getAllEntries(modelName)

	client := &http.Client{}

	for _, v := range entries.Results {
		id := v.Id
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", BUILDER_WRITE_API_ENDPOINT, modelName, id), nil)

		if err != nil {
			log.Fatal(err)
		}

		addBuilderToken(req)

		resp, err := client.Do(req)

		if err != nil {
			log.Fatalln(err)
		}
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("Deletion finished")
}

func GetDownloads(limit string, offset string) []BuilderEntry[BuilderDownloadData] {
	resp, err := http.Get(fmt.Sprintf(
		"%s/%s?apiKey=8f355b298ee34d67a83efdd69982f6d8&limit=%s&offset=%s",
		BUILDER_CONTENT_API_ENDPOINT,
		MODEL_DOWNLOADS,
		limit,
		offset,
	))

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var downloads BuilderResults[BuilderDownloadData]

	json.Unmarshal(body, &downloads)

	return downloads.Results
}

func GetResourcePosts(limit string, offset string) []BuilderEntry[BuilderContentData] {
	resp, err := http.Get(fmt.Sprintf(
		"%s/page-resource-post?apiKey=8f355b298ee34d67a83efdd69982f6d8&limit=%s&offset=%s",
		BUILDER_CONTENT_API_ENDPOINT,
		limit,
		offset,
	))

	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var resources BuilderResults[BuilderContentData]

	json.Unmarshal(body, &resources)

	return resources.Results
}

func SendAuthorToBuilder(authorData wp.UserData, generateBuilderData convertAuthor) {
	ad := generateBuilderData(authorData)

	bs, err := json.Marshal(ad)

	if err != nil {
		log.Fatal(err)
	}

	SendToBuilder(bs, MODEL_BLOG_AUTHORS, false)
	fmt.Printf("finished uploading user: %s", authorData.Name)
	fmt.Println()
}

func SendDownloadToBuilder(download BuilderEntry[BuilderDownloadData]) {
	assetUrl := UploadAssetToBuilder(download.Data.FileUpload, BLOG_RESOURCE_ASSET_FOLDER, download.Name)
	downloadToSend := CreateDownloadEntry(download, assetUrl)

	bs, err := json.Marshal(downloadToSend)

	if err != nil {
		log.Fatal(err)
	}

	SendToBuilder(bs, MODEL_DOWNLOADS, false)
	msg := fmt.Sprintf("%s Uploaded", downloadToSend.Name)
	fmt.Println(msg)
}

func SendMigrationAssetToBuilder(asset MigrationAsset) {
	as, err := json.Marshal(asset)

	if err != nil {
		log.Fatal(err)
	}

	SendToBuilder(as, MODEL_MIGRATION_ASSETS, true)
}

func SendPostToBuilder(wpData wp.PostData, generateBuilderData convertPost) {
	bd := generateBuilderData(wpData)

	bs, err := json.Marshal(bd)

	if err != nil {
		log.Fatal(err)
	}

	SendToBuilder(bs, MODEL_BLOG_POST, false)
	msg := fmt.Sprintf("Finished uploading post: %s", wpData.Title.Rendered)
	fmt.Println(msg)
}

func SendResourceToBuilder(resource BuilderEntry[BuilderResourceData]) {
	bs, err := json.Marshal(resource)
	if err != nil {
		log.Fatal(err)
	}
	SendToBuilder(bs, MODEL_RESOURCE_POSTS, false)
	msg := fmt.Sprintf("Finished uploading resource: %s", resource.Name)
	fmt.Println(msg)
}

func SendToBuilder(bs []byte, model string, isAsset bool) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", BUILDER_WRITE_API_ENDPOINT, model), bytes.NewBuffer(bs))

	if err != nil {
		log.Fatalln(err)
	}

	if isAsset {
		token := getMainToken()
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")
	} else {
		addBuilderToken(req)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

func UploadAssetToBuilder(asset string, assetFolder string, fileNameBackup string) string {
	resp, err := http.Get(asset)

	if err != nil {
		log.Fatal(err)
	}

	var fileName string

	contentType := resp.Header.Values("Content-Type")[0]

	if len(resp.Header.Values("Content-Disposition")) > 0 {
		fileNameSlice := strings.Split(resp.Header.Values("Content-Disposition")[0], "filename=")
		if len(fileNameSlice) > 1 {
			fileName = strings.Trim(fileNameSlice[1], "\"")
		} else {
			fileName = fmt.Sprintf("%s.%s", fileNameBackup, strings.Split(contentType, "/")[1])
			fileName = strings.ReplaceAll(fileName, " ", "-")
			fileName = strings.ToLower(fileName)
		}
	} else if len(strings.Split(asset, "/")) < 8 {
		fileName = fmt.Sprintf("%s.%s", fileNameBackup, strings.Split(contentType, "/")[1])
		fileName = strings.ReplaceAll(fileName, " ", "-")
		fileName = strings.ToLower(fileName)
	} else {
		fileName = strings.Split(asset, "/")[7]
	}

	fileName = strings.ReplaceAll(fileName, " ", "")

	existingUrl, assetExists := doesAssetExist(fileName)

	if assetExists {
		return existingUrl
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	uploadEndpoint := fmt.Sprintf("%s?name=%s%s", BUILDER_UPLOAD_API_ENDPOINT, fileName, assetFolder)

	req, err := http.NewRequest("POST", uploadEndpoint, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}

	token := getMainToken()
	req.Header.Add("Content-Type", contentType)
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

	migrationAsset := CreateMigrationAssetEntry(fileName, builderImgUrl.ImageUrl)

	SendMigrationAssetToBuilder(migrationAsset)

	return builderImgUrl.ImageUrl
}
