// Copyright (c) 2018 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questa libreria è descritta nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"

// Il package todo implementa tipi e funzioni per gestire note testuali con un database SQLite.
package todo

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	//inizializza il driver sqlite3
	_ "github.com/mattn/go-sqlite3"
)

//Nota rappresenta una nota.
type Nota struct {
	id    int64
	testo string
	Fatto bool
}

//GetID restituisce l'id della nota.
func (nt Nota) GetID() int64 {
	return nt.id
}

//GetTesto restituisce il testo della nota.
func (nt Nota) GetTesto() string {
	return nt.testo
}

//Testo imposta il testo della nota.
func (nt *Nota) Testo(str string) {
	nt.testo = strings.TrimSpace(str)
}

//Valida indica se la nota è valida.
func (nt *Nota) Valida() bool {
	return (len(strings.TrimSpace(nt.testo)) > 0)
}

//String restituisce il testo della nota.
func (nt Nota) String() string {
	return nt.testo
}

//Gestore gestisce le note.
type Gestore struct {
	base *sql.DB
}

//ErrGestoreNonPronto è restituito quando il gestore non è pronto.
var ErrGestoreNonPronto error = errors.New("gestore non pronto")

//ErrNotaNonValida è restituito quando una nota non è valida.
var ErrNotaNonValida error = errors.New("nota non valida")

//ErrNotaNonTrovata è restituito quando una nota non è disponibile.
var ErrNotaNonTrovata error = errors.New("nota non trovata")

const createStmt string = "CREATE TABLE note (id INTEGER PRIMARY KEY ASC AUTOINCREMENT, testo VARCHAR(200) NOT NULL, fatto BOOLEAN NOT NULL);"

//FiltroElenco rappresenta il filtro di selezione delle note.
type FiltroElenco int

const (
	//NessunFiltro non applica filtri all'elenco delle note.
	NessunFiltro FiltroElenco = 0
	//NoteDaFare seleziona le note da fare.
	NoteDaFare FiltroElenco = 1
	//NoteFatte seleziona le note fatte.
	NoteFatte FiltroElenco = 2
)

//Tutte restituisce true se il filtro seleziona tutte le note.
func (f FiltroElenco) Tutte() bool {
	return (f == NessunFiltro) || (f == (NoteDaFare | NoteFatte))
}

//Fatte restituisce true se il filtro seleziona solo le note fatte.
func (f FiltroElenco) Fatte() bool {
	return (f == NoteFatte)
}

//DaFare restituisce true se il filtro seleziona solo le note da fare.
func (f FiltroElenco) DaFare() bool {
	return (f == NoteDaFare)
}

//NewGestore apre o crea un file SQLite in cui salvare le note.
//Se la connessione al database non riesce, il metodo Pronto restituisce false
//e i vari metodi per accedere o modificare le note restituiscono l'errore ErrGestoreNonPronto.
func NewGestore(filePath string) (gn *Gestore, err error) {
	var db *sql.DB
	gn = &Gestore{base: nil}
	createTable := false

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		// il file non esiste
		createTable = true
	}

	// apre il database
	if db, err = sql.Open("sqlite3", filePath); err != nil {
		return
	}

	if createTable {
		// crea la tabella delle note
		if _, err = db.Exec(createStmt); err != nil {
			db.Close()
			return
		}
	}

	err = nil
	gn.base = db
	return
}

//Chiudi chiude il database sottostante se inizializzato.
func (gn *Gestore) Chiudi() {
	if gn.base != nil {
		gn.base.Close()
	}
}

//Pronto restituisce true se il gestore è pronto.
func (gn *Gestore) Pronto() bool {
	return (gn.base != nil)
}

//Elenco restituisce uno slice di note selezionate dal database oppure nil
//se il gestore non è pronto o in caso di errori nell'interrogazione del database.
//Il parametro filtro indica quali note devono essere selezionate.
func (gn *Gestore) Elenco(filtro FiltroElenco) (note []Nota) {
	if !gn.Pronto() {
		return nil
	}

	var slc []interface{}
	var query string
	var rws *sql.Rows
	var err error

	query = " WHERE fatto = ?"
	switch {
	case filtro.Fatte():
		slc = append(slc, true)
	case filtro.DaFare():
		slc = append(slc, false)
	default:
		query = ""
	}

	query = "SELECT id, testo, fatto FROM note" + query + ";"

	if rws, err = gn.base.Query(query, slc...); err != nil {
		return nil
	}

	note = make([]Nota, 0, 5)

	var valID int64
	var valTesto string
	var valFatto bool

	defer rws.Close()

	for rws.Next() {
		if rws.Scan(&valID, &valTesto, &valFatto) == nil {
			note = append(note, Nota{id: valID, testo: valTesto, Fatto: valFatto})
		}
	}

	return
}

