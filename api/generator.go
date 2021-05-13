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

const (
	username = "lalbert"
	password = "ib8jaeyahTei"
	hostname = "hn105529-001.dbaas.ovh.net:35552"
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
	StopsAt int    `json:"stopsAt,omitempty"`
	Social  string `json:"social,omitempty"`
}

type Action struct {
	Action       string        `json:"action,omitempty"`
	Index        int           `json:"index,omitempty"`
	Url          string        `json:"url,omitempty"`
	ToSearch     string        `json:"toSearch,omitempty"`
	WatchContext *WatchContext `json:"watchContext,omitempty"`
}

func WriteScript(infos *map[string]string, order *[]string) (string, error) {

	log.Println("ici")
	rand.Seed(time.Now().UnixNano())

	social := (*infos)["social"]
	stopsAt, _ := strconv.Atoi((*infos)["stopsAt"])

	interactionPercent, _ := strconv.Atoi((*infos)["interactionPercent"])

	n, _ := strconv.Atoi((*infos)["watchFromURL"])
	urls, err := getWatchURL(n, social, stopsAt, interactionPercent)
	if err != nil {
		fmt.Println("Error while getting info from DB")
	}

	next, _ := strconv.Atoi((*infos)["watchNext"])
	nexts := getWatchNext(next, social, stopsAt, interactionPercent)

	recommended, _ := strconv.Atoi((*infos)["watchRecommended"])
	recommendeds := getWatchRecommended(recommended, social, stopsAt, interactionPercent)

	home, _ := strconv.Atoi((*infos)["watchFromHome"])
	homes := getWatchFromHome(home, social, stopsAt, interactionPercent)

	channel, _ := strconv.Atoi((*infos)["watchFromChannel"])
	channels := getWatchFromChannel(channel, social, stopsAt, interactionPercent)

	searchType := (*infos)["search"]
	search, _ := strconv.Atoi((*infos)["watchFromSearch"])
	searches := getSearchAndWatch(search, searchType, social, stopsAt, interactionPercent)

	log.Println("All data formated !")

	json, err := write(order, urls, nexts, recommendeds, homes, channels, searches)
	if err != nil {
		log.Println(err)
	}
	res := `{ "actions": ` + string(json) + ` }`
	return res, nil
}

// TODO test with URL (connection to DB is needed
func getWatchURL(n int, social string, stopsAt int, interactionPercent int) (*[]Action, error) {

	if n == 0 {
		return &[]Action{}, nil
	}

	urls := make([]Action, n)
	count := 0

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", username, password, hostname))
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

		w := getWatchContext(social, stopsAt, interactionPercent)

		urls[count] = Action{
			Action:       "watch",
			Url:          "https://www.youtube.com/watch?v=" + id,
			WatchContext: w}
		count++
	}

	return &urls, nil
}

func getSearchAndWatch(n int, search string, social string, stopsAt int, interactionPercent int) *[]Action {

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
		w := getWatchContext(social, stopsAt, interactionPercent)

		searches[i] = Action{
			Action:   "search",
			ToSearch: str}
		searches[i+1] = Action{
			Action:       "watch",
			Index:        index,
			WatchContext: w}
	}

	return &searches
}

func getWatchFromChannel(n int, social string, stopsAt int, interactionPercent int) *[]Action {

	channels := make([]Action, 2*n)
	for i := 0; i < 2*n; i += 2 {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interactionPercent)

		channels[i] = Action{
			Action: "goToChannel"}
		channels[i+1] = Action{
			Action:       "watch",
			Index:        index,
			WatchContext: w}
	}
	return &channels
}

func getWatchFromHome(n int, social string, stopsAt int, interactionPercent int) *[]Action {

	homes := make([]Action, 2*n)
	for i := 0; i < 2*n; i += 2 {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interactionPercent)

		homes[i] = Action{
			Action: "goToHome"}
		homes[i+1] = Action{
			Action:       "watch",
			Index:        index,
			WatchContext: w}
	}
	return &homes
}

func getWatchNext(n int, social string, stopsAt int, interactionPercent int) *[]Action {

	nexts := make([]Action, n)

	for i := 0; i < n; i++ {
		w := getWatchContext(social, stopsAt, interactionPercent)
		nexts[i] = Action{
			Action:       "watch",
			Index:        1,
			WatchContext: w}
	}

	return &nexts
}

// TODO coorect typo in recommended
func getWatchRecommended(n int, social string, stopsAt int, interactionPercent int) *[]Action {

	recommendeds := make([]Action, n)
	for i := 0; i < n; i++ {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interactionPercent)

		recommendeds[i] = Action{
			Action:       "watch",
			Index:        index,
			WatchContext: w}
	}
	return &recommendeds
}

// TODO Improve watch context to handle several "stopsAt" variable
func getWatchContext(social string, stopsAt int, interactionPercent int) *WatchContext {

	r := rand.Intn(101)
	w := WatchContext{}
	if r < interactionPercent {
		w = WatchContext{Social: social, StopsAt: stopsAt}
	} else {
		w = WatchContext{StopsAt: stopsAt}
	}
	return &w
}

func writeOrdered(order *[]string, urls *[]Action, nexts *[]Action, recommendeds *[]Action, homes *[]Action, channels *[]Action, searches *[]Action) ([]byte, error) {

	actions := []Action{}
	for _, o := range *order {
		switch o {
		case "url":
			actions = append(actions, *urls...)
		case "next":
			actions = append(actions, *nexts...)
		case "recommended":
			actions = append(actions, *recommendeds...)
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

// TODO write without specified order
func writeUnordered(urls *[]Action, nexts *[]Action, recommendeds *[]Action, homes *[]Action, channels *[]Action, searches *[]Action) ([]byte, error) {

	urlIndex := 0
	nextIndex := 0
	recommendedIndex := 0
	homeIndex := 0
	channelIndex := 0
	searchIndex := 0

	n := len(*urls) + len(*nexts) + len(*recommendeds) + len(*homes) + len(*channels) + len(*searches)
	log.Println(n)
	actions := []Action{}

	for i := 0; i < n; i++ {
		r := rand.Intn(6)
		switch r {
		case 0:
			if urlIndex < len(*urls) {
				actions = append(actions, (*urls)[urlIndex])
				urlIndex++
			}
		case 1:
			if nextIndex < len(*nexts) {
				actions = append(actions, (*nexts)[nextIndex])
				nextIndex++
			}
		case 2:
			if recommendedIndex < len(*recommendeds) {
				actions = append(actions, (*recommendeds)[recommendedIndex])
				recommendedIndex++
			}
		case 3:
			if homeIndex < len(*homes) {
				actions = append(actions, (*homes)[homeIndex])
				homeIndex++
			}
		case 4:
			if channelIndex < len(*channels) {
				actions = append(actions, (*channels)[channelIndex])
				channelIndex++
			}
		case 5:
			if searchIndex < len(*searches) {
				actions = append(actions, (*searches)[searchIndex])
				searchIndex++
			}
		}
	}
	return json.Marshal(actions)

}

func write(order *[]string, urls *[]Action, nexts *[]Action, recommendeds *[]Action, homes *[]Action, channels *[]Action, searches *[]Action) ([]byte, error) {
	if len(*order) > 0 {
		return writeOrdered(order, urls, nexts, recommendeds, homes, channels, searches)
	}
	return writeUnordered(urls, nexts, recommendeds, homes, channels, searches)
}
