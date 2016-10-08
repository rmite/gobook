// Copyright (c) 2016 Renato Mite. Tutti i diritti riservati. All rights reserved.
// Questa libreria è descritta nella guida
// "Programmare in Linguaggio Go - La guida italiana per muovere i primi passi"
// come esercizio di scrittura del codice Go.

/*
Il package strutil implementa funzioni per l'acquisizione e
la convalida dell'input dalla console.

Funzioni di convalida

Le funzioni di convalida che soddisfano i tipi
  type ValidateStr func(s string) error
  type ValidateInt func(n int) error
devono restituire un errore in caso di fallimento della convalida altrimenti nil.
*/
package strutil

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//ValidateStr è il tipo funzione per convalidare le stringhe.
type ValidateStr func(s string) error

//ValidateInt è il tipo funzione per convalidare i valori int.
type ValidateInt func(n int) error

//GetStr acquisisce una stringa dalla console.
func GetStr(inputMsg string, validate ValidateStr) string {
	var s string
	var err error

	reader := bufio.NewReader(os.Stdin)

	if validate == nil {
		validate = ValidateStrAll
	}

	for {
		fmt.Println(inputMsg)
		s, err = reader.ReadString('\n')
		if err == nil {
			s = strings.TrimRight(s, "\r\n")
			if err = validate(s); err == nil {
				return s
			}
		}
		fmt.Println(err)
	}
}

//GetInt acquisisce un int dalla console.
func GetInt(inputMsg string, validate ValidateInt) int {
	var s string
	var v int
	var err error

	reader := bufio.NewReader(os.Stdin)

	if validate == nil {
		validate = ValidateIntAll
	}

	for {
		fmt.Println(inputMsg)
		s, err = reader.ReadString('\n')
		if err == nil {
			s = strings.TrimRight(s, "\r\n")
			v, err = strconv.Atoi(s)
			if err == nil {
				if err = validate(v); err == nil {
					return v
				}
			} else {
				err = errors.New("Inserisci un numero intero.")
			}
		}
		fmt.Println(err)
	}
}

//IsBlank restituisce true se la stringa è vuota o solo spazi.
func IsBlank(s string) bool {
	return (len(strings.TrimSpace(s)) == 0)
}

//ValidateStrAll convalida qualsiasi stringa.
func ValidateStrAll(s string) error {
	return nil
}

//ValidateStrNoBlank convalida le stringhe che non sono vuote.
func ValidateStrNoBlank(s string) error {
	if IsBlank(s) {
		return errors.New("La stringa non può essere vuota.")
	}
	return nil
}

//ValidateIntAll convalida qualsiasi int.
func ValidateIntAll(n int) error {
	return nil
}

//ValidateIntGreaterZero convalida qualsiasi int maggiore di zero.
func ValidateIntGreaterZero(n int) error {
	if n < 1 {
		return errors.New("Il numero deve essere maggiore di zero.")
	}
	return nil
}
