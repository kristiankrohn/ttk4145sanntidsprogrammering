package main

import (


)

func Beregn_kostnad(){
  //gjør beregninger - eksterne og interne ordre
    //det enkleste er å se på det totale antall ordre og etasjer den skal kjøre
  //del resultat på nettverk
}

func Ekstern_ordre(){
  //check om ordre allerede er i køen
  //legg til hvis ikke
  //del ordre på nettverket
  //del på nytt hvis timeout
}

func Inkommende_ordre(){
  //check om ordre allerede er i køen
  //legg til hvis ikke
  //fjern intern ordre fra kø når ordren er fullført
}

func Vurder_kostnad(){
  //sammenlign innkommende resultat
  //Vurder om vi skal ta ordre og legge den til i intern ordrekø
}

func Kvitter_ordre(){
  //send kvittering for ekstern ordre på nettverket
  //motta kvittering
  //fjern ordre fra kø
}

func main(){

}
