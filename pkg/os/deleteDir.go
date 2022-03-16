package os

import (
	"fmt"
	"os"
	"time"
)

// DeleteDirSomeTimeLater は一定時間後にディレクトリを削除する
func DeleteDirSomeTimeLater(dirPath string, seconds int) error {
	// wait some seconds
	time.Sleep(time.Second * time.Duration(seconds))
	err := os.RemoveAll(dirPath)
	if err != nil {
		return fmt.Errorf("DeleteDirSomeTimeLater: %v", err)
	}
	return nil
}
