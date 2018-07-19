// Copyright (c) 2018 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questa applicazione web è descritta nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"

package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rmite/todo"
	web "rmite/webman"
)

//NotaAPI rappresenta una nota per le api.
type NotaAPI struct {
	ID     int64  `json:"id"`
	Testo  string `json:"nota"`
	Fatto  bool   `json:"fatto"`
	Valida bool   `json:"valida"`
}

//RisultatoAPI descrive il risultato di un'operazione via api.
type RisultatoAPI struct {
	OK        bool   `json:"ok"`
	Messaggio string `json:"msg"`
}

//leggiNotaAPI legge la nota in formato json contenuta nel corpo di una richiesta.
func leggiNotaAPI(r *http.Request) (napi NotaAPI) {
	if (r.Body != nil) && (r.ContentLength > 0) {
		dati := make([]byte, r.ContentLength)
		if _, err := r.Body.Read(dati); err != nil {
			if json.Unmarshal(dati, &napi) != nil {
				napi.Valida = false
			}
		}
	}
	return
}

//apiMostraNota restituisce i dati di una nota.
func apiMostraNota(w http.ResponseWriter, r *http.Request) {
	if !web.CheckMethod(r, []string{http.MethodGet}, true, w) {
		return
	}

	idstr := r.FormValue("id")

	var nt *todo.Nota = &todo.Nota{}

	if id, err := strconv.ParseInt(idstr, 10, 64); err == nil {
		nt, _ = gn.Recupera(id)
	}

	napi := NotaAPI{ID: nt.GetID(), Testo: nt.GetTesto(), Fatto: nt.Fatto, Valida: nt.Valida()}
	web.ServeJSON(r, napi, http.StatusOK, w)
}
