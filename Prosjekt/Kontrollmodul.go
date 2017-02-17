package main

import (


)

func Beregn_kostnad(){
  //gjør beregninger - eksterne og interne ordre
  //del resultat på nettverk
}

func Ekstern_ordre(){
  //check om ordre allerede er i køen
  //legg til hvis ikke
  //del ordre på nettverket
}

func Inkommende_ordre(){
  //check om ordre allerede er i køen
  //legg til hvis ikke
  //fjern intern ordre fra kø når ordren er fullført
}

func Vurder_kostnad(){
  //sammenlign innkommende resultat
  //ta ordre eller la være hvis noen andre har bedre kost
}

func Kvitter_ordre(){
  //send kvittering for ekstern ordre på nettverket
  //motta kvittering
  //fjern ordre fra kø
}

func main(){

}
