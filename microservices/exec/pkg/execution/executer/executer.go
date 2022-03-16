/*
executer.go
ContextManagerの値をベースにコマンドを実行する。
コマンドを実行した標準出力、出力ファイルなどはOutputManagerにセットする。
実行した後に使用したディレクトリは一定時間後、削除する。
*/

package executer

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"webapi/microservices/exec/pkg/execution/contextManager"
	"webapi/microservices/exec/pkg/execution/outputManager"
	"webapi/microservices/exec/pkg/msgs"
	pkgHttp "webapi/pkg/http/url"
	pkgOs "webapi/pkg/os"
)

// Executer はコマンドを実行する構造体のインタフェース
// 返り値はOutputManager(インターフェース) エラーもoutputManagerの中に入れる、
type Executer interface {
	Execute(contextManager.ContextManager) outputManager.OutputManager
}

// New はfileExecuter構造体を返す。
func New() Executer {
	return &executer{}
}

type executer struct{}

// errorOutWrap は中の３行は頻繁に使用するので行数削減と見やすくするため
// OutputManagerの中にセットする
func errorOutWrap(out outputManager.OutputManager, err error, status string) outputManager.OutputManager {
	out.SetErrorMsg(err.Error())
	out.SetStatus(status)
	return out
}

func (f *executer) Execute(ctx contextManager.ContextManager) (out outputManager.OutputManager) {

	// コマンド実行
	out = Exec(ctx.Command(), ctx.Env().ExecuteTimeoutSec, ctx.Env().StdoutBufferSize, ctx.Env().StderrBufferSize)
	if out.Status() != msgs.OK {
		return
	}

	// 出力ファイルたちはまだ通常のパスなのでそれを
	// CURLで取得するためにURLパスに変換する。
	outFileURLs, err := GetOutFileURLs(ctx.OutputDir(), ctx.Env().ProgramServerIP, ctx.Env().ProgramServerPort, ctx.Env().FileServer.Dir)
	if err != nil {
		out.SetStatus(msgs.SERVERERROR)
		out.SetErrorMsg(err.Error())
		return
	}

	// 時間経過後ファイルを削除
	go func() {
		err := pkgOs.DeleteDirSomeTimeLater(ctx.ProgramTempDir(), ctx.Env().WorkedDirKeepSec)
		if err != nil {
			fmt.Printf("Execute: %v \n", err)
		}
	}()

	out.SetOutURLs(outFileURLs)

	return
}

// Exec は実行するためのコマンド, 時間制限をもらい、OutputManagerインタフェースを返す
// エラーがでた場合もoutputInfoのエラーメッセージの中に格納する。
func Exec(command string, timeOut int, stdOutBufferSize, stdErrBufferSize int) (out outputManager.OutputManager) {
	out = outputManager.New()

	var timeoutError *pkgOs.TimeoutError
	stdout, stderr, err := pkgOs.ExecuteWithTimeout(command, timeOut)

	if err1 := out.SetStdOut(&stdout, stdOutBufferSize); err1 != nil {
		out.SetStatus(msgs.SERVERERROR)
		out.SetErrorMsg(fmt.Sprintf("err: %v, err1: %v", err.Error(), err1.Error()))
		return
	}

	if err2 := out.SetStdErr(&stderr, stdErrBufferSize); err2 != nil {
		out.SetStatus(msgs.SERVERERROR)
		out.SetErrorMsg(fmt.Sprintf("err: %v, err2: %v", err.Error(), err2.Error()))
		return
	}

	if err != nil {
		if errors.As(err, &timeoutError) {
			// プログラムがタイムアウトした場合
			out.SetStatus(msgs.PROGRAMTIMEOUT)
			out.SetErrorMsg(err.Error())
			return
		} else {
			// プログラムがエラーで終了した場合
			out.SetStatus(msgs.PROGRAMERROR)
			out.SetErrorMsg(err.Error())
		}
	} else {
		// 正常終了した場合
		out.SetStatus(msgs.OK)
	}

	return
}

// GetOutFileURLs はコマンドを実行した後に使用する。
// プログラム出力ディレクトリの全てのファイルを取得するURLのリストを返す。
func GetOutFileURLs(outputDir string, serverIP, serverPort, fileServerDir string) ([]string, error) {
	// 出力されたディレクトリの複数ファイルをglobで取得
	pattern := outputDir + "/*"
	outFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("GetOutFileURLs: %v ", err)
	}

	outFileURLs := make([]string, 0, 20)
	for _, outfile := range outFiles {
		outFileURL, err := pkgHttp.GetURLFromFilePath(outfile, serverIP, serverPort, fileServerDir)
		if err != nil {
			return nil, fmt.Errorf("GetOutFileURLs: %v", err)
		}
		// サーバがwindowsだった場合、出力パス区切りを¥から/に変更する。
		outFileURL = strings.Replace(outFileURL, "¥", "/", -1)
		outFileURLs = append(outFileURLs, outFileURL)
	}

	return outFileURLs, nil
}
