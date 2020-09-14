package funciones

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

var size, path, name, unit string
var tipo, fit, eliminar, agregar, idMkfs, idRep, ruta, usr, pwd, idLogin, idMkgrp, idMkdir string
var mbr Mbr
var logueado bool = false
var arrayLinea []string
var listadoMount []ParticionMontada
var abecedario = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var incremento int
var id []string
var totalPart int64
var ebr, ebrf Ebr
var part Partition

//*********************************** STRUCTS **********************************************************
type Mbr struct {
	Mbr_tamano         int64     //8 bytes
	Mbr_fecha_creacion [10]byte  // 20 bytes
	Mbr_disk_signature int64     // 8 bytes
	Mbr_partition_1    Partition // 35 bytes
	Mbr_partition_2    Partition // 35 bytes
	Mbr_partition_3    Partition // 35 bytes
	Mbr_partition_4    Partition // 35 bytes
}
type Ebr struct {
	Part_status byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_next   int64
	Part_name   [16]byte
}

//Struct Superboot
type Sb struct {
	Sb_nombre_hd                         [16]byte
	Sb_arbol_virtual_count               int64
	Sb_detalle_directo_count             int64
	Sb_inodos_count                      int64
	Sb_bloques_count                     int64
	Sb_arbol_virtual_free                int64
	Sb_detalle_directorio_free           int64
	Sb_inodos_free                       int64
	Sb_bloques_free                      int64
	Sb_date_creacion                     [16]byte
	Sb_date_ultimo_montaje               [16]byte
	Sb_montajes_count                    int64
	Sb_ap_bitmap_arbol_directorio        int64
	Sb_ap_arbol_directorio               int64
	Sb_ap_bitmap_detalle_directorio      int64
	Sb_ap_detalle_directorio             int64
	Sb_ap_bitmap_tabla_inodo             int64
	Sb_ap_tabla_inodo                    int64
	Sb_ap_bitmap_bloques                 int64
	Sb_ap_bloques                        int64
	Sb_ap_log                            int64
	Sb_size_struct_arbol_directorio      int64
	Sb_size_struct_detalle_directorio    int64
	Sb_size_struct_inodo                 int64
	Sb_size_struct_bloque                int64
	Sb_first_free_bit_arbol_directorio   int64
	Sb_first_free_bit_detalle_directorio int64
	Sb_first_free_bit_tabla_inodo        int64
	Sb_first_free_bit_bit_bloques        int64
	Sb_magic_num                         int64
}

//********************************************

type Avd struct {
	Avd_fecha_creacion              [16]byte
	Avd_nombre_directorio           [20]byte
	Avd_ap_array_subdirectorios     [5]int64
	Avd_ap_detalle_directorio       int64
	Avd_ap_arbol_virtual_directorio int64
	Avd_proper                      [9]byte
}

//********************************************
//Struct Detalle de Directorio
type Dd struct {
	Dd_array_files           [5]SubDetalle
	Dd_ap_detalle_directorio int64
}

//********************************************
//Struct Detalle de Directorio
type SubDetalle struct {
	DdArray_file_nombre            [19]byte
	DdArray_file_ap_inodo          int64
	DdArray_file_date_creacion     [16]byte
	DdArray_file_date_modificacion [16]byte
}

//********************************************
//Struct Tabla I-Nodo
type I struct {
	I_count_inodo   int64
	I_size_archivo  int64
	I_array_bloques [4]int64
	I_ap_indirecto  int64
	I_id_proper     [9]byte
}

//********************************************
//Struct Bloque de Datos
type Db struct {
	Db_data [24]byte
}

//********************************************
//Struct Arbol Virtual de Directorio
type Log struct {
	Log_tipo_operacion [19]byte
	Log_tipo           byte
	Log_nombre         [19]byte
	Log_contenido      [24]byte
	Log_fecha          [16]byte
}

//********************************************
//Struct Arbol Virtual de Directorio
type Partition struct {
	Part_status byte     // 1 byte
	Part_type   byte     // 1 byte
	Part_fit    byte     // 1 byte
	Part_start  int64    // 8 bytes
	Part_size   int64    // 8 bytes
	Part_name   [16]byte // 16 bytes
}
type ParticionMontada struct {
	Path  string
	Name  string
	IDn   string
	Letra string
}

//******************************************************************************************************

func AsignarArray(arraInicio []string) {
	arrayLinea = arraInicio
	/*for i := 0; i < len(arrayLinea); i++ {
		fmt.Println(arrayLinea[i])
	}*/
}
func EliminarBarra(atributo string) string {
	if strings.Contains(atributo, "*") {
		natributo := strings.ReplaceAll(atributo, "\\*", "")
		return natributo
	}
	return atributo
}

var comT bool = false

