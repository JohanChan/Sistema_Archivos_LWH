package estructura

import (
	"proyecto/analizador"
)

type nodo struct {
	Particion analizador.Partition
	IDn       int
	siguiente *nodo
}
type particionMontada struct {
}

var inicio *nodo
var indice int

func Insertar(particion analizador.Partition, IDn int) {
	var nuevo *nodo
	nuevo.Particion = particion
	nuevo.IDn = IDn
	if EstaVacia() {
		inicio = nuevo
		indice++
	} else {
		nuevo.siguiente = inicio
		inicio = nuevo
		indice++
	}
}
func Elimianr(IDn int) {
	aux := inicio
	auxS := inicio.siguiente
	if IDn == inicio.IDn {
		inicio = nil
	} else {
		for i := 0; i < Tama(); i++ {
			if IDn == auxS.IDn {
				aux = auxS.siguiente
			}
			aux = aux.siguiente
			auxS = auxS.siguiente
		}
	}

}
func Tama() int {
	return indice
}
func EstaVacia() bool {
	return inicio == nil
}
