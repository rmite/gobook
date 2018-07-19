// Copyright (c) 2018 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questa applicazione web è descritta nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"

//File principale dell'applicazione

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"rmite/todo"
	web "rmite/webman"
)

var app *web.ServerManager
var server *http.Server
var modelli *template.Template

var gn *todo.Gestore
var filtro todo.FiltroElenco
var uiMsg string

func main() {
	var err error

	//inizializza il gestore note
	if gn, err = todo.NewGestore("note.db"); err != nil {
		log.Fatalln(err)
	}
	filtro = todo.NessunFiltro

	//crea la mappa delle funzioni per i template
	fm := template.FuncMap{
		"msg":    usaMessaggio,
		"filtro": recuperaFiltro}

	//inizializza i template
	if modelli, err = template.New("").Funcs(fm).ParseFiles("privato\\modelli\\home.html", "privato\\modelli\\modifica.html", "privato\\modelli\\elimina.html"); err != nil {
		log.Fatalln(err)
	}

	//crea il gestore dell'applicazione
	app = web.NewServerManager(nil)

	//imposta i percorsi
	app.EnlistFuncOK("/", mostraHomepage)
	app.EnlistFuncOK("/inserisci", aggiungiNota)
	app.EnlistFuncOK("/modifica", modificaNota)
	app.EnlistFuncOK("/aggiorna", aggiornaNota)
	app.EnlistFuncOK("/cambia", cambiaStato)
	app.EnlistFuncOK("/avviso/rimuovi", avvisoRimuovi)
	app.EnlistFuncOK("/conferma/rimuovi", rimuoviNota)
	app.EnlistFuncOK("/chiudi", chiudiApp)
	app.EnlistFuncOK("/api/mostra/nota", apiMostraNota)

	//imposta il gestore dei file
	fs := http.FileServer(http.Dir(".\\pubblico"))
	app.EnlistOK("/img/", fs)
	app.EnlistOK("/files/", fs)

	//crea il server
	server = &http.Server{Addr: ":8080", Handler: app}

	log.Fatalln(server.ListenAndServe())
}

//recuperaFiltro restituisce il filtro per l'elenco delle note.
//Funzione usata nei template.
func recuperaFiltro() (f todo.FiltroElenco) {
	return filtro
}

//usaMessaggio restituisce il messaggio impostato nelle funzioni di gestione e lo cancella.
//Funzione usata nei template.
func usaMessaggio() (m string) {
	m = uiMsg
	uiMsg = "" //messaggio temporaneo, svuota la stringa
	return
}

//inviaMessaggio invia un messaggio all'utente.
//Restituisce true se il messaggio è mostrato con reindirizzamento.
func inviaMessaggio(w http.ResponseWriter, r *http.Request, redirectHome bool, code int, msg string) bool {
	if web.CheckAccept(r, []string{"application/json"}, false, w) {
		//vuole risposta in JSON
		risultato := RisultatoAPI{OK: (code == http.StatusOK), Messaggio: msg}
		web.ServeJSON(r, risultato, code, w)
		return false
	}

	uiMsg = msg
	if code < 300 || code > 399 {
		code = http.StatusFound
	}
	if redirectHome {
		http.Redirect(w, r, "/", code)
	}
	return true
}

//mostraPagina risponde ad una richiesta eseguendo il template specificato.
func mostraPagina(nome string, dati interface{}, w http.ResponseWriter, r *http.Request) bool {
	err := modelli.ExecuteTemplate(w, nome+".html", dati)

	if err != nil {
		app.ReplyStatus(http.StatusInternalServerError, err.Error(), w, r)
	}

	return (err == nil)
}

//mostraPaginaNota recupera la nota con id specificato in query string ed esegue il template specificato.
func mostraPaginaNota(nome string, w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")

	var err error
	var id int64

	id, err = strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		inviaMessaggio(w, r, true, http.StatusBadRequest, fmt.Sprintf("ID nota '%s' non valido.", idstr))
		return
	}

	var nt *todo.Nota

	nt, err = gn.Recupera(id)

	if err != nil {
		inviaMessaggio(w, r, true, http.StatusNotFound, fmt.Sprintf("Nota con ID '%s' non trovata.", idstr))
		return
	}

	mostraPagina(nome, nt, w, r)
}

//mostraHomepage gestisce la pagina iniziale.
func mostraHomepage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	switch path {
	case "/note/tutte":
		filtro = todo.NessunFiltro
	case "/note/fatte":
		filtro = todo.NoteFatte
	case "/note/dafare":
		filtro = todo.NoteDaFare
	}

	mostraPagina("home", gn, w, r)
}

//aggiungiNota gestisce l'aggiunta di una nota e reindirizza alla homepage.
func aggiungiNota(w http.ResponseWriter, r *http.Request) {
	if !web.CheckMethod(r, []string{http.MethodPost}, true, w) {
		return
	}

	testo := strings.TrimSpace(r.FormValue("nota"))

	if len(testo) == 0 {
		inviaMessaggio(w, r, true, http.StatusBadRequest, "Specifica il testo della nota.")
		return
	}

	if _, err := gn.Aggiungi(testo); err == nil {
		inviaMessaggio(w, r, true, http.StatusOK, "Nota aggiunta con successo.")
	} else {
		inviaMessaggio(w, r, true, http.StatusBadRequest, fmt.Sprintf("Operazione non riuscita: %s", err))
	}
}