func FuncionComando(arreglo []string) {
	for i := 0; i <= len(arreglo)-1; i++ {
		if strings.EqualFold(arreglo[i], "mkdisk") {
			comT = true
			AtributosMKDISK(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "rmdisk") {
			comT = true
			RMDISK(arreglo[i+2])
		} else if strings.EqualFold(arreglo[i], "fdisk") {
			comT = true
			AtributoFDISK(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "mount") {
			comT = true
			AtributosMount(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "unmount") {
			comT = true
			AtributosUnmount(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "pause") {
			comT = true
			CapturarPantalla("Presione Enter para continuar")
		} else if strings.EqualFold(arreglo[i], "particiones") {
			read()
		} else if strings.EqualFold(arreglo[i], "mkfs") {
			comT = true
			AtributoMKFS(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "rep") {
			comT = true
			AtributosREP(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "login") {
			comT = true
			AtributosLogin(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "logout") {
			comT = true
			Logout()
		} else if strings.EqualFold(arreglo[i], "mkgrp") {
			comT = true
			AtributosMKGRP(arreglo[i+1], i+1)
		} else {
			if comT == false {
				fmt.Println("Comando no especificado")
				return
			}
		}
	}
}

//*************************** Atributos Funciones ******************************************************
func AtributosMKDISK(atrributo string, indice int) {
	if strings.EqualFold(atrributo, "-size") {
		size = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMKDISK(arrayLinea[indice], indice)

	} else if strings.EqualFold(atrributo, "-path") {
		path = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMKDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atrributo, "-name") {
		name = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMKDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atrributo, "-unit") {
		unit = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMKDISK(arrayLinea[indice], indice)
	} else {
		CrearDisco(size, path, name, unit)
		size = ""
		name = ""
		unit = ""
	}

}
func RMDISK(path string) {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("Disco a eliminar no Existe!")
		return
	}
	captura := CapturarPantalla("Estas seguro que deseas eliminar? " + path + " y/n")
	if strings.EqualFold(captura, "y") {
		path = strings.Replace(path, "\"", "", -1)
		os.Remove(path)
		fmt.Println("Disco " + path + " ha sido eliminado")
	}
}
func CapturarPantalla(mensaje string) string {
	fmt.Println(mensaje)
	bf := bufio.NewReader(os.Stdin)
	entrada, _ := bf.ReadString('\n')
	cadena := strings.TrimRight(entrada, "\r\n")
	return cadena
}
func AtributoFDISK(atributo string, indice int) {
	if strings.EqualFold(atributo, "-size") {
		size = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoFDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-unit") {
		unit = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoFDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-path") {
		path = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoFDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-type") || strings.EqualFold(atributo, "-tipo") {
		tipo = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoFDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-fit") {
		fit = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoFDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-delete") {
		eliminar = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoFDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-name") {
		name = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoFDISK(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-add") {
		agregar = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoFDISK(arrayLinea[indice], indice)
	} else {
		if eliminar != "" {
			EliminarParticion(eliminar, name, path)
			eliminar = ""
		} else if agregar != "" {
			fmt.Println("Aqui no hacemos eso")
		} else {
			mbr = LeerMBR(path)
			if strings.EqualFold(tipo, "l") {
				CrearParticionLogica(size, unit, path, fit, name)
			} else {
				CrearParticion(size, unit, path, tipo, fit, name)
			}
		}
		size = ""
		unit = ""
		tipo = ""
		fit = ""
		name = ""
		agregar = ""
		path = ""
		eliminar = ""
	}

}
func AtributosMount(atributo string, indice int) {
	if strings.EqualFold(atributo, "-name") {
		name = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMount(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-path") {
		path = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMount(arrayLinea[indice], indice)
	} else {
		if name == "" && path == "" {
			ImprimirMontadas(listadoMount)
		} else {

			MontarParticion(path, name)
			path = ""
			name = ""
		}
	}
}
func AtributosUnmount(atributo string, indice int) {
	if strings.EqualFold(atributo, "-idn") {
		id = append(id, EliminarBarra(arrayLinea[indice+1]))
		indice += 2
		AtributosUnmount(arrayLinea[indice], indice)
	} else {
		DesmontarParticion(id)
	}
}
func AtributoMKFS(atributo string, indice int) {
	if strings.EqualFold(atributo, "-id") {
		idMkfs = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoMKFS(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-tipo") {
		tipo = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoMKFS(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-add") {

	} else if strings.EqualFold(atributo, "-unit") {
		unit = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributoMKFS(arrayLinea[indice], indice)
	} else {
		FormateoLWH(idMkfs)
		idMkfs = ""
		tipo = ""
		unit = ""

	}
}
func AtributosREP(atributo string, indice int) {
	if strings.EqualFold(atributo, "-nombre") || strings.EqualFold(atributo, "-name") {
		name = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosREP(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-path") {
		path = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosREP(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-id") {
		idRep = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosREP(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-ruta") {
		ruta = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosREP(arrayLinea[indice], indice)
	} else {
		CrearReportes(name, path, idRep, ruta)
		path = ""
		name = ""
		ruta = ""
		idRep = ""
	}
}
func AtributosLogin(atributo string, indice int) {
	if strings.EqualFold(atributo, "-usr") {
		usr = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosLogin(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-pwd") {
		pwd = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosLogin(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-id") {
		idLogin = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosLogin(arrayLinea[indice], indice)
	} else {
		IniciarSesion(usr, pwd, idLogin)
	}
}
func AtributosMKGRP(atributo string, indice int) {
	if strings.EqualFold(atributo, "-id") {
		idMkgrp = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMKGRP(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-name") {
		name = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMKGRP(arrayLinea[indice], indice)
	} else {
		if logueado {
			CrearGrupo(idMkgrp, name)
		} else {
			fmt.Println("No hay usuario logueado")
		}
	}
}

var p string

func AtributosMKDIR(atributo string, indice int) {
	if strings.EqualFold(atributo, "-id") {
		idMkdir = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMKDIR(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-path") {
		path = EliminarBarra(arrayLinea[indice+1])
		indice += 2
		AtributosMKDIR(arrayLinea[indice], indice)
	} else if strings.EqualFold(atributo, "-p") {
		p = EliminarBarra(arrayLinea[indice])
		indice += 1
		AtributosMKDIR(arrayLinea[indice], indice)
	} else {
		if logueado {
			Mkdir(idMkdir, path, p)
		} else {
			fmt.Println("No hay usuario logueado")
		}
	}
}

//******************************************************************************************************
var punteroBloques []int64
var punteroInodos []int64

//*************************** MKGRP ********************************************************************
func CrearGrupo(idp string, nombre string) {
	if len(nombre) < 10 {
		pPart := SearchPath(idp)
		nPart := SearchNombre(idp)
		mbr = LeerMBR(pPart)
		par := ObtenerParticion(nPart, mbr)
		sb = LeerSb(pPart, par.Part_start)
		a := LeerBitmap(pPart, sb.Sb_ap_bitmap_bloques, sb.Sb_bloques_count)
		bInodo := LeerBitmap(pPart, sb.Sb_ap_bitmap_tabla_inodo, sb.Sb_inodos_count)
		avd = LeerAvd(pPart, sb.Sb_ap_arbol_directorio)
		dd = LeerDD(pPart, sb.Sb_ap_detalle_directorio)
		i = LeerTInodo(pPart, sb.Sb_ap_tabla_inodo)
		var bloqueActual, iActual int64
		var cadena string = ""
		var tot, totalInodos int
		nBit := CantidadBitOcupados(a)
		for in := 0; in < len(i.I_array_bloques); in++ {
			if i.I_array_bloques[in] != 0 {
				tot++
				punteroBloques = append(punteroBloques, i.I_array_bloques[in])
				bloqueActual = i.I_array_bloques[in]
				db = LeerBloque(pPart, i.I_array_bloques[in])
				for j := 0; j < len(db.Db_data); j++ {
					if db.Db_data[j] > 0 {
						cadena += string(db.Db_data[j])
					}
				}
			}
		}
		cadena += "1,G," + nombre + "_"
		fmt.Println(len(cadena))
		n := SplitSubN(cadena)
		var bloquesNuevos float64 = float64(len(n) - tot)
		var iNodosNuevos float64 = float64(len(n) / 5)
		fmt.Println(iNodosNuevos)
		if tot%4 == 0 {
			for j := 0; j < len(bInodo); j++ {
				if bInodo[j] == '1' {
					totalInodos++
				}
			}
			iActual = dd.Dd_array_files[0].DdArray_file_ap_inodo
			fmt.Println(iActual)
			var NewInodo I
			NewInodo.I_ap_indirecto = 0
			NewInodo.I_count_inodo = 0
			NewInodo.I_size_archivo = 0

			fmt.Println("Se crea nuevo inodo")

		} else {
			for i := 0; i < int(bloquesNuevos); i++ {
				bloqueActual += sb.Sb_size_struct_bloque
				punteroBloques = append(punteroBloques, bloqueActual)
			}
			for ii := 0; ii < len(punteroBloques); ii++ {
				var dbNew Db
				fmt.Println("i " + strconv.Itoa(ii))
				i.I_array_bloques[ii] = punteroBloques[ii]
				copy(dbNew.Db_data[:], n[ii])
				fmt.Println(string(dbNew.Db_data[:]))
				EscribirBloques(pPart, punteroBloques[ii], dbNew)
			}
			EscribirTInodos(pPart, sb.Sb_ap_tabla_inodo, i)
			punteroBloques = nil
			fmt.Println(punteroBloques)
			EscribirBitmap(pPart, sb.Sb_ap_bitmap_bloques+int64(nBit), int64(bloquesNuevos), '1')
		}
	} else {
		fmt.Println("nombre 10+ caracteres")
	}
}
func clean(s []int64, tam int) []int64 {
	for i := tam; i >= 0; i-- {
		s = LimpiarArreglo(s, i)
	}
	return s
}
func LimpiarArreglo(s []int64, index int) []int64 {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}
func FullBlock(a Db) bool {
	if Isfull(string(a.Db_data[:])) == 24 {
		return true
	}
	return false
}
func Isfull(a string) int {
	var n int = 0
	for _, r := range a {
		if unicode.IsLetter(r) {
			n++
		} else if unicode.IsDigit(r) {
			n++
		} else if r == ',' {
			n++
		} else if r == '\n' {
			n++
		}
	}
	return n
}
func CantidadBitOcupados(a []byte) int {
	var CantBit int = 0
	for i := 0; i < len(a); i++ {
		if a[i] != '0' {
			CantBit++
		}
	}
	return CantBit
}
func SplitSubN(s string) []string {
	var c string = ""
	var ac []string

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		c = c + string(r)
		if (i+1)%25 == 0 {
			ac = append(ac, c)
			c = ""
		} else if (i + 1) == l {
			ac = append(ac, c)
		}
	}
	return ac
}

//******************************************************************************************************
//*************************** MKDIR ********************************************************************
func Mkdir(idm string, patch string, p string) {
	if idm == "" || patch == "" {
		fmt.Println("Faltan parametros obligatorios")
	} else {
		pPart := SearchPath(idm)
		nPart := SearchNombre(idm)
		mbr = LeerMBR(pPart)
		par := ObtenerParticion(nPart, mbr)
		sb = LeerSb(pPart, par.Part_start)
	}
}

//******************************************************************************************************

//*************************** LOGOUT ********************************************************************
func Logout() {
	if logueado {
		fmt.Println("Adios")
		logueado = false
		usr = ""
		pwd = ""
		idLogin = ""
	} else {
		fmt.Println("Ninguna sesión iniciada")
	}
}

//******************************************************************************************************
var nGrupo string

//*************************** LOGIN ********************************************************************
func IniciarSesion(usuario string, password string, idn string) {
	if !EstaMontada(idn) {
		fmt.Println("Particion no montada")
		return
	} else {
		pfile := SearchPath(idn)
		nfile := SearchNombre(idn)
		mbr = LeerMBR(pfile)
		parfile := ObtenerParticion(nfile, mbr)
		sb = LeerSb(pfile, parfile.Part_start)
		avd = LeerAvd(pfile, sb.Sb_ap_arbol_directorio)
		dd = LeerDD(pfile, sb.Sb_ap_detalle_directorio)
		i = LeerTInodo(pfile, sb.Sb_ap_tabla_inodo)
		var cadena string = ""
		for in := 0; in < len(i.I_array_bloques); in++ {
			if i.I_array_bloques[in] != 0 {
				db = LeerBloque(pfile, i.I_array_bloques[in])
				for j := 0; j < len(db.Db_data); j++ {
					if db.Db_data[j] > 0 {
						cadena += string(db.Db_data[j])
					}
				}

			}
		}
		cadena = strings.ReplaceAll(cadena, "_", ",")
		slcLogin := strings.Split(cadena, ",")
		for i := 0; i < len(slcLogin)-1; i++ {
			if slcLogin[i] == "U" {
				//rr := slcLogin[i+1]
				pa := slcLogin[i+2]
				fmt.Println(pa)
				if slcLogin[i-2] != "0" {
					if slcLogin[i+1] == usuario && slcLogin[i+2] == password {
						fmt.Println("Bienvenido " + usuario)
						logueado = true
						nGrupo = slcLogin[i-2]

					} else {
						fmt.Println("Datos incorrectos o usuario no existe")
					}
				}
			}
		}
	}

}
func EstaMontada(idn string) bool {
	var bandera bool = false
	for i := 0; i < len(listadoMount); i++ {
		if strings.EqualFold(listadoMount[i].IDn, idn) {
			bandera = true
		}
	}
	return bandera
}

//******************************************************************************************************

//*************************** FORMATEO LWH **************************************************************

var avd Avd
var dd Dd
var i I
var log Log
var sb Sb
var db Db

func read() {
	sb = LeerSb("/home/johan/MisDiscos/johan.dsk", 1048743)
	fmt.Println("********* SuperBoot **********")
	fmt.Println("Nombre disco " + string(sb.Sb_nombre_hd[:]))
	fmt.Println("Carne " + strconv.FormatInt(sb.Sb_magic_num, 10))
	fmt.Println("Fecha Creacion " + string(sb.Sb_date_creacion[:]))
	fmt.Println("*******************************")

	avd = LeerAvd("/home/johan/MisDiscos/johan.dsk", sb.Sb_ap_arbol_directorio)
	fmt.Println("************* AVD ************")
	fmt.Println("Nombre directorio " + string(avd.Avd_nombre_directorio[:]))
	fmt.Println("apuntador dd " + strconv.FormatInt(avd.Avd_ap_detalle_directorio, 10))
	fmt.Println("Fecha Creacion " + string(avd.Avd_fecha_creacion[:]))
	fmt.Println("*******************************")

	dd = LeerDD("/home/johan/MisDiscos/johan.dsk", sb.Sb_ap_detalle_directorio)
	fmt.Println("************* DD ************")
	fmt.Println("Nombre archivo " + string(dd.Dd_array_files[0].DdArray_file_nombre[:]))
	fmt.Println("apuntador inodo " + strconv.FormatInt(dd.Dd_array_files[0].DdArray_file_ap_inodo, 10))
	fmt.Println("Fecha Creacion " + string(dd.Dd_array_files[0].DdArray_file_date_creacion[:]))
	fmt.Println("*******************************")

	i = LeerTInodo("/home/johan/MisDiscos/johan.dsk", sb.Sb_ap_tabla_inodo)
	fmt.Println("************* INODO ************")
	fmt.Println("Nombre propietario " + string(i.I_id_proper[:]))
	for ii := 0; ii < len(i.I_array_bloques); ii++ {
		fmt.Println("apntador ibloque " + strconv.FormatInt(i.I_array_bloques[ii], 10))
	}
	fmt.Println("Fecha Creacion " + string(dd.Dd_array_files[0].DdArray_file_date_creacion[:]))
	fmt.Println("*******************************")

	for ii := 0; ii < len(i.I_array_bloques); ii++ {
		if i.I_array_bloques[ii] != 0 {
			db = LeerBloque("/home/johan/MisDiscos/johan.dsk", i.I_array_bloques[ii])
			fmt.Println("************* BLOQUE " + strconv.Itoa(ii) + " ************")
			fmt.Println("Contenido " + string(db.Db_data[:]))
			fmt.Println("*******************************")
		}
	}

}
func FormateoLWH(id string) {
	pathArchivo := SearchPath(id)
	nombreParticion := SearchNombre(id)
	if strings.EqualFold(pathArchivo, "") {
		fmt.Println("Particion no montada ")
		return
	}
	mbr := LeerMBR(pathArchivo)

	par := ObtenerParticion(nombreParticion, mbr)
	var inicioPart int64 = par.Part_start
	var sizePart int64 = par.Part_size

	var sizeAvd int64 = int64(binary.Size(avd))
	var sizeDd int64 = int64(binary.Size(dd))
	var sizeInodo int64 = int64(binary.Size(i))
	var sizeBitacora int64 = int64(binary.Size(log))
	var sizeSb int64 = int64(binary.Size(sb))
	var sizeDb int64 = int64(binary.Size(db))

	nEstructuras := (sizePart - (2 * sizeSb)) / (27 + sizeAvd + sizeDd + (5*sizeInodo + (20 * sizeDb) + sizeBitacora))
	fmt.Println("Total estructuras " + strconv.FormatInt(nEstructuras, 10))
	cantidadAvd := nEstructuras
	cantidadDD := nEstructuras
	cantidadInodos := nEstructuras
	cantidadBloques := nEstructuras
	//cantidadBitacoras := nEstructuras
	fmt.Println("Inicio particion " + strconv.FormatInt(inicioPart, 10))

	sb.Sb_ap_bitmap_arbol_directorio = inicioPart + sizeSb
	sb.Sb_ap_arbol_directorio = sb.Sb_ap_bitmap_arbol_directorio + cantidadAvd
	sb.Sb_ap_bitmap_detalle_directorio = sb.Sb_ap_arbol_directorio + (cantidadAvd * sizeAvd)
	sb.Sb_ap_detalle_directorio = sb.Sb_ap_bitmap_detalle_directorio + cantidadDD
	sb.Sb_ap_bitmap_tabla_inodo = sb.Sb_ap_detalle_directorio + (sizeDd * cantidadDD)
	sb.Sb_ap_tabla_inodo = sb.Sb_ap_bitmap_tabla_inodo + cantidadInodos
	sb.Sb_ap_bitmap_bloques = sb.Sb_ap_tabla_inodo + (sizeInodo * cantidadInodos)
	sb.Sb_ap_bloques = sb.Sb_ap_bitmap_bloques + cantidadBloques
	sb.Sb_ap_log = sb.Sb_ap_bloques + (sizeDb * cantidadBloques)
	espacios := "               "
	copy(sb.Sb_nombre_hd[:], espacios)
	copy(sb.Sb_nombre_hd[:], nombreParticion)
	sb.Sb_arbol_virtual_count = cantidadAvd
	sb.Sb_detalle_directo_count = cantidadDD
	sb.Sb_inodos_count = cantidadInodos
	sb.Sb_bloques_count = cantidadBloques
	copy(sb.Sb_date_creacion[:], time.Now().Format("02-01-2006 15:04:05"))
	copy(sb.Sb_date_ultimo_montaje[:], time.Now().Format("02-01-2006 15:04:05"))
	sb.Sb_montajes_count = 0
	sb.Sb_size_struct_arbol_directorio = sizeAvd
	sb.Sb_size_struct_bloque = sizeDb
	sb.Sb_size_struct_detalle_directorio = sizeDd
	sb.Sb_size_struct_inodo = sizeInodo
	sb.Sb_magic_num = 201603052
	sb.Sb_inodos_free = cantidadInodos
	sb.Sb_bloques_free = cantidadBloques
	sb.Sb_arbol_virtual_free = cantidadAvd
	sb.Sb_detalle_directorio_free = cantidadDD
	sb.Sb_first_free_bit_arbol_directorio = sb.Sb_ap_bitmap_arbol_directorio
	sb.Sb_first_free_bit_bit_bloques = sb.Sb_ap_bitmap_bloques
	sb.Sb_first_free_bit_detalle_directorio = sb.Sb_ap_bitmap_detalle_directorio
	sb.Sb_first_free_bit_tabla_inodo = sb.Sb_ap_bitmap_tabla_inodo

	EscribirSb(pathArchivo, inicioPart, sb)
	sb = LeerSb(pathArchivo, par.Part_start)
	avd = DatosDirectorioRoot()
	dd = DatosDDRoot()
	i = DatosInodoRoot()
	db1 := DatosBloquesRoot("1,G,root_1,root,U,root,2")
	db2 := DatosBloquesRoot("01603052_")

	fmt.Println("Inicio bavd " + strconv.FormatInt(sb.Sb_ap_bitmap_arbol_directorio, 10))
	EscribirBitmap(pathArchivo, sb.Sb_ap_bitmap_arbol_directorio, cantidadAvd, '0')

	fmt.Println("Inicio bdd " + strconv.FormatInt(sb.Sb_ap_bitmap_detalle_directorio, 10))
	EscribirBitmap(pathArchivo, sb.Sb_ap_bitmap_detalle_directorio, cantidadDD, '0')

	fmt.Println("Inicio btinodo " + strconv.FormatInt(sb.Sb_ap_bitmap_tabla_inodo, 10))
	EscribirBitmap(pathArchivo, sb.Sb_ap_bitmap_tabla_inodo, cantidadInodos, '0')

	fmt.Println("Inicio bbloques " + strconv.FormatInt(sb.Sb_ap_bitmap_bloques, 10))
	EscribirBitmap(pathArchivo, sb.Sb_ap_bitmap_bloques, cantidadBloques, '0')

	EscribirAvd(pathArchivo, sb.Sb_ap_arbol_directorio, avd)
	EscribirDd(pathArchivo, sb.Sb_ap_detalle_directorio, dd)
	EscribirTInodos(pathArchivo, sb.Sb_ap_tabla_inodo, i)
	EscribirBloques(pathArchivo, i.I_array_bloques[0], db1)
	EscribirBloques(pathArchivo, i.I_array_bloques[1], db2)

	EscribirBitmap(pathArchivo, sb.Sb_ap_bitmap_arbol_directorio, 1, '1')
	EscribirBitmap(pathArchivo, sb.Sb_ap_bitmap_detalle_directorio, 1, '1')
	EscribirBitmap(pathArchivo, sb.Sb_ap_bitmap_tabla_inodo, 1, '1')
	EscribirBitmap(pathArchivo, sb.Sb_ap_bitmap_bloques, 2, '1')

}
func DatosBloquesRoot(texto string) Db {
	var root Db
	copy(root.Db_data[:], texto)
	return root
}
func DatosInodoRoot() I {
	var root I
	root.I_count_inodo = 2
	espacios := "         "
	copy(root.I_id_proper[:], espacios)
	copy(root.I_id_proper[:], "root")
	root.I_count_inodo = 2
	root.I_ap_indirecto = -1
	root.I_array_bloques[0] = sb.Sb_ap_bloques
	root.I_array_bloques[1] = sb.Sb_ap_bloques + sb.Sb_size_struct_bloque
	root.I_size_archivo = 2 * sb.Sb_size_struct_bloque
	return root
}
func DatosDDRoot() Dd {
	var root Dd
	espacios := "                    "
	copy(root.Dd_array_files[0].DdArray_file_date_creacion[:], time.Now().Format("02-01-2006 15:04:05"))
	copy(root.Dd_array_files[0].DdArray_file_date_modificacion[:], time.Now().Format("02-01-2006 15:04:05"))
	copy(root.Dd_array_files[0].DdArray_file_nombre[:], espacios)
	copy(root.Dd_array_files[0].DdArray_file_nombre[:], "users.txt")
	root.Dd_array_files[0].DdArray_file_ap_inodo = sb.Sb_ap_tabla_inodo
	return root
}
func DatosDirectorioRoot() Avd {
	var root Avd
	espacios := "                    "
	copy(root.Avd_nombre_directorio[:], espacios)
	copy(root.Avd_nombre_directorio[:], "/")
	root.Avd_ap_detalle_directorio = sb.Sb_ap_detalle_directorio
	copy(root.Avd_fecha_creacion[:], time.Now().Format("02-01-2006 15:04:05"))
	espacios = "         "
	copy(root.Avd_proper[:], espacios)
	copy(root.Avd_proper[:], "root")
	for i := 0; i < len(avd.Avd_ap_array_subdirectorios); i++ {
		root.Avd_ap_array_subdirectorios[i] = 0
	}
	return root
}

//*****************************LEER LWH**********************************************************
func LeerAvd(path string, comienzo int64) Avd {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Ruta incorrecta")
	}
	archivo.Seek(comienzo, 0)
	s := Avd{}
	tama := int(binary.Size(s))
	datos := LeerByteArchivo(archivo, tama)
	buffer := bytes.NewBuffer(datos)
	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		panic(err)
	}
	return s
}
func LeerDD(path string, comienzo int64) Dd {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Ruta incorrecta")
	}
	archivo.Seek(comienzo, 0)
	s := Dd{}
	tama := int(binary.Size(s))
	datos := LeerByteArchivo(archivo, tama)
	buffer := bytes.NewBuffer(datos)
	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		panic(err)
	}
	return s
}
func LeerTInodo(path string, comienzo int64) I {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Ruta incorrecta")
	}
	archivo.Seek(comienzo, 0)
	s := I{}
	tama := int(binary.Size(s))
	datos := LeerByteArchivo(archivo, tama)
	buffer := bytes.NewBuffer(datos)
	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		panic(err)
	}
	return s
}
func LeerBloque(path string, comienzo int64) Db {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Ruta incorrecta")
	}
	archivo.Seek(comienzo, 0)
	s := Db{}
	tama := int(binary.Size(s))
	datos := LeerByteArchivo(archivo, tama)
	buffer := bytes.NewBuffer(datos)
	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		panic(err)
	}
	return s
}
func LeerSb(path string, comienzo int64) Sb {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Ruta incorrecta")
	}
	archivo.Seek(comienzo, 0)
	s := Sb{}
	tama := int(binary.Size(s))
	datos := LeerByteArchivo(archivo, tama)
	buffer := bytes.NewBuffer(datos)
	err = binary.Read(buffer, binary.BigEndian, &s)
	if err != nil {
		panic(err)
	}
	return s
}
func LeerBitmap(path string, comienzo int64, cantidadBits int64) []byte {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Ruta incorrecta")
	}
	archivo.Seek(comienzo, 0)

	datos := LeerByteArchivo(archivo, int(cantidadBits))

	return datos
}

//***********************************************************************************************
func EscribirSb(path string, comienzo int64, sb Sb) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	e := &sb
	archivo.Seek(comienzo, 0)
	var binario1 bytes.Buffer
	binary.Write(&binario1, binary.BigEndian, e)
	EscribirArchivo(archivo, binario1.Bytes())
}
func EscribirBitmap(path string, comienzo int64, cantidadBit int64, numero byte) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	e := &numero
	archivo.Seek(comienzo, 0)

	var binario1 bytes.Buffer
	for i := 0; i < int(cantidadBit); i++ {
		binary.Write(&binario1, binary.BigEndian, e)
	}
	EscribirArchivo(archivo, binario1.Bytes())
}
func EscribirAvd(path string, comienzo int64, avd Avd) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	e := &avd

	archivo.Seek(comienzo, 0)
	var binario1 bytes.Buffer
	binary.Write(&binario1, binary.BigEndian, e)

	EscribirArchivo(archivo, binario1.Bytes())
}
func EscribirDd(path string, comienzo int64, dd Dd) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	e := &dd

	archivo.Seek(comienzo, 0)
	var binario1 bytes.Buffer

	binary.Write(&binario1, binary.BigEndian, e)

	EscribirArchivo(archivo, binario1.Bytes())
}
func EscribirTInodos(path string, comienzo int64, i I) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	e := &i

	archivo.Seek(comienzo, 0)
	var binario1 bytes.Buffer

	binary.Write(&binario1, binary.BigEndian, e)

	EscribirArchivo(archivo, binario1.Bytes())
}
func EscribirBloques(path string, comienzo int64, db Db) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	e := &db

	archivo.Seek(comienzo, 0)
	var binario1 bytes.Buffer

	binary.Write(&binario1, binary.BigEndian, e)

	EscribirArchivo(archivo, binario1.Bytes())
}
func SearchPath(id string) string {
	pathF := ""
	for i := 0; i < len(listadoMount); i++ {
		if strings.EqualFold(listadoMount[i].IDn, id) {
			pathF = listadoMount[i].Path
		}
	}
	return pathF
}
func SearchNombre(id string) string {
	nombreF := ""
	for i := 0; i < len(listadoMount); i++ {
		if strings.EqualFold(listadoMount[i].IDn, id) {
			nombreF = listadoMount[i].Name
		}
	}
	return nombreF
}
func ObtenerParticion(nombre string, mbr Mbr) Partition {
	var part Partition
	var arrNombre [16]byte
	var espacios string = "                "
	copy(arrNombre[:], espacios)
	copy(arrNombre[:], nombre)
	nString := string(arrNombre[:])
	if strings.EqualFold(string(mbr.Mbr_partition_1.Part_name[:]), nString) {
		part = mbr.Mbr_partition_1
	} else if strings.EqualFold(string(mbr.Mbr_partition_2.Part_name[:]), nString) {
		part = mbr.Mbr_partition_2
	} else if strings.EqualFold(string(mbr.Mbr_partition_3.Part_name[:]), nString) {
		part = mbr.Mbr_partition_3
	} else if strings.EqualFold(string(mbr.Mbr_partition_4.Part_name[:]), nString) {
		part = mbr.Mbr_partition_4
	}
	return part
}

