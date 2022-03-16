/*
contextManager.go
入力ファイルや出力ディレクトリ、パラメータなど登録プログラムの実行に必要な
情報等を保持、用意し、executer.Execute()に渡し、実行してもらう。
本来ならばたくさんのパラメータをExecute()に渡さなければいけないがその
たくさんのパラメータを全てこのパッケージのContextManagerインタフェースが保持する。
*/

package contextManager

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	execConfig "webapi/microservices/exec/config"
	execEnv "webapi/microservices/exec/env"
	pkgHttp "webapi/pkg/http/upload"
	pkgInt "webapi/pkg/int"
	pkgOs "webapi/pkg/os"
	pkgString "webapi/pkg/string"
)

// ContextManager はコマンド実行に必要な要素を持つ構造体のインタフェース
// コマンド実行に必要なパラメータ等を一括で管理する構造体のインタフェース
// コマンド実行に必要なものは全てこの中に入れる。
type ContextManager interface {
	// InputFilePath 登録プログラムに処理させる入力ファイルを返す。
	InputFilePath() string
	// SetInputFilePath move fileserver/upload/a.txt fileserver/programOut/(program name)/(random str)/a.txt
	SetInputFilePath() error

	// Command 登録してあるコマンドを返す。
	Command() string
	// SetCommand programConfigHolderを受け取り、登録してあるコマンドを設定する。
	SetCommand() error

	// ProgramName 実行するプログラム名を返す
	ProgramName() string
	// SetProgramName 実行するプログラム名を設定する。
	SetProgramName(string)

	// ProgramTempDir プログラムが作業するテンポラリディレクトリ
	ProgramTempDir() string
	// SetProgramTempDir プログラムが作業するテンポラリディレクトリを取得する
	SetProgramTempDir(proConf execConfig.ProgramConfigHolder, cfg *execEnv.Env) error

	// OutputDir プログラムがファイルを出力するディレクトリを返す。
	OutputDir() string
	// SetOutputDir プログラムがファイルを出力するディレクトリを設定する。
	SetOutputDir(string)

	// UploadedFilePath アップロードされたファイルパスを返す。
	UploadedFilePath() string
	// SetUploadedFilePath アップロードされたファイルパスを設定する
	SetUploadedFilePath(string) // webからの場合に使用する

	// Parameta コマンド実行する際に使用するパラメータを返す。
	Parameta() string
	// SetParameta コマンド実行する際に使用するパラメータを設定する。
	SetParameta(string)

	// ProgramConfig 登録プログラムの情報インターフェースを返す。
	ProgramConfig() execConfig.ProgramConfigHolder
	// SetProgramConfig 登録プログラムの情報インターフェースを設定する。
	SetProgramConfig(holder execConfig.ProgramConfigHolder)

	// Env サーバの設定値を保持する
	Env() *execEnv.Env
	SetEnv(cfg *execEnv.Env)
}

// New はcontextManagerを返す。
// プログラム名とプログラム出力ディレクトリはセットする。
// それ以外に必要な要素は定義した後で設定し、executerに渡す感じ
func New(w http.ResponseWriter, r *http.Request, cfg *execEnv.Env) (ContextManager, error) {
	fName := "New"
	ctx := &contextManager{}

	// file(multi-data)をこのサーバのfileserver/upload/ランダム文字列(被らないように)にアップロードする。
	uploadDir := filepath.Join(cfg.FileServer.Dir, cfg.FileServer.UploadDir, pkgString.GetRandomString(20))
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		return nil, err
	}
	// 時間経過後ファイルを削除
	go func() {
		err := pkgOs.DeleteDirSomeTimeLater(uploadDir, 20)
		if err != nil {
			l := log.New(os.Stdout, "ERR", log.LstdFlags|log.Lshortfile)
			l.Printf("ERR: %v \n", err)
		}
	}()

	maxUploadSize := int64(pkgInt.MBToByte(int(cfg.MaxUploadSizeMB)))
	uploadFilePath, err := pkgHttp.UploadHelper(w, r, uploadDir, maxUploadSize)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", fName, err)
	}

	if !pkgOs.FileExists(uploadFilePath) {
		return nil, fmt.Errorf("%v: %v", fName, errors.New(uploadFilePath+" is not found."))
	}
	ctx.SetUploadedFilePath(uploadFilePath)

	ctx.SetParameta(r.FormValue("parameta"))
	programName := r.FormValue("proName")

	// プログラム名を取得する処理
	// webからの場合はFormValueにプログラム名が乗ってくるが
	// cliからの場合はURLのエンドポイントに乗ってくるため処理を分ける必要がある。
	if programName == "" {
		programName = r.URL.Path[len("/exec/"):]
	}
	ctx.SetProgramName(programName)

	proConf, err := execConfig.GetProConfByName(programName)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", fName, err)
	}
	ctx.SetProgramConfig(proConf)

	err = ctx.SetProgramTempDir(proConf, cfg)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}

	err = ctx.SetProgramOutDir()
	if err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}

	if err = ctx.SetInputFilePath(); err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}
	if err = ctx.SetCommand(); err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}

	ctx.SetEnv(cfg)

	return ctx, nil
}

