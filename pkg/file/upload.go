package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"

	storage_go "github.com/supabase-community/storage-go"
)

var (
	bucketId  string = "gotest-images"
	projectId string = "jaurcadmqgccmhzllhxr"
	rawUrl    string = fmt.Sprintf("https://%s.supabase.co/storage/v1", projectId)
)

const MaxUploadSize = 1024 * 1024

func UploadHandler(c echo.Context) error {
	var apiKey string = os.Getenv("SUPABASE_API_KEY")
	storageClient := storage_go.NewClient(rawUrl, apiKey, nil)

	// ファイルの読み込み
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("error")
		errJson := map[string]string{
			"statusCode":   strconv.Itoa(http.StatusInternalServerError),
			"errorMessage": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errJson)
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// ファイルのmimetypeの取得
	// ファイルを512バイト分のみ読み込み
	buff := make([]byte, 512)
	_, err = src.Read(buff)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// 読み込んだバッファからmimetypeを推定
	mimeType := http.DetectContentType(buff)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		return c.JSON(http.StatusInternalServerError, "許可されていないファイルタイプです。JPEGかPNGをアップロードしてください")
	}
	// 読み込んだバッファ分を戻す
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fileOptions := storage_go.FileOptions{ContentType: &mimeType}
	_, err = storageClient.UploadFile(bucketId, file.Filename, src, fileOptions)
	if err != nil {
		errJson := map[string]string{
			"statusCode":   strconv.Itoa(http.StatusInternalServerError),
			"errorMessage": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errJson)
	}

	return c.String(http.StatusOK, "uploaded success")
}
