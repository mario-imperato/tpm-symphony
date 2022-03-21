## TPM-Symphony

Obiettivo del progetto la realizzazione di un componente in grado di orchestrare le invocaioni a servizi Rest/Json ed esporre degli endpoint invocabili dall'esterno per l'attivazione di queste orchestrazioni.
Il progetto ha un obiettivo di genericità e non vuole essere legato ad uno specifico scenario ma soddisfare, entro certi limiti, un numero diverso di use-case.
Nella sua implementazione corrente ha l'obiettivo di coprire tutti i requisiti espressi in __come si chiama l'ambito attuale__ attraverso un approccio molto configurabile. Le esigenze espresse
prevedono la necessità di 

- invocare un certo numero di microservizi in serie ed in parallelo (per esigenze di latenza complessiva evidentemente)
- usare le response di alcuni per costruire le request dei successivi
- decidere le chiamate effettive e l'output in base ai parametri in ingresso e alle risposte ottenute.

Il diagramma successivo mostra, nell'ambito del contesto corrente, come il microservizio (a sx) debba invocare in  parallelo 5 microservizi e poi eseguire successivamente le invocazioni di due microservizi ulteriori.
Dalla figura, si evidenzia come alcuni di questi microservizi si appoggino su procedure mainframe e siano sostanzialmente dei wrapper.

![Processo di verifica ricarica PostePay](REAME-files/processo-verifica-ricarica.png)

In alcuni degli scenari in ambito sono previste delle operazioni __dispositive__. Dal punto di vista dell'orchestrazione nulla cambia in quanto, anche gli elementi dispositivi saranno __wrappati__ e quindi riconducibili 
all'interno di chiamate ad endpoint rest/json. Eventuali esigenze di transattività saranno gestite all'interno del singolo __wrapper__ come si chiarisce nella figura successiva nella quale un componente con
un JBoss containerizzato chiama tre procedure mainframe in modalità XA.

![Processo ricarica PostePay](REAME-files/processo-dispositiva.png)

### Architettura

![Symphony architecture](REAME-files/symphony-architecture.png)

L'obiettivo del progetto è dunque quello di  mettere a disposizione un artefatto che possa adattarsi a use-case diversi dai processi in ambito mediante operazioni di configurazione.
Da un punto di vista generale, il componente permette di esporre endpoint a partire da una descrizione degli stessi e di __eseguire__ un grafo orientato aciclico nel quale:

- i nodi sono insiemi di endpoint da invocare (se piu' di uno nello stesso nodo, in parallelo)
- gli archi permettono di transire da un nodo al successivo. 

Gli archi uscenti da un nodo possono avere molteplicità superiore ad uno ma, in ogni caso, solamente uno sarà percorribile. Non sono previsti __parallel gateway__ ma la scelta di uno degli archi uscenti sarà
determinata da espressioni booleane calcolate sullo stato della esecuzione. Solo un arco potrà essere percorso.

La figura precedente mostra quattro sotto componenti principali

- Orchestration registry: questo componente si occupa di caricare la descrizione delle orchestrazioni a partire da una qualche forma di repository che viene 
acceduto da un crawler (es. cartelle messe a disposizione del microservizio piuttosto che repo git-hub).
- API: tra le varie informazioni configurate e disponibili nella registry si troverà anche la descrizione del/degli enpoint che debbono essere esposti tramite API.
- Execution: all'invocazione di uno degli endpoint esposti il componente provvederà ad attraversare il grafo che costituisce l'orchestrazione invocando i vari servizi censiti nelle informazioni di configurazione; l'esecuzione termina con un nodo foglia 
(senza piu' archi uscenti) e l'output verrà restituito in uscita.
- Observability: a parte le normali informazioni di metriche e tracce previste, il sistema loggerà anche l'esecuzione con un livello informativo simile a quello previsto per il formato [HAR](https://en.wikipedia.org/wiki/HAR_(file_format)).