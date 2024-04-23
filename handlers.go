package main

import (
	"dictionaryAi/restDb"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type GeneratedSentenceResult struct {
	GeneratedRussianSentence string `json:"generatedRussianSentence"`
	GeneratedEngSentence     string `json:"generatedEngSentence"`
}

var words struct {
	Words []string `json:"words"`
}

type WordsToDelete struct {
	WordsId []int `json:"wordsIdToDelete"`
}

var word restDb.WordData

func timeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is: " + t + "\n"
	fmt.Fprintf(w, "%s", Body)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Error:", http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s\n", "Method not allowed!")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(d, &word)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}

	if word.Word != "" {
		result := restDb.AddWord(word)
		if result == -1 {
			http.Error(w, "Error:", http.StatusBadRequest)
			return
		}
		log.Println(word.Word, " was successfully added to DB")
		fmt.Fprintf(w, "%s was successfully added to DB", word.Word)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Error:", http.StatusBadRequest)
		return
	}
}

func GenerateSentenceHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host, r.Method)
	if r.Method != http.MethodPost {
		http.Error(rw, "Error:", http.StatusMethodNotAllowed)
		fmt.Fprintf(rw, "%s\n", "Method not allowed!")
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error:", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(d, &words)

	russianMessage, err := GenerateRussianMessageRequest(words.Words)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}
	engMessage, err := TranslateMessageRequest(russianMessage)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}
	result := GeneratedSentenceResult{
		GeneratedRussianSentence: russianMessage,
		GeneratedEngSentence:     engMessage,
	}

	e := json.NewEncoder(rw)

	err = e.Encode(result)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}
}

// SliceToJSON encodes a slice with JSON records
func SliceToJSON(slice interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(slice)
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	//указываем путь к нужному файлу
	path := filepath.Join("pages", "main.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//выводим шаблон клиенту в браузер
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

type notAllowedHandler struct{}

func (h notAllowedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	MethodNotAllowedHandler(rw, r)
}

func DefaultHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("DefaultHandler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	rw.WriteHeader(http.StatusNotFound)
	Body := r.URL.Path + " is not supported. Thanks for visiting!\n"
	fmt.Fprintf(rw, "%s", Body)
}

// MethodNotAllowedHandler is executed when the HTTP method is incorrect
func MethodNotAllowedHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	rw.WriteHeader(http.StatusNotFound)
	Body := "Method not allowed!\n"
	fmt.Fprintf(rw, "%s", Body)
}

// WordsListHandler is for getting all data from the worldList database
func WordsListHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("GetAllHandler Serving:", r.URL.Path, "from", r.Host)

	data, err := restDb.ListWords()
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = SliceToJSON(data, rw)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
}

// DeleteHandler is for deleting an existing user + DELETE
func DeleteHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("DeleteHandler Serving:", r.URL.Path, "from", r.Host)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error:", http.StatusBadRequest)
		return
	}

	var wordsListToDelete WordsToDelete

	err = json.Unmarshal(d, &wordsListToDelete)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Error:", http.StatusBadRequest)
		return
	}

	fmt.Println(wordsListToDelete)

	err = restDb.DeleteWords(wordsListToDelete.WordsId)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

}
