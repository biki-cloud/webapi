/*
ユーザにレスポンスするメッセージを定義したパッケージ
ここにわかりやすく書いておいて、他で使用する。
*/

package msgs

import (
	"fmt"
)

var (
	NotAllowMethodError = "許可されていないメソッドです."
	UploadSuccess       = "アップロードが成功しました。"
	UploadSizeIsTooBig  = "アップロードされたファイルが大きすぎます"
)

func UploadFileSizeExceedError(maxUploadFileSize int64) string {
	return fmt.Sprintf("%v %vMB以下のファイルを指定してください", UploadSizeIsTooBig, maxUploadFileSize)
}
