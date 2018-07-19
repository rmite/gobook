// Copyright (c) 2018 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questa libreria è descritta nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"

/*
Package webman implementa tipi e funzioni per la gestione di un server web.
*/
package webman

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// ===== Tipo serverAction =====

type serverAction struct {
	path        string
	statusReply bool
	statusCode  int
	handler     http.Handler
}

// isStatusReply restituisce true se l'oggetto serverAction gestisce il codice di stato specificato.
func (sa serverAction) isStatusReply(code int) bool {
	return (sa.statusReply && sa.statusCode == code)
}

// statusKey restituisce la chiave per i gestori di stato.
func statusKey(code int) string {
	return fmt.Sprintf("#STATUS_%d#", code)
}

// ===== Tipo actionMap =====

type actionMap map[string]serverAction

// ===== Tipo ServerManager =====

/*
ServerManager rappresenta il gestore server a cui possono essere
associati gestori di richiesta per percorsi e per codici di stato HTTP.

I gestori di codice sono associati con il metodo EnlistStatusReply, mentre con
i metodi Enlist e EnlistFunc puoi associare un gestore a un percorso.

Un percorso che termina con lo slash indica che il relativo gestore può replicare a tutte le richieste il cui percorso ha quella stessa radice.
*/
type ServerManager struct {
	actions actionMap
	log     *log.Logger
}

//ErrServerManagerNotReady è l'errore restituito quando il server manager non è stato inizializzato.
var ErrServerManagerNotReady error = errors.New("server manager not ready")

/*
NewServerManager restituisce un oggetto ServerManager inizializzato.

Nell'oggetto sono impostate le seguenti risposte di default:

  DefaultBadRequestReply per il codice http.StatusBadRequest (400)
  DefaultNotFoundReply per il codice http.StatusNotFound (404)
  DefaultServerErrorReply per il codice http.StatusInternalServerError (500)

per cambiare questi gestori di stato o aggiungere altri usa il metodo EnlistStatusReply.

Per aggiungere o cambiare gestori di percorso usa i metodi Enlist e EnlistFunc.

Il log interno è scritto sull'oggetto log.Logger specificato. Se il parametro è nil, il gestore scrive su os.Stderr.
*/
func NewServerManager(logger *log.Logger) *ServerManager {
	// crea l'oggetto
	sm := &ServerManager{log: logger}
	// imposta il logger default se necessario
	if sm.log == nil {
		sm.log = log.New(os.Stderr, "", (log.LstdFlags | log.LUTC))
	}
	// inizializza la mappa azioni
	sm.actions = make(actionMap, 5)
	// imposta i gestori di default
	sm.EnlistStatusReply(http.StatusBadRequest, http.HandlerFunc(DefaultBadRequestReply))
	sm.EnlistStatusReply(http.StatusNotFound, http.HandlerFunc(DefaultNotFoundReply))
	sm.EnlistStatusReply(http.StatusInternalServerError, http.HandlerFunc(DefaultServerErrorReply))
	return sm
}

//IsReady restituisce true quando il gestore è stato inizializzato con NewServerManager.
func (sm *ServerManager) IsReady() bool {
	return (sm.actions != nil)
}

/*
Enlist imposta il gestore per un determinato percorso di richiesta e con un codice di stato specificati.

Il metodo genera un panic con l'errore ErrServerManagerNotReady se il metodo IsReady restituisce false.
*/
func (sm *ServerManager) Enlist(path string, hler http.Handler, code int) {
	if !sm.IsReady() {
		panic(ErrServerManagerNotReady)
	}
	var pattern string
	if path[0] == '/' {
		pattern = path
	} else {
		pattern = "/" + path
	}
	sm.actions[pattern] = serverAction{path: pattern, statusReply: false, statusCode: code, handler: hler}
}

/*
EnlistOK imposta il gestore per un determinato percorso di richiesta specificato e con il codice di stato http.StatusOK.

Il metodo genera un panic con l'errore ErrServerManagerNotReady se il metodo IsReady restituisce false.
*/
func (sm *ServerManager) EnlistOK(path string, hler http.Handler) {
	sm.Enlist(path, hler, http.StatusOK)
}

/*
EnlistFunc imposta una funzione come gestore per un determinato percorso di richiesta e con un codice di stato specificati.

Il metodo genera un panic con l'errore ErrServerManagerNotReady se il metodo IsReady restituisce false.
*/
func (sm *ServerManager) EnlistFunc(path string, f http.HandlerFunc, code int) {
	sm.Enlist(path, f, code)
}

/*
EnlistFuncOK imposta una funzione come gestore per un determinato percorso di richiesta specificato e con il codice di stato http.StatusOK.

Il metodo genera un panic con l'errore ErrServerManagerNotReady se il metodo IsReady restituisce false.
*/
func (sm *ServerManager) EnlistFuncOK(path string, f http.HandlerFunc) {
	sm.Enlist(path, f, http.StatusOK)
}

