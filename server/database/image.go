package database

// Import Cloudinary and other necessary libraries
//===================
import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func Credentials() (*cloudinary.Cloudinary, context.Context) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, _ := cloudinary.NewFromParams(GetItemFromENV("CLOUDINARY_CLOUDNAME"), GetItemFromENV("CLOUDINARY_APIKEY"), GetItemFromENV("CLOUDINARY_APISECRET"))
	// cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx
}

func UploadImage(cld *cloudinary.Cloudinary, ctx context.Context, image_pId string, image_url string) string {

	// Upload the image.
	// Set the asset's public ID and allow overwriting the asset with new versions
	resp, err := cld.Upload.Upload(ctx, image_url, uploader.UploadParams{
		PublicID:       image_pId, //"quickstart_butterfly"
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true)})
	if err != nil {
		fmt.Println("error")
		return ""
	}

	// Log the delivery URL
	fmt.Println("****2. Upload an image****\nDelivery URL:", resp.SecureURL)
	return resp.SecureURL
}

func GetAssetInfo(cld *cloudinary.Cloudinary, ctx context.Context, image_public_id string) string {
	resp, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: image_public_id})
	if err != nil {
		fmt.Println("error")
		return ""
	}
	fmt.Println("****3. Get and use details of the image****\nDetailed response:\n", resp)
	return resp.SecureURL

	// Assign tags to the uploaded image based on its width. Save the response to the update in the variable 'update_resp'.
	// if resp.Width > 900 {
	// 	update_resp, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
	// 		PublicID: "quickstart_butterfly",
	// 		Tags:     []string{"large"}})
	// 	if err != nil {
	// 		fmt.Println("error")
	// 	} else {
	// 		// Log the new tag to the console.
	// 		fmt.Println("New tag: ", update_resp.Tags, "\n")
	// 	}
	// } else {
	// 	update_resp, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
	// 		PublicID: "quickstart_butterfly",
	// 		Tags:     []string{"small"}})
	// 	if err != nil {
	// 		fmt.Println("error")
	// 	} else {
	// 		// Log the new tag to the console.
	// 		fmt.Println("New tag: ", update_resp.Tags)
	// 	}
	// }
}

func TransformImage(cld *cloudinary.Cloudinary, ctx context.Context, image_public_id string) {
	// Instantiate an object for the asset with public ID "my_image"
	img, err := cld.Image(image_public_id) //"quickstart_butterfly"
	if err != nil {
		fmt.Println("error")
	}

	// Generate and log the delivery URL
	url, err := img.String()
	if err != nil {
		fmt.Println("error")
	} else {
		print("****4. Transform the image****\nTransfrmation URL: ", url)
	}
}
