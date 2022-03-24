package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	gg "webapi/microservices/apigw/pkg/getAllPrograms"
	mg "webapi/microservices/apigw/pkg/memoryGetter"
	"webapi/microservices/apigw/pkg/minimumServerSelector"
	gp "webapi/microservices/apigw/pkg/programHasServers"
	sc "webapi/microservices/apigw/pkg/serverAliveConfirmer"
)

// GetMinimumMemoryServerHandler は実際に疎通できるサーバの中から使用メモリが最小の
// サーバのURLをJSONで表示するAPI
func (app *Application) GetMinimumMemoryServerHandler(w http.ResponseWriter, r *http.Request) {
	minimumMemoryServerSelector := minimumServerSelector.New()
	memoryGetter := mg.New()
	serverAliveConfirmer := sc.New()
	url, err := minimumMemoryServerSelector.Select(app.Env.ExecServers, serverAliveConfirmer, memoryGetter, "/health/memory", "/health")
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
	programName := r.URL.Path[len("/program-server/minimumMemory-and-hasProgram/"):]
	app.InfoLog.Printf("programName: %v ", programName)
	serverAliveConfirmer := sc.New()
	aliveServers, err := sc.GetAliveServers(app.Env.ExecServers, "/health", serverAliveConfirmer)
	if err != nil {
		app.ServerError(w, err)
	}
	app.InfoLog.Printf("aliveservers: %v", aliveServers)

	programHasServersGetter := gp.New()
	programHasServers, err := programHasServersGetter.Get(aliveServers, "/program/all", programName)
	if err != nil {
		app.ServerError(w, err)
	}

	app.InfoLog.Printf("programHasServers: %v ", programHasServers)

	minimumMemoryServerSelector := minimumServerSelector.New()
	memoryGetter := mg.New()
	url, err := minimumMemoryServerSelector.Select(programHasServers, serverAliveConfirmer, memoryGetter, "/health/memory", "/health")
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
	serverAliveConfirmer := sc.New()
	aliveServers, err := sc.GetAliveServers(app.Env.ExecServers, "/health", serverAliveConfirmer)
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
	// allServerMapsはキーにプログラム名が入る。値はプログラム情報のmapが入る。
	allServerMaps := map[string]interface{}{}
	serverAliveConfirmer := sc.New()
	aliveServers, err := sc.GetAliveServers(app.Env.ExecServers, "/health", serverAliveConfirmer)
	if err != nil {
		app.ServerError(w, err)
	}

	allProgramGetter := gg.New()
	endPoint := "/program/all"
	allServerMaps, err = allProgramGetter.Get(aliveServers, endPoint)

	app.RenderJSON(w, allServerMaps)
}
