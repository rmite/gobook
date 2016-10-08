// Copyright (c) 2016 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questo programma è descritto nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"
// come esercizio di scrittura del codice Go.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	soc "rmite/gobook/lib/socutil"
	su "rmite/gobook/lib/strutil"
	"sort"
	"strings"
)

const numGironi int = 2
const numSquadreGirone int = 4
const maxSquadre = (numGironi * numSquadreGirone)
const maxPartiteGiorno = (numGironi * 2)

var squadre [numGironi][numSquadreGirone]string
var partite soc.ListaPartite

func main() {

	fmt.Println("Benvenuto in campo con Go!")

	var inputConsole bool = true
	var fileSquadre string

	if len(os.Args) > 1 {
		fileSquadre = os.Args[1]

		if _, err := os.Stat(fileSquadre); os.IsNotExist(err) {
			fmt.Printf("Il file '%s' non esiste.\n", fileSquadre)
		} else {
			if err = caricaSquadre(fileSquadre); err != nil {
				fmt.Printf("Impossibile leggere dal file '%s'.\n", fileSquadre)
				fmt.Println(err)
			} else {
				fmt.Println("Squadre caricate.")
				inputConsole = false
			}
		}

	}

	// Acquisisce i nomi delle squadre dalla console.
	if inputConsole {
		getSquadre()
		salvaSquadre()
	}

	mostraSquadre()

	creaPartite()

	mostraPartite()

	salvaPartite()

}

// ----- FUNZIONI Squadre -----

func mostraSquadre() {
	fmt.Println("\n=========== SQUADRE ============")
	for g := 0; g < numGironi; g++ {
		for s := 0; s < numSquadreGirone; s++ {
			fmt.Println(squadre[g][s])
		}
	}
	fmt.Println("================================")
}

func getSquadre() {
	// Acquisisce le squadre con l'input.
	var trovato bool
	var nome string
	var lower string
	var tmp map[string]string

	tmp = make(map[string]string)

	for g := 1; g <= numGironi; g++ {
		fmt.Printf("\n--- Squadre GIRONE %d ---\n", g)
		for s := 1; s <= numSquadreGirone; s++ {
			trovato = true
			for trovato {
				nome = su.GetStr(fmt.Sprintf("Inserisci il nome della squadra %d", s), su.ValidateStrNoBlank)
				lower = strings.ToLower(nome)
				if _, trovato = tmp[lower]; trovato {
					fmt.Printf("La squadra %s è stata già inserita!\n", nome)
				} else {
					tmp[lower] = nome
					squadre[g-1][s-1] = nome
				}
			}
		}
	}
}

func caricaSquadre(filePath string) error {
	// Acquisisce le squadre da file.
	var f *os.File
	var err error

	// apre il file
	f, err = os.Open(filePath)
	if err != nil {
		// esce per errori durante l'apertura, es. il file non esiste
		return err
	}

	defer f.Close() // differisce la chiusura del file

	// crea map squadre
	var tmp map[string]string
	tmp = make(map[string]string)

	// crea un reader
	r := bufio.NewReader(f)

	// legge i nomi delle squadre
	var trovato bool
	var ln, nome, lower string
	var count = 0

	for g := 0; g < numGironi; g++ {
		for s := 0; s < numSquadreGirone; s++ {

			ln, err = r.ReadString('\n')
			if err != nil {
				// esce dalla funzione se non è possibile leggere righe dal file
				return err
			}

			nome = strings.TrimSuffix(ln, "\r\n")
			if su.IsBlank(nome) {
				continue // ignora le righe vuote
			}

			lower = strings.ToLower(nome)
			if _, trovato = tmp[lower]; !trovato {
				tmp[lower] = nome
				squadre[g][s] = nome
				count++
			}
		}
	}

	if count != maxSquadre {
		return errors.New("Il file non contiene il numero esatto di squadre.")
	}
	return nil
}

