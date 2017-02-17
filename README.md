# Gitrepo for gruppe 67

## Prosedyre ved ny ekstern ordre:
1. Noden som mottar ordren publiserer den på nettverket
2. Alle noder returnerer med kostnad og lagrer ordren lokalt
3. Den med lavest kostnad og lavest IP tar ordren og skrur på lyset
4. Noden som utfører ordren kvitterer når den er ferdig
5. Dersom kvittering ikke kommer innen 15sek publiseres ordren på nytt.

## Prosedyre ved interne ordre:
1. Sjekk om ordren finnes fra før eller vi allerede er i den etasjen ordren er til
2. Legg ordren til i intern ordrekø og inkrementer ordreteller og skrur på lyset
3. Sorter køen og lagre på fil, bakerste element er første ordre som blir utført
4. Når ordren er utført dekrementers ordretelleren og oppdaterer filen og skrur av lyset

## Nettverksmodul
* Broadcast_TCP(melding)()
  * Oppdag andre noder og legg til i nodeliste
  * Send melding til alle på nodeliste
  * Fjern fra nodeliste dersom død
* Recieve_TCP()(melding, IP)
  * Les inn meldinger og returner melding og IP

## Heismodul
* Kjør_heis(intern kø og antall ordrer)(intern kø og antall ordrer)
  * Åpne dør, lukkedør, lys, kjører automatisk til siste element i ordrekøen dersom antall ordrer > 0
  * Fjerne utførte ordrer
* Intern_ordre(intern kø og antall ordrer)(intern kø og antall ordrer)
  * Se "Prosedyre ved intern ordre"
* Etasje_indikator()()
  * Vis hvilken etasje vi er i
## Kontrollmodul
* Beregn_kostnad(intern kø og antall ordrer)
  * Se på intern ordrekø og beregn kostnad og send den ut på nettverket
* Ekstern_ordre
  * Publiser ordren på nettverket, publiser på nytt dersom timeout
* Inkommende_ordre
  * Legg ordren til i lokalt ordrearray
* Vurder_kostnad
  * Se på innkommende kostnader og vurder og vi skal utføre ordren og legg til i intern ordrekø
* Kvitter_ordre
  * Dersom utført ekstern ordre, send kvittering ut på nettverket
