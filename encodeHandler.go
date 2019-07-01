package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/dineshappavoo/basex"
)

func encodeHandler(response http.ResponseWriter, request *http.Request, db Database, baseURL string) {
	decoder := json.NewDecoder(request.Body)
	var data struct {
		URL      string `json:"url"`
		ShortURL string `json:"shorturl"`
	}
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(response, `{"error": "Unable to parse json"}`, http.StatusBadRequest)
		return
	}

	if !govalidator.IsURL(data.URL) {
		http.Error(response, `{"error": "Not a valid URL"}`, http.StatusBadRequest)
		return
	}

	var shorturl string
	if data.ShortURL != "" {
		_, err := db.saveShort(data.URL, data.ShortURL)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		shorturl = data.ShortURL
	} else {
		id, _shorturl, err := db.GetID(data.URL)

		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		if id == 0 {
			exists := true

			lastID, err = db.getLastID()
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
			}

			for exists {
				lastID++
				s := strconv.FormatInt(lastID, 10)
				_shorturl, _ = basex.Encode(s)
				exists, _ = db.shortExists(_shorturl)
			}

			_, err := db.saveShort(data.URL, _shorturl)
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
				return
			}

			_, err = db.setLastID(lastID)
			if err != nil {
				http.Error(response, err.Error(), http.StatusInternalServerError)
			}
		}

		shorturl = _shorturl
	}

	resp := map[string]string{"url": baseURL + shorturl, "short_url": shorturl, "error": ""}

	jsonData, _ := json.Marshal(resp)
	response.Write(jsonData)
}
