package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

var linksMap = map[string]string{}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sEnc := base64.StdEncoding.EncodeToString(body)
	linksMap[sEnc] = string(body)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http" + "://" + r.Host + r.URL.String() + sEnc))
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	linkID := chi.URLParam(r, "linkID")

	link, ok := linksMap[linkID]

	fmt.Println(linkID, link)

	if !ok {
		http.Error(w, `Link not found`, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link, http.StatusTemporaryRedirect)
}