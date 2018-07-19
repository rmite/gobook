// Copyright (c) 2018 Renato Mite. Tutti i diritti riservati. All rights reserved.

//mostraInfoNota richiama l'API Mostra Nota per avere le informazioni sulla nota con id specificato.
function mostraInfoNota(id) {
   // crea la richiesta
   var req = new XMLHttpRequest();
   // imposta la funzione di analisi della risposta
   req.onreadystatechange = function() {
      if (this.readyState == 4) {
         //imposta info su messaggio di errore default
         var info = "Operazione non riuscita - Risposta Server: " + this.status;
         if (this.status == 200) {
            //operazione riuscita, decodifica la risposta in JSON per estrarre i dati
            try {
               var ntAPI = JSON.parse(this.responseText);
               if (ntAPI.valida == true) {
                  //nota valida, inserisce in info i dati della nota
                  info = ntAPI.id + ") " + ntAPI.nota + " [" + ntAPI.fatto + "]";
               }
            }
            catch (err) {
               //decodifica JSON fallita, imposta info su un messaggio di errore
               info = "Nessuna informazione disponibile.";
            }
         }
         //mostra le informazioni o il messaggio di errore
         alert(info);
      }
   };
   // imposta la richiesta
   req.open("GET", "/api/mostra/nota?id="+id, true);
   req.setRequestHeader("Accept", "application/json");
   // invia la richiesta
   req.send();
}


//analizzaRisposta analizza la risposta ricevuta da una richiesta per riportare eventuali errori.
//Verifica il formato della risposta e mostra all'utente il messaggio di errore in JSON,
//oppure testo e codice di stato se diverso da OK (200) per altri formati.
function analizzaRisposta(req) {
   var info = "";
   //recupera il formato della risposta
   var ct = req.getResponseHeader("Content-Type").split(";",1)[0];
   //verifica il formato
   switch (ct) {
      case "application/json":
         //formato JSON, cerca di estrarre i dati di un oggetto RisultatoAPI
         try {
            var risAPI = JSON.parse(req.responseText);
            if (risAPI.ok == false) {
               //operazione non riuscita
               info = risAPI.msg;
               if (info.length < 1) {
                  //messaggio vuoto, imposta un messaggio di errore di default
                  info = "Operazione non riuscita. Risposta Server: " + req.status;
               }
            }
         }
         catch(err) {
            //decodifica JSON fallita
            if (req.status != 200) {
               //stato diverso da OK, imposta un messaggio di errore con il testo della risposta
               info = "Risposta Server: " + req.responseText + " - Errore client: " + err.message;
            }
         }
         break;
      case "text/plain":
         //formato testo
         if (req.status != 200) {
            //stato diverso da OK, imposta un messaggio di errore con il testo della risposta
            info = "Risposta Server: " + req.responseText;
         }
         break;
      default:
         //altri formati
         if (req.status != 200) {
            //stato diverso da OK, imposta un messaggio di errore con il codice di stato
            info = "Risposta Server: " + req.status;
         }
   }
   //restituisce il messaggio di errore
   return info;
}


//cambiaTestoNota permette all'utente di modificare il testo di una nota.
//link rappresenta il link che racchiude il testo della nota, id e fatto sono i dati della nota
function cambiaTestoNota(link, id, fatto) {
   //chiede il nuovo testo all'utente
   var txt = prompt("Inserisci il testo della nota:", link.innerText);
   if (txt == null) {return;}
   // crea la richiesta
   var req = new XMLHttpRequest();
   // imposta la funzione di analisi della risposta
   req.onreadystatechange = function() {
      if (this.readyState == 4) {
         var info;
         //analizza la risposta
         info = analizzaRisposta(req);
         if (info.length > 0) {
            //mostra il messaggio di errore
            alert(info);
         } else {
            //sostituisce il testo nel link con il nuovo testo
            link.innerText = txt;
         }
      }
   };
   //crea l'oggetto napi da codificare in JSON per la richiesta
   var napi = {"id":id, "nota":txt, "fatto":fatto, "valida":true};
   //imposta la richiesta che invia i dati al percorso di modifica
   req.open("POST", "/aggiorna", true);
   req.setRequestHeader("Content-Type", "application/json");
   req.setRequestHeader("Accept", "application/json");
   //invia la richiesta con l'oggetto napi codificato
   req.send(JSON.stringify(napi));
}


//cambiaStatoNota permette all'utente di modificare lo stato di una nota.
//link rappresenta il link che racchiude l'immagine dello stato Fatto\Non Fatto, id e fatto sono i dati della nota
function cambiaStatoNota(link, id, fatto) {
   //recupera il persorso del link
   var path = link.href;
   var imgHTML;
   //verifica lo stato da impostare
   if (path.endsWith("true")) {
      //imposta true, la nuova immagine e il percorso cambieranno in Non Fatto
      imgHTML = '<img class="icon" alt="Cambia in Non Fatto" title="Cambia in Non Fatto" src="/img/fatto.png">';
      path = path.substr(0, path.length-4) + "false";
   } else if (path.endsWith("false")) {
      //imposta false, la nuova immagine e il percorso cambieranno in Fatto
      imgHTML = '<img class="icon" alt="Cambia in Fatto" title="Cambia in Fatto" src="/img/non-fatto.png">';
      path = path.substr(0, path.length-5) + "true";
   } else {
      return; //percorso non valido
   }
   // crea la richiesta
   var req = new XMLHttpRequest();
   // imposta la funzione di analisi della risposta
   req.onreadystatechange = function() {
      if (this.readyState == 4) {
         var info;
         //analizza la risposta
         info = analizzaRisposta(req);
         if (info.length > 0) {
            //mostra il messaggio di errore
            alert(info);
         } else {
            //sostituisce il percorso del link e l'immagine
            link.href = path;
            link.innerHTML = imgHTML;
         }
      }
   };
   //imposta la richiesta che punta al percorso di cambio stato
   req.open("POST", link.href, true);
   req.setRequestHeader("Content-Type", "application/json");
   req.setRequestHeader("Accept", "application/json");
   //invia la richiesta
   req.send();
}