//******************************************************************************************************

//*************************** CREAR DISCO **************************************************************

func CrearDisco(size string, path string, name string, unit string) {
	if size == "" || path == "" || name == "" {
		fmt.Println("Atributo obligatorio")
	} else {
		if s, _ := strconv.Atoi(size); s <= 0 {
			panic("No se crea disco, por número invalido")
		} else {
			name = strings.ReplaceAll(name, " ", "")
			path = strings.ReplaceAll(path, "\"", "")
			CrearCarpeta(path)
			nsize, _ := strconv.ParseInt(size, 10, 64)
			if strings.EqualFold(unit, "k") {
				CrearArchivo(path+name, nsize*1024)
			} else if strings.EqualFold(unit, "m") || unit == "" {
				CrearArchivo(path+name, nsize*1024*1024)
			}
		}

	}
}
func CrearCarpeta(pathn string) {
	_, err := os.Stat(pathn)
	if os.IsNotExist(err) {
		tr := exec.Command("mkdir", "-p", pathn)
		err = tr.Run()
		if err != nil {
			panic(err)
		} else {

		}
	} else {
	}
}
func CrearArchivo(path string, tama int64) {
	archivo, err := os.Create(path)
	if err != nil {
		fmt.Println("Error")
		panic(err)
		return
	}

	var i int8 = 0
	s := &i

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	EscribirArchivo(archivo, binario.Bytes())

	archivo.Seek(tama-1, 0)

	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	EscribirArchivo(archivo, binario2.Bytes())

	archivo.Seek(0, 0)
	mbr.Mbr_tamano = tama
	mbr.Mbr_disk_signature = rand.Int63()
	fecha := FormatoFecha()
	letraFecha := []byte(fecha)
	for i := 0; i < len(letraFecha); i++ {
		mbr.Mbr_fecha_creacion[i] = letraFecha[i]
	}
	mbr.Mbr_partition_1.Part_status = '0'
	mbr.Mbr_partition_2.Part_status = '0'
	mbr.Mbr_partition_3.Part_status = '0'
	mbr.Mbr_partition_4.Part_status = '0'

	mbr.Mbr_partition_1.Part_fit = ' '
	mbr.Mbr_partition_2.Part_fit = ' '
	mbr.Mbr_partition_3.Part_fit = ' '
	mbr.Mbr_partition_4.Part_fit = ' '

	mbr.Mbr_partition_1.Part_type = ' '
	mbr.Mbr_partition_2.Part_type = ' '
	mbr.Mbr_partition_3.Part_type = ' '
	mbr.Mbr_partition_4.Part_type = ' '

	mbr.Mbr_partition_1.Part_size = 0
	mbr.Mbr_partition_2.Part_size = 0
	mbr.Mbr_partition_3.Part_size = 0
	mbr.Mbr_partition_4.Part_size = 0
	var espacios string = "                "
	copy(mbr.Mbr_partition_1.Part_name[:], espacios)
	copy(mbr.Mbr_partition_2.Part_name[:], espacios)
	copy(mbr.Mbr_partition_3.Part_name[:], espacios)
	copy(mbr.Mbr_partition_4.Part_name[:], espacios)

	s1 := &mbr
	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, s1)
	EscribirArchivo(archivo, binario3.Bytes())
}
func EscribirArchivo(archivo *os.File, bytes []byte) {
	_, err := archivo.Write(bytes)
	if err != nil {
		panic(err)

	}
}
func FormatoFecha() string {
	t := time.Now()
	return t.Format("01-02-2006")
}

