package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"parser/bechmarker"
	"parser/core"
	"parser/serp"
)

const baseYandexURL = "https://yandex.ru/search/touch/?service=www.yandex&ui=webmobileapp.yandex&numdoc=50&lr=213&p=0&text=%s"

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	resp, err := http.Get(baseYandexURL + search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	hostMap := make(map[string]int)
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		res := serp.ParseYandexResponse(bodyBytes)
		if res.Error == nil {
			bechmarker.Benchmark(res.Items)
			for _, item := range res.Items {
				val, err := core.Get(item.Host)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				hostMap[item.Host] = val
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	j, err := json.Marshal(hostMap)
	if err == nil {
		w.Write(j)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Handle() {
	r := mux.NewRouter()
	r.HandleFunc("/sites", keyValueGetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
