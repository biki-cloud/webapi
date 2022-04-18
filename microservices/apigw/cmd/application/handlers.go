package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	"webapi/microservices/apigw/env"
	"webapi/microservices/apigw/pkg/getAllPrograms"
	"webapi/microservices/apigw/pkg/memoryGetter"
	"webapi/microservices/apigw/pkg/minimumServerSelector"
	"webapi/microservices/apigw/pkg/programHasServers"
	"webapi/microservices/apigw/pkg/serverAliveConfirmer"
)

// GetMinimumMemoryServerHandler は実際に疎通できるサーバの中から使用メモリが最小の
// サーバのURLをJSONで表示するAPI
func (app *Application) GetMinimumMemoryServerHandler(w http.ResponseWriter, r *http.Request) {
	app.Env = env.New()
	minimumMemoryServerSelector := minimumServerSelector.New()
	mg := memoryGetter.New()
	sac := serverAliveConfirmer.New()
	url, err := minimumMemoryServerSelector.Select(app.Env.ExecServers, sac, mg, "/health/memory", "/health")
	if err != nil {
		app.ServerError(w, err)
	}

	type j struct {
		Url string `json:"url"`
	}

	jsonStr := j{Url: url}

	app.RenderJSON(w, jsonStr)
}

// GetMinimumMemoryAndHasProgram は実際にプログラムがあるサーバかつ、使用メモリが最小の
// サーバのURLをJSONで表示するAPI
func (app *Application) GetMinimumMemoryAndHasProgram(w http.ResponseWriter, r *http.Request) {
	app.Env = env.New()
	programName := r.URL.Path[len("/program-server/minimumMemory-and-hasProgram/"):]
	app.InfoLog.Printf("programName: %v ", programName)
	sac := serverAliveConfirmer.New()
	aliveServers, err := serverAliveConfirmer.GetAliveServers(app.Env.ExecServers, "/health", sac)
	if err != nil {
		app.ServerError(w, err)
	}
	app.InfoLog.Printf("aliveservers: %v", aliveServers)

	programHasServersGetter := programHasServers.New()
	programHasServers, err := programHasServersGetter.Get(aliveServers, "/program/all", programName)
	if err != nil {
		app.ServerError(w, err)
	}

	app.InfoLog.Printf("programHasServers: %v ", programHasServers)

	minimumMemoryServerSelector := minimumServerSelector.New()
	mg := memoryGetter.New()
	url, err := minimumMemoryServerSelector.Select(programHasServers, sac, mg, "/health/memory", "/health")
	if err != nil {
		app.ServerError(w, err)
	}

	type j struct {
		Url string `json:"url"`
	}

	var jsonStr j
	if len(programHasServers) == 0 {
		jsonStr.Url = fmt.Sprintf("%v is not found in all exec.", programName)
	} else {
		jsonStr.Url = url
	}

	app.RenderJSON(w, jsonStr)
}

// mapToStruct はmapからstructに変換する。
func mapToStruct(m interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("mapToStruct: %v", err)
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return fmt.Errorf("mapToStruct: %v", err)
	}
	return nil
}

func (app *Application) GetAliveServersHandler(w http.ResponseWriter, r *http.Request) {
	app.Env = env.New()
	sac := serverAliveConfirmer.New()
	aliveServers, err := serverAliveConfirmer.GetAliveServers(app.Env.ExecServers, "/health", sac)
	if err != nil {
		app.ServerError(w, err)
	}

	type data struct {
		AliveServers []string `json:"AliveServers"`
	}

	jsonStr := data{AliveServers: aliveServers}

	app.RenderJSON(w, jsonStr)
}

func (app *Application) GetAllProgramsHandler(w http.ResponseWriter, r *http.Request) {
	app.Env = env.New()
	// allServerMapsはキーにプログラム名が入る。値はプログラム情報のmapが入る。
	allServerMaps := map[string]interface{}{}
	sac := serverAliveConfirmer.New()
	aliveServers, err := serverAliveConfirmer.GetAliveServers(app.Env.ExecServers, "/health", sac)
	if err != nil {
		app.ServerError(w, err)
	}

	allProgramGetter := getAllPrograms.New()
	endPoint := "/program/all"
	allServerMaps, err = allProgramGetter.Get(aliveServers, endPoint)

	app.RenderJSON(w, allServerMaps)
}