//******************************************************************************************************

//*************************** CREAR PARTICION **************************************************************

func CrearParticion(size string, unit string, path string, tip string, fit string, name string) {
	if size == "" || path == "" || name == "" {
		fmt.Println("Faltan campos Obligatorios")
	} else {
		p := ObtenerParticion(name, mbr)
		if p.Part_size != 0 {
			fmt.Println("Particion con ese nombre ya existe")
			return
		}
		nsize := GetUnit(unit, size)
		totalPart += nsize
		if ParticionvsMbr(nsize) {
			fmt.Println("Particion mayor a Tamaño Disco")
			totalPart -= nsize
			return
		} else if SumaMayorDisco() {
			fmt.Println("Sin espacio en Disco")
			totalPart -= nsize
			return
		}
		fit = GetFit(fit)
		if strings.EqualFold(tip, "") {
			tip = "p"
		}
		espacios := "                "
		nFit := []byte(fit)
		nTipo := []byte(tip)
		if mbr.Mbr_partition_1.Part_status == '0' {
			if strings.EqualFold(tip, "e") {
				tipo = ""
				if EsExtendida(mbr.Mbr_partition_2) || EsExtendida(mbr.Mbr_partition_3) || EsExtendida(mbr.Mbr_partition_4) {
					fmt.Println("Ya se creo particion Ex")
					totalPart -= nsize
				} else {
					mbr.Mbr_partition_1.Part_status = '1'
					mbr.Mbr_partition_1.Part_fit = nFit[0]
					mbr.Mbr_partition_1.Part_type = nTipo[0]
					mbr.Mbr_partition_1.Part_size = nsize
					mbr.Mbr_partition_1.Part_start = int64(binary.Size(mbr))
					copy(mbr.Mbr_partition_1.Part_name[:], espacios)
					copy(mbr.Mbr_partition_1.Part_name[:], name)
					EscribirParticion(path, '1')
				}
			} else {
				mbr.Mbr_partition_1.Part_status = '1'
				mbr.Mbr_partition_1.Part_fit = nFit[0]
				mbr.Mbr_partition_1.Part_type = nTipo[0]
				mbr.Mbr_partition_1.Part_size = nsize
				mbr.Mbr_partition_1.Part_start = int64(binary.Size(mbr))
				copy(mbr.Mbr_partition_1.Part_name[:], espacios)
				copy(mbr.Mbr_partition_1.Part_name[:], name)
				EscribirParticion(path, '1')
			}
		} else if mbr.Mbr_partition_2.Part_status == '0' {
			if strings.EqualFold(tip, "e") {
				tipo = ""
				if EsExtendida(mbr.Mbr_partition_1) || EsExtendida(mbr.Mbr_partition_3) || EsExtendida(mbr.Mbr_partition_4) {
					fmt.Println("Ya se creo particion")
					totalPart -= nsize
				} else {
					mbr.Mbr_partition_2.Part_status = '1'
					mbr.Mbr_partition_2.Part_fit = nFit[0]
					mbr.Mbr_partition_2.Part_type = nTipo[0]
					mbr.Mbr_partition_2.Part_size = nsize
					mbr.Mbr_partition_2.Part_start = int64(mbr.Mbr_partition_1.Part_start+mbr.Mbr_partition_1.Part_size) + 1
					copy(mbr.Mbr_partition_2.Part_name[:], espacios)
					copy(mbr.Mbr_partition_2.Part_name[:], name)
					EscribirParticion(path, '2')
				}
			} else {
				mbr.Mbr_partition_2.Part_status = '1'
				mbr.Mbr_partition_2.Part_fit = nFit[0]
				mbr.Mbr_partition_2.Part_type = nTipo[0]
				mbr.Mbr_partition_2.Part_size = nsize
				mbr.Mbr_partition_2.Part_start = int64(mbr.Mbr_partition_1.Part_start+mbr.Mbr_partition_1.Part_size) + 1
				copy(mbr.Mbr_partition_2.Part_name[:], espacios)
				copy(mbr.Mbr_partition_2.Part_name[:], name)
				EscribirParticion(path, '2')
			}
		} else if mbr.Mbr_partition_3.Part_status == '0' {
			if strings.EqualFold(tip, "e") {
				tipo = ""
				if EsExtendida(mbr.Mbr_partition_1) || EsExtendida(mbr.Mbr_partition_2) || EsExtendida(mbr.Mbr_partition_4) {
					fmt.Println("Ya se creo particion")
					totalPart -= nsize

				} else {
					mbr.Mbr_partition_3.Part_status = '1'
					mbr.Mbr_partition_3.Part_fit = nFit[0]
					mbr.Mbr_partition_3.Part_type = nTipo[0]
					mbr.Mbr_partition_3.Part_size = nsize
					mbr.Mbr_partition_3.Part_start = int64(mbr.Mbr_partition_2.Part_start+mbr.Mbr_partition_2.Part_size) + 1
					copy(mbr.Mbr_partition_3.Part_name[:], name)
					EscribirParticion(path, '3')
				}
			} else {
				mbr.Mbr_partition_3.Part_status = '1'
				mbr.Mbr_partition_3.Part_fit = nFit[0]
				mbr.Mbr_partition_3.Part_type = nTipo[0]
				mbr.Mbr_partition_3.Part_size = nsize
				mbr.Mbr_partition_3.Part_start = int64(mbr.Mbr_partition_2.Part_start+mbr.Mbr_partition_2.Part_size) + 1
				copy(mbr.Mbr_partition_3.Part_name[:], name)
				EscribirParticion(path, '3')
			}
		} else if mbr.Mbr_partition_4.Part_status == '0' {
			if strings.EqualFold(tip, "e") {
				tipo = ""
				if EsExtendida(mbr.Mbr_partition_1) || EsExtendida(mbr.Mbr_partition_3) || EsExtendida(mbr.Mbr_partition_2) {
					fmt.Println("Ya se creo particion")
					totalPart -= nsize

				} else {
					mbr.Mbr_partition_4.Part_status = '1'
					mbr.Mbr_partition_4.Part_fit = nFit[0]
					mbr.Mbr_partition_4.Part_type = nTipo[0]
					mbr.Mbr_partition_4.Part_size = nsize
					mbr.Mbr_partition_4.Part_start = int64(mbr.Mbr_partition_3.Part_start+mbr.Mbr_partition_3.Part_size) + 1
					copy(mbr.Mbr_partition_4.Part_name[:], name)
					EscribirParticion(path, '4')
				}
			} else {
				mbr.Mbr_partition_4.Part_status = '1'
				mbr.Mbr_partition_4.Part_fit = nFit[0]
				mbr.Mbr_partition_4.Part_type = nTipo[0]
				mbr.Mbr_partition_4.Part_size = nsize
				mbr.Mbr_partition_4.Part_start = int64(mbr.Mbr_partition_3.Part_start+mbr.Mbr_partition_3.Part_size) + 1
				copy(mbr.Mbr_partition_4.Part_name[:], name)
				EscribirParticion(path, '4')
			}
		} else {
			fmt.Println("Sin particiones Disponibles")
		}
	}
}
func ParticionvsMbr(size int64) bool {
	return size >= mbr.Mbr_tamano
}
func SumaMayorDisco() bool {
	return totalPart >= mbr.Mbr_tamano
}
func EscribirParticion(path string, number byte) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	archivo.Seek(0, 0)
	var binario3 bytes.Buffer
	ebrN := Ebr{}
	s := &mbr
	binary.Write(&binario3, binary.BigEndian, s)
	EscribirArchivo(archivo, binario3.Bytes())

	if strings.EqualFold(string(mbr.Mbr_partition_1.Part_type), "e") && number == '1' {
		ebrN.Part_next = -1
		ebrN.Part_status = '0'
		ebrN.Part_size = 0
		ebrN.Part_start = 0
		s := &ebrN
		archivo.Seek(mbr.Mbr_partition_1.Part_start, 0)
		var binario1 bytes.Buffer
		binary.Write(&binario1, binary.BigEndian, s)
		EscribirArchivo(archivo, binario1.Bytes())
	} else if strings.EqualFold(string(mbr.Mbr_partition_2.Part_type), "e") && number == '2' {
		ebrN.Part_next = -1
		ebrN.Part_status = '0'
		ebrN.Part_size = 0
		ebrN.Part_start = 0
		s := &ebrN
		archivo.Seek(mbr.Mbr_partition_2.Part_start, 0)
		var binario1 bytes.Buffer
		binary.Write(&binario1, binary.BigEndian, s)
		EscribirArchivo(archivo, binario1.Bytes())
	} else if strings.EqualFold(string(mbr.Mbr_partition_3.Part_type), "e") && number == '3' {
		ebrN.Part_next = -1
		ebrN.Part_status = '0'
		ebrN.Part_size = 0
		ebrN.Part_start = 0
		s := &ebrN
		archivo.Seek(mbr.Mbr_partition_3.Part_start, 0)
		var binario1 bytes.Buffer
		binary.Write(&binario1, binary.BigEndian, s)
		EscribirArchivo(archivo, binario1.Bytes())
	} else if strings.EqualFold(string(mbr.Mbr_partition_4.Part_type), "e") && number == '4' {
		ebrN.Part_next = -1
		ebrN.Part_status = '0'
		ebrN.Part_size = 0
		ebrN.Part_start = 0
		s := &ebrN
		archivo.Seek(mbr.Mbr_partition_4.Part_start, 0)
		var binario1 bytes.Buffer
		binary.Write(&binario1, binary.BigEndian, s)
		EscribirArchivo(archivo, binario1.Bytes())
	} else {
		if number == '1' {
			archivo.Seek(mbr.Mbr_partition_1.Part_start, 0)
			var binario1 bytes.Buffer
			for i := 0; i < int(mbr.Mbr_partition_1.Part_size); i++ {
				binary.Write(&binario1, binary.BigEndian, &number)
			}
			EscribirArchivo(archivo, binario1.Bytes())
		} else if number == '2' {
			archivo.Seek(mbr.Mbr_partition_2.Part_start, 0)
			var binario2 bytes.Buffer
			for i := 0; i < int(mbr.Mbr_partition_2.Part_size); i++ {
				binary.Write(&binario2, binary.BigEndian, &number)
			}
			EscribirArchivo(archivo, binario2.Bytes())
		} else if number == '3' {
			archivo.Seek(mbr.Mbr_partition_3.Part_start, 0)
			var binario3f bytes.Buffer
			for i := 0; i < int(mbr.Mbr_partition_3.Part_size); i++ {
				binary.Write(&binario3f, binary.BigEndian, &number)
			}
			EscribirArchivo(archivo, binario3f.Bytes())
		} else if number == '4' {
			archivo.Seek(mbr.Mbr_partition_4.Part_start, 0)
			var binario4 bytes.Buffer
			for i := 0; i < int(mbr.Mbr_partition_4.Part_size); i++ {
				binary.Write(&binario4, binary.BigEndian, &number)
			}
			EscribirArchivo(archivo, binario4.Bytes())
		}
	}

}
func EsExtendida(partition Partition) bool {
	return strings.EqualFold(string(partition.Part_type), "e")
}
func GetFit(fit string) string {
	if strings.EqualFold(fit, "bf") {
		return "b"
	} else if strings.EqualFold(fit, "ff") {
		return "f"
	} else if strings.EqualFold(fit, "wf") || fit == "" {
		return "w"
	}
	return ""
}
func GetUnit(unit string, size string) int64 {
	nsize, _ := strconv.ParseInt(size, 10, 64)
	if strings.EqualFold(unit, "k") || unit == "" {
		nsize = nsize * 1024
		return nsize
	} else if strings.EqualFold(unit, "m") {
		nsize = nsize * 1024 * 1024
		return nsize
	} else if strings.EqualFold(unit, "b") {

	}
	return nsize
}

