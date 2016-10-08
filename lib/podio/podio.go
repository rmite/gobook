// Copyright (c) 2016 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questa libreria è descritta nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"
// per un esempio di test del codice Go.

// Il package podio implementa la funzione DividiPunti per sfide multiple.
package podio

import (
	"errors"
)

// PuntiPodio contiene i punti per i primi tre classificati.
type PuntiPodio struct {
	Primo   int
	Secondo int
	Terzo   int
}

// DividiPunti divide i punti fra i primi tre classificati.
func DividiPunti(punti int) (p PuntiPodio, err error) {
	if punti <= 0 {
		err = errors.New("punti deve essere maggiore di zero")
		return
	}

	resto := (punti % 6)
	p.Terzo = (punti - resto) / 6
	p.Secondo = (p.Terzo * 2)
	p.Primo = (p.Terzo * 3)

	if resto > 0 {
		switch {
		case resto < 4:
			p.Primo += resto
		case resto == 4:
			p.Primo += 3
			p.Secondo++
		case resto == 5:
			p.Primo += 3
			p.Secondo += 2
		}
	}
	return
}
