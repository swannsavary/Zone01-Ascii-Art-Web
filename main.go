package main

import (
	asciiart "ascii-art-web/ascii-art"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// donnees de chaque page
type PageData struct {
	Titre string
	Texte string
	Choix string
	Ecrit string
}

var texteSaisi string = ""
var choixSaisi string = "standard"
var textecrit string = ""
var tmpl *template.Template

// page d'accueille du site
func accueil(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "error: 404 ", 404)
		return
	}
	data := PageData{Titre: "Ascii-Art-Web", Texte: texteSaisi, Choix: choixSaisi, Ecrit: textecrit}
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "error: 500 ", 500)
		return
	}
}

// prend l'input puis l'affiche en mode ascii-art
func texte(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "cette route n'accepte que les requêtes POST", 400)
		return
	}
	choixSaisi = r.FormValue("choix")
	textecrit = r.FormValue("monTexte")
	texteNormalise := strings.ReplaceAll(textecrit, "\r\n", "\n")
	resultat, err := asciiart.AsciiArt(texteNormalise, choixSaisi)

	if err != nil {
		if err.Error() == "500" {
			w.WriteHeader(500)
		} else if err.Error() == "400" {
			w.WriteHeader(400)
		}
		data := PageData{Titre: "Ascii-Art-Web", Texte: "Erreur : " + err.Error(), Choix: choixSaisi, Ecrit: textecrit}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "error: 500 ", 500)
			return
		}
		return
	}

	texteSaisi = resultat
	data := PageData{Titre: "Ascii-Art-Web", Texte: texteSaisi, Choix: choixSaisi, Ecrit: textecrit}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "error: 500 ", 500)
		return
	}
}

func main() {
	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", accueil)
	http.HandleFunc("/ascii-art", texte)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