//******************************************************************************************************

//****************************** LEER MBR **************************************************************

func LeerMBR(path string) Mbr {
	path = strings.ReplaceAll(path, "\"", "")
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Ruta incorrecta")

	}
	mbr := Mbr{}
	tam := int(unsafe.Sizeof(mbr))
	datos := LeerByteArchivo(archivo, tam)
	buffer := bytes.NewBuffer(datos)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		panic(err)
	}
	return mbr
}
func LeerByteArchivo(archivo *os.File, numero int) []byte {
	bytes := make([]byte, numero)
	_, err := archivo.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

//******************************************************************************************************

//****************************** REPORTES **************************************************************

func CrearReportes(nombre string, pathN string, id string, ruta string) {
	if nombre == "" || path == "" || id == "" {
		fmt.Println("Faltan datos obligatorios")
	} else {
		CrearCarpeta(pathN)
		if strings.EqualFold(nombre, "mbr") {
			fmt.Println("SMBR")
			pathArchivo := SearchPath(id)
			ReporteMBR(nombre, pathArchivo)
		} else if strings.EqualFold(nombre, "disk") {
			fmt.Println("SDISK")
			pathArchivo := SearchPath(id)
			ReporteParticiones(nombre, pathArchivo)
		} else if strings.EqualFold(nombre, "sb") {
			fmt.Println("SSB")
		} else if strings.EqualFold(nombre, "bm_arbdir") {
			fmt.Println("SBM-ARBDIR")
		} else if strings.EqualFold(nombre, "bm_detdir") {
			fmt.Println("SBM-DETDIR")
		} else if strings.EqualFold(nombre, "bm_inode") {
			fmt.Println("SBM-INODE")
		} else if strings.EqualFold(nombre, "bm_block") {
			fmt.Println("SBM-BLOCK")
		} else if strings.EqualFold(nombre, "bitacora") {
			fmt.Println("BITACORA")
		} else if strings.EqualFold(nombre, "directorio") {
			fmt.Println("DIRECTORIO")
		} else if strings.EqualFold(nombre, "tree_file") {
			pathArchivo := SearchPath(id)
			nn := SearchNombre(id)
			if ruta == "" {
				ruta = "/"
			}
			ReporteTreeFile(nombre, pathArchivo, id, ruta, nn)
			fmt.Println("TREE-FILE")
		} else if strings.EqualFold(nombre, "tree_directorio") {
			fmt.Println("TREE-DIRECTORIO")
		} else if strings.EqualFold(nombre, "tree_complete") {
			fmt.Println("TREE-COMPLETE")
		} else if strings.EqualFold(nombre, "ls") {
			fmt.Println("SLS")
		} else {
			fmt.Println("Parametro incorrecto")
		}
	}
}
func ReporteMBR(nombre string, p string) {
	if strings.EqualFold(p, "") {
		fmt.Println("Particion no montada ")
		return
	}
	mbr := LeerMBR(p)
	pp, fn := filepath.Split(path)
	archivo, err := os.Create(pp + fn + ".dot")
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	var c int
	reportMBR := "digraph rMBR {\n"
	reportMBR += " abc[shape=none, margin=0,label=<\n<TABLE BORDER =\"0\" CELLBORDER =\"1\" CELLSPACING =\"0\" CELLPADDING=\"4\">\n"
	reportMBR += "<TR><TD BGCOLOR=\"lightblue\">Nombre</TD>\n<TD BGCOLOR=\"lightblue\">Descripcion</TD></TR>\n"
	reportMBR += "<TR><TD>mbr_tamaño</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_tamano, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>mbr_fecha_creacion</TD>" + "<TD>" + string(mbr.Mbr_fecha_creacion[:]) + "</TD></TR>\n"
	reportMBR += "<TR><TD>mbr_disk_signature</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_disk_signature, 10) + "</TD></TR>\n"

	c = contadorNombre(mbr.Mbr_partition_1)
	reportMBR += "<TR><TD>part_Nombre_1</TD>" + "<TD>" + string(mbr.Mbr_partition_1.Part_name[:c]) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Fit_1</TD>" + "<TD>" + string(mbr.Mbr_partition_1.Part_fit) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Size_1</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_partition_1.Part_size, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Start_1</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_partition_1.Part_start, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Status_1</TD>" + "<TD>" + string(mbr.Mbr_partition_1.Part_status) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Type_1</TD>" + "<TD>" + string(mbr.Mbr_partition_1.Part_type) + "</TD></TR>\n"

	c = contadorNombre(mbr.Mbr_partition_2)
	reportMBR += "<TR><TD>part_Nombre_2</TD>" + "<TD>" + string(mbr.Mbr_partition_2.Part_name[:c]) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Fit_2</TD>" + "<TD>" + string(mbr.Mbr_partition_2.Part_fit) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Size_2</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_partition_2.Part_size, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Start_2</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_partition_2.Part_start, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Status_2</TD>" + "<TD>" + string(mbr.Mbr_partition_2.Part_status) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Type_2</TD>" + "<TD>" + string(mbr.Mbr_partition_2.Part_type) + "</TD></TR>\n"

	c = contadorNombre(mbr.Mbr_partition_3)
	reportMBR += "<TR><TD>part_Nombre_3</TD>" + "<TD>" + string(mbr.Mbr_partition_3.Part_name[:c]) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Fit_3</TD>" + "<TD>" + string(mbr.Mbr_partition_3.Part_fit) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Size_3</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_partition_3.Part_size, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Start_3</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_partition_3.Part_start, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Status_3</TD>" + "<TD>" + string(mbr.Mbr_partition_3.Part_status) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Type_3</TD>" + "<TD>" + string(mbr.Mbr_partition_3.Part_type) + "</TD></TR>\n"

	c = contadorNombre(mbr.Mbr_partition_4)
	reportMBR += "<TR><TD>part_Nombre_4</TD>" + "<TD>" + string(mbr.Mbr_partition_4.Part_name[:c]) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Fit_4</TD>" + "<TD>" + string(mbr.Mbr_partition_4.Part_fit) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Size_4</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_partition_4.Part_size, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Start_4</TD>" + "<TD>" + strconv.FormatInt(mbr.Mbr_partition_4.Part_start, 10) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Status_4</TD>" + "<TD>" + string(mbr.Mbr_partition_4.Part_status) + "</TD></TR>\n"
	reportMBR += "<TR><TD>part_Type_4</TD>" + "<TD>" + string(mbr.Mbr_partition_4.Part_type) + "</TD></TR>\n"

	reportMBR += "</TABLE>>];\n}"
	archivo.WriteString(reportMBR)
	archivo.Sync()
	archivo.Close()

	tr := exec.Command("dot", "-Tjpg", pp+fn+".dot", "-o", pp+fn+".jpg")
	out, m := tr.CombinedOutput()
	if m != nil {
		fmt.Println(fmt.Sprint(m) + ": " + string(out))
	}
}
func contadorNombre(particion Partition) int {
	var cuenta int
	for i := 0; i < 16; i++ {
		if particion.Part_name[i] != ' ' {
			cuenta++
		}

	}
	return cuenta
}
func ReporteSb() {

}

var Extendida bool = false

func ReporteParticiones(nombre string, pathN string) {
	if strings.EqualFold(pathN, "") {
		fmt.Println("Particion no montada ")
		return
	}
	pp, fn := filepath.Split(path)
	archivo, err := os.Create(pp + fn + ".dot")
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	mbr = LeerMBR(pathN)
	var colSpan int = 0
	var particionActual Partition
	var nombrePart string = "LIBRE"
	reporteParticiones := "digraph rParticiones {\n"
	reporteParticiones += "abd[shape=none, margin =0, label=<\n"
	reporteParticiones += "<TABLE BORDER=\"1\" CELLBORDER =\"1\" CELLSPACING = \"3\" CELLPADDING=\"4\">\n"
	reporteParticiones += "<TR><TD ROWSPAN=\"3\">MBR</TD>\n"
	if mbr.Mbr_partition_1.Part_status != '0' {
		if strings.EqualFold(string(mbr.Mbr_partition_1.Part_type), "p") {
			reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#AACADC\">PRIMARIA</TD>\n"
		} else if strings.EqualFold(string(mbr.Mbr_partition_1.Part_type), "e") {
			ebr = LeerEBR(pathN, mbr.Mbr_partition_1.Part_start)
			colSpan = TotalEbr(ebr, pathN)
			Extendida = true
			particionActual = mbr.Mbr_partition_1
			reporteParticiones += "<TD COLSPAN=\"" + strconv.Itoa(colSpan) + "\" BGCOLOR=\"#EF775D\">EXTENDIDA</TD>\n"
		}
	} else {
		reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#96FF33\">" + nombrePart + "</TD>\n"
	}
	if mbr.Mbr_partition_2.Part_status != '0' {
		if strings.EqualFold(string(mbr.Mbr_partition_2.Part_type), "p") {
			reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#AACADC\">PRIMARIA</TD>\n"
		} else if strings.EqualFold(string(mbr.Mbr_partition_2.Part_type), "e") {
			ebr = LeerEBR(pathN, mbr.Mbr_partition_2.Part_start)
			colSpan = TotalEbr(ebr, pathN)
			Extendida = true
			particionActual = mbr.Mbr_partition_2
			reporteParticiones += "<TD COLSPAN=\"" + strconv.Itoa(colSpan) + "\" BGCOLOR=\"#EF775D\">EXTENDIDA</TD>\n"
		}

	} else {
		reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#96FF33\">" + nombrePart + "</TD>\n"
	}
	if mbr.Mbr_partition_3.Part_status != '0' {
		if strings.EqualFold(string(mbr.Mbr_partition_3.Part_type), "p") {
			reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#AACADC\">PRIMARIA</TD>\n"
		} else if strings.EqualFold(string(mbr.Mbr_partition_3.Part_type), "e") {
			ebr = LeerEBR(pathN, mbr.Mbr_partition_3.Part_start)
			colSpan = TotalEbr(ebr, pathN)
			Extendida = true
			particionActual = mbr.Mbr_partition_3
			reporteParticiones += "<TD COLSPAN=\"" + strconv.Itoa(colSpan) + "\" BGCOLOR=\"#EF775D\">EXTENDIDA</TD>\n"
		}

	} else {
		reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#96FF33\">" + nombrePart + "</TD>\n"
	}
	if mbr.Mbr_partition_4.Part_status != '0' {
		if strings.EqualFold(string(mbr.Mbr_partition_4.Part_type), "p") {
			reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#AACADC\">PRIMARIA</TD>\n"
		} else if strings.EqualFold(string(mbr.Mbr_partition_4.Part_type), "e") {
			ebr = LeerEBR(pathN, mbr.Mbr_partition_4.Part_start)
			colSpan = TotalEbr(ebr, pathN)
			Extendida = true
			particionActual = mbr.Mbr_partition_4
			reporteParticiones += "<TD COLSPAN=\"" + strconv.Itoa(colSpan) + "\" BGCOLOR=\"#EF775D\">EXTENDIDA</TD>\n"
		}
	} else {
		reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#96FF33\">" + nombrePart + "</TD>\n"
	}
	if mbr.Mbr_partition_1.Part_size+mbr.Mbr_partition_2.Part_size+mbr.Mbr_partition_3.Part_size+mbr.Mbr_partition_4.Part_size < mbr.Mbr_tamano {
		reporteParticiones += "<TD ROWSPAN=\"3\" BGCOLOR=\"#96FF33\">" + nombrePart + "</TD>"
	}
	reporteParticiones += "</TR>\n"

	if Extendida {
		ebr = LeerEBR(pathN, particionActual.Part_start)
		if ebr.Part_status == '1' {
			reporteParticiones += "<TR>\n"
			reporteParticiones += RecorrerEbr(ebr, pathN)
			reporteParticiones += "</TR>\n"
		}
	}

	reporteParticiones += "</TABLE>>];\n}"
	cadena = ""
	total = 1
	archivo.WriteString(reporteParticiones)
	archivo.Sync()
	archivo.Close()

	tr := exec.Command("dot", "-Tjpg", pp+fn+".dot", "-o", pp+fn+".jpg")
	out, m := tr.CombinedOutput()
	if m != nil {
		fmt.Println(fmt.Sprint(m) + ": " + string(out))
	}
}

var cadena string = ""

func RecorrerEbr(e Ebr, pathN string) string {
	if e.Part_status != '0' {
		if e.Part_next == -1 {
			cadena += "<TD BGCOLOR=\"#F0E23A\">LOGICA</TD>\n"
		} else {
			cadena += "<TD BGCOLOR=\"#F0E23A\">LOGICA</TD>\n"
			e = LeerEBR(pathN, e.Part_next)
			RecorrerEbr(e, pathN)
		}

	}
	return cadena
}

var total int = 1

func TotalEbr(e Ebr, path string) int {
	if e.Part_next != -1 {
		e = LeerEBR(path, e.Part_next)
		total++
		TotalEbr(e, path)
	}
	return total
}
func ReporteTreeFile(nombre string, p string, idR string, r string, nomp string) {
	//mbr := LeerMBR(p)
	fmt.Println("NN " + nombre)
	fmt.Println("P " + p)
	fmt.Println("IDR " + idR)
	fmt.Println("R " + r)
	archivo, err := os.Create(path + nombre + ".dot")
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	if nombre == "" || p == "" || idR == "" || r == "" {
		fmt.Println("Faltan datos")
	} else {
		if strings.EqualFold(p, "") {
			fmt.Println("Particion no montada ")
			return
		}
		fmt.Println(path + nombre)
		mbr := LeerMBR(p)
		archivo, err := os.Create(path + nombre + ".dot")
		defer archivo.Close()
		if err != nil {
			panic(err)
		}
		partAct := ObtenerParticion(nomp, mbr)
		sb := LeerSb(p, partAct.Part_start)
		fmt.Println("Inicio part " + strconv.FormatInt(partAct.Part_start, 10))
		avd = LeerAvd(p, sb.Sb_ap_arbol_directorio)
		dd = LeerDD(p, sb.Sb_ap_detalle_directorio)
		reportTree := "digraph rMBR {\n"
		reportTree += "node[shape=plaintext]\nstruct1[\nlabel=<"
		reportTree += "<TABLE>\n<TR><TD BGCOLOR=\"lightblue\" COLSPAN=\"2\">" + string(avd.Avd_nombre_directorio[:]) + "</TD></TR>\n"
		reportTree += "<TR><TD>Fecha Creacion</TD><TD>" + string(avd.Avd_fecha_creacion[:]) + "</TD></TR>\n"
		var i int

		reportTree += "<TR><TD PORT=\"" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[0], 10) + "\">APD1</TD><TD>" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[0], 10) + "</TD></TR>\n"
		reportTree += "<TR><TD PORT=\"" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[1], 10) + "\">APD2</TD><TD>" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[1], 10) + "</TD></TR>\n"
		reportTree += "<TR><TD PORT=\"" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[2], 10) + "\">APD3</TD><TD>" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[2], 10) + "</TD></TR>\n"
		reportTree += "<TR><TD PORT=\"" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[3], 10) + "\">APD4</TD><TD>" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[3], 10) + "</TD></TR>\n"
		reportTree += "<TR><TD PORT=\"" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[4], 10) + "\">APD5</TD><TD>" + strconv.FormatInt(avd.Avd_ap_array_subdirectorios[4], 10) + "</TD></TR>\n"

		reportTree += "<TR><TD PORT=\"f" + strconv.FormatInt(avd.Avd_ap_detalle_directorio, 10) + "\">Detalle Directorio</TD><TD>" + strconv.FormatInt(avd.Avd_ap_detalle_directorio, 10) + "</TD></TR>\n"
		i++
		reportTree += "<TR><TD PORT=\"" + strconv.Itoa(i) + "\">API </TD><TD>" + strconv.FormatInt(avd.Avd_ap_arbol_virtual_directorio, 10) + "</TD></TR>\n"
		reportTree += "<TR><TD>Proper</TD><TD>" + string(avd.Avd_proper[:]) + "</TD></TR>\n</TABLE>>];\n"
		reportTree += "struct2[\nlabel=<"
		reportTree += "<TABLE>\n<TR><TD BGCOLOR=\"lightblue\">DETALLE DIRECTORIO</TD><TD></TD></TR>"
		reportTree += "<TR><TD>" + string(dd.Dd_array_files[0].DdArray_file_nombre[:]) + "</TD><TD>" + strconv.FormatInt(dd.Dd_array_files[0].DdArray_file_ap_inodo, 10) + "</TD></TR>\n"
		reportTree += "<TR><TD>APD2</TD><TD>" + strconv.FormatInt(dd.Dd_array_files[1].DdArray_file_ap_inodo, 10) + "</TD></TR>\n"
		reportTree += "<TR><TD>APD3</TD><TD>" + strconv.FormatInt(dd.Dd_array_files[2].DdArray_file_ap_inodo, 10) + "</TD></TR>\n"
		reportTree += "<TR><TD>APD4</TD><TD>" + strconv.FormatInt(dd.Dd_array_files[3].DdArray_file_ap_inodo, 10) + "</TD></TR>\n"
		reportTree += "<TR><TD>APD5</TD><TD>" + strconv.FormatInt(dd.Dd_array_files[4].DdArray_file_ap_inodo, 10) + "</TD></TR>\n"
		reportTree += "<TR><TD>APDI</TD><TD>" + strconv.FormatInt(dd.Dd_ap_detalle_directorio, 10) + "</TD></TR>\n"
		reportTree += "</TABLE>>];\n}"
		archivo.WriteString(reportTree)
		archivo.Sync()
		archivo.Close()

		tr := exec.Command("dot", "-Tjpg", path+nombre+".dot", "-o", path+nombre+".jpg")
		out, m := tr.CombinedOutput()
		if m != nil {
			fmt.Println(fmt.Sprint(m) + ": " + string(out))
		}
	}
}

