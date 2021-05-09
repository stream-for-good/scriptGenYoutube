package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var conspiList = []string{
	"5G",
	"Weapon",
	"Conspiracy",
	"Lies",
	"Reveal",
	"Truth",
	"Man-made",
	"Human-made",
	"Wuhan-Virus",
}

var covidList = []string{
	"Sars-Cov2",
	"Covid-19",
	"Corona Virus",
	"Corona",
	"Covid",
}

type WatchContext struct {
	StopsAt string `json:"stopsAt,omitempty"`
	Social  string `json:"social,omitempty"`
}

type Action struct {
	Action       string        `json:"action,omitempty"`
	Index        string        `json:"index,omitempty"`
	Url          string        `json:"url,omitempty"`
	ToSearch     string        `json:"toSearch,omitempty"`
	WatchContext *WatchContext `json:"watchContext,omitempty"`
}

func WriteScript(infos *map[string]string, order *[]string) (string, error) {

	log.Println("In WriteScript function")

	social := (*infos)["social"]
	stopsAt := (*infos)["stopsAt"]

	interactionPercent, _ := strconv.Atoi((*infos)["interactionPercent"])

	//n, _ := strconv.Atoi((*infos)["watchFromURL"])
	/*urls, err := getWatchURL(n, social, stopsAt, interractionPercent)
	if err != nil {
		fmt.Println("Error while getting info from DB")
	}*/
	urls := &[]Action{}

	next, _ := strconv.Atoi((*infos)["watchNext"])
	nexts := getWatchNext(next, social, stopsAt, interactionPercent)

	recommanded, _ := strconv.Atoi((*infos)["watchRecommanded"])
	recommandeds := getWatchRecommanded(recommanded, social, stopsAt, interactionPercent)

	home, _ := strconv.Atoi((*infos)["watchFromHome"])
	homes := getWatchFromHome(home, social, stopsAt, interactionPercent)

	channel, _ := strconv.Atoi((*infos)["watchFromChannel"])
	channels := getWatchFromChannel(channel, social, stopsAt, interactionPercent)

	searchType := (*infos)["search"]
	search, _ := strconv.Atoi((*infos)["watchFromSearch"])
	searches := getSearchAndWatch(search, searchType, social, stopsAt, interactionPercent)

	log.Println("All data formated")

	json, err := writeOrder(order, urls, nexts, recommandeds, homes, channels, searches)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(json))
	res := `{ "actions": ` + string(json) + ` }`
	return res, nil
}

func getWatchURL(n int, social string, stopsAt string, interractionPercent int) (*[]Action, error) {

	urls := make([]Action, n)
	count := 0

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/gocron")
	if err != nil {
		log.Fatal("Error connecting to the database ...")
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	ids, err := db.Query(`SELECT TOP (*infos)["watchFromURL"] id FROM video WHERE label='(*infos)["type"]' ORDER BY RAND`)
	if err != nil {
		return nil, err
	}
	defer ids.Close()

	id := ""
	for ids.Next() {
		if err := ids.Scan(&id); err != nil {
			log.Fatal(err)
			return nil, err
		}

		w := getWatchContext(social, stopsAt, interractionPercent)

		urls[count] = Action{
			Action:       "watch",
			Url:          "https://www.youtube.com/watch?v=" + id,
			WatchContext: w}
		count++
	}

	return &urls, nil
}

func getSearchAndWatch(n int, search string, social string, stopsAt string, interractionPercent int) *[]Action {

	log.Println("In search & watch function with %d actions", n)
	rand.Seed(time.Now().UnixNano())
	searches := make([]Action, 2*n)

	for i := 0; i < 2*n; i += 2 {

		str := ""
		if search == "conspi" {

			r := rand.Intn(len(covidList))
			str = covidList[r] + " "

			r = rand.Intn(len(conspiList))
			str += conspiList[r]
		}

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interractionPercent)

		searches[i] = Action{
			Action:   "search",
			ToSearch: str}
		searches[i+1] = Action{
			Action:       "watch",
			Index:        strconv.Itoa(index),
			WatchContext: w}
	}

	return &searches
}

func getWatchFromChannel(n int, social string, stopsAt string, interractionPercent int) *[]Action {

	rand.Seed(time.Now().UnixNano())
	channels := make([]Action, 2*n)
	for i := 0; i < 2*n; i += 2 {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interractionPercent)

		channels[i] = Action{
			Action: "goToChannel"}
		channels[i+1] = Action{
			Action:       "watch",
			Index:        strconv.Itoa(index),
			WatchContext: w}
	}
	return &channels
}

func getWatchFromHome(n int, social string, stopsAt string, interractionPercent int) *[]Action {

	rand.Seed(time.Now().UnixNano())
	homes := make([]Action, 2*n)
	for i := 0; i < 2*n; i += 2 {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interractionPercent)

		homes[i] = Action{
			Action: "goToHome"}
		homes[i+1] = Action{
			Action:       "watch",
			Index:        strconv.Itoa(index),
			WatchContext: w}
	}
	return &homes
}

func getWatchNext(n int, social string, stopsAt string, interractionPercent int) *[]Action {

	nexts := make([]Action, n)

	for i := 0; i < n; i++ {
		w := getWatchContext(social, stopsAt, interractionPercent)
		nexts[i] = Action{
			Action:       "watch",
			Index:        "1",
			WatchContext: w}
	}

	return &nexts
}

func getWatchRecommanded(n int, social string, stopsAt string, interractionPercent int) *[]Action {

	rand.Seed(time.Now().UnixNano())
	recommandeds := make([]Action, n)
	for i := 0; i < n; i++ {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interractionPercent)

		recommandeds[i] = Action{
			Action:       "watch",
			Index:        strconv.Itoa(index),
			WatchContext: w}
	}
	return &recommandeds
}

func getWatchContext(social string, stopsAt string, interractionPercent int) *WatchContext {

	log.Println(social, stopsAt, interractionPercent)
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(101)
	w := WatchContext{}
	if r < interractionPercent {
		w = WatchContext{Social: social, StopsAt: stopsAt}
		log.Println("ici ", w)
	} else {
		w = WatchContext{StopsAt: stopsAt}
	}
	log.Println(w)
	return &w
}

func writeOrder(order *[]string, urls *[]Action, nexts *[]Action, recommandeds *[]Action, homes *[]Action, channels *[]Action, searches *[]Action) ([]byte, error) {
	log.Println("In write order func")

	actions := []Action{}
	for _, o := range *order {
		switch o {
		case "url":
			actions = append(actions, *urls...)
		case "upNext":
			actions = append(actions, *nexts...)
		case "recommanded":
			actions = append(actions, *recommandeds...)
		case "home":
			actions = append(actions, *homes...)
		case "channel":
			actions = append(actions, *channels...)
		case "search":
			actions = append(actions, *searches...)
		}
	}
	return json.Marshal(actions)

}