func salvaSquadre() {
	// Salva le squadre in un file.
	var f *os.File
	var err error

	// crea il file
	f, err = os.Create("squadre.txt")
	if mostraErrore(err) {
		// esce per errori durante la creazione
		return
	}

	defer f.Close() // differisce la chiusura del file

	// crea un writer
	w := bufio.NewWriter(f)

	// scrive i nomi delle squadre
	errCount := 0
	for g := 0; g < numGironi; g++ {
		for s := 0; s < numSquadreGirone; s++ {
			_, err = w.WriteString(fmt.Sprintf("%s\r\n", squadre[g][s]))
			if mostraErrore(err) {
				errCount++
				break
			}
		}
		w.Flush()
	}

	if errCount == 0 {
		fmt.Println("Squadre salvate nel file squadre.txt")
	}
}

// ----- FUNZIONI Partite -----

func mostraPartite() {
	fmt.Println("\n====== CALENDARIO PARTITE ======")
	for _, p := range partite {
		fmt.Println(p)
	}
	fmt.Println("================================")
}

func validateMaxPG(n int) error {
	if n < 1 || n > maxPartiteGiorno || (n != 1 && (n%2) != 0) {
		return fmt.Errorf("Numero partite al giorno non valido. Deve essere 1 oppure un numero pari fino a %d.", maxPartiteGiorno)
	}
	return nil
}

func creaPartite() {
	var partiteAlGiorno int
	partiteAlGiorno = su.GetInt("Inserisci il numero di partite al giorno:", validateMaxPG)

	// Crea le partite dei gironi
	var g int
	var nomeA, nomeB string
	var pID string
	var p soc.Partita
	var err error

	partite = make(soc.ListaPartite, 0)

	// Ciclo dei gironi
	for g = 0; g < numGironi; g++ {
		fmt.Printf("\n----- Partite Girone %d -----\n", g+1)

		// Ciclo della prima squadra fino a numSquadreGirone-1 perché deve esserci almeno una seconda squadra
		for s := 0; s < (numSquadreGirone - 1); s++ {
			nomeA = squadre[g][s]

			// Ciclo della seconda squadra da s+1 perché le combinazioni sono solo con squadre successive
			for o := (s + 1); o < numSquadreGirone; o++ {
				nomeB = squadre[g][o]
				pID = fmt.Sprintf("g%d-%d-%d", g, s, o)
				p, err = soc.NewPartita(pID, nomeA, nomeB)
				if err == nil {
					partite = append(partite, p)
					fmt.Println(p.ShortString())
				} else {
					fmt.Println(err)
					return
				}
			}
		}
	}
	fmt.Println("--------------------------")

	// Imposta la giornata in cui si gioca
	var totComb int
	var i, n int
	var counter = 0 // contatore partite
	var giorno = 1  // contatore giorni

	// Calcola le combinazioni in un girone
	totComb = (numSquadreGirone * (numSquadreGirone - 1)) / 2

	for comb := 1; comb <= totComb; comb++ {
		for g = 0; g < numGironi; g++ {

			i = (g * totComb) // Indice della prima partita del girone
			// Alterna le partite in modo che non giochino le stesse squadre
			if (comb % 2) == 0 {
				//comb pari, prende partite dalla fine dell'elenco
				n = (comb / 2)
				i += (totComb - n)
			} else {
				//comb dispari, prende partite dall'inizio dell'elenco
				n = (comb - 1) / 2
				i += n
			}

			partite[i].Giornata = giorno

			// Incrementa il contatore delle partite
			counter++
			if (counter % partiteAlGiorno) == 0 {
				giorno++ // Incrementa il giorno
			}

		}
	}

	// Ordina le partite in base alla giornata
	sort.Sort(partite)
}

func salvaPartite() {
	// Salva le partite in un file.
	var f *os.File
	var err error

	// crea il file
	f, err = os.Create("partite.txt")
	if mostraErrore(err) {
		// esce per errori durante la creazione
		return
	}

	defer f.Close() // differisce la chiusura del file

	// crea un writer
	w := bufio.NewWriter(f)

	// scrive le partite
	errCount := 0
	for _, p := range partite {
		_, err = w.WriteString(fmt.Sprintf("%s\r\n", p))
		if mostraErrore(err) {
			errCount++
			break
		}
	}
	w.Flush()

	if errCount == 0 {
		fmt.Println("Partite salvate nel file partite.txt")
	}
}

// ----- FUNZIONI Errori -----

func mostraErrore(e error) bool {
	if e != nil {
		fmt.Println(e)
		return true
	}
	return false
}
