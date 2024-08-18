package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/imjap/pkg/response"
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

	uid, _ := uuid.NewRandom()
	now := time.Now().Format("2006010215405")
	filename := now + "-" + uid.String()
	fileOptions := storage_go.FileOptions{ContentType: &mimeType}
	_, err = storageClient.UploadFile(bucketId, filename, src, fileOptions)
	if err != nil {
		errJson := map[string]string{
			"statusCode":   strconv.Itoa(http.StatusInternalServerError),
			"errorMessage": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errJson)
	}

	return c.String(http.StatusOK, "uploaded success")
}

func GetFilesHandler(c echo.Context) error {
	var apiKey string = os.Getenv("SUPABASE_API_KEY")
	storageClient := storage_go.NewClient(rawUrl, apiKey, nil)

	files, err := storageClient.ListFiles(bucketId, "", storage_go.FileSearchOptions{
		Limit:  10,
		Offset: 0,
	})

	if err != nil {
		errJson := map[string]string{
			"statusCode":   strconv.Itoa(http.StatusInternalServerError),
			"errorMessage": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errJson)
	}

	resFiles := make([]response.File, 0, len(files))
	for _, v := range files {
		var file response.File = response.File{
			ID:        v.Id,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		resFiles = append(resFiles, file)
	}

	return c.JSON(http.StatusOK, resFiles)
}