//******************************************************************************************************

//****************************** MONTAR Y DESMONTAR PARTICION ***************************************************
func DesmontarParticion(id []string) {
	for i := 0; i < len(listadoMount); i++ {
		for j := 0; j < len(id); j++ {
			if listadoMount[i].IDn == id[j] {
				listadoMount = Desmontar(listadoMount, i)
			}
		}

	}
}
func Desmontar(s []ParticionMontada, index int) []ParticionMontada {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func MontarParticion(pathN string, name string) {
	fmt.Println(pathN)
	_, e := os.Stat(pathN)
	if e != nil {
		fmt.Println("Disco no existe ")
		return
	}
	mbr = LeerMBR(pathN)
	var arrnombre [16]byte
	espacios := "                "
	copy(arrnombre[:], espacios)
	copy(arrnombre[:], name)
	nombreE := string(arrnombre[:])
	if string(mbr.Mbr_partition_1.Part_name[:]) != nombreE && string(mbr.Mbr_partition_2.Part_name[:]) != nombreE && string(mbr.Mbr_partition_3.Part_name[:]) != nombreE && string(mbr.Mbr_partition_4.Part_name[:]) != nombreE {
		fmt.Println("Particion " + nombreE + " no existe en disco, pruebe nuevamente")
	} else {
		for i := 0; i < len(listadoMount); i++ {
			if name == listadoMount[i].Name {
				fmt.Println("Particion " + name + " ya esta montada")
				return
			}
		}
		montar := ParticionMontada{}
		l := AsignarLiteral(path)
		incremento = AsignarNumero(&incremento)
		montar.Letra = l
		montar.Name = name
		montar.Path = path
		montar.IDn = "vd" + l + strconv.Itoa(incremento)
		listadoMount = append(listadoMount, montar)
	}
}
func ImprimirMontadas(lista []ParticionMontada) {
	for i := 0; i < len(lista); i++ {
		fmt.Println("-path->"+lista[i].Path+" -name->"+lista[i].Name, " -id->"+lista[i].IDn)
	}
}
func AsignarNumero(incremento *int) int {
	return *incremento + 1

}

var al int = 0

func AsignarLiteral(path string) string {

	for i := 0; i < len(listadoMount); i++ {
		if path == listadoMount[i].Path {
			return listadoMount[i].Letra
		}
	}
	letra := abecedario[al]
	al++
	return letra
}

//******************************************************************************************************

func EliminarParticion(eliminar string, name string, path string) {
	mbrEvaluar := &mbr
	if name != "" || path != "" {
		p := ParticionAEliminar(mbrEvaluar, name)
		fmt.Println("Particion a eliminar " + strconv.Itoa(p))
		if strings.EqualFold(eliminar, "fast") {
			if p == 1 || p == 2 || p == 3 || p == 4 {
				EscribirParticionFast(path, mbr)
			} else if p == 0 {
				fmt.Println("Particion no pertenece a disco " + path)
			}
		} else if strings.EqualFold(eliminar, "full") {
			EscribirParticionFast(path, mbr)
			switch p {
			case 1:
				fmt.Println("Particion a eliminar 1")
				EliminarFull(mbrEvaluar.Mbr_partition_1, path)
			case 2:
				fmt.Println("Particion a eliminar 2")
				EliminarFull(mbrEvaluar.Mbr_partition_2, path)
			case 3:
				fmt.Println("Particion a eliminar 3")
				EliminarFull(mbrEvaluar.Mbr_partition_3, path)
			case 4:
				fmt.Println("Particion a eliminar 4")
				EliminarFull(mbrEvaluar.Mbr_partition_4, path)
			default:
				fmt.Println("Particion no pertenece a disco " + path)
				break
			}
		}
	} else {

	}
}
func EliminarFull(particion Partition, path string) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	var bite byte = 0
	n := &bite
	archivo.Seek(particion.Part_start, 0)

	var binario3 bytes.Buffer
	for i := 0; i < int(particion.Part_size); i++ {
		binary.Write(&binario3, binary.BigEndian, n)
	}
	EscribirArchivo(archivo, binario3.Bytes())

}
func EscribirParticionFast(path string, mbr Mbr) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	archivo.Seek(0, 0)
	var binario3 bytes.Buffer
	s := &mbr
	//for i := 0; i < int(mbr.Part_size); i++ {
	binary.Write(&binario3, binary.BigEndian, s)
	EscribirArchivo(archivo, binario3.Bytes())
	//}
}
func ParticionAEliminar(mbr *Mbr, nombre string) int {
	var arrnombre [16]byte
	copy(arrnombre[:], "                ")
	copy(arrnombre[:], nombre)
	nombreE := string(arrnombre[:])
	var del int
	if strings.EqualFold(nombreE, string(mbr.Mbr_partition_1.Part_name[:])) && mbr.Mbr_partition_1.Part_status == '1' {
		mbr.Mbr_partition_1.Part_status = '0'
		mbr.Mbr_partition_1.Part_size = 0
		mbr.Mbr_partition_1.Part_start = 0
		mbr.Mbr_partition_1.Part_fit = ' '
		mbr.Mbr_partition_1.Part_type = ' '
		var espacio string = "                "
		copy(mbr.Mbr_partition_1.Part_name[:], espacio)
		del = 1
	} else if strings.EqualFold(nombreE, string(mbr.Mbr_partition_2.Part_name[:])) && mbr.Mbr_partition_2.Part_status == '1' {
		mbr.Mbr_partition_2.Part_status = '0'
		mbr.Mbr_partition_2.Part_size = 0
		mbr.Mbr_partition_2.Part_start = 0
		mbr.Mbr_partition_2.Part_fit = ' '
		mbr.Mbr_partition_2.Part_type = ' '
		var espacio string = "                "
		copy(mbr.Mbr_partition_2.Part_name[:], espacio)
		del = 2
	} else if strings.EqualFold(nombreE, string(mbr.Mbr_partition_3.Part_name[:])) && mbr.Mbr_partition_3.Part_status == '1' {
		mbr.Mbr_partition_3.Part_status = '0'
		mbr.Mbr_partition_3.Part_size = 0
		mbr.Mbr_partition_3.Part_start = 0
		mbr.Mbr_partition_3.Part_fit = ' '
		mbr.Mbr_partition_3.Part_type = ' '
		var espacio string = "                "
		copy(mbr.Mbr_partition_3.Part_name[:], espacio)
		del = 3
	} else if strings.EqualFold(nombreE, string(mbr.Mbr_partition_4.Part_name[:])) && mbr.Mbr_partition_4.Part_status == '1' {
		mbr.Mbr_partition_4.Part_status = '0'
		mbr.Mbr_partition_4.Part_size = 0
		mbr.Mbr_partition_4.Part_start = 0
		mbr.Mbr_partition_4.Part_fit = ' '
		mbr.Mbr_partition_4.Part_type = ' '
		var espacio string = "                "
		copy(mbr.Mbr_partition_4.Part_name[:], espacio)
		del = 4
	}
	return del
}

