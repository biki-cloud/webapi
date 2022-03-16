package executer_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"webapi/microservices/exec/env"
	executer2 "webapi/microservices/exec/pkg/execution/executer"
	"webapi/microservices/exec/pkg/execution/outputManager"
	"webapi/microservices/exec/pkg/msgs"
	"webapi/test"
)

var (
	executer executer2.Executer
	cfg      *env.Env
)

func init() {
	cfg = env.New()
	executer = executer2.New()
}

func TestFileExecuter_Execute(t *testing.T) {
	// コンテキストを用意する→実行→出力を判定、テストする

	tests := test.GetMaterials()

	for _, tt := range tests {
		if tt.IsSkip {
			continue
		}

		ctx, err := tt.Setup()
		if err != nil {
			// 帰ってきたエラーが期待しているものと違ったら
			if !errors.Is(err, tt.ExpectedError) {
				t.Errorf("test name: %v, got: %v, want: %v", tt.TestName, err, tt.ExpectedError)
				os.Remove(tt.UploadFilePath)
				continue
			} else {
				// 期待しているエラーがきた場合はここにくる
				// 準備していたアップロードファイルを削除し、このループは終了する。
				os.Remove(tt.UploadFilePath)
				continue
			}
		}

		out := executer.Execute(ctx)
		t.Run(tt.TestName, func(j *testing.T) {
			testExecute(t, out, tt, cfg)
		})
	}
}

func testExecute(t *testing.T, out outputManager.OutputManager, tt test.Struct, cfg *env.Env) {
	t.Helper()

	if len(tt.ExpectedOutFileNames) > 0 {
		for _, ExpectedOutFileName := range tt.ExpectedOutFileNames {
			isExists := false
			for _, outPath := range out.OutURLs() {
				if filepath.Base(outPath) == ExpectedOutFileName {
					isExists = true
				}
			}
			if !isExists {
				t.Errorf("test name: %v,\"%v\" is not exists of %v", tt.TestName, ExpectedOutFileName, out.OutURLs())
			}
		}
	}

	if (out.StdOut() == "") != tt.ExpectedStdOutIsEmpty {
		t.Errorf("test name: %v, out.StdOut() is empty. stdout: %v", tt.TestName, out.StdOut())
	}
	if (out.StdErr() == "") != tt.ExpectedStdErrIsEmpty {
		t.Errorf("test name: %v, out.StdErr() is not empty. stdout: %v", tt.TestName, out.StdErr())
	}
	if out.Status() != tt.ExpectedStatus {
		t.Errorf("test name: %v, out.ExpectedStatus() is not %v. got: %v", tt.TestName, msgs.OK, out.Status())
	}
	if (out.ErrorMsg() == "") != tt.ExpectedErrMsgIsEmpty {
		t.Errorf("test name: %v, out.ErrorMsg is not empty. got: %v", tt.TestName, out.ErrorMsg())
	}

	// out.Stdout,errはcfg（設定ファイル）の値より小さくなくてはならない。設定値がマックスなので。
	if len(out.StdOut()) > cfg.StdoutBufferSize {
		t.Errorf("test name: %v, len(out.StdOut()):%v is not more less cfg.StdoutBufferSize: %v \n", tt.TestName, len(out.StdOut()), cfg.StdoutBufferSize)
	}
	if len(out.StdErr()) > cfg.StderrBufferSize {
		t.Errorf("test name: %v, len(out.StdErr()):%v is not more less cfg.StderrBufferSize: %v. ", tt.TestName, len(out.StdErr()), cfg.StderrBufferSize)
	}

	t.Cleanup(func() {
		os.RemoveAll("fileserver")
		os.Remove(tt.UploadFilePath)
	})
}
