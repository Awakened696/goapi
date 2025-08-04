package main

import (
	"net/http"
	"fmt"
	"strings"
	"encoding/json"
)

const jsonContentType = "application/json"

type HeroStore interface {
	GetHeroId(id string) string
	GetHeroPowerstat() []HeroPowerStat
}

type HeroServer struct {
	store HeroStore
	http.Handler
}

type HeroPowerStat struct {
	Id 				int
	Name 			string
	Intelligence 	int
	Strength 		int
	Speed 			int
	Durability 		int
	Power 			int
	Combat 			int
}

func NewHeroServer(store HeroStore) *HeroServer {
	h := new(HeroServer)
	
	h.store = store
	
	router := http.NewServeMux()	
	router.Handle("/api/4b3e7de93f96e6c75ce7e09a504a7c6b//powerstats", http.HandlerFunc(h.powerstatsHandler))
	router.Handle("/api/4b3e7de93f96e6c75ce7e09a504a7c6b/", http.HandlerFunc(h.nameheroHandler))
	
	h.Handler = router
	
	return h
}

func (h *HeroServer) powerstatsHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(h.store.GetHeroPowerstat())

}

func (h *HeroServer) nameheroHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/4b3e7de93f96e6c75ce7e09a504a7c6b/")
	
	h.showHeroName(w, id)
}

func (h *HeroServer) showHeroName(w http.ResponseWriter, id string) {
	name := h.store.GetHeroId(id)
	
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	
	fmt.Fprint(w, name)
}