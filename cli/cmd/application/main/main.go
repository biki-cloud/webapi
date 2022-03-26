/*
WebAPIのコマンドラインツールとして動作するプログラム

処理の流れ

１、APIゲートウェイサーバ数台の中で生きているサーバかつ消費メモリが一番少ないサーバを選択する。

２、１で選択したAPIゲートウェイサーバに入力されたプログラム名でアクセスし、execサーバの中でプログラム名を保持しているかつ消費メモリが一番少ないexecサーバを選択する。

３、２で選択したexecサーバに入力されたファイルを処理させる。処理させた後は入力された出力ディレクトリに出力する。
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"webapi/cli/env"
	"webapi/cli/pkg/download"
	"webapi/cli/pkg/post"
	"webapi/cli/pkg/selectServer"
	"webapi/microservices/apigw/pkg/memoryGetter"
	"webapi/microservices/apigw/pkg/minimumServerSelector"
	"webapi/microservices/apigw/pkg/serverAliveConfirmer"
	pkgLog "webapi/pkg/log"
	pkgOs "webapi/pkg/os"
)

func main() {
	// 環境変数APIGW_SERVER_URISがセットされていなければデフォルトの値をセットする
	err := pkgOs.SetEnvIfNotExists("APIGW_SERVER_URIS", "http://127.0.0.1:8001,http://127.0.0.1:8002,http://127.0.0.1:8003")
	if err != nil {
		log.Fatalf("Err: %v \n", err.Error())
	}
	e := env.New()

	var (
		proName               string
		inputFile             string
		outputDir             string
		parameta              string
		LogFlag               bool
		outJSONFLAG           bool
		displayAllProgramFlag bool
	)
	flag.StringVar(&proName, "name", "", "(required) 登録プログラムの名称を入れてください。 登録されているプログラムは-aで参照できます。例 -> -name convertToJson")
	flag.StringVar(&inputFile, "i", "", "(required) 登録プログラムに処理させる入力ファイルのパスを指定してください。 例 -> -i ./input/test.txt")
	flag.StringVar(&outputDir, "o", "", "(required) 登録プログラムの出力ファイルを出力するディレクトリを指定してください。 例 -> -o ./proOut")
	parametaUsage := "(option) 登録プログラムに渡すパラメータを指定してください。例 -> -p " + strconv.Quote("-name mike")
	flag.StringVar(&parameta, "p", "", parametaUsage)
	flag.BoolVar(&LogFlag, "l", false, "(option) -lを付与すると詳細なログを出力します。通常は使用しません。")
	flag.BoolVar(&displayAllProgramFlag, "a", false, fmt.Sprintf("(option) -aを付与するとwebサーバに登録されているプログラムのリストを表示します。使用例 -> %s -a", flag.CommandLine.Name()))
	jsonExample := `
	{
	  "status": "program timeout or program error or server error or ok",
	  "stdout": "作成プログラムの標準出力",
	  "stderr": "作成プログラムの標準エラー出力",
	  "outURLs": [作成プログラムの出力ファイルのURLのリスト(この値は気にしなくて大丈夫です。)],
	  "errmsg": "サーバ内のプログラムで起きたエラーメッセージ"
	}
	statusの各項目
	program timeout -> 登録プログラムがサーバー内で実行された際にタイムアウトになった場合
	program error   -> 登録プログラムがサーバー内で実行された際にエラーになった場合
	server error    -> サーバー内のプログラムがエラーを起こした場合
	ok              -> エラーを起こさなかった場合
	`
	flag.BoolVar(&outJSONFLAG, "j", false, "(option, but recommend) -j を付与するとコマンド結果の出力がJSON形式になり、次のように出力します。"+jsonExample)

	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "\nUsage: \n  %s <option> -name <プログラム名> -i <入力ファイル> -o <出力ディレクトリ>\n", flag.CommandLine.Name())
		fmt.Fprintf(o, "\n\n"+
			"Description:  \n  "+
			"プログラムサーバに登録してあるプログラムを起動し、サーバ上で処理させ出力を返す。\n  "+
			"サーバで処理され、出力されたファイルは-oで指定した出力ディレクトリに出力されます。\n  "+
			"-oで指定する出力ディレクトリは存在しなければ作成されます。\n  "+
			"\n  実行する前にAPIゲートウェイサーバのアドレスを環境変数にセットしてください。値は環境に応じて変更してください。\n  "+
			"Linux  : export APIGW_SERVER_URIS=http://127.0.0.1:8001,http://127.0.0.1:8002,http://127.0.0.1:8003 \n  "+
			"Windows: SET APIGW_SERVER_URIS=http://127.0.0.1:8001,http://127.0.0.1:8002,http://127.0.0.1:8003 \n\n  ")
		fmt.Fprintf(o, " \n"+
			"Examples: \n")
		fmt.Fprintf(o,
			"# 一番シンプルな実行方法 \n  "+
				"%s -name <プログラム名> -i <入力ファイル> -o <出力ファイル> \n  ", flag.CommandLine.Name())
		fmt.Fprintf(o, " \n"+
			"# パラメータを付加させる場合, -pの後の文字列をダブルクォーテーションで囲む必要がある。中の文字列の構成は登録プログラムの仕様に依存する。 \n  "+
			"%s -name <プログラム名> -i <入力ファイル> -o <出力ファイル> -p \"<パラメータ１,パラメータ２>\" \n  ", flag.CommandLine.Name())
		fmt.Fprintf(o, " \n"+
			"# 実行結果をJSONで受け取る場合 \n  "+
			"%s -j -name <プログラム名> -i <入力ファイル> -o <出力ファイル> \n", flag.CommandLine.Name())
		fmt.Fprintf(o, " \n"+
			"# プログラムの処理過程を表示しながら実行する場合 \n  "+
			"%s -l -name <プログラム名> -i <入力ファイル> -o <出力ファイル> \n", flag.CommandLine.Name())
		fmt.Fprintf(o, "\n")
		fmt.Fprintf(o, "\nOptions: \n")
		flag.PrintDefaults()
		fmt.Fprintf(o, "\nUpdated date 2022.3.25 by morituka. \n\n")
	}
	flag.Parse()

	// 引数がなければヘルプを表示する
	if len(os.Args) == 1 {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	// ロガーをセットする
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// ロガーフラグがない場合はログを出さない
	if !LogFlag {
		logger.SetOutput(new(pkgLog.NullWriter))
	}

	// config.jsonのAPIゲートウェイサーバからメモリ消費が一番少ないサーバを選択する。
	// 生きているサーバのリストを取得
	confirmer := serverAliveConfirmer.New()
	aliveServers, err := serverAliveConfirmer.GetAliveServers(e.APIGWServerURIs, "/health", confirmer)
	if len(aliveServers) == 0 {
		log.Fatalln("生きているAPIゲートウェイサーバはありませんでした。")
	}
	mg := memoryGetter.New()
	if err != nil {
		log.Fatalf("err: %v \n", err)
	}

	// 生きているサーバにアクセスしていき、メモリ状況を取得、一番消費メモリが少ないサーバを取得する
	serverMemoryMap, err := minimumServerSelector.GetServerMemoryMap(aliveServers, "/health/memory", mg)
	if err != nil {
		log.Fatalf("APIゲートウェイサーバにてエラーが発生しました。err: %v\n", err)
	}

	minUrl := minimumServerSelector.GetMinimumMemoryServer(serverMemoryMap)

	apiGateWayServerAddr := minUrl

	logger.Printf("selected APIGateWay address: %v \n", apiGateWayServerAddr)

	// 全てのプログラム情報を取得する。allProgramInfoはjsonで出力される。
	command := fmt.Sprintf("curl %v/program-server/program/all", apiGateWayServerAddr)
	programsJSON, stderr, err := pkgOs.SimpleExec(command)
	if err != nil || programsJSON == "" {
		fmt.Printf("err from SimpleExec(command: %v), err msg: %v. stderr: %v, \n", command, err.Error(), stderr)
		os.Exit(1)
	}

	// プログラム一覧を確認する。 -a があれば実行
	if displayAllProgramFlag {
		fmt.Println(programsJSON)
		os.Exit(1)
	}

	// ---------- プログラムの実行に必要なパラメータが適切に準備されているか ----------
	requiredArgsAreOK := proName == "" || inputFile == "" || outputDir == ""

	if requiredArgsAreOK {
		fmt.Println("必須のパラメータが不足しています。")
		fmt.Println("-------------------------------------------")
		fmt.Printf("- name: %v\n", proName)
		fmt.Printf("- i   : %v\n", inputFile)
		fmt.Printf("- o   : %v\n", outputDir)
		fmt.Println("-------------------------------------------")
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	// 入力ファイル存在確認
	if !pkgOs.FileExists(inputFile) {
		fmt.Printf("no such file or directory: %v\n", inputFile)
		os.Exit(1)
	}

	// 出力ディレクトリなければ作成
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.Mkdir(outputDir, os.ModePerm); err != nil {
			fmt.Printf("err from Mkdir(%v): %v \n", outputDir, err.Error())
			os.Exit(1)
		}
	}

	// ---------- 入力されたプログラム名は存在するのか確認する --------------
	// json str からmapに変換
	var proMaps map[string]interface{}
	if err := json.Unmarshal([]byte(programsJSON), &proMaps); err != nil {
		fmt.Println(err)
	}

	// プログラム情報に入力されてプログラム名はあるのか、ない場合はエラーを出す。
	if _, ok := proMaps[proName]; !ok {
		fmt.Printf("%v is not found in all program server.", proName)
		os.Exit(1)
	}

	// ---------- プログラムサーバの中で入力ファイルを処理させる。 ----------
	// programServerAddr は目的のプログラムが登録してあるプログラムサーバの中で使用メモリが最小のサーバが入る。
	// APIゲートウェイサーバにアクセスして取得する。
	var programServerAddr string

	selector := selectServer.New()
	selectURL := apiGateWayServerAddr + "/program-server/minimumMemory-and-hasProgram/" + proName
	programServerAddr, err = selector.Select(selectURL)
	if err != nil {
		fmt.Printf("err from Select(): %v \n", err.Error())
		os.Exit(1)
	}

	logger.Printf("selected program exec address: %v \n", programServerAddr)

	// inputfile,parametaをサーバへ送信しサーバー上で処理する。
	proURL := fmt.Sprintf("%v/api/exec/%v", programServerAddr, proName)
	poster := post.New()
	res, err := poster.Post(proName, proURL, inputFile, parameta)
	if err != nil {
		fmt.Printf("サーバ上でエラーが発生しました。 err msg: %v \n", err.Error())
		os.Exit(1)
	}

	// サーバでの実行結果を表示する。
	if outJSONFLAG {
		b, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Printf("プログラムサーバからのレスポンスをJSONに変換するのを失敗しました。レスポンス: %v \n", res)
			os.Exit(1)
		}

		// サーバでの処理に成功し、正常にJSONが標準出力された場合
		fmt.Print(string(b))

	} else {
		//サーバでの処理に成功し、正常に標準出力、エラー出力のみを表示する場合
		fmt.Println(res.StdOut())
		fmt.Println(res.StdErr())
	}

	// サーバから出力されたJSONにファイルをダウンロードするためのURLが記述されているので
	// ファイルをカレントディレクトリにダウンロードし、出力ディレクトリへ移動させる
	done := make(chan error, len(res.OutURLs()))
	downloader := download.New()
	var wg sync.WaitGroup
	for _, getOutFileURL := range res.OutURLs() {
		wg.Add(1) // ゴルーチン起動のたびにインクリメント
		go downloader.Download(getOutFileURL, outputDir, done, &wg, pkgOs.NewMover())
	}
	wg.Wait()   // ゴルーチンでAddしたものが全てDoneされたら次に処理がいく
	close(done) // ゴルーチンが全て終了したのでチャネルをクローズする。

	for e := range done {
		if e != nil {
			fmt.Printf("プログラムサーバで処理が完了したファイルをダウンロードする際にエラーが発生しました。err msg: %v \n", err.Error())
			os.Exit(1)
		}
	}
}
