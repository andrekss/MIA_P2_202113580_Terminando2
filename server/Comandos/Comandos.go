package Comandos

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"main.go/Reportes"
	"main.go/build"

	"main.go/global"
	"main.go/structs"
)

func MKDisk(size *int, fit *string, unit *string) { // mkdisk -size=30 -fit=FF -unit=K
	var aviso bool = true

	Ajuste := [3]string{"BF", "FF", "WF"}
	// Verificar si se proporciona el tamaño
	if *size <= 0 {
		fmt.Println("Error: Debes proporcionar un tamaño válido para el disco.")
		global.Salida = "Error: Debes proporcionar un tamaño válido para el disco."
		flag.PrintDefaults()
		//os.Exit(1)
	}

	for _, valor := range Ajuste {
		if strings.EqualFold(valor, *fit) {
			aviso = false
			break
		}
	}
	if aviso {
		fmt.Println("Error no existe este ajuste")
		global.Salida = "Error no existe este ajuste"
		flag.PrintDefaults()
		//os.Exit(1)
	}

	*size = build.Conversion(*unit, *size)

	fmt.Printf("Creando disco con tamaño: %d bytes\n", *size)
	global.Salida += "Creando disco con tamaño: " + strconv.Itoa(*size) + "\n"
	// Aquí puedes agregar lógica adicional para crear el disco con el tamaño y la unidad proporcionados.

	Arch := build.CrearArchivo(*size)

	var fitSlice []byte
	fitSlice = append(fitSlice, (*fit)[0])

	var tecnica structs.MBR
	//tecnica.Fecha = time.Now()
	//tecnica.fit = (*fit)[0]

	tecnica.Tamaño = int32(*size)
	tecnica.Signature = int32(build.Letra)

	copy(tecnica.Fit[:], fitSlice)
	copy(tecnica.Fecha[:], time.Now().String())

	// Write object in bin file
	if err := build.Escribir(Arch, tecnica, 0); err != nil {
		return
	}
	defer Arch.Close()
}

func RMDisk() { // go run *.go rmdisk -driveletter=A

	driveletter := flag.String("driveletter", "", "Borrar disco")

	flag.CommandLine.Parse(os.Args[2:])
	var nombre string = "MIA/P1/" + *driveletter + ".dsk"
	var Verificación string

	fmt.Print("¿Quiere confirmar esta accion? (tabule V para confirmar, si no tabule cualquier otro): ")
	fmt.Scan(&Verificación)

	if strings.EqualFold(Verificación, "V") {
		build.Existencia((*driveletter)[0])
		build.Eliminar(nombre)
	} else {
		fmt.Print("No se eliminó ningun archivo")
	}
}

func FDisk(size *int, driveletter *string, name *string, unit *string, types *string, fit *string) { //go run *.go fdisk -size=300 -driveletter=A -name=Particion1 -unit=B

	ruta := "./MIA/P1/" + strings.ToUpper(*driveletter) + ".dsk"
	// Verificar si se proporciona el tamaño
	if *size <= 0 {
		fmt.Println("Error: Debes proporcionar un tamaño válido para el disco.")
		global.Salida = "Error: Debes proporcionar un tamaño válido para el disco."
		flag.PrintDefaults()
		//os.Exit(1)
	}
	build.Existencia((*driveletter)[0]) // verificación

	*size = build.Conversion(*unit, *size)

	// partición nueva
	var Partición structs.Partition
	Partición.Size = int32(*size)

	var name1 [16]byte
	copy(name1[:], *name)
	copy(Partición.Name[:], name1[:])

	tip := []byte(*types)
	copy(Partición.Tipo[:], tip)

	fits := []byte(*fit)
	copy(Partición.Fit[:], fits)
	copy(Partición.Status[:], "0") // verifica que esta montada
	// fin del llenado de la nueva partición

	// aqui ejecutamos todo el comando
	build.Funcionalidades(*driveletter, name1, Partición, ruta, *unit)

	fmt.Printf("Creando una partición con tamaño: %d bytes\n", *size)
	global.Salida = "Creando una partición con tamaño: " + strconv.Itoa(*size)
}

