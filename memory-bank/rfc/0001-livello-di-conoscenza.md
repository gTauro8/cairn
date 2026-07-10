# RFC-0001: Perché il software ha bisogno di un livello di conoscenza

*   **Stato:** Draft
*   **Autore:** Gemini CLI (in collaborazione con il team di Cairn)
*   **Data:** 2026-07-10
*   **Area:** Architettura / Filosofia di Prodotto

---

## 1. Stato dell'Arte: Il Divario tra Codice e Conoscenza

In qualsiasi progetto software moderno, la base di codice rappresenta l'artefatto finale e definitivo di ciò che il sistema *fa*. Tuttavia, il codice esprime esclusivamente il **come** (la sintassi, l'algoritmo, l'istruzione di macchina). Il **perché** — i compromessi architetturali accettati, le strade scartate dopo ore di brainstorming, le assunzioni sul business, le scoperte empiriche fatte durante il debug o l'impatto di un limite tecnologico esterno — rimane disperso in una costellazione di canali effimeri.

Oggi, questa conoscenza vitale è frammentata tra:
*   **Sistemi di tracciamento esterni** (Notion, Confluence, Jira, Slack, email): canali distanti dal codice che marciscono rapidamente perché aggiornarli richiede uno sforzo cosciente estraneo al flusso di lavoro quotidiano dello sviluppatore.
*   **Commenti inline e commit log**: sebbene vicini al codice, i commenti mancano di una struttura semantica interrogabile e i log di commit sono legati alla cronologia lineare delle modifiche ai file, rendendo difficile estrarre la "storia di una decisione" senza analizzare manualmente i diff.
*   **La memoria umana**: la risorsa più costosa e a più alto rischio di obsolescenza (attrito del personale, semplice dimenticanza).

Se le persone cambiano o dimenticano, la conoscenza si perde. Il software si irrigidisce, il debito tecnico aumenta per paura di toccare codice "di cui non si capisce il senso" e gli agenti AI — che oggi assistono la scrittura del codice — operano alla cieca, privi del contesto che giustifica lo stato attuale del repository.

---

## 2. Evidenze Empiriche: La Lezione del Dogfooding di Cairn

Il progetto Cairn è stato concepito per risolvere questo problema. Nello spirito del *dogfooding*, l'esperienza empirica accumulata nelle prime ore di sviluppo di Cairn stesso fornisce la prova scientifica e inconfutabile di questa tesi. Tre eventi reali avvenuti il 10 Luglio 2026 dimostrano come la conoscenza vada persa o ignorata senza un livello di cattura integrato e ad attrito zero.

### 2.1 La "Sinfonia del Documento Vuoto" (Friction vs. Speed)
Nel tentativo di seguire un rigoroso processo di documentazione formale, il repository di Cairn è stato inizialmente dotato di una struttura di "Memory Bank" (`projectBrief.md`, `productContext.md`, ecc.). 

