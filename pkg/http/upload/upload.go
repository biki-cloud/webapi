package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	msg "webapi/microservices/exec/pkg/msgs"
	"webapi/pkg/int"
	utils2 "webapi/pkg/os"
	utilString "webapi/pkg/string"
)

var FileSizeTooBigError = errors.New("upload file size is too big.")

// UploadHelper はファイルをアップロードするためのハンドラー。
// サーバの中で動作する。
// w,rとサーバ内のアップロードさせるディレクトリを指定することでそのディレクトリに保存させる。
func UploadHelper(w http.ResponseWriter, r *http.Request, uploadDir string, maxUploadSizeByte int64) (string, error) {
	//maxUploadSize := int64(int2.MBToByte(int(cfg.MaxUploadSizeMB)))
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("UploadHelper: %v ", errors.New(r.Method+" is not allowed."))
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSizeByte)
	if err := r.ParseMultipartForm(maxUploadSizeByte); err != nil {
		return "", fmt.Errorf("%w, file size: %v", FileSizeTooBigError, msg.UploadFileSizeExceedError(int.ByteToMB(maxUploadSizeByte)))
	}

	//FormFileの引数はHTML内のform要素のnameと一致している必要があります
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("UploadHelper: %v", err)
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
		return
	}(file)

	// 存在していなければ、保存用のディレクトリを作成します。
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("UploadHelper: %v", err)
	}

	// 保存用ディレクトリ内に新しいファイルを作成します。
	// アップロードファイルに半角や全角のスペースがある場合は削除する。
	spaceRemovedUploadFileName := utilString.RemoveSpace(fileHeader.Filename)
	uploadFilePath := filepath.Join(uploadDir, spaceRemovedUploadFileName)
	dst, err := os.Create(uploadFilePath)
	if err != nil {
		return "", fmt.Errorf("UploadHelper: %v", err)
	}

	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			return
		}
		return
	}(dst)

	// アップロードされたファイルを先程作ったファイルにコピーします。
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("UploadHelper: %w", err)
	}

	return uploadFilePath, nil

}

type Uploader interface {
	// Upload アップロードするファイルをURLを受け取り、アップロードする。
	Upload(url string, uploadFilePath string) error
}

func NewUploader() Uploader {
	return &uploader{}
}

type uploader struct{}

// Upload クライアントサイドのGoプログラムで動作する。
// 指定したサーバURLに指定したファイルをアップロードする。
func (u *uploader) Upload(url string, uploadFilePath string) error {
	command := fmt.Sprintf("curl -X POST -F file=@%v %v", uploadFilePath, url)
	stdout, stderr, err := utils2.SimpleExec(command)
	if strings.Contains(stdout, "request body too large") || err != nil {
		return fmt.Errorf("Upload: stdout: %v \n stderr: %v ", stdout, stderr)
	}
	return nil
}
