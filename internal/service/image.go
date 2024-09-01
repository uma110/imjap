package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/imjap/internal/model"
	"github.com/labstack/echo/v4"

	storage_go "github.com/supabase-community/storage-go"
)

var (
	bucketId  string = "images"
	projectId string = "pxoevrjzbjwbcqmbhrcf"
	rawUrl    string = fmt.Sprintf("https://%s.supabase.co/storage/v1", projectId)
)

const MaxUploadSize = 1024 * 1024

type ImageService struct{}

func (is *ImageService) UploadFile(c echo.Context) (bool, error) {
	var apiKey string = os.Getenv("SUPABASE_API_KEY")
	storageClient := storage_go.NewClient(rawUrl, apiKey, nil)

	// ファイルの読み込み
	file, err := c.FormFile("file")
	if err != nil {
		return false, err
	}
	src, err := file.Open()
	if err != nil {
		return false, err
	}
	defer src.Close()

	// ファイルのmimetypeの取得
	// ファイルを512バイト分のみ読み込み
	buff := make([]byte, 512)
	_, err = src.Read(buff)
	if err != nil {
		return false, err
	}
	// 読み込んだバッファからmimetypeを推定
	mimetype := http.DetectContentType(buff)
	if mimetype != "image/jpeg" && mimetype != "image/png" {
		return false, err
	}
	// 読み込んだバッファ分を戻す
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}

	uid, _ := uuid.NewRandom()
	now := time.Now().Format("2006010215405")
	filename := now + "-" + uid.String()
	fileOptions := storage_go.FileOptions{ContentType: &mimetype}
	_, err = storageClient.UploadFile(bucketId, filename, src, fileOptions)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (is *ImageService) GetFiles() ([]model.Image, error) {
	var apiKey string = os.Getenv("SUPABASE_API_KEY")
	storageClient := storage_go.NewClient(rawUrl, apiKey, nil)

	files, err := storageClient.ListFiles(bucketId, "", storage_go.FileSearchOptions{
		Limit:  10,
		Offset: 0,
	})

	if err != nil {
		return nil, err
	}

	resFiles := make([]model.Image, 0, len(files))
	for _, v := range files {
		var file model.Image = model.Image{
			ID:        v.Id,
			Name:      v.Name,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		resFiles = append(resFiles, file)
	}

	return resFiles, nil
}

func (is *ImageService) GetFile(filename string) ([]byte, string, error) {
	var apiKey string = os.Getenv("SUPABASE_API_KEY")
	storageClient := storage_go.NewClient(rawUrl, apiKey, nil)

	// ファイルのダウンロード
	data, err := storageClient.DownloadFile(bucketId, filename)

	if err != nil {
		return nil, "", err
	}

	// ファイルのmimetypeの取得
	// 読み込んだバッファからmimetypeを推定
	mimetype := http.DetectContentType(data[:512])
	return data, mimetype, nil
}
