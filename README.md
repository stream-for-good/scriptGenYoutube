# scriptGenYoutube

This repo contains the code needed to generate scripts for the bots (state machine) watching Youtube videos for the projet Stream4Good (and my M2 master thesis).


It is using Go module (to ease deployment) therefore verify your GO111MODULE is set to "on" (you can access your go env variable with `go env` and set the GO111MODULE thanks to `SETX GO111MODULE "on"`
To use it, start the program :

`go run main.go` or

`go build main.go` && `.\main`


The API starts on port `10001`

Send a request to the API (`http://127.0.0.1:10001/generate`)
With a JSON body of that form :

`	"type": "conspi",
	"watchNext": "15",
	"watchFromURL": "0",
	"watchFromHome": "10",
	"search": "conspi",
	"watchFromSearch": "5",
	"watchFromChannel": "5",
	"watchRecommended": "15",
	"stopsAt": "5",
	"social": "like",
 "interactionPercent": "50",
	"order": ["home", "next", "search", "channel", "recommended"]
 `
 
The parameters `watchNext`, `watchFromURL`, `watchFromSearch`, `watchFromChannel`, `watchRecommended` specifies the numer of video to watch from each one of the categories. (Note watchFromURL has not been tested yet as I can't curenetly connect to the DB)

The `type` parameter specifies the profile to chosse from the DB (NOT IMPLEMENTED YET)

The `search` parameter indicates what type of search we'd like to perform (choose word from dictionnaries in the code) (FOR NOW ONLY SUPPORTING 'conspi' and 'not conspi')

The `stopsAt` parameter indicates how long (in seconds) to watch each videos. A empty or omited field means that the videos will be fully watched

The `social` parameter indicate what action to do (either 'like' or 'dislike'). The `interactionPercent` indicates how often (in percent) this social action will be performed.

Finaly, the `order` specifies in which order the actions will be performed. An empty array (or ommited) will create a random order.


NOTE THAT:
You can omit parameters. All data are exepected as string (execpt for order which is an array of strings)

## Usage

```bash
curl -d ' {"type": "conspi", \
   "watchNext": "15",\
   "watchFromURL": "0",\
   "watchFromHome": "10",\
   "search": "conspi",\
   "watchFromSearch": "5",\
   "watchFromChannel": "5",\
   "watchRecommended": "15",\
   "stopsAt": "5",\
   "social": "like",\
   "interactionPercent": "50",
   \"order": ["home", "next", "search", "channel", "recommended"]}'\
   -H "Content-Type: application/json"\
   https://scriptgenyoutube.miage.dev/generate
```


## Building

docker image can be generated, pushed and run from the makefile

### creating docker container

```make build```

### pushing to repo

```make push```

### running the local docker image on port 8080

```make run```