type contextManager struct {
	programName      string
	programTempDir   string
	outputDir        string
	parameta         string
	uploadedFilePath string
	inputFilePath    string
	command          string
	programConfig    execConfig.ProgramConfigHolder
	stdOutBufferSize int
	stdErrBufferSize int
	env              *execEnv.Env
}

// SetProgramOutDir はプログラムが出力するディレクトリを作成、セットする。
// SetProgramTempDirの後に実行されなければいけない
func (c *contextManager) SetProgramOutDir() error {
	if c.programTempDir == "" {
		return fmt.Errorf("c.programTempDir is not set. you shoud SetProgrtamTempDir before me.")
	}
	randomName := pkgString.GetRandomString(20) + "_out"
	programOutDir := filepath.Join(c.programTempDir, randomName)
	c.outputDir = programOutDir
	err := os.MkdirAll(programOutDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("SetProgramOutDir: %v", err)
	}
	return nil
}

// SetProgramTempDir 登録プログラムが使用する作業ディレクトリとして使用する日付とランダム文字列から作成されたテンポラリディレクトリを作成、セットする
func (c *contextManager) SetProgramTempDir(proConf execConfig.ProgramConfigHolder, cfg *execEnv.Env) error {
	outDirName := pkgString.GetNowTimeStringWithHyphen() + "-" + pkgString.GetRandomString(20)
	programDir := filepath.Join(cfg.FileServer.Dir, cfg.FileServer.WorkDir, proConf.Name(), outDirName)
	c.programTempDir = programDir
	err := os.MkdirAll(programDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("SetProgramTempDir: %v", err)
	}
	return nil
}

func (c *contextManager) InputFilePath() string { return c.inputFilePath }

func (c *contextManager) Command() string { return c.command }

func (c *contextManager) ProgramName() string     { return c.programName }
func (c *contextManager) SetProgramName(s string) { c.programName = s }

func (c *contextManager) ProgramTempDir() string { return c.programTempDir }

func (c *contextManager) OutputDir() string     { return c.outputDir }
func (c *contextManager) SetOutputDir(s string) { c.outputDir = s }

func (c *contextManager) Parameta() string     { return c.parameta }
func (c *contextManager) SetParameta(s string) { c.parameta = s }

func (c *contextManager) UploadedFilePath() string     { return c.uploadedFilePath }
func (c *contextManager) SetUploadedFilePath(s string) { c.uploadedFilePath = s }

func (c *contextManager) ProgramConfig() execConfig.ProgramConfigHolder      { return c.programConfig }
func (c *contextManager) SetProgramConfig(pc execConfig.ProgramConfigHolder) { c.programConfig = pc }

func (c *contextManager) Env() *execEnv.Env       { return c.env }
func (c *contextManager) SetEnv(cfg *execEnv.Env) { c.env = cfg }

// SetInputFilePath uploadディレクトリの中のファイルをfileserver/programOut/convertToJson/xxxxxx/の中に移動させる。
// move fileserver/upload/a.txt fileserver/programOut/(program name)/(random str)/a.txt
func (c *contextManager) SetInputFilePath() error {
	inputFilePath := filepath.Join(filepath.Dir(c.OutputDir()), filepath.Base(c.UploadedFilePath()))
	if err := os.Rename(c.UploadedFilePath(), inputFilePath); err != nil {
		return fmt.Errorf("SetInputFilePath: %v", err)
	}
	c.inputFilePath = inputFilePath
	if !pkgOs.FileExists(c.inputFilePath) {
		return fmt.Errorf("SetInputFilePath: %v", errors.New(c.inputFilePath+"is not found."))
	}
	return nil
}

// SetCommand templateコマンドからINPUTFILE,OUTPUTDIR, PARAMETAなどをreplaceして正規コマンドを作成する。
func (c *contextManager) SetCommand() error {
	if c.inputFilePath == "" {
		return fmt.Errorf("SetCommand: %v", errors.New("c.inputFilePath is empty. should SetInputFilePath() before SetCommand()."))
	}
	c.command = c.programConfig.ReplacedCmd(c.inputFilePath, c.OutputDir(), c.Parameta())
	return nil
}