//Totale restituisce il numero di note gestite oppure 0 in caso di errori
//nell'interrogazione del database o se il gestore non è pronto.
func (gn *Gestore) Totale(filtro FiltroElenco) (tot int) {
	if !gn.Pronto() {
		return 0
	}

	var slc []interface{}
	var query string

	query = " WHERE fatto = ?"
	switch {
	case filtro.Fatte():
		slc = append(slc, true)
	case filtro.DaFare():
		slc = append(slc, false)
	default:
		query = ""
	}

	query = "SELECT COUNT(*) FROM note" + query + ";"

	if err := gn.base.QueryRow(query, slc...).Scan(&tot); err != nil {
		tot = 0
	}

	return
}

//Aggiungi inserisce una nuova nota col testo specificato e stato false.
//Se l'inserimento riesce, restituisce l'identificativo numerico della nota e nil.
//Se l'inserimento non riesce, restituisce -1 e l'errore SQL avvenuto oppure l'errore ErrNotaNonTrovata
//se non è stato possibile recuperare l'identificativo della nota dopo l'inserimento.
func (gn *Gestore) Aggiungi(testoNota string) (id int64, err error) {
	id = -1

	if !gn.Pronto() {
		err = ErrGestoreNonPronto
		return
	}

	if len(strings.TrimSpace(testoNota)) == 0 {
		err = ErrNotaNonValida
		return
	}

	var res sql.Result
	if res, err = gn.base.Exec("INSERT INTO note (testo, fatto) values(?, 0);", testoNota); err != nil {
		return
	}
	if id, err = res.LastInsertId(); err != nil {
		id = -1
		err = ErrNotaNonTrovata
	}
	return
}

/*
Aggiorna applica le modifiche a una nota nel database sottostante.
Se l'aggiornamento riesce, restituisce nil.
Negli altri casi, restituisce ErrGestoreNonPronto se il gestore non è pronto,
ErrNotaNonValida se la nota passata non è valida, oppure l'eventuale errore SQL.

La nota deve essere recuperata dal gestore affinché abbia il suo identificativo.

Ad esempio dopo
  n, _ := g.Recupera(1)
puoi modificare testo e stato prima di chiamare Aggiorna:
  n.Testo("Comprare altri 2 litri di latte")
  n.Fatto = false
  g.Aggiorna(n)
*/
func (gn *Gestore) Aggiorna(nt *Nota) (err error) {
	if !gn.Pronto() {
		err = ErrGestoreNonPronto
		return
	}

	if !nt.Valida() {
		err = ErrNotaNonValida
		return
	}

	_, err = gn.base.Exec("UPDATE note SET testo = ?, fatto = ? WHERE id = ?;", nt.testo, nt.Fatto, nt.id)

	return
}

//CambiaStato modifica lo stato di una nota nel database sottostante.
//Se la modifica riesce, restituisce nil.
//Negli altri casi, restituisce ErrGestoreNonPronto se il gestore non è pronto,
//oppure l'eventuale errore SQL.
func (gn *Gestore) CambiaStato(IDNota int64, valoreFatto bool) (err error) {
	if !gn.Pronto() {
		err = ErrGestoreNonPronto
		return
	}

	_, err = gn.base.Exec("UPDATE note SET fatto = :valore WHERE id = :idn AND fatto <> :valore;", sql.Named("valore", valoreFatto), sql.Named("idn", IDNota))

	return
}

//Recupera restituisce la nota con identificativo specificato e nil.
//Se il recupero non riesce, restituisce una nota vuota (non valida)
//e ErrGestoreNonPronto se il gestore non è pronto, l'eventuale errore SQL
//oppure ErrNotaNonTrovata se non c'è una nota con l'identificativo specificato.
func (gn *Gestore) Recupera(IDNota int64) (nt *Nota, err error) {
	nt = &Nota{id: -1, testo: "", Fatto: false}

	if !gn.Pronto() {
		err = ErrGestoreNonPronto
		return
	}

	var valTesto string
	var valFatto bool

	err = gn.base.QueryRow("SELECT testo, fatto FROM note WHERE id = ?", IDNota).Scan(&valTesto, &valFatto)

	if err == nil {
		nt.id = IDNota
		nt.testo = valTesto
		nt.Fatto = valFatto
	}

	if err == sql.ErrNoRows {
		err = ErrNotaNonTrovata
	}

	return
}

//Elimina rimuove la nuova nota con id specificato e restituisce nil in caso di successo.
//Restituisce ErrGestoreNonPronto se il gestore non è pronto,
//altrimenti l'errore SQL se l'eliminazione non riesce.
func (gn *Gestore) Elimina(IDNota int64) (err error) {
	if !gn.Pronto() {
		err = ErrGestoreNonPronto
		return
	}

	_, err = gn.base.Exec("DELETE FROM note WHERE id = ?", IDNota)

	return
}
