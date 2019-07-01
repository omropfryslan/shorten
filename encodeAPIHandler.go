package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/dineshappavoo/basex"
)

type data struct {
	Success  int    `xml:"succes"`
	Error    string `xml:"error"`
	LongURL  string `xml:"longUrl"`
	ShortURL string `xml:"shortUrl"`
}

func encodeAPIHandler(response http.ResponseWriter, request *http.Request, db Database, baseURL string, apiKey string) {
	RequestURI := request.RequestURI

	u, _ := url.Parse(RequestURI)
	q := u.Query()

	if q.Get("apiKey") != apiKey {
		http.Error(response, "Invalid API key", http.StatusBadRequest)
		return
	}

	if !govalidator.IsURL(q.Get("longUrl")) {
		http.Error(response, "Invalid URL", http.StatusBadRequest)
		return
	}

	longURL := q.Get("longUrl")
	format := q.Get("format")

	var shortURL string
	id, shortURL, err := db.GetID(longURL)
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
			shortURL, _ = basex.Encode(s)
			exists, _ = db.shortExists(shortURL)
		}

		_id, err := db.saveShort(longURL, shortURL)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.setLastID(lastID)
		if err != nil {
			http.Error(response, err.Error(), http.StatusInternalServerError)
		}

		id = _id
	}

	xmlresp := &data{}
	xmlresp.ShortURL = baseURL + shortURL
	xmlresp.LongURL = longURL
	xmlresp.Success = 1
	xmlresp.Error = "0"

	var f func(interface{}, string, string) ([]byte, error)

	switch format {
	case "json":
		f = json.MarshalIndent
	default:
		f = xml.MarshalIndent
	}

	x, err := f(xmlresp, " ", "    ")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	switch format {
	case "json":
		response.Header().Set("Content-Type", "application/json")
	default:
		response.Header().Set("Content-Type", "application/xml")
		response.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"))
	}

	response.Write(x)
}