func Mount(driveletter *string, name *string) { // go run *.go mount -driveletter=A -name=Particion1

	// Abrir archivo
	ruta := "./MIA/P1/" + strings.ToUpper(*driveletter) + ".dsk"
	file, err := build.AbrirArchivo(ruta)
	if err != nil {
		return
	}

	var MBR structs.MBR

	if err := build.LeerArchivo(file, &MBR, 0); err != nil {
		return
	}

	var index int = -1
	var count = 0
	// buscamos la partición especifica
	for i := 0; i < 4; i++ {
		if MBR.Partitions[i].Size != 0 {
			count++
			if strings.Contains(string(MBR.Partitions[i].Name[:]), *name) {
				index = i
				break
			}
		}
	}

	// id = DriveLetter + Correlative + 80

	id := strings.ToUpper(*driveletter) + strconv.Itoa(count) + "80"

	copy(MBR.Partitions[index].Status[:], "1") // verifica que esta montada
	copy(MBR.Partitions[index].Id[:], id)

	for i := 0; i < len(global.Particiones); i++ {
		if global.Particiones[i].Nombre == string(*name) && global.Particiones[i].Disco == string(*driveletter) {
			global.Particiones[i].Mount = id
		}

	}

	if err := build.Escribir(file, MBR, 0); err != nil {
		return
	}

	/*
	      solo para imprimir
	   	var MBR2 structs.MBR
	   	if err := build.LeerArchivo(file, &MBR2, 0); err != nil {
	   		return
	   	}*/

	fmt.Print("Se a montado la partición" + *name)
	global.Salida = "Se a montado la partición " + string(*name)

	defer file.Close()
}

func Unmount() { // go run *.go Unmount -id=A180
	id := flag.String("id", "", "id de la particioón")
	flag.CommandLine.Parse(os.Args[2:])

	i := 0
	for {
		ruta := "./MIA/P1/" + string(build.Alfabeto[i]) + ".dsk"
		var revision structs.MBR
		file, err := build.AbrirArchivo(ruta)
		if err != nil {
			break
		}

		if err := build.LeerArchivo(file, &revision, 0); err != nil {
			return
		}

		for j := 0; j < 4; j++ { // recorremos las particiones
			var idd [4]byte
			copy(idd[:], *id)

			if revision.Partitions[j].Id == idd {
				copy(revision.Partitions[j].Status[:], "0") // indicamos que se desmontó
				build.Escribir(file, revision, 0)

				defer file.Close()
				fmt.Print("Se a desmontó la partición " + string(revision.Partitions[j].Name[:]))
				global.Salida = "Se a desmontó la partición " + string(revision.Partitions[j].Name[:])
				break
			}
		}
		i += 1
	}
}

func MKfs() {
	id := flag.String("id", "", "id partición montada")
	//types := flag.String("type", "Full", "tipo de formateo")
	//fs := flag.String("fs", "2fs", "Formateo a otro sistema")

	flag.CommandLine.Parse(os.Args[2:])
	i := 0
	for {
		ruta := "./MIA/P1/" + string(build.Alfabeto[i]) + ".dsk"
		var revision structs.MBR
		file, err := build.AbrirArchivo(ruta)
		if err != nil {
			break
		}

		if err := build.LeerArchivo(file, &revision, 0); err != nil {
			return
		}

		for j := 0; j < 4; j++ { // particiones
			var idd [4]byte
			copy(idd[:], *id)

			if revision.Partitions[j].Id == idd {
			}
		}
		i += 1
	}
}

func Rep(name *string, path *string, id *string) { //go run *.go rep -name=MBR -path=./reporte -id=A080
	tabla := Reportes.CreandoEstructura(*name, *id)
	Reportes.Reportes(tabla, *path, *id)
	//*ruta = *ruta
}

func Login() {
	usuario := "Andres"
	contraseña := "HOLAmundo"

	user := flag.String("user", "", "Usuario")
	pass := flag.String("pass", "", "contraseña")
	//	id := flag.String("id", "", "id")

	flag.CommandLine.Parse(os.Args[2:])

	if *user == usuario && *pass == contraseña {
	} else {
		fmt.Print("Usuario desconocido")
	}

}
