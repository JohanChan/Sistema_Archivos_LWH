package analizador

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var size, path, name, unit string
var tipo, fit, eliminar, agregar string
var mbr Mbr
var arrayLinea []string
var listadoMount []ParticionMontada
var abecedario = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var incremento int
var id []string

const (
	Informacion = "\033[1;34m"
	Noticia     = "\033[1;36m"
	Advertencia = "\033[1;33m"
	Error       = "\033[1;31m"
	DebugColor  = "\033[0;36m"
	last        = "\033[0m"
)

type Mbr struct {
	Mbr_tamano         int64     //8 bytes
	Mbr_fecha_creacion [20]byte  // 20 bytes
	Mbr_disk_signature int64     // 8 bytes
	Mbr_partition_1    Partition // 35 bytes
	Mbr_partition_2    Partition // 35 bytes
	Mbr_partition_3    Partition // 35 bytes
	Mbr_partition_4    Partition // 35 bytes
}

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

func AsignarArray(arraInicio []string) {
	arrayLinea = arraInicio
	/*for i := 0; i < len(arrayLinea); i++ {
		fmt.Println(arrayLinea[i])
	}*/
}
func FuncionComando(arreglo []string) {
	for i := 0; i <= len(arreglo)-1; i++ {
		if strings.EqualFold(arreglo[i], "mkdisk") {
			AtributosMKDISK(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "rmdisk") {
			RMDISK(arreglo[i+2])
		} else if strings.EqualFold(arreglo[i], "fdisk") {
			AtributoFDISK(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "mount") {
			AtributosMount(arreglo[i+1], i+1)
		} else if strings.EqualFold(arreglo[i], "unmount") {
			AtributosUnmount(arreglo[i+1], i+1)
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
			mbr = LeerMBR(path)
			MontarParticion(path, name)
			path = ""
			name = ""
		}

	}
}

func MontarParticion(path string, name string) {
	var arrnombre [16]byte
	copy(arrnombre[:], name)
	nombreE := string(arrnombre[:])
	if string(mbr.Mbr_partition_1.Part_name[:]) != nombreE && string(mbr.Mbr_partition_2.Part_name[:]) != nombreE && string(mbr.Mbr_partition_3.Part_name[:]) != nombreE && string(mbr.Mbr_partition_4.Part_name[:]) != nombreE {
		fmt.Println("Particion no existe en disco, pruebe nuevamente")
	} else {
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
func AsignarLiteral(path string) string {
	longitud := len(abecedario)
	for i := 0; i < len(listadoMount); i++ {
		if path == listadoMount[i].Path {
			return listadoMount[i].Letra
		}
	}
	random := rand.Intn(longitud)
	return abecedario[random]
}
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

	}

}
func EliminarBarra(atributo string) string {
	if strings.Contains(atributo, "*") {
		fmt.Println("tiene slash")
		natributo := strings.ReplaceAll(atributo, "\\*", "")

		fmt.Println("Sin slash si contiene " + natributo)
		return natributo
	}
	return atributo
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
		mbr = LeerMBR(path)
		if eliminar != "" {
			EliminarParticion(eliminar, name, path)
			eliminar = ""
		} else {
			CrearParticion(size, unit, path, tipo, fit, name)

		}
	}

}
func EliminarParticion(eliminar string, name string, path string) {
	mbrEvaluar := &mbr
	if name != "" || path != "" {
		p := ParticionAEliminar(mbrEvaluar, name)
		fmt.Println("P ", p)
		if strings.EqualFold(eliminar, "fast") {
			if p == 1 || p == 2 || p == 3 || p == 4 {
				EscribirParticionFast(path, mbr)
				fmt.Println(Noticia + "Particion eliminada correctamente" + last)
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
		DatosMbr(path)
	} else {
		fmt.Println(Error + "Elija un valor correcto")
	}
}
func EliminarFull(particion Partition, path string) {
	fmt.Println(string(particion.Part_name[:]))
	fmt.Println(particion.Part_size)
	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	var bite int8 = 0
	n := &bite
	archivo.Seek(particion.Part_start, 0)

	var binario3 bytes.Buffer
	for i := 0; i < int(particion.Part_size); i++ {
		binary.Write(&binario3, binary.BigEndian, n)
	}
	EscribirArchivo(archivo, binario3.Bytes())

}
func ParticionAEliminar(mbr *Mbr, nombre string) int {
	var arrnombre [16]byte
	copy(arrnombre[:], nombre)
	nombreE := string(arrnombre[:])
	var del int
	if strings.EqualFold(nombreE, string(mbr.Mbr_partition_1.Part_name[:])) && mbr.Mbr_partition_1.Part_status == '1' {
		mbr.Mbr_partition_1.Part_status = '0'
		del = 1
	} else if strings.EqualFold(nombreE, string(mbr.Mbr_partition_2.Part_name[:])) && mbr.Mbr_partition_2.Part_status == '1' {
		mbr.Mbr_partition_2.Part_status = '0'
		del = 2
	} else if strings.EqualFold(nombreE, string(mbr.Mbr_partition_3.Part_name[:])) && mbr.Mbr_partition_3.Part_status == '1' {
		mbr.Mbr_partition_3.Part_status = '0'
		del = 3
	} else if strings.EqualFold(nombreE, string(mbr.Mbr_partition_4.Part_name[:])) && mbr.Mbr_partition_4.Part_status == '1' {
		mbr.Mbr_partition_4.Part_status = '0'
		del = 4
	}
	return del
}
func CrearParticion(size string, unit string, path string, tipo string, fit string, name string) {
	if size == "" || path == "" || name == "" {
		fmt.Println("Faltan campos obligatorios")
	} else {
		nsize := GetUnit(unit, size)
		fit = GetFit(fit)
		if strings.EqualFold(tipo, "") {
			tipo = "p"
		}
		nFit := []byte(fit)
		nTipo := []byte(tipo)
		if mbr.Mbr_partition_1.Part_status == '0' {
			mbr.Mbr_partition_1.Part_status = '1'
			mbr.Mbr_partition_1.Part_fit = nFit[0]
			mbr.Mbr_partition_1.Part_type = nTipo[0]
			mbr.Mbr_partition_1.Part_size = nsize
			mbr.Mbr_partition_1.Part_start = int64(binary.Size(mbr))
			copy(mbr.Mbr_partition_1.Part_name[:], name)
			EscribirParticion(path, mbr, '1')
		} else if mbr.Mbr_partition_2.Part_status == '0' {
			fmt.Println(string(mbr.Mbr_partition_1.Part_type))
			if !EsExtendida(mbr.Mbr_partition_1) {
				mbr.Mbr_partition_2.Part_status = '1'
				mbr.Mbr_partition_2.Part_fit = nFit[0]
				mbr.Mbr_partition_2.Part_type = nTipo[0]
				mbr.Mbr_partition_2.Part_size = nsize
				mbr.Mbr_partition_2.Part_start = int64(mbr.Mbr_partition_1.Part_start + int64(binary.Size(mbr.Mbr_partition_1)))
				copy(mbr.Mbr_partition_2.Part_name[:], name)
				EscribirParticion(path, mbr, '2')
			} else {
				fmt.Println(Advertencia, "Extendida repetida")
			}
		} else if mbr.Mbr_partition_3.Part_status == '0' {
			if !EsExtendida(mbr.Mbr_partition_2) {
				mbr.Mbr_partition_3.Part_status = '1'
				mbr.Mbr_partition_3.Part_fit = nFit[0]
				mbr.Mbr_partition_3.Part_type = nTipo[0]
				mbr.Mbr_partition_3.Part_size = nsize
				mbr.Mbr_partition_3.Part_start = int64(mbr.Mbr_partition_2.Part_start + int64(binary.Size(mbr.Mbr_partition_2)))
				copy(mbr.Mbr_partition_3.Part_name[:], name)
				EscribirParticion(path, mbr, '3')
			} else {
				fmt.Println(Advertencia, "Extendida repetida")
			}

		} else if mbr.Mbr_partition_4.Part_status == '0' {
			if !EsExtendida(mbr.Mbr_partition_3) {
				mbr.Mbr_partition_4.Part_status = '1'
				mbr.Mbr_partition_4.Part_fit = nFit[0]
				mbr.Mbr_partition_4.Part_type = nTipo[0]
				mbr.Mbr_partition_4.Part_size = nsize
				mbr.Mbr_partition_4.Part_start = int64(mbr.Mbr_partition_3.Part_start + int64(binary.Size(mbr.Mbr_partition_3)))
				copy(mbr.Mbr_partition_4.Part_name[:], name)
				EscribirParticion(path, mbr, '4')
			} else {
				fmt.Println(Advertencia, "Extendida repetida")
			}

		}
	}
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
	fmt.Println("Obteniendo nSize Gunit ", nsize)
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
func EsExtendida(partition Partition) bool {
	return strings.EqualFold(string(partition.Part_type), "e")
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
func EscribirParticion(path string, mbr Mbr, number byte) {
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

	if number == '1' {
		archivo.Seek(mbr.Mbr_partition_1.Part_start, 0)
		var binario1 bytes.Buffer
		for i := 0; i < int(mbr.Mbr_partition_1.Part_size); i++ {
			binary.Write(&binario1, binary.BigEndian, &number)

		}
		EscribirArchivo(archivo, binario1.Bytes())
	} else if number == '2' {
		archivo.Seek(mbr.Mbr_partition_2.Part_start, 0)
		fmt.Print("Comienzo particion 2 ")
		fmt.Println(mbr.Mbr_partition_2.Part_start)
		var binario2 bytes.Buffer
		for i := 0; i < int(mbr.Mbr_partition_2.Part_size); i++ {
			binary.Write(&binario2, binary.BigEndian, &number)

		}
		EscribirArchivo(archivo, binario2.Bytes())
	} else if number == '3' {
		archivo.Seek(mbr.Mbr_partition_3.Part_start, 0)
		fmt.Print("Comienzo particion 3 ")
		fmt.Println(mbr.Mbr_partition_3.Part_start)
		var binario3f bytes.Buffer
		for i := 0; i < int(mbr.Mbr_partition_3.Part_size); i++ {
			binary.Write(&binario3f, binary.BigEndian, &number)

		}
		EscribirArchivo(archivo, binario3f.Bytes())
	} else if number == '4' {
		archivo.Seek(mbr.Mbr_partition_4.Part_start, 0)
		fmt.Print("Comienzo particion 4 ")
		fmt.Println(mbr.Mbr_partition_4.Part_start)
		var binario4 bytes.Buffer
		for i := 0; i < int(mbr.Mbr_partition_4.Part_size); i++ {
			binary.Write(&binario4, binary.BigEndian, &number)

		}
		EscribirArchivo(archivo, binario4.Bytes())
	}

	//}
}
func DatosMbr(path string) {

	archivo, err := os.OpenFile(path, os.O_RDWR, 0755)
	defer archivo.Close()
	if err != nil {
		panic(err)
	}
	mbr := Mbr{}
	tam := int(unsafe.Sizeof(mbr))
	datos := LeerByteArchivo(archivo, tam)
	buffer := bytes.NewBuffer(datos)
	err = binary.Read(buffer, binary.BigEndian, &mbr)
	if err != nil {
		panic(err)
	}
	//fecha := string(mbr.Mbr_fecha_creacion[:])
	fmt.Println("************************************")
	fmt.Println(string(mbr.Mbr_partition_1.Part_status))
	fmt.Println(string(mbr.Mbr_partition_2.Part_status))
	fmt.Println(string(mbr.Mbr_partition_3.Part_status))
	fmt.Println(string(mbr.Mbr_partition_4.Part_status))
	fmt.Println("************************************")
}
func RMDISK(path string) {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("Archivo a eliminar no Existe!")
		return
	}
	captura := CapturarPantalla("Estas seguro que deseas eliminar? y/n")
	if strings.EqualFold(captura, "y") {
		path = strings.Replace(path, "\"", "", -1)
		os.Remove(path)
		fmt.Println("Archivo " + path + " ha sido eliminado")
	}
}
func Imprimir(arreglo []string) {
	for i := 0; i < len(arreglo); i++ {
		fmt.Println(arreglo[i])
	}
}
func CrearCarpeta(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		path = strings.Replace(path, "\"", "", -1)
		tr := exec.Command("mkdir", "-p", path)
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
	defer archivo.Close()
	if err != nil {
		panic(err)
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
	mbr := Mbr{}
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
func CrearDisco(size string, path string, name string, unit string) {
	if size == "" || path == "" || name == "" {
		fmt.Println("Atributo obligatorio")
	} else {
		if s, _ := strconv.Atoi(size); s <= 0 {
			panic("No se crea disco, por número invalido")
		} else {
			CrearCarpeta(path)
			nsize, _ := strconv.ParseInt(size, 10, 64)
			if unit == "k" {
				CrearArchivo(path+name, nsize*1024)
			} else if strings.EqualFold(unit, "m") || unit == "" {
				CrearArchivo(path+name, nsize*1024*1024)
			}
		}
	}
	//	LeerMBR(path + name)

}
func LeerMBR(path string) Mbr {
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
func CapturarPantalla(mensaje string) string {
	fmt.Println(mensaje)
	bf := bufio.NewReader(os.Stdin)
	entrada, _ := bf.ReadString('\n')
	cadena := strings.TrimRight(entrada, "\r\n")
	return cadena
}
func FormatoFecha() string {
	t := time.Now()
	return t.Format("01-02-2006")
}
