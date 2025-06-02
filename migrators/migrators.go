package migrators

import (
	"sync"

	"simplepractice.com/wp-post-migrator/builder"
	"simplepractice.com/wp-post-migrator/utils"
	wp "simplepractice.com/wp-post-migrator/wordpress"
)

func MigrateAuthors(perFetch int, pageNum int, wg *sync.WaitGroup) {
	defer wg.Done()
	if pageNum == 1 {
		builder.DeleteAllEntries(builder.MODEL_BLOG_AUTHORS)
	}
	userSlice := wp.GetWordpressUsers(perFetch, pageNum)

	for _, v := range userSlice {

		isAuthor := wp.CheckIfIsAuthor(v.Id)

		if isAuthor {
			if len(v.YoastHeadJson.OgImage) > 0 {
				img := v.YoastHeadJson.OgImage[0].Url
				v.YoastHeadJson.OgImage[0].Url = builder.UploadAssetToBuilder(img, builder.BLOG_AUTHOR_IMAGE_FOLDER, "")
			} else {
				yoastData := append(v.YoastHeadJson.OgImage, wp.SeoImage{
					Url: "https://cdn.builder.io/api/v1/image/assets%2F6a01583f2cca435aa4d0f2561c553e0c%2F0c25f9072b2c482982d69195fa10f19b?quality=60&width=200&height=200",
				})
				v.YoastHeadJson.OgImage = yoastData
			}
			builder.SendAuthorToBuilder(v, builder.CreateAuthorEntry)
		}
	}
}

func MigrateDownloads(limit string, offset string, wg *sync.WaitGroup) {
	defer wg.Done()
	downloads := builder.GetDownloads(limit, offset)

	if offset == "0" {
		builder.DeleteAllEntries(builder.MODEL_DOWNLOADS)
	}

	for _, v := range downloads {
		builder.SendDownloadToBuilder(v)
	}

}

func MigratePosts(perFetch int, pageNum int, wg *sync.WaitGroup) {
	defer wg.Done()

	if pageNum == 1 {
		utils.DeleteFile("error.log")
		builder.DeleteAllEntries(builder.MODEL_BLOG_POST)
	}
	postSlice := wp.GetWordpressPosts(perFetch, pageNum)
	for _, v := range postSlice {
		builder.SendPostToBuilder(v, builder.CreatePostEntry)
	}
}

func MigrateResources(limit string, offset string, wg *sync.WaitGroup) {
	defer wg.Done()
	if offset == "0" {
		builder.DeleteAllEntries(builder.MODEL_RESOURCE_POSTS)
	}
	resources := builder.GetResourcePosts(limit, offset)

	for _, v := range resources {
		resourceEntry := builder.CreateResourceEntry(v)
		builder.SendResourceToBuilder(resourceEntry)
	}
}