//************************************ PARTICIONES LOGICAS *******************************************************
func CrearParticionLogica(size string, unit string, path string, fit string, name string) {
	mbr = LeerMBR(path)
	tipo = ""
	if size == "" || path == "" || name == "" {
		fmt.Println("Faltan campos Obligatorios")
	} else {
		nsize := GetUnit(unit, size)
		fit = GetFit(fit)
		nFit := []byte(fit)
		var n [16]byte
		espacios := "                "
		copy(n[:], espacios)
		copy(n[:], name)

		if ParticionExtendida(mbr) == 1 {
			part = mbr.Mbr_partition_1
			CrearEBR(part, path, name, nsize, nFit)
		} else if ParticionExtendida(mbr) == 2 {
			part = mbr.Mbr_partition_2
			CrearEBR(part, path, name, nsize, nFit)
		} else if ParticionExtendida(mbr) == 3 {
			part = mbr.Mbr_partition_3
			CrearEBR(part, path, name, nsize, nFit)
		} else if ParticionExtendida(mbr) == 4 {
			part = mbr.Mbr_partition_4
			CrearEBR(part, path, name, nsize, nFit)
		} else {
			fmt.Println("No existe particion extendida")
		}

	}
}
func ParticionExtendida(mbrExt Mbr) int {
	var retorno int = 0
	if strings.EqualFold(string(mbrExt.Mbr_partition_1.Part_type), "e") {
		retorno = 1
	} else if strings.EqualFold(string(mbrExt.Mbr_partition_2.Part_type), "e") {
		retorno = 2
	} else if strings.EqualFold(string(mbrExt.Mbr_partition_3.Part_type), "e") {
		retorno = 3
	} else if strings.EqualFold(string(mbrExt.Mbr_partition_4.Part_type), "e") {
		retorno = 4
	}
	return retorno
}

