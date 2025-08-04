package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"fmt"
	"encoding/json"
	"reflect"
	"io"
)

type StubHeroStore struct {
	names map[string]string
	powerstat []HeroPowerStat
}

func TestGetId(t *testing.T) {
	store := StubHeroStore{
		map[string]string{
			"1": "A-Bomb",
			"100": "Black Flash",
			"247": "Evil Deadpool",
			"517": "Phoenix",
		}, nil,
	}
	server := NewHeroServer(&store)
	
	t.Run("returns name hero id=247", func(t *testing.T) {
		
		request := newGetNameRequest("247")
		response := httptest.NewRecorder()
				
		server.ServeHTTP(response, request)
		
		assertStatus(t, response.Code, http.StatusOK)
		assertResponceBody(t, response.Body.String(), "Evil Deadpool")
	})
	
	t.Run("returns name hero id=517", func(t *testing.T) {
		request := newGetNameRequest("517")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		
		assertStatus(t, response.Code, http.StatusOK)
		assertResponceBody(t, response.Body.String(), "Phoenix")
	})
	
	t.Run("returns name hero id=1", func(t *testing.T) {
		request := newGetNameRequest("1")
		response := httptest.NewRecorder()
				
		server.ServeHTTP(response, request)
		
		assertStatus(t, response.Code, http.StatusOK)
		assertResponceBody(t, response.Body.String(), "A-Bomb")
	})
	
	t.Run("returns name hero id=100", func(t *testing.T) {
		request := newGetNameRequest("100")
		response := httptest.NewRecorder()
				
		server.ServeHTTP(response, request)
		
		assertStatus(t, response.Code, http.StatusOK)
		assertResponceBody(t, response.Body.String(), "Black Flash")
	})
	
	t.Run("returns 404 on missing id hero", func(t *testing.T) {
		request := newGetNameRequest("1000")
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)
		
		assertStatus(t, response.Code, http.StatusNotFound)
		
	})
}

func TestPowerstats(t *testing.T){
	
	t.Run("Get powerstat of Bane", func(t  *testing.T){
		wantedPowerstat := []HeroPowerStat {
			{60, "Bane", 88, 38, 23, 56, 51, 95},	
		}
		
		store := StubHeroStore{nil, wantedPowerstat}
		server := NewHeroServer(&store)
		
		request := newPowerstatRequest(60)
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)
		
		got := getPowerstatFromResponse(t, response.Body)
		
		assertStatus(t, response.Code, http.StatusOK)
		assertPowerstat(t, got, wantedPowerstat)
		
		assertContentType(t, response, jsonContentType)
	})
	
	t.Run("Get powerstat of Batman", func(t  *testing.T){
		wantedPowerstat := []HeroPowerStat {
			{70, "Batman", 100, 26, 27, 50, 47, 100},	
		}
		
		store := StubHeroStore{nil, wantedPowerstat}
		server := NewHeroServer(&store)
		
		request := newPowerstatRequest(70)
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)
		
		got := getPowerstatFromResponse(t, response.Body)
		
		assertStatus(t, response.Code, http.StatusOK)
		assertPowerstat(t, got, wantedPowerstat)
		
		assertContentType(t, response, jsonContentType)
	})
	
	t.Run("Get powerstat of Tiger Shark", func(t  *testing.T){
		wantedPowerstat := []HeroPowerStat {
			{666, "Tiger Shark", 38, 72, 46, 70, 51, 28},	
		}
		
		store := StubHeroStore{nil, wantedPowerstat}
		server := NewHeroServer(&store)
		
		request := newPowerstatRequest(666)
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)
		
		got := getPowerstatFromResponse(t, response.Body)
		
		assertStatus(t, response.Code, http.StatusOK)
		assertPowerstat(t, got, wantedPowerstat)
		
		assertContentType(t, response, jsonContentType)
	})
}

func newGetNameRequest(id string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://superheroapi.com/api/4b3e7de93f96e6c75ce7e09a504a7c6b/%s", id), nil)
	return req
}

func newPowerstatRequest(id int) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://superheroapi.com/api/4b3e7de93f96e6c75ce7e09a504a7c6b/%q/powerstats", id), nil)
	return req
}

func getPowerstatFromResponse(t testing.TB, body io.Reader) (powerstat []HeroPowerStat) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&powerstat)
	
	if err != nil{
		t.Fatalf("Unable to parse response from server %q into slice of HeroPowerStat, '%v'", body, err)
	}
	
	return
}

func (s *StubHeroStore) GetHeroId(id string) string {
	name := s.names[id]
	return name
}

func (s *StubHeroStore) GetHeroPowerstat() []HeroPowerStat {
	return s.powerstat
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
	fmt.Println(got)
}

func assertResponceBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
	
}

func assertPowerstat(t testing.TB, got, want []HeroPowerStat) {
	t.Helper()
	
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}