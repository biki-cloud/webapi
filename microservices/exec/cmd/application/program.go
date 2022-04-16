/*
クライアントプログラムからAPIにアクセスし、プログラムを実行するハンドラを定義している。
*/

package application

import (
	"errors"
	"fmt"
	"net/http"

	"webapi/microservices/exec/config"
	"webapi/microservices/exec/pkg/execution/contextManager"
	"webapi/microservices/exec/pkg/execution/executer"
	"webapi/microservices/exec/pkg/execution/outputManager"
	"webapi/microservices/exec/pkg/msgs"
)

// APIExec はプログラムを実行するためのハンドラー。処理結果をJSON文字列で返す。
// サーバの中でエラーが起こってもステータスコード200でJSONを返し、サーバエラーの詳細をJSONに記述する。
// cliからアクセスされる。cliの場合はこのハンドラにリクエストがくる前にファイルのアップロードは
// 完了し,アップロードディレクトリに格納されている。bodyにファイルベース名とパラメタを格納し、リクエストとしてこのハンドラ
// に送られる。
// アップロードファイルやパラメータ等を使用し、コマンド実行する。
func (app *Application) APIExec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	programName := r.URL.Path[len("/exec/"):]

	var out outputManager.OutputManager = outputManager.New()
	var newExecuter executer.Executer = executer.New()

	ctx, err := contextManager.New(w, r, app.Cfg)
	// プログラムがこのサーバになかった場合,もしくは他のエラーの場合
	if err != nil {
		if errors.Is(err, config.ProgramNotFoundError) {
			msg := fmt.Sprintf("%v is not found.", programName)
			out.SetErrorMsg(msg)
		} else {
			out.SetErrorMsg(err.Error())
		}
		app.InfoLog.Printf("err: %v \n", err.Error())
		out.SetStatus(msgs.SERVERERROR)
		app.RenderJSON(w, out)
		return
	}

	out = newExecuter.Execute(ctx)
	app.InfoLog.Println("---------------------------------------------------")
	app.InfoLog.Printf("EXEC: Name    : %v", ctx.ProgramName())
	app.InfoLog.Printf("EXEC: Infile  : %v", ctx.UploadedFilePath())
	app.InfoLog.Printf("EXEC: Parameta: %v", ctx.Parameta())
	app.InfoLog.Printf("EXEC: Command : %v", ctx.Command())
	app.InfoLog.Printf("EXEC: Status  : %v", out.Status())
	app.InfoLog.Printf("EXEC: ErrMsg  : %v", out.ErrorMsg())
	app.InfoLog.Println("---------------------------------------------------")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	app.RenderJSON(w, out)
}

// AllHandler は登録されているプログラムの全てをJSONで返す。
func (app *Application) AllHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.Body

	proConfList, err := config.GetPrograms()
	if err != nil {
		app.ServerError(w, err)
	}

	m := map[string]interface{}{}
	// proConfListから代入していく
	for _, ele := range proConfList {
		m1 := map[string]string{}
		m[ele.Name()] = m1
		//m1["command"] = ele.Command()
		help, err := ele.Help()
		if err != nil {
			app.ServerError(w, err)
		}
		m1["help"] = help
	}

	app.RenderJSON(w, m)
}