var intLogica int64

func CrearEBR(particion Partition, path string, nombre string, size int64, fit []byte) {
	ebr = LeerEBR(path, particion.Part_start)
	var ebrNuevo Ebr
	var arr [16]byte
	copy(arr[:], nombre)
	if ebr.Part_status == '0' {
		ebr.Part_fit = fit[0]
		ebr.Part_size = size
		ebr.Part_start = particion.Part_start
		ebr.Part_status = '1'
		ebr.Part_name = arr
		EscribirEbr(path, particion.Part_start, ebr)
	} else {
		ebr = LeerEBR(path, particion.Part_start)
		AgregarMasEbr(ebr, ebrNuevo, path, arr, fit[0], size)
	}
}

var tamebr int64

func AgregarMasEbr(eAnterior Ebr, ebrNuevo Ebr, path string, arr [16]byte, f byte, tam int64) {
	if eAnterior.Part_next == -1 && eAnterior.Part_status == '1' {
		ebrNuevo.Part_next = -1
		ebrNuevo.Part_fit = f
		ebrNuevo.Part_name = arr
		ebrNuevo.Part_status = '1'
		ebrNuevo.Part_size = tam
		ebrNuevo.Part_start = eAnterior.Part_start + eAnterior.Part_size
		eAnterior.Part_next = ebrNuevo.Part_start
		EscribirEbr(path, eAnterior.Part_start, eAnterior)
		EscribirEbr(path, ebrNuevo.Part_start, ebrNuevo)
	} else {
		eAnterior = LeerEBR(path, eAnterior.Part_next)
		AgregarMasEbr(eAnterior, ebrNuevo, path, arr, f, tam)
	}
}
func TamEbr(e Ebr) {
	if e.Part_status != '0' {
		if e.Part_next == -1 {
			cadena += "<TD BGCOLOR=\"#F0E23A\">LOGICA</TD>\n"
		} else {
			tamebr += e.Part_size
			TamEbr(e)
		}
	}
}
func EscribirEbr(path string, comienzo int64, eb Ebr) {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	// estado 0 libre, 1 ocupado
	e := &eb
	archivo.Seek(comienzo, 0)
	var binario1 bytes.Buffer
	binary.Write(&binario1, binary.BigEndian, e)
	EscribirArchivo(archivo, binario1.Bytes())

	/*archivo.Seek(eb.Part_start, 0)
	var binario3 bytes.Buffer
	ss := &s
	n := strconv.FormatInt(ebr.Part_size, 10)
	n1, _ := strconv.Atoi(n)
	for i := 0; i < n1; i++ {
		binary.Write(&binario3, binary.BigEndian, ss)
	}

	EscribirArchivo(archivo, binario3.Bytes())*/
}
func LeerEBR(path string, comienzo int64) Ebr {
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		fmt.Println("Ruta incorrecta")
	}
	archivo.Seek(comienzo, 0)
	ebr := Ebr{}
	tama := int(binary.Size(ebr))
	datos := LeerByteArchivo(archivo, tama)
	buffer := bytes.NewBuffer(datos)
	err = binary.Read(buffer, binary.BigEndian, &ebr)
	if err != nil {
		panic(err)
	}
	return ebr
}

//******************************************************************************************************
