package main

import (
	asciiart "ascii-art-web/ascii-art"
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Titre string
	Texte string
	Choix string
	Ecrit string
}

var texteSaisi string = ""
var choixSaisi string = "standard"
var textecrit string = ""

func accueil(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := PageData{Titre: "Mon input", Texte: texteSaisi, Choix: choixSaisi, Ecrit: textecrit}
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "error: 500 ", 500)
		return
	}
}

func texte(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "cette route n'accepte que les requêtes POST", 400)
		return
	}
	choixSaisi = r.FormValue("choix")
	textecrit = r.FormValue("monTexte")

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "template introuvable", 404)
		return
	}
	resultat, err := asciiart.AsciiArt(textecrit, choixSaisi)
	if err != nil {
		if err.Error() == "404" {
			w.WriteHeader(404)
		} else if err.Error() == "400" {
			w.WriteHeader(400)
		}
		data := PageData{Titre: "Mon input", Texte: "Erreur : " + err.Error(), Choix: choixSaisi, Ecrit: textecrit}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "error: 500 ", 500)
			return
		}
		return
	}

	texteSaisi = resultat
	data := PageData{Titre: "Mon input", Texte: texteSaisi, Choix: choixSaisi, Ecrit: textecrit}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "error: 500 ", 500)
		return
	}
}

func main() {
	http.HandleFunc("/", accueil)
	http.HandleFunc("/ascii-art", texte)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