//modificaNota gestisce la pagina per modificare una nota.
func modificaNota(w http.ResponseWriter, r *http.Request) {
	if !web.CheckMethod(r, []string{http.MethodGet}, true, w) {
		return
	}

	mostraPaginaNota("modifica", w, r)
}

//aggiornaNota gestisce l'aggiornamento di una nota.
func aggiornaNota(w http.ResponseWriter, r *http.Request) {
	if !web.CheckMethod(r, []string{http.MethodPost}, true, w) {
		return
	}

	var err error
	var id int64
	var testo string
	var fatto bool

	switch strings.Split(r.Header.Get("Content-Type"), ";")[0] {
	case "application/json":
		napi := leggiNotaAPI(r)
		if napi.Valida {
			id = napi.ID
			testo = strings.TrimSpace(napi.Testo)
			fatto = napi.Fatto
		} else {
			inviaMessaggio(w, r, true, http.StatusBadRequest, "Dati nota non validi.")
			return
		}
	case "application/x-www-form-urlencoded":
		idstr := r.FormValue("id")
		testo = strings.TrimSpace(r.FormValue("nota"))
		fatto = (r.FormValue("fatto") == "true")
		id, err = strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			inviaMessaggio(w, r, true, http.StatusBadRequest, fmt.Sprintf("ID nota '%s' non valido.", idstr))
			return
		}
	default:
		inviaMessaggio(w, r, true, http.StatusBadRequest, "Formato richiesta non valido.")
		return
	}

	var nt *todo.Nota

	nt, err = gn.Recupera(id)

	if err != nil {
		inviaMessaggio(w, r, true, http.StatusNotFound, fmt.Sprintf("Nota con ID '%d' non trovata.", id))
		return
	}

	if len(testo) == 0 {
		if inviaMessaggio(w, r, false, http.StatusBadRequest, "Specifica il testo della nota.") {
			mostraPagina("modifica", nt, w, r)
		}
		return
	}

	nt.Testo(testo)
	nt.Fatto = fatto

	if err = gn.Aggiorna(nt); err == nil {
		inviaMessaggio(w, r, true, http.StatusOK, "Nota aggiornata con successo.")
	} else {
		inviaMessaggio(w, r, true, http.StatusInternalServerError, fmt.Sprintf("Errore %s", err))
	}
}

//cambiaStato gestisce la modifica dello stato di una nota.
func cambiaStato(w http.ResponseWriter, r *http.Request) {
	if !web.CheckMethod(r, []string{http.MethodGet, http.MethodPost}, true, w) {
		return
	}

	valori := r.URL.Query()
	idstr := valori.Get("id")

	var err error
	var id int64

	id, err = strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		inviaMessaggio(w, r, true, http.StatusBadRequest, fmt.Sprintf("ID nota '%s' non valido.", idstr))
		return
	}

	var nt *todo.Nota

	nt, err = gn.Recupera(id)

	if err != nil {
		inviaMessaggio(w, r, true, http.StatusNotFound, fmt.Sprintf("Nota con ID '%s' non trovata.", idstr))
		return
	}

	fatto := (valori.Get("fatto") == "true")

	if nt.Fatto != fatto {
		if err = gn.CambiaStato(id, fatto); err != nil {
			inviaMessaggio(w, r, true, http.StatusInternalServerError, fmt.Sprintf("Errore %s", err))
			return
		}
	}

	inviaMessaggio(w, r, true, http.StatusOK, "Nota aggiornata con successo.")
}

//avvisoRimuovi chiede conferma di rimuovere una nota.
func avvisoRimuovi(w http.ResponseWriter, r *http.Request) {
	if !web.CheckMethod(r, []string{http.MethodGet}, true, w) {
		return
	}

	mostraPaginaNota("elimina", w, r)
}

//rimuoviNota gestisce la rimozione di una nota.
func rimuoviNota(w http.ResponseWriter, r *http.Request) {
	if !web.CheckMethod(r, []string{http.MethodGet, http.MethodDelete}, true, w) {
		return
	}

	idstr := r.URL.Query().Get("id")

	var err error
	var id int64

	id, err = strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		inviaMessaggio(w, r, true, http.StatusBadRequest, fmt.Sprintf("ID nota '%s' non valido.", idstr))
		return
	}

	if err = gn.Elimina(id); err != nil {
		inviaMessaggio(w, r, true, http.StatusInternalServerError, fmt.Sprintf("Errore %s", err))
		return
	}

	inviaMessaggio(w, r, true, http.StatusOK, "Nota eliminata.")
}

//chiudiApp avvia la chiusura e mostra una pagina per informare l'utente.
func chiudiApp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<html><head></head><body>La connessione &egrave; terminata.<br/>Puoi chiudere il browser.<br/>Arrivederci.</body></html>")
	go chiudi()
}

//chiude server, gestore note ed esce dall'applicazione
func chiudi() {
	time.Sleep(3 * time.Second)
	server.Close()
	gn.Chiudi()
	os.Exit(0)
}