/*
EnlistStatusReply imposta il gestore per un codice di stato.

Il metodo genera un panic con l'errore ErrServerManagerNotReady se il metodo IsReady restituisce false.

Con questo metodo puoi impostare ad esempio il gestore per il codice 404 (Pagina non trovata).
*/
func (sm *ServerManager) EnlistStatusReply(code int, hler http.Handler) {
	if !sm.IsReady() {
		panic(ErrServerManagerNotReady)
	}
	pattern := statusKey(code)
	sm.actions[pattern] = serverAction{path: "", statusReply: true, statusCode: code, handler: hler}
}

/*
ReplyStatus risponde alla richiesta con il gestore associato al codice HTTP specificato.

Il metodo genera un panic con l'errore ErrServerManagerNotReady se il metodo IsReady restituisce false.

Un'annotazione contenente il codice HTTP della risposta è inserita all'interno del log interno.

Se nessun gestore è stato associato al codice specificato, ReplyStatus affida la risposta
alla funzione WriteStatus passando il codice e il messaggio specificati insieme al percorso della richiesta.
*/
func (sm *ServerManager) ReplyStatus(code int, message string, w http.ResponseWriter, r *http.Request) {
	if !sm.IsReady() {
		panic(ErrServerManagerNotReady)
	}
	// scrive nel log lo stato e il messaggio
	sm.log.Printf("RISPOSTA: STATO %d\n", code)
	sm.log.Printf("MESSAGGIO: [%s]\n", message)
	// trova l'azione con il gestore risposta stato
	a := sm.actions[statusKey(code)]
	if a.isStatusReply(code) && a.handler != nil {
		// se c'è, lo usa
		a.handler.ServeHTTP(w, r)
	} else {
		// se non c'è, risponde direttamente
		WriteStatus(code, message, r.URL.Path, w)
	}
}

// pathBegins restituisce true se path richiesta è uguale a path azione o ha la stessa radice.
func pathBegins(reqPath, actPath string) bool {
	if len(actPath) == 0 || len(reqPath) == 0 {
		return false
	}
	actPure := strings.TrimSuffix(actPath, "/")
	reqPure := strings.TrimSuffix(reqPath, "/")
	if actPure == reqPure {
		return true
	}
	actPure += "/"
	aLen := len(actPure)
	return ((len(reqPath) >= aLen) && (reqPath[:aLen] == actPure))
}

// getAction cerca l'azione per il percorso di richiesta specificato.
func (sm *ServerManager) getAction(reqPath string) (a serverAction) {
	// imposta il ritorno su reply StatusNotFound
	a = serverAction{path: "", statusReply: true, statusCode: http.StatusNotFound}
	// esce se ServerManager non è stato inizializzato
	if !sm.IsReady() {
		return
	}
	// ricerca azione con path identica
	if action, ok := sm.actions[reqPath]; ok {
		return action
	}
	// ricerca azione con path identica a radice path richiesta
	var lastFoundLen = -1
	for _, act := range sm.actions {
		if !pathBegins(reqPath, act.path) {
			continue
		}
		if len(act.path) > lastFoundLen {
			lastFoundLen = len(act.path)
			a = act
		}
	}
	return
}

// -- Implementazione dell'interfaccia http.Handler --

/*
ServeHTTP smista le richieste dal server ai vari gestori associati a ServerManager.

Il metodo scrive nel log interno i dettagli della richiesta e cerca il gestore a cui affidare la richiesta.

Se trova il gestore che può replicare al percorso di richiesta, usa il relativo metodo ServeHTTP.

Se il gestore associato a un percorso è nil, risponde con ReplyStatus inviando il codice HTTP associato al percorso.

Se nessun gestore può replicare al percorso di richiesta, il metodo invia la risposta 404 (Pagina non trovata) con il gestore di stato associato al codice.
Questo gestore corrisponde a DefaultNotFoundReply se non è stato sostituito con il metodo EnlistStatusReply.
*/
func (sm *ServerManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// scrive nel log i dettagli della richiesta
	sm.log.Print("====================")
	sm.log.Printf("host: %s\n", r.Host)
	sm.log.Printf("rURI: %s\n", r.RequestURI)
	for k, v := range r.Header {
		sm.log.Printf(" %s = %s\n", k, v)
	}
	// recupera l'azione associata al percorso di richiesta
	action := sm.getAction(r.URL.Path)
	// azione risposta di stato
	if action.statusReply {
		// l'azione rappresenta una risposta di stato
		// es. codice 404 - Pagina non trovata
		sm.ReplyStatus(action.statusCode, "", w, r)
		return
	}
	// altrimenti azione associata ad un percorso
	if action.handler != nil {
		// gestore disponibile, lo usa
		action.handler.ServeHTTP(w, r)
	} else {
		// gestore non disponibile
		// risponde con il codice di stato dell'azione
		sm.ReplyStatus(action.statusCode, "", w, r)
	}
}

