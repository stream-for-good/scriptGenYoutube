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

	log.Println("Request Received...")

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
	getDataFromJSON(jsonData, &infos)

	order := []string{}
	getSliceFromJSON(jsonData, &order, &w)

	result, err := WriteScript(&infos, &order)
	if err != nil {
		respond(&w, http.StatusInternalServerError, err.Error())
	}
	respond(&w, http.StatusOK, result)
}

// Processes the JSON received by the server and fills a map to be used by the Mesh API
// Parameters :
// - json   : JSON representation of the task received of the server
// - infos  : Contains infos passed in the script and that are going to be used to generate scripts
// - w      : Used to respond error messages to sender in case of malformed JSON
func getDataFromJSON(jsonData *gabs.Container, infos *map[string]string) bool {

	scriptType, ok := jsonData.Path("type").Data().(string)
	if !ok {
		return false
	}
	(*infos)["type"] = strings.ToLower(scriptType)

	search, ok := jsonData.Path("search").Data().(string)
	if !ok {
		(*infos)["search"] = "0"
	}
	(*infos)["search"] = strings.ToLower(search)

	stopsAt, ok := jsonData.Path("stopsAt").Data().(string)
	if !ok {
		(*infos)["stopsAt"] = ""
	}
	(*infos)["stopsAt"] = strings.ToLower(stopsAt)

	social, ok := jsonData.Path("social").Data().(string)
	if !ok {
		(*infos)["social"] = ""
	}
	(*infos)["social"] = strings.ToLower(social)

	interactionPercent, ok := jsonData.Path("interactionPercent").Data().(string)
	if !ok {
		(*infos)["interactionPercent"] = "0"
	}
	(*infos)["interactionPercent"] = strings.ToLower(interactionPercent)

	watchNext, ok := jsonData.Path("watchNext").Data().(string)
	if !ok {
		(*infos)["watchNext"] = "0"
	}
	(*infos)["watchNext"] = watchNext

	watchFromURL, ok := jsonData.Path("watchFromURL").Data().(string)
	if !ok {
		(*infos)["watchFromURL"] = "0"
	}
	(*infos)["watchFromURL"] = watchFromURL

	watchFromHome, ok := jsonData.Path("watchFromHome").Data().(string)
	if !ok {
		(*infos)["watchFromHome"] = "0"
	}
	(*infos)["watchFromHome"] = watchFromHome

	watchFromSearch, ok := jsonData.Path("watchFromSearch").Data().(string)
	if !ok {
		(*infos)["watchFromSearch"] = "0"
	}
	(*infos)["watchFromSearch"] = watchFromSearch

	watchFromChannel, ok := jsonData.Path("watchFromChannel").Data().(string)
	if !ok {
		(*infos)["watchFromChannel"] = "0"
	}
	(*infos)["watchFromChannel"] = watchFromChannel

	watchRecommended, ok := jsonData.Path("watchRecommended").Data().(string)
	if !ok {
		(*infos)["watchRecommended"] = "0"
	}
	(*infos)["watchRecommended"] = watchRecommended
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
