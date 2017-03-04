
## Codice della guida "Programmare in Linguaggio Go"


[![Go Report Card](https://goreportcard.com/badge/github.com/rmite/gobook)](https://goreportcard.com/report/github.com/rmite/gobook) &nbsp; [![GoDoc](https://godoc.org/github.com/rmite/gobook/lib/strutil?status.svg)](https://godoc.org/github.com/rmite/gobook/lib/strutil)


Questo repository contiene il **codice di esempio** descritto nella guida **"Programmare in Linguaggio Go"**.

<p align="center"><img src="go-guide-cover.jpg" /></p>

La guida **"Programmare in Linguaggio Go"** è disponibile in ebook su [Amazon][guide-amazon], [GooglePlay][guide-goplay] e [Kobo][guide-kobo], in libro cartaceo su [Amazon][guide-amazon].

Go è un linguaggio dai notevoli [punti di forza][go-strengths] e la guida è pensata per tutti coloro che vogliono apprendere questo nuovo linguaggio di programmazione, siano aspiranti programmatori o programmatori esperti. Descrive **le specifiche e le peculiarità di Go** in modo graduale, in un percorso che va dai valori ai packages, passando per array, slice, map, puntatori, goroutines e canali.

Il capitolo **"Go all'opera"** accompagna il lettore nella scrittura del programma [soccer][code-program] e delle librerie [strutil][code-strutil] e [socutil][code-socutil] contenuti in questo repository. Questi file di codice mostrano come avvalersi delle **librerie standard Go** per conoscere gli argomenti della linea di comando, acquisire l'input dalla console, convalidare i dati, leggere e scrivere stringhe su flussi e files, verificare se un file esiste, e ancora come definire strutture con relativi metodi e come implementare i metodi delle interfacce standard per la rappresentazione testuale nella console e per l'ordinamento degli slice con criteri personali.
Il programma _soccer_ non è un semplice "Hello World!", è un programma più articolato le cui istruzioni, insieme alle due librerie, sono spiegate nella guida passo passo, approfondendo allo stesso tempo i temi della progettazione e dell'ottimizzazione del codice.

Il capitolo **"Rudimenti di test"** descrive sia i files [podio.go][test-lib] e [podio_test.go][test-file] contenuti nella cartella podio di questo repository sia le basi per scrivere e eseguire test.

## Usare il codice

Per visionare il codice durante la lettura della guida, dalla console **clona il repository nella cartella src del tuo workspace**:

```
$ cd %GOPATH%\src
$ git clone https://github.com/rmite/gobook.git rmite
```

e per eseguirlo spostati nella sottocartella soccer, importa i packages e crea l'eseguibile

```
$ cd rmite\soccer\
$ go get && go build
```

che puoi avviare sempre da console:

(Windows)
```
$ soccer.exe
```

(Linux)
```
$ ./soccer
```

Vedi anche [License][license] e [Contributing][contribute].

[guide-cover]: go-guide-cover.jpg
[guide-amazon]: https://www.amazon.it/dp/B01M2URIVX
[guide-goplay]: https://play.google.com/store/books/details/Renato_Mite_Programmare_in_Linguaggio_Go?id=4Ag6DQAAQBAJ
[guide-kobo]: https://store.kobobooks.com/it-it/ebook/programmare-in-linguaggio-go
[code-program]: soccer/soccer.go
[code-strutil]: lib/strutil/strutil.go
[code-socutil]: lib/socutil/partita.go
[test-lib]: lib/podio/podio.go
[test-file]: lib/podio/podio_test.go
[license]: LICENSE.md
[contribute]: CONTRIBUTING.md
[go-strengths]: https://medium.com/@renato.mite/punti-di-forza-del-linguaggio-go-2905d698740e
