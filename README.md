
## Codice della guida "Programmare in Linguaggio Go"


[![Go Report Card](https://goreportcard.com/badge/github.com/rmite/gobook)](https://goreportcard.com/report/github.com/rmite/gobook)


Questo repository contiene il **codice di esempio** descritto nella guida **"Programmare in Linguaggio Go"**.

<p align="center"><img src="go-guide-cover.jpg" /></p>

La guida **"Programmare in Linguaggio Go"** è disponibile in ebook su [Amazon][guide-amazon], [GooglePlay][guide-goplay] e [Kobo][guide-kobo], in libro cartaceo su [Amazon][guide-amazon].

Go è un linguaggio dai notevoli [punti di forza][go-strengths] e la guida è pensata per tutti coloro che vogliono apprendere questo nuovo linguaggio di programmazione, siano aspiranti programmatori o programmatori esperti. Descrive **le specifiche e le peculiarità di Go** in modo graduale, in un percorso che va dai valori ai packages, passando per array, slice, map, puntatori, goroutines e canali.

Il capitolo **"Go all'opera"** accompagna il lettore nella scrittura del programma [soccer][code-program] e delle librerie [strutil][code-strutil] e [socutil][code-socutil] contenuti in questo repository. Questi file di codice mostrano come avvalersi delle **librerie standard Go** per conoscere gli argomenti della linea di comando, acquisire l'input dalla console, convalidare i dati, leggere e scrivere stringhe su flussi e files, verificare se un file esiste, e ancora come definire strutture con relativi metodi e come implementare i metodi delle interfacce standard per la rappresentazione testuale nella console e per l'ordinamento degli slice con criteri personali.
Il programma _soccer_ non è un semplice "Hello World!", è un programma più articolato le cui istruzioni, insieme alle due librerie, sono spiegate nella guida passo passo, approfondendo allo stesso tempo i temi della progettazione e dell'ottimizzazione del codice.

Il capitolo **"Rudimenti di test"** descrive sia i files [podio.go][test-lib] e [podio_test.go][test-file] contenuti nella cartella podio di questo repository sia le basi per scrivere e eseguire test.

## Usare il codice di esempio

[![GoDoc](https://godoc.org/github.com/rmite/gobook/lib/strutil?status.svg)](https://godoc.org/github.com/rmite/gobook/lib/strutil)

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

## Database sql e applicazioni web

La nuova edizione della guida si arricchisce di tre capitoli in cui è descritto come usare il package standard per interagire con database sql e come sviluppare applicazioni web.

Nel repository trovi la libreria todo per gestire note con un database SQLite, la libreria webman per gestire server web e l'applicazione RicordaLista descritte nei rispettivi capitoli.

### La libreria todo

[![GoDoc](https://godoc.org/github.com/rmite/gobook/lib/todo?status.svg)](https://godoc.org/github.com/rmite/gobook/lib/todo)

La libreria [todo][code-todo] usa il **package standard database/sql** e il [driver per SQLite 3 di mattn][mattn-sqlite-driver] disponibile qui su GitHub.

La libreria permette di creare e gestire un database di note con tipi e metodi che non espongono il database.

### La libreria webman

[![GoDoc](https://godoc.org/github.com/rmite/gobook/lib/webman?status.svg)](https://godoc.org/github.com/rmite/gobook/lib/webman)

La libreria [webman][code-webman] raccoglie tipi e funzioni per **gestire le richieste a un server web**.

### L'applicazione RicordaLista

RicordaLista è **un'applicazione** web che fornisce un'interfaccia grafica per gestire un elenco di note.

Le note sono salvate in un database SQLite e gestite con la libreria [todo][code-todo].

<p align="center"><img src="webapp/pubblico/img/titolo.png" /></p>

L'applicazione si avvale dei seguenti **file esterni**:
 * il database delle note, nella stessa cartella
 * modelli html, nella sottocartella "\\privato\\modelli"
 * immagini, nella sottocartella "\\pubblico\\img"
 * stili e script, nella sottocartella "\\pubblico\\files"

In alternativa al modello semplice della homepage "home.html", nel repository c'è il modello "home2.html" insieme al file javascript "apilib.js" che permettono di vedere come l'applicazione risponde a richieste asincrone e API.

Nella **cartella webapp** c'è l'eseguibile dell'applicazione "webapp.exe" per Windows a 64bit.

Per usare questo eseguibile o quello compilato da te, scarica la cartella webapp con i file esterni che servono all'applicazione.

L'applicazione legge il modello "home.html" quindi sarà necessario rinominare i due files della homepage per usare il modello alternativo al posto di quello semplice.

---

Vedi anche [License][license] e [Contributing][contribute].

[guide-cover]: go-guide-cover.jpg
[guide-amazon]: https://www.amazon.it/Programmare-Linguaggio-Go-italiana-muovere-ebook/dp/B01M2URIVX
[guide-goplay]: https://play.google.com/store/books/details/Renato_Mite_Programmare_in_Linguaggio_Go?id=4Ag6DQAAQBAJ
[guide-kobo]: https://www.kobo.com/it/it/ebook/programmare-in-linguaggio-go
[code-program]: soccer/soccer.go
[code-strutil]: lib/strutil/strutil.go
[code-socutil]: lib/socutil/partita.go
[test-lib]: lib/podio/podio.go
[test-file]: lib/podio/podio_test.go
[license]: LICENSE.md
[contribute]: CONTRIBUTING.md
[go-strengths]: https://medium.com/@renato.mite/punti-di-forza-del-linguaggio-go-2905d698740e
[code-todo]: lib/todo/todo.go
[mattn-sqlite-driver]: https://github.com/mattn/go-sqlite3
[code-webman]: lib/webman/webman.go
[webapp-title]: webapp/pubblico/img/titolo.png
