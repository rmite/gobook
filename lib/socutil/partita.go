// Copyright (c) 2016 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questa libreria è descritta nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"
// come esercizio di scrittura del codice Go.

// Il package socutil implementa tipi e funzioni per la gestione di partite di calcio.
package socutil

import (
	"errors"
	"fmt"
	su "rmite/gobook/lib/strutil"
)

//Partita rappresenta una partita.
type Partita struct {
	id       string
	Giornata int
	squadraA string
	squadraB string
}

//ListaPartite è un elenco di partite.
type ListaPartite []Partita

//NewPartita restituisce un nuovo oggetto Partita.
func NewPartita(id, squadraA, squadraB string) (p Partita, err error) {
	p = Partita{id: "", squadraA: "", squadraB: ""}
	if su.IsBlank(id) {
		err = errors.New("L'ID partita non è valido.")
		return
	}
	if su.IsBlank(squadraA) {
		err = errors.New("La squadra A non è valida.")
		return
	}
	if su.IsBlank(squadraB) {
		err = errors.New("La squadra B non è valida.")
		return
	}
	p.id = id
	p.squadraA = squadraA
	p.squadraB = squadraB
	err = nil
	return
}

//IsValid indica se un oggetto Partita è valido.
func (p *Partita) IsValid() bool {
	return (len(p.id) > 0)
}

//ID restituisce l'identificativo della partita.
func (p *Partita) ID() string {
	return p.id
}

//SquadraA restituisce il nome della squadra A.
func (p *Partita) SquadraA() string {
	return p.squadraA
}

//SquadraB restituisce il nome della squadra B.
func (p *Partita) SquadraB() string {
	return p.squadraB
}

//ShortString restituisce la stringa breve che descrive la partita.
func (p Partita) ShortString() string {
	return fmt.Sprintf("%s - %s", p.squadraA, p.squadraB)
}

// ----- FUNZIONI INTERFACCIA Stringer -----

//String restituisce la stringa che descrive la partita.
func (p Partita) String() string {
	return fmt.Sprintf("Giorno %d) %s - %s", p.Giornata, p.squadraA, p.squadraB)
}

// ----- FUNZIONI INTERFACCIA Sort -----

func (lp ListaPartite) Len() int {
	return len(lp)
}

func (lp ListaPartite) Less(i, j int) bool {
	return lp[i].Giornata < lp[j].Giornata
}

func (lp ListaPartite) Swap(i, j int) {
	lp[i], lp[j] = lp[j], lp[i]
}
