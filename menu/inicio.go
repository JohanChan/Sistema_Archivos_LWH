package menu

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"proyecto/analizador"

	"strings"
)

var arrayLinea []string
var cadenaArchivo string

func Menu() {
	fmt.Println("Ingrese comando")
	bf := bufio.NewReader(os.Stdin)
	entrada, _ := bf.ReadString('\n')
	cadena := strings.TrimRight(entrada, "\r\n")

	SplitSaltos(cadena)

	//arrayLinea = ValidaPath(arrayLinea)
	Exec(arrayLinea)
	analizador.Imprimir(arrayLinea)
	fmt.Println("Desea ingresar otro comando?, s/n")
	bf = bufio.NewReader(os.Stdin)
	entrada, _ = bf.ReadString('\n')
	chain := strings.TrimRight(entrada, "\r\n")
	if strings.EqualFold(chain, "s") {
		Menu()
	} else if strings.EqualFold(chain, "n") {
		fmt.Println("Un gusto, vuelva pronto")
	}
}
func SplitSaltos(cadena string) {
	sSaltos := strings.SplitN(cadena, "\n", -1)
	for i := 0; i < len(sSaltos); i++ {
		SplitEspacios(sSaltos[i])
	}

}
func SplitEspacios(cadena string) {
	sEspacio := strings.SplitN(cadena, " ", -1)
	for i := 0; i < len(sEspacio); i++ {
		SplitFlecha(sEspacio[i])
	}
}
func SplitFlecha(cadena string) {
	SplitFlecha := strings.SplitN(cadena, "->", -1)
	for i := 0; i < len(SplitFlecha); i++ {
		arrayLinea = append(arrayLinea, SplitFlecha[i])
	}
}
func LeerArchivo(file string) {
	datos, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Ruta no encontrada")
	} else {
		cadenaArchivo = string(datos)
		SplitSaltos(cadenaArchivo)
		arrayLinea = ValidaPath(arrayLinea)
	}
}
func Exec(arreglo []string) {
	for i := 0; i <= len(arreglo)-1; i++ {
		if strings.EqualFold(arreglo[i], "exec") {
			path := arreglo[i+2]
			LimpiarArreglo()
			LeerArchivo(path)
			EjecutarArchivo()
		} else {
			EjecutarArchivo()
			return
		}
	}
}

func LimpiarArreglo() {
	arrayLinea = nil
}
func EjecutarArchivo() {
	arrayLinea = append(arrayLinea, " ")
	analizador.AsignarArray(arrayLinea)
	analizador.FuncionComando(arrayLinea)
	LimpiarArreglo()
}
func ValidaPath(arreglo []string) []string {
	var arregloOrdenado []string
	bandera := 0
	var cadena string
	for i := 0; i < len(arreglo); i++ {
		if bandera == 0 && strings.Contains(arreglo[i], "\"") {
			bandera = 1
			cadena += arreglo[i] + " "
		} else if bandera == 1 && strings.Contains(arreglo[i], "\"") {
			bandera = 0
			cadena += arreglo[i] + " "
			arregloOrdenado = append(arregloOrdenado, cadena)
			cadena = ""
		} else if bandera == 1 && !strings.Contains(arreglo[i], "\"") {
			cadena += arreglo[i] + " "
		} else {
			arregloOrdenado = append(arregloOrdenado, arreglo[i])
		}
	}
	return arregloOrdenado

}
