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

type watchContext struct {
	stopsAt string `json:stopsAt, omitempty`
	social  string `json:social, omitempty`
}

type action struct {
	action       string
	index        string       `json:index, omitempty`
	url          string       `json:url, omitempty`
	toSearch     string       `json:toSearch, omitempty`
	watchContext watchContext `json:watchContext, omitempty`
}

func WriteScript(infos *map[string]string, order *[]string) error {

	social := (*infos)["social"]
	stopsAt := (*infos)["stopsAt"]

	interractionPercent, _ := strconv.Atoi((*infos)["interractionPercent"])

	n, _ := strconv.Atoi((*infos)["watchFromURL"])

	urls, err := getWatchURL(n, social, stopsAt, interractionPercent)
	if err != nil {
		fmt.Println("Error while getting info from DB")
	}

	next, _ := strconv.Atoi((*infos)["watchNext"])
	nexts := getWatchNext(next, social, stopsAt, interractionPercent)

	recommanded, _ := strconv.Atoi((*infos)["watchRecommanded"])
	recommandeds := getWatchRecommanded(recommanded, social, stopsAt, interractionPercent)

	home, _ := strconv.Atoi((*infos)["watchFromHome"])
	homes := getWatchFromHome(home, social, stopsAt, interractionPercent)

	channel, _ := strconv.Atoi((*infos)["watchFromChannel"])
	channels := getWatchFromChannel(channel, social, stopsAt, interractionPercent)

	searchType := (*infos)["search"]
	search, _ := strconv.Atoi((*infos)["watchFromSearch"])
	searches := getSearchAndWatch(search, searchType, social, stopsAt, interractionPercent)

	json, err := writeOrder(order, urls, nexts, recommandeds, homes, channels, searches)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(json)

	return nil
}

func getWatchURL(n int, social string, stopsAt string, interractionPercent int) (*[]action, error) {

	rand.Seed(time.Now().UnixNano())
	urls := []action{}
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

		urls[count] = action{
			action:       "watch",
			url:          "https://www.youtube.com/watch?v=" + id,
			watchContext: w}
		count++
	}

	return &urls, nil
}

func getSearchAndWatch(n int, search string, social string, stopsAt string, interractionPercent int) *[]action {

	rand.Seed(time.Now().UnixNano())
	searches := []action{}

	for i := 0; i < 2*n; i += 2 {

		str := ""
		if search == "conspi" {
			r := rand.Intn(len(conspiList))
			str = conspiList[r] + " "

			r = rand.Intn(len(covidList))
			str += covidList[r]
		}

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interractionPercent)

		searches[i] = action{
			action:   "search",
			toSearch: str}
		searches[i+1] = action{
			action:       "watch",
			index:        strconv.Itoa(index),
			watchContext: w}
	}

	return &searches
}

func getWatchFromChannel(n int, social string, stopsAt string, interractionPercent int) *[]action {

	channels := []action{}
	for i := 0; i < 2*n; i += 2 {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interractionPercent)

		channels[i] = action{
			action: "goToChannel"}
		channels[i+1] = action{
			action:       "watch",
			index:        strconv.Itoa(index),
			watchContext: w}
	}
	return &channels
}

func getWatchFromHome(n int, social string, stopsAt string, interractionPercent int) *[]action {

	homes := []action{}
	for i := 0; i < 2*n; i += 2 {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interractionPercent)

		homes[i] = action{
			action: "goToHome"}
		homes[i+1] = action{
			action:       "watch",
			index:        strconv.Itoa(index),
			watchContext: w}
	}
	return &homes
}

func getWatchNext(n int, social string, stopsAt string, interractionPercent int) *[]action {

	nexts := []action{}

	for i := 0; i < n; i++ {
		w := getWatchContext(social, stopsAt, interractionPercent)
		nexts[i] = action{
			action:       "watch",
			index:        "1",
			watchContext: w}
	}
	return &nexts
}

func getWatchRecommanded(n int, social string, stopsAt string, interractionPercent int) *[]action {

	recommandeds := []action{}
	for i := 0; i < n; i++ {

		index := rand.Intn(20) + 1
		w := getWatchContext(social, stopsAt, interractionPercent)

		recommandeds[i] = action{
			action:       "watch",
			index:        strconv.Itoa(index),
			watchContext: w}
	}
	return &recommandeds
}

func getWatchContext(social string, stopsAt string, interractionPercent int) watchContext {

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(101)
	w := watchContext{}
	if r < interractionPercent {
		w = watchContext{social: social, stopsAt: stopsAt}
	} else {
		w = watchContext{stopsAt: stopsAt}
	}
	return w
}

func writeOrder(order *[]string, urls *[]action, nexts *[]action, recommandeds *[]action, homes *[]action, channels *[]action, searches *[]action) ([]byte, error) {
	actions := []action{}
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
		case "seraches":
			actions = append(actions, *searches...)
		}
	}
	return json.Marshal(actions)

}
