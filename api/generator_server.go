package api

import (
	"io"
	"strings"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
)

func Generate(w http.ResponseWriter, r *http.Request) {

	log.Println("Service Started...")

	req, err := ioutil.ReadAll(r.Body)
	defer bodyCloser(r.Body)
	if err != nil {
		log.Println(err)
		respond(&w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonData, err := gabs.ParseJSON(req)
	if err != nil {
		respond(&w, http.StatusInternalServerError, "You must provide a valid JSON as body")
		return
	}

	infos := make(map[string]string)
	getDataFromJSON(jsonData, &infos, &w)

	order := []string{}
	getSliceFromJSON(jsonData, &order, &w)

	err = WriteScript(&infos, &order)

	if err != nil {
		respond(&w, http.StatusInternalServerError, err.Error())
	}
}

// Processes the JSON received by the server and fills a map to be used by the Mesh API
// Parameters :
// - json   : JSON representation of the task received of the server
// - infos  : Contains infos passed in the script and that are going to be used to generate scripts
// - w      : Used to respond error messages to sender in case of malformed JSON
func getDataFromJSON(jsonData *gabs.Container, infos *map[string]string, w *http.ResponseWriter) bool {
	log.Println("In getDataFromJSON function")

	scriptType, ok := jsonData.Path("type").Data().(string)
	if !ok {
		return false
	}
	(*infos)["type"] = strings.ToLower(scriptType)

	search, ok := jsonData.Path("search").Data().(string)
	if !ok {
		return false
	}
	log.Println("%s search actions", strings.ToLower(search))
	(*infos)["search"] = strings.ToLower(search)

	timeWatching, ok := jsonData.Path("timeWatching").Data().(string)
	if !ok {
		return false
	}
	(*infos)["timeWatching"] = strings.ToLower(timeWatching)

	social, ok := jsonData.Path("social").Data().(string)
	if !ok {
		return false
	}
	(*infos)["social"] = strings.ToLower(social)

	watchNext, ok := jsonData.Path("watchNext").Data().(string)
	if !ok {
		return false
	}
	(*infos)["watchNext"] = watchNext

	watchFromURL, ok := jsonData.Path("watchFromURL").Data().(string)
	if !ok {
		return false
	}
	(*infos)["watchFromURL"] = watchFromURL

	watchFromHome, ok := jsonData.Path("watchFromHome").Data().(string)
	if !ok {
		return false
	}
	(*infos)["watchFromHome"] = watchFromHome

	watchFromSearch, ok := jsonData.Path("watchFromSearch").Data().(string)
	if !ok {
		return false
	}
	(*infos)["watchFromSearch"] = watchFromSearch

	watchFromChannel, ok := jsonData.Path("watchFromChannel").Data().(string)
	if !ok {
		return false
	}
	(*infos)["watchFromChannel"] = watchFromChannel

	watchRecommanded, ok := jsonData.Path("watchRecommanded").Data().(string)
	if !ok {
		return false
	}
	(*infos)["watchRecommanded"] = watchRecommanded
	return true
}

func getSliceFromJSON(jsonData *gabs.Container, infos *[]string, w *http.ResponseWriter) bool {
	children := jsonData.S("order").Children()
	/*if err != nil {
		log.Println(err)
		return false
	}*/
	for _, child := range children {
		*infos = append(*infos, child.Data().(string))
	}
	return true
}

// Wrapper to easily send an HTTP status code and message on given ResponseWriter
// Parameters:
// - w       : ResponseWriter to use
// - status  : status code to send (http.StatusOK for ex)
// - msg     : message to be sent as body
func respond(w *http.ResponseWriter, status int, msg string) {
	(*w).WriteHeader(status)
	io.Copy(*w, strings.NewReader(msg))
}

// close the given body
func bodyCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}
