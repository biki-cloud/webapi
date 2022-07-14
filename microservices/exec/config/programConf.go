package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	os2 "webapi/pkg/os"

	"webapi/microservices/exec/env"
)

var (
	cfg env.Env = *env.New()
)

var proFilePath string

// SetProConfPath programConfig.jsonのパスをグローバル変数のproFilePathにセットする
func SetProConfPath(proConfName string) {
	if strings.Contains(proConfName, string(filepath.Separator)) {
		panic("proConfName want a file base name!!!!")
	}
	currentDir, err := os2.GetCurrentDir()
	if err != nil {
		panic("err msg: " + err.Error())
	}
	proFilePath = filepath.Join(currentDir, "config", proConfName)
}

// ProgramConfigHolder はプログラムの設定値を保持する構造体のインターフェース
type ProgramConfigHolder interface {
	Name() string
	Command() string
	Help() (string, error)
	DetailedHelp() (string, error)
	ReplacedCmd(string, string, string) string
	ToProperPath()
	SetCommand(string)
	SetHelpPath(string)
	SetDetailedHelpPath(string)
}

// programConfig はユーザーが入力したプログラムの情報を保持する構造体。
// ProgramConfigHolderインタフェースを実装している。
type programConfig struct {
	ProName     string `json:"name"`
	ProCommand  string `json:"command"`
	ProHelpPath string `json:"helpPath"`
	ProDetailedHelpPath string `json:"detailedHelpPath"`
	ProLogName  string `json:"logName"`
}

// NewProgramConfigHolder programConfigHolderインターフェースを返す
func NewProgramConfigHolder() ProgramConfigHolder {
	return &programConfig{}
}

// ToProperPath はパスを適切な感じに変換する。
func (n *programConfig) ToProperPath() {
	// windowsで実行する場合、programConfig.jsonのパスをlinuxのように
	// 記述したいので以下のことをする。windowsのパス区切りをlinuxに変換する。
	PathStrList := []string{"¥¥", "\\", "/"}
	for _, p := range PathStrList {
		n.ProCommand = strings.Replace(n.ProCommand, p, string(filepath.Separator), -1)
		n.ProName = strings.Replace(n.ProName, p, string(filepath.Separator), -1)
	}

	// コマンド、helpなどにcurrentDirを追加し、フルパスにする。
	// currentDirをProCommand,やProHelpPathなどに追加しないとパスがうまくつながらず、エラーが出る。
	currentDir, err := os2.GetCurrentDir()
	if err != nil {
		panic("err msg: " + err.Error())
	}
	n.ProCommand = strings.Replace(n.ProCommand, cfg.ProgramsDir, filepath.Join(currentDir, cfg.ProgramsDir), -1)
	n.ProHelpPath = strings.Replace(n.ProHelpPath, cfg.ProgramsDir, filepath.Join(currentDir, cfg.ProgramsDir), -1)
	n.ProDetailedHelpPath = strings.Replace(n.ProDetailedHelpPath, cfg.ProgramsDir, filepath.Join(currentDir, cfg.ProgramsDir), -1)
}

// Name は構造体のNameを返す
func (n *programConfig) Name() string {
	return n.ProName
}

// Command は構造体のcommandを返す
func (n *programConfig) Command() string {
	return n.ProCommand
}

// SetCommand コマンドをセットする
func (n *programConfig) SetCommand(cmd string) {
	n.ProCommand = cmd
}

// Help は構造体のhelpを返す
func (n *programConfig) Help() (string, error) {
	bytes, err := ioutil.ReadFile(n.ProHelpPath)
	if err != nil {
		return "", fmt.Errorf("Help: %v", err)
	}
	return string(bytes), nil
}

// SetHelpPath ヘルプパスをセットする
func (n *programConfig) SetHelpPath(help string) {
	n.ProHelpPath = help
}

// DetailedHelp は構造体のhelpを返す
func (n *programConfig) DetailedHelp() (string, error) {
	bytes, err := ioutil.ReadFile(n.ProDetailedHelpPath)
	if err != nil {
		return "", fmt.Errorf("DetailedHelp: %v", err)
	}
	return string(bytes), nil
}

// DetailedHelpPath
func (n *programConfig) SetDetailedHelpPath(path string) {
	n.ProDetailedHelpPath = path
}

// ReplacedCmd はprogram_config.jsonに記述されているコマンド文字列のINPUTFILE, OUTPUTDIR, PARAMETA を入力パラメータで置換する。
// 置換後のコマンドを返す。
func (n *programConfig) ReplacedCmd(infile string, outputDir string, parameta string) string {
	tmp1 := strings.Replace(n.ProCommand, "INPUTFILE", infile, 1)
	tmp2 := strings.Replace(tmp1, "OUTPUTDIR", outputDir, 1)
	cmd := strings.Replace(tmp2, "PARAMETA", parameta, 1)
	return cmd
}

// ErrProgramNotFound 指定されたプログラムが見つからなかった場合に返すエラーの構造体
var ErrProgramNotFound = errors.New("program not found error")

// GetProConfByName はプログラムの名前を受け取り、programConfig.jsonの中を検索ヒットした
// ものをProgramConfigHolder(インターフェース)として返す。
func GetProConfByName(programName string) (ProgramConfigHolder, error) {
	proConfigs, err := GetPrograms()
	if err != nil {
		return nil, fmt.Errorf("GetProConfByName: %v", err)
	}
	for _, program := range proConfigs {
		if programName == program.Name() {
			return program, nil
		}
	}
	return nil, fmt.Errorf("%v: %w", programName, ErrProgramNotFound)
}

// GetPrograms はprogramConfigHolderのインターフェースを返す。
func GetPrograms() ([]ProgramConfigHolder, error) {
	if proFilePath == "" {
		SetProConfPath("programConfig.json")
	}

	inputPath := proFilePath
	_, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("GetPrograms: %v", err)
	}

	byteArray, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("GetPrograms: %v", err)
	}

	var j map[string][]map[string]interface{}

	err = json.Unmarshal(byteArray, &j)
	if err != nil {
		return nil, fmt.Errorf("GetPrograms: %v", err)
	}

	proConfList := make([]ProgramConfigHolder, 0, 20)
	for _, v := range j["programs"] {
		var p programConfig
		err := MapToStruct(v, &p)
		if err != nil {
			return nil, fmt.Errorf("GetPrograms: %v", err)
		}
		p.ToProperPath()
		proConfList = append(proConfList, &p)
	}
	return proConfList, nil
}

// MapToStruct はmapからstructに変換する。
func MapToStruct(m map[string]interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("MapToStruct: %v", err)
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return fmt.Errorf("MapToStruct: %v", err)
	}
	return nil
}
