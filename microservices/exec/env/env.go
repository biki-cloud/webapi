package env

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"webapi/pkg/k8s"
	pkgOs "webapi/pkg/os"
)

// FileServerConf is relevant fileserver information struct.
type FileServerConf struct {
	// コンテンツ等を格納するファイルサーバのディレクトリ名
	Dir string
	// ファイルサーバディレクトリの中に入力ファイルをアップロードさせるためのディレクトリ名
	UploadDir string
	// ファイルサーバディレクトリの中に処理ディレクトリなどを一時的に作成し、動作させるが
	// その一時ディレクトリのルートディレクトリの名前
	WorkDir string
}

// Env has all config that are LogConf, FileServerConf,,etc.
type Env struct {
	// 登録プログラムを入れるためのディレクトリ名
	ProgramsDir string
	// コンテンツ等を格納するファイルサーバ等の情報を保持した構造体
	FileServer FileServerConf
	// プログラム実行が終了したディレクトリを何秒保持するか
	WorkedDirKeepSec int
	// このサーバを立てるIP
	ProgramServerIP string
	// このサーバを立てるポート
	ProgramServerPort string
	// プログラム実行が始まって何秒経過したらタイムアウトとするか
	ExecuteTimeoutSec int
	// プログラムの標準出力を何バイト出力するか
	StdoutBufferSize int
	// プログラムの標準エラー出力を何バイト出力するか
	StderrBufferSize int
	// 入力ファイルサイズの上限を何メガにするか
	MaxUploadSizeMB int64
}

// Print 設定値を表示する。
func Print(w io.Writer) {
	e := New()
	fmt.Fprintf(w, "FILESERVER_DIRNAME : %v \n", e.FileServer.Dir)
	fmt.Fprintf(w, "UPLOAD_DIRNAME     : %v \n", e.FileServer.UploadDir)
	fmt.Fprintf(w, "WORK_DIRNAME       : %v \n", e.FileServer.WorkDir)
	fmt.Fprintf(w, "PROGRAMS_DIRNAME   : %v \n", e.ProgramsDir)
	fmt.Fprintf(w, "WORKED_DIR_KEEP_SEC: %v \n", e.WorkedDirKeepSec)
	fmt.Fprintf(w, "MY_IP              : %v \n", e.ProgramServerIP)
	fmt.Fprintf(w, "MY_PORT            : %v \n", e.ProgramServerPort)
	fmt.Fprintf(w, "EXECUTE_TIMEOUT_SEC: %v \n", e.ExecuteTimeoutSec)
	fmt.Fprintf(w, "STDOUT_BUFFER_SIZE : %v \n", e.StdoutBufferSize)
	fmt.Fprintf(w, "STDERR_BUFFER_SIZE : %v \n", e.StderrBufferSize)
	fmt.Fprintf(w, "MAX_UPLOAD_SIZE_MB : %v \n", e.MaxUploadSizeMB)
}

// New はconfig.jsonの中身をstructに入れたものを返す
func New() *Env {
	// 構造体を初期化
	e := new(Env)

	if os.Getenv("ENV") == "k8s" {
		// for k8s
		e.FileServer.Dir = os.Getenv("FILESERVER_DIRNAME")
		e.FileServer.UploadDir = os.Getenv("UPLOAD_DIRNAME")
		e.FileServer.WorkDir = os.Getenv("WORK_DIRNAME")

		e.ProgramsDir = os.Getenv("PROGRAMS_DIRNAME")

		n, err := strconv.Atoi(os.Getenv("WORKED_DIR_KEEP_SEC"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}

		e.WorkedDirKeepSec = n

		// K8S_WORKER_NODE_IPSからランダムで設定することによりNODEPORTサービスを利用しているので
		// outURLsにURLをセットするときにどれかのノードポートにアクセスすればダウンロードできる
		li := pkgOs.ListEnvToSlice(os.Getenv("K8S_WORKER_NODE_IPS"))
		e.ProgramServerIP, err = k8s.LoadBalance(li, os.Getenv("K8S_WEBSITE_NODEPORT_PORT"))
		if err != nil {
			log.Fatalf("Env.New(): %v \n", err.Error())
		}
		e.ProgramServerPort = os.Getenv("EXEC_NODEPORT_PORT")

		n, err = strconv.Atoi(os.Getenv("EXECUTE_TIMEOUT_SEC"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}
		e.ExecuteTimeoutSec = n

		n, err = strconv.Atoi(os.Getenv("STDOUT_BUFFER_SIZE"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}
		e.StdoutBufferSize = n

		n, err = strconv.Atoi(os.Getenv("STDERR_BUFFER_SIZE"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}
		e.StderrBufferSize = n

		n, err = strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE_MB"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}
		e.MaxUploadSizeMB = int64(n)

	} else {
		// ローカルでの動作の場合はこの環境変数を設定を使用する
		// 以下の設定は環境変数がセットされていなければセットし、
		// セットされていれば何もしない
		m := make(map[string]string)
		m["FILESERVER_DIRNAME"] = "fileserver"
		m["UPLOAD_DIRNAME"] = "upload"
		m["WORK_DIRNAME"] = "work"
		m["PROGRAMS_DIRNAME"] = "programs"
		m["WORKED_DIR_KEEP_SEC"] = "600"
		m["MY_IP"] = pkgOs.GetLocalIP()
		m["EXECUTE_TIMEOUT_SEC"] = "100"
		m["STDOUT_BUFFER_SIZE"] = "1000000"
		m["STDERR_BUFFER_SIZE"] = "1000000"
		m["MAX_UPLOAD_SIZE_MB"] = "300"

		for k, v := range m {
			err := pkgOs.SetEnvIfNotExists(k, v)
			if err != nil {
				log.Fatalf("New: %v", err.Error())
			}
		}

		e.FileServer.Dir = os.Getenv("FILESERVER_DIRNAME")
		e.FileServer.UploadDir = os.Getenv("UPLOAD_DIRNAME")
		e.FileServer.WorkDir = os.Getenv("WORK_DIRNAME")

		e.ProgramsDir = os.Getenv("PROGRAMS_DIRNAME")

		n, err := strconv.Atoi(os.Getenv("WORKED_DIR_KEEP_SEC"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}

		e.WorkedDirKeepSec = n
		e.ProgramServerIP = os.Getenv("MY_IP")
		e.ProgramServerPort = os.Getenv("MY_PORT")

		n, err = strconv.Atoi(os.Getenv("EXECUTE_TIMEOUT_SEC"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}
		e.ExecuteTimeoutSec = n

		n, err = strconv.Atoi(os.Getenv("STDOUT_BUFFER_SIZE"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}
		e.StdoutBufferSize = n

		n, err = strconv.Atoi(os.Getenv("STDERR_BUFFER_SIZE"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}
		e.StderrBufferSize = n

		n, err = strconv.Atoi(os.Getenv("MAX_UPLOAD_SIZE_MB"))
		if err != nil {
			log.Fatalf("New: %v \n", err.Error())
		}
		e.MaxUploadSizeMB = int64(n)

	}

	return e
}
