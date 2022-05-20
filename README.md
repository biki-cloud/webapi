## EMS (Easy MicroService)

EMS is a platform where everybody can develop microservices.
Users can concentrate on developing their microservices.
Usually, EMS works on the k8s.

## Contents
- [Let's Get started](#lets-get-started)
- [Why I created EMS](#why-i-created-ems)
- [Web Servers of EMS](#web-servers-of-ems)
    - [website](#website)
    - [apigw](#apigw)
      - [apigw REST API](#apigw-REST-API)
    - [exec](#exec)
      - [exec REST API](#exec-REST-API)
    - [cli](#cli)
- [Test](#Test)
- [How to create microservices](#How-to-create-microservices)

## Let's Get started
how to work as local environment.

1. deploy to local
```shell
cd deployment/local 

./deploy.sh 
apigw1
apigw2
server1
server2
website1
cli
process kill
run all server

ps
PID TTY           TIME CMD
22639 ttys000    0:00.01 files/apigw1/apigw -port 8001
22640 ttys000    0:00.01 files/apigw2/apigw -port 8002
22641 ttys000    0:00.01 files/exec1/exec -port 9001
22642 ttys000    0:00.01 files/exec2/exec -port 9002
22643 ttys000    0:00.01 files/website1/website -port 7001
```

2. Verify <br>
Access to [http://localhost:7001/user/top](http://localhost:7001/user/top)

## Why I created EMS
description about that why i created EMS.

## Web Servers of EMS

## website
website is a web application that users use registered microservices of EMS.

## apigw
apigw works load balancing to exec service and manage registered microservices of exec service.

### apigw REST API
```go
package application

import (
  "net/http"

  pkgHttpHandlers "webapi/pkg/http/handlers"
)

func (app *Application) Routes() *http.ServeMux {
  router := http.NewServeMux()

  // Return exec server URL that uses lowest using memory.
  router.HandleFunc("/program-server/memory/minimum", app.GetMinimumMemoryServerHandler)

  // Return exec server URL that has program and uses lowest using program.
  router.HandleFunc("/program-server/minimumMemory-and-hasProgram/", app.GetMinimumMemoryAndHasProgram)

  // Return exec servers URL that is working in real time.
  router.HandleFunc("/program-server/alive", app.GetAliveServersHandler)

  // Get program info that all exec server has.
  // Return display it as JSON.
  router.HandleFunc("/program-server/program/all", app.GetAllProgramsHandler)

  // This handler uses that determine this server is alive. 
  router.HandleFunc("/health", pkgHttpHandlers.HealthHandler)

  // This handler displays the status of the memory of this server as JSON.
  router.HandleFunc("/health/memory", pkgHttpHandlers.GetRuntimeHandler)

  return router
}
```

## exec
exec works registering created microservices and executing microservices.

### exec REST API
```go
package application

import (
  "net/http"

  pkgHttpHandlers "webapi/pkg/http/handlers"
)

func (app *Application) Routes() *http.ServeMux {
  r := http.NewServeMux()

  // This handler works as file server handler.
  fileServer := "/" + app.Cfg.FileServer.Dir + "/"
  r.Handle(fileServer, http.StripPrefix(fileServer, http.FileServer(http.Dir(app.Cfg.FileServer.Dir))))

  // This API that executes registered service.
  r.HandleFunc("/api/exec/", app.APIExec)

  // This API that uploads file.
  r.HandleFunc("/api/upload", app.APIUpload)

  // This handler displays the status of the memory of this server as JSON.
  r.HandleFunc("/health/memory", pkgHttpHandlers.GetRuntimeHandler)

  // This API displays the program which is registered on this server.
  r.HandleFunc("/program/all", app.AllHandler)

  // This handler uses that determine this server is alive. 
  r.HandleFunc("/health", pkgHttpHandlers.HealthHandler)

  // This API works for downloading contents of the file server.
  r.HandleFunc("/download/", app.Download)

  return r
}
```

## cli
The user uses microservices that are registered EMS on the command line.

### How to use
```shell
# The simplest way to execute.
cli -name <program name> -i <input file> -o <output directory> 
   
# If add some parameta, need to surround letters were after -p with a double quotation.
# The letters depend on the content of microservices.
cli -name <program name> -i <input file> -o <output directory> -p "<parameter1,parameter2>" 
   
# In case that receives the result of executing with JSON.
cli -j -name <program name> -i <input file> -o <output directory> 
 
# In case that executing while outputting a log.
cli -l -name <program name> -i <input file> -o <output directory>
```


## Test
```shell
go test ./...
```

## How to create microservices
1. Determine microservices that we want to create <br>
Content of implement: We want to create microservices that take file, then output file that is added extension of ".json" to them.

<br>

2. Create project
```shell
mkdir ConvertToJson
cd ConvertToJson
touch convert_to_json.py
```

<br>

3. Determine command to execute
```shell
python3 convert_to_json.py <input file> <output dir> 
```

<br>

4. Coding
```python
import os
import shutil
import sys

infile = sys.argv[1]
output_dir = sys.argv[2]

outfile = os.path.join(output_dir, os.path.basename(infile)) + ".json"
shutil.move(infile, outfile)
```

<br>

5. Write help about coded program
```shell
cat help.txt
take file, then output file that is added extension of ".json" to them.
```

## How to register microservices to exec service.
1. Edit exec/config/programConfig.json
```json
{
  "programs": [
    {
      "name": "ConvertToJson",
      "command": "python3 programs/ConvertToJson/convert_to_json.py INPUTFILE OUTPUTDIR",
      "helpPath": "programs/ConvertToJson/help.txt"
    }
  ]
}
```

<br>

2. Move created project directory to exec/programs directory
```shell
mv ConvertToJson exec/programs/ConvertToJson
```