// ===== FUNZIONI =====

/*
WriteStatus risponde con un messaggio di testo e un codice di stato specificati.
La risposta contiene anche un percorso se specificato.

Nell'intestazione della risposta, il valore "Content-Type" è impostato su "text/plain; charset=utf-8".
*/
func WriteStatus(code int, message string, path string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	var reply string
	if len(path) > 0 {
		reply = fmt.Sprintf("%d - %s - Path: %s", code, message, path)
	} else {
		reply = fmt.Sprintf("%d - %s", code, message)
	}
	fmt.Fprintf(w, reply)
}

/*
DefaultBadRequestReply invia la risposta di default per il codice 400.

  400 - Richiesta non valida. - Path: [percorso di richiesta]

Questo metodo è usato da ServerManager per default. Per sostituirlo, imposta un gestore
per questo codice con il metodo EnlistStatusReply.
*/
func DefaultBadRequestReply(w http.ResponseWriter, r *http.Request) {
	WriteStatus(http.StatusBadRequest, "Richiesta non valida.", r.URL.Path, w)
}

/*
DefaultNotFoundReply invia la risposta di default per il codice 404.

  404 - Pagina non trovata. - Path: [percorso di richiesta]

Questo metodo è usato da ServerManager per default. Per sostituirlo, imposta un gestore
per questo codice con il metodo EnlistStatusReply.
*/
func DefaultNotFoundReply(w http.ResponseWriter, r *http.Request) {
	WriteStatus(http.StatusNotFound, "Pagina non trovata.", r.URL.Path, w)
}

/*
DefaultServerErrorReply invia la risposta di default per il codice 500.

  500 - Errore interno del server. - Path: [percorso di richiesta]

Questo metodo è usato da ServerManager per default. Per sostituirlo, imposta un gestore
per questo codice con il metodo EnlistStatusReply.
*/
func DefaultServerErrorReply(w http.ResponseWriter, r *http.Request) {
	WriteStatus(http.StatusInternalServerError, "Errore interno del server.", r.URL.Path, w)
}

/*
CheckMethod verifica se il metodo di una richiesta rientra fra quelli ammessi.

Il parametro replyNotAllowed indica se replicare con 405 MethodNotAllowed quando il metodo differisce.
*/
func CheckMethod(r *http.Request, allowedMethods []string, replyNotAllowed bool, w http.ResponseWriter) bool {
	for _, method := range allowedMethods {
		if r.Method == method {
			return true
		}
	}

	if replyNotAllowed {
		allowedList := strings.Join(allowedMethods, ", ")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Allow", allowedList)
		w.WriteHeader(http.StatusMethodNotAllowed)
		var reply string
		reply = fmt.Sprintf("%d - %s invece di %s", http.StatusMethodNotAllowed, r.Method, allowedList)
		fmt.Fprintf(w, reply)
	}
	return false
}

/*
CheckAccept verifica se la richiesta accetta uno dei mimetype specificati.

Il parametro replyNotAcceptable indica se replicare con 406 NotAcceptable quando la richiesta non accetta nessuno dei mimetype.
*/
func CheckAccept(r *http.Request, allowedMime []string, replyNotAcceptable bool, w http.ResponseWriter) bool {
	accepted := r.Header.Get("Accept")
	for _, mime := range allowedMime {
		if strings.Contains(accepted, mime) {
			return true
		}
	}

	if replyNotAcceptable {
		allowedList := strings.Join(allowedMime, ", ")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotAcceptable)
		var reply string
		reply = fmt.Sprintf("%d - la richiesta non accetta il formato della risposta: %s", http.StatusNotAcceptable, allowedList)
		fmt.Fprintf(w, reply)
	}
	return false
}

/*
ServeJSON risponde alla richiesta specificata con i dati nel parametro reply e il codice di stato specificato.

Nell'intestazione della risposta, il valore "Content-Type" è impostato su "application/json; charset=utf-8".

Il contenuto di reply è scritto nella risposta direttamente se di tipo []byte o string,
altrimenti è prima serializzato in formato JSON con la funzione json.Marshal.
La risposta sarà vuota se la funzione json.Marshal restituisce un errore.
*/
func ServeJSON(r *http.Request, reply interface{}, code int, w http.ResponseWriter) {
	if !CheckAccept(r, []string{"*/*", "application/*", "application/json"}, true, w) {
		return
	}

	var dati []byte
	switch rv := reply.(type) {
	case []byte:
		dati = rv
	case string:
		dati = []byte(rv)
	default:
		var err error
		dati, err = json.Marshal(reply)
		if err != nil {
			WriteStatus(http.StatusInternalServerError, err.Error(), r.URL.Path, w)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(dati)
}
