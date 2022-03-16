package env

import (
	"log"
	"os"
	"strconv"
	os2 "webapi/pkg/os"
	"webapi/pkg/random"
)

// FileServerConf is relevant fileserver information struct.
type FileServerConf struct {
	Dir       string `json:"dir"`
	UploadDir string `json:"uploadDir"`
	WorkDir   string `json:"workDir"`
}

// Env has all config that are LogConf, FileServerConf,,etc.
type Env struct {
	ProgramsDir       string
	ProgramsJSON      string
	FileServer        FileServerConf
	WorkedDirKeepSec  int
	ProgramServerIP   string
	ProgramServerPort string
	ExecuteTimeoutSec int
	StdoutBufferSize  int
	StderrBufferSize  int
	MaxUploadSizeMB   int64
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
		li := os2.ListEnvToSlice(os.Getenv("K8S_WORKER_NODE_IPS"))
		e.ProgramServerIP = random.Choice(li)
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
		// for local
		// set environment variable here, if it is not setting
		m := make(map[string]string)
		m["FILESERVER_DIRNAME"] = "fileserver"
		m["UPLOAD_DIRNAME"] = "upload"
		m["WORK_DIRNAME"] = "work"
		m["PROGRAMS_DIRNAME"] = "programs"
		m["WORKED_DIR_KEEP_SEC"] = "600"
		m["MY_IP"] = os2.GetLocalIP()
		m["EXECUTE_TIMEOUT_SEC"] = "10"
		m["STDOUT_BUFFER_SIZE"] = "1000000"
		m["STDERR_BUFFER_SIZE"] = "1000000"
		m["MAX_UPLOAD_SIZE_MB"] = "300"

		for k, v := range m {
			err := os2.SetEnvIfNotExists(k, v)
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
