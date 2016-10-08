// Copyright (c) 2016 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questa libreria è descritta nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"
// per un esempio di test del codice Go.

package podio

import (
	"fmt"
	"testing"
)

type datiPodio struct {
	punti        int
	podio        PuntiPodio
	generaErrore bool
}

var dtPodio = []datiPodio{
	{-6, PuntiPodio{0, 0, 0}, true},
	{0, PuntiPodio{0, 0, 0}, true},
	{1, PuntiPodio{1, 0, 0}, false},
	{2, PuntiPodio{2, 0, 0}, false},
	{3, PuntiPodio{3, 0, 0}, false},
	{4, PuntiPodio{3, 1, 0}, false},
	{5, PuntiPodio{3, 2, 0}, false},
	{6, PuntiPodio{3, 2, 1}, false},
	{7, PuntiPodio{4, 2, 1}, false},
	{12, PuntiPodio{6, 4, 2}, false}}

var risPodio PuntiPodio
var err error

func TestDividiPunti(t *testing.T) {

	for _, dp := range dtPodio {

		risPodio, err = DividiPunti(dp.punti)

		switch {
		case err != nil:
			if dp.generaErrore {
				t.Logf("MSG : Errore previsto. La divisione di %d punti genera l'errore '%v' \n", dp.punti, err)
			} else {
				t.Errorf("ERR : Errore non previsto. La divisione di %d punti genera l'errore '%v' quando non dovrebbe \n", dp.punti, err)
			}

		case dp.generaErrore:
			t.Errorf("ERR : La divisione di %d punti dovrebbe generare un errore e non lo fa \n", dp.punti)

		case risPodio != dp.podio:
			t.Errorf("ERR : La divisione di %d punti produce %v invece di %v \n", dp.punti, risPodio, dp.podio)

		default:
			t.Logf("MSG : Divisione di %d punti eseguita con successo: %v \n", dp.punti, risPodio)
		}

	}

}

func TestPropPodio(t *testing.T) {

	for p := 1; p < 101; p++ {

		risPodio, err = DividiPunti(p)

		switch {
		case err != nil:
			t.Errorf("ERR : Errore non previsto. La divisione di %d punti genera l'errore '%v' quando non dovrebbe \n", p, err)

		case (risPodio.Terzo * 2) > risPodio.Secondo:
			t.Errorf("ERR : La divisione di %d punti produce %v con proporzione errata fra secondo e terzo posto \n", p, risPodio)

		case (risPodio.Terzo * 3) > risPodio.Primo:
			t.Errorf("ERR : La divisione di %d punti produce %v con proporzione errata fra primo e terzo posto \n", p, risPodio)

		case (risPodio.Secondo * 3) > (risPodio.Primo * 2):
			t.Errorf("ERR : La divisione di %d punti produce %v con proporzione errata fra primo e secondo posto \n", p, risPodio)

		case (risPodio.Primo + risPodio.Secondo + risPodio.Terzo) != p:
			t.Errorf("ERR : La divisione di %d punti produce %v e la somma dei punti non corrisponde \n", p, risPodio)

		default:
			t.Logf("MSG : Divisione di %d punti eseguita con successo: %v \n", p, risPodio)
		}

	}

}

// Examples

func ExampleDividiPunti_zero() {
	_, err = DividiPunti(0)
	fmt.Println(err)
	// Output:
	// punti deve essere maggiore di zero
}

func ExampleDividiPunti_neg() {
	_, err = DividiPunti(-1)
	fmt.Println(err)
	// Output:
	// punti deve essere maggiore di zero
}

// Benchmarks

func BenchmarkDividiPunti(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DividiPunti(100)
	}
}
