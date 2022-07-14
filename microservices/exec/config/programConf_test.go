package config_test

import (
	"log"
	"strings"
	"testing"
	"webapi/pkg/os"

	"webapi/microservices/exec/config"
)

var currentDir string

func init() {
	c, err := os.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	config.SetProConfPath("programConfig_test.json")
}

func TestProgramConfig_ToProperPath(t *testing.T) {
	// 実際にヘルプやコマンドが記述してあるテキストファイル
	// を読み込んで、文字列などを表示できるかをテストする
	// 読み込むテストファイルはこのファイルがあるカレントディレクトリにある
	// testsディレクトリに入っているファイルたち

	p := config.NewProgramConfigHolder()
	p.SetHelpPath("tests/help.txt")
	p.SetCommand("python3 tests/xxxxx/eg.py")
	p.SetDetailedHelpPath("tests/detailedHelp.html")

	p.ToProperPath()

	// きちんとHelpが読み取れているか
	help, err := p.Help()
	if err != nil {
		t.Errorf("Help() : %v, want: %v \n", err, nil)
	}

	if !strings.Contains(help, "this is the help.txt") {
		t.Errorf("doesn't Contain %v of %v", "this is the help.txt", help)
	}

	// きちんとDetailedHelpが読み取れているか
	detailedHelp, err := p.DetailedHelp()
	if err != nil {
		t.Errorf("DetailedHelp() : %v, want: %v \n", err, nil)
	}

	if !strings.Contains(detailedHelp, "<h1>detailed help</h1>") {
		t.Errorf("doesn't Contain %v of %v", "<h1>detailed help</h1>", detailedHelp)
	}
}

func TestProgramConfig_ReplacedCmd(t *testing.T) {
	p := config.NewProgramConfigHolder()
	p.SetCommand("INPUTFILE OUTPUTDIR PARAMETA")
	// p.ProCommand = "INPUTFILE OUTPUTDIR PARAMETA"
	cmd := p.ReplacedCmd("i", "o", "p")

	if cmd != "i o p" {
		t.Errorf("ReplaceCmd(): %v, want: i o p \n", cmd)
	}
}

func TestGetProConfByName(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		testGetProConfByName(t, "convertToJson", true)
	})
	t.Run("test 2", func(t *testing.T) {
		testGetProConfByName(t, "dummy", false)
	})
}

func testGetProConfByName(t *testing.T, programName string, want bool) {
	_, err := config.GetProConfByName(programName)
	if (err != nil) == want {
		t.Errorf(err.Error())
	}
}

func TestGetPrograms(t *testing.T) {
	programConfigHolders, err := config.GetPrograms()
	if err != nil {
		t.Errorf("err from GetPrograms(): %v \n", err.Error())
	}
	for _, p := range programConfigHolders {
		if p.Name() != "convertToJson" {
			t.Errorf("p.Name() is not %v \n", "convertToJson")
		}
	}
}