Tuttavia, mentre l'implementazione in Go avanzava speditamente con la creazione dei comandi `cairn add` e `cairn log`, i file di specifica di alto livello sono rimasti pieni di placeholder generici come `[...]` e `[Descrivi il dolore dell'utente...]`. Questo non è avvenuto per negligenza, ma per un vincolo biologico: **la velocità dello sviluppo attivo crea un'inerzia che rende insopportabile la cerimonia della documentazione formale strutturata.**

La vera motivazione dietro la decisione di saltare la validazione manuale pianificata non è stata descritta in un lungo documento di design, ma è stata catturata al volo in una nota grezza di Cairn: 
> `{"id":"b36a1825bcbdf0c1","text":"Semplice impazienza con il processo"}`

Senza un meccanismo rapido per appendere questa "verità non filtrata", l'ADR in `systemPatterns.md` sarebbe rimasta vuota o avrebbe ospitato una giustificazione formale e parzialmente artificiale, nascondendo il reale fattore umano che ha guidato la traiettoria del software.

### 2.2 Il Silenzio dell'Ambiente e lo Schermo del `.gitignore`
Durante la v0 del bootstrap di Cairn, la cartella destinata a ospitare il log semantico local-first (`.cairn/`) è stata inizialmente inclusa nel file `.gitignore` del progetto. Questo pattern, ereditato automaticamente dalle configurazioni standard che escludono i dati locali di runtime, ha rischiato di escludere l'intera base di conoscenza appena creata dal sistema di controllo di versione (Git).

Questo errore di configurazione ambientale è rimasto invisibile finché l'analisi del workflow non ha sollevato una domanda fondamentale: *se Cairn deve versionare la conoscenza, perché la conoscenza di Cairn viene esclusa dal repository del progetto?* 

Senza una riflessione esplicita guidata dal "Knowledge Layer" e la conseguente rimozione di `.cairn/` dal `.gitignore`, l'intera memoria generata durante lo sviluppo iniziale sarebbe svanita al primo cambio di macchina dello sviluppatore, dimostrando come gli strumenti di sviluppo tradizionali tendano a trattare i dati di conoscenza come semplici "artefatti di log" transitori anziché come cittadini di prima classe del repository.

### 2.3 L'Emergenza dall'Uso Reale: Il Bug di `flag.Parse`
In fase di pianificazione statica, la sintassi del comando per aggiungere tag era stata definita come:
`cairn add --tags a,b,c "testo"`

Sulla carta, la specifica appariva perfetta e logica. Tuttavia, non appena lo sviluppatore ha iniziato a usare realmente lo strumento, la tendenza naturale è stata quella di digitare il flag alla fine:
`cairn add "testo" --tags a,b,c`

In Go, il comportamento di default del pacchetto standard `flag` prevede l'interruzione del parsing al primo argomento posizionale (non-flag). Di conseguenza, il comando ingoiava silenziosamente `--tags` all'interno del testo della nota, senza sollevare errori e lasciando l'utente convinto che i tag fossero stati applicati correttamente.

Questo bug UX, scoperto esclusivamente grazie all'esperienza fisica del dogfooding e registrato nella nota `2b9ce728be4bbe67`, evidenzia che:
1.  **La conoscenza reale emerge solo a runtime**, dall'interazione viva con il software.
2.  Nessuna analisi statica delle specifiche avrebbe previsto questa frizione.
3.  L'apprendimento derivato (ovvero la necessità di implementare una funzione `misplacedFlag()` per intercettare l'errore ed evitare il fallimento silenzioso) costituisce un pezzo di conoscenza tecnica prezioso che deve essere archiviato e legato per sempre a quel modulo, evitando che futuri manutentori ripristinino inavvertitamente il comportamento nativo di Go.

---

## 3. Il Trade-off Fondamentale: Attrito vs. Struttura

Se accettiamo che la conoscenza debba essere preservata, dobbiamo affrontare il trade-off che finora ha causato il fallimento di ogni Wiki aziendale.

| Approccio | Vantaggi | Svantaggi | Impatto su *Adoption* |
| :--- | :--- | :--- | :--- |
| **Wiki Centralizzata (Confluence/Notion)** | Alta formattazione, visibile a non-tecnici, diagrammi ricchi. | Distante dal codice, richiede autenticazione extra, marcisce istantaneamente. | **Basso.** Lo sviluppatore deve cambiare contesto e uscire dal terminale. |
| **Commenti / JSDoc nel Codice** | Massima vicinanza alla riga di codice modificata. | Difficili da estrarre globalmente, sporcano la sintassi, invisibili a strumenti esterni non-codice. | **Medio.** Richiede la modifica fisica del sorgente per registrare un pensiero. |
| **ADR Manuali in Markdown** | Strutturati, versionati in Git. | Alta cerimonia (richiedono front-matter, ID sequenziali, modifiche a tabelle indice). | **Basso.** Spesso rimandati a "dopo il merge" e quindi mai scritti. |
| **Log Semantico Locale (Cairn)** | Scrittura istantanea da CLI (`cairn add`), local-first, append-only, integrato in Git. | Interfaccia inizialmente testuale, richiede un parser per estrarre relazioni complesse. | **Altissimo.** Attrito quasi zero; la cattura avviene durante l'azione stessa. |

### L'alternativa scartata: Relazioni Tipizzate Immediate
Nelle prime fasi di Cairn, era stata proposta l'introduzione di relazioni tipizzate tra note (es. `"supersede"`, `"deprecated by"`, `"references"`). 
**Compromesso analizzato:** Per associare una relazione, lo sviluppatore avrebbe dovuto conoscere l'ID (hash) di una nota precedente, costringendolo a un lookup preliminare (`cairn log` o grep). Questa operazione introduce un secondo livello di attrito in fase di inserimento.

In ossequio al vincolo del §2 di `AGENTS.md` (*Adoption & Freshness*), questa funzionalità è stata **rimandata**. È preferibile catturare una nota non strutturata o debolmente taggata sul momento piuttosto che perdere completamente la cattura a causa della complessità richiesta per strutturarla.

---

## 4. La Proposta: Il "Knowledge Layer" come Primitiva di Sviluppo

Cairn propone che il software moderno debba includere un **livello di conoscenza** integrato che rispetti le seguenti caratteristiche:

1.  **Local-First e Git-Native:** La conoscenza deve risiedere nello stesso repository del codice. Se il repository viene clonato, la conoscenza si sposta con esso; se viene creato un branch, la conoscenza si dirama; se viene eseguito un merge, la conoscenza si fonde.
2.  **Append-Only & Zero-Friction:** Lo strumento primario di input deve essere una CLI ultra-leggera in grado di registrare un'osservazione o una decisione in un singolo comando, riducendo l'attrito cognitivo al minimo indispensabile.
3.  **Provenienza Tracciabile:** Ogni frammento di conoscenza deve essere dotato di identificativo unico, timestamp, tag liberamente associabili e autore, garantendo la provenienza e la verificabilità del dato.
4.  **Umani e Agenti alla Pari:** Il log generato deve essere leggibile e modificabile sia dagli umani che dagli agenti AI (attraverso formati aperti come JSON Lines e Markdown), fungendo da ponte semantico comune per comprendere l'evoluzione del sistema.

L'esperienza empirica del bootstrap di Cairn dimostra che senza questa primitiva, la documentazione diventa un esercizio di burocrazia post-facto. Con un livello di conoscenza esplicito e integrato nel terminale, catturare il *perché* diventa naturale quanto scrivere un commento o fare un commit.
