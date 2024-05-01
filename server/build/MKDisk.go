package build

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

	"main.go/global"
)

func IndiceAlfabeto(letra string) int {
	for i := 0; i < 26; i++ {
		if Alfabeto[i] == byte(letra[0]) {
			return i
		}
	}
	return 0

}

func CrearArchivo(TamanhoArchivo int) *os.File {
	for { // cliclo sirve para rotular con letras los discos
		// Verificar si el archivo ya existe
		_, err := os.Stat(NombreArchivo)
		if !os.IsNotExist(err) {
			Letra += 1
			NombreArchivo = "MIA/P1/" + string(Alfabeto[Letra]) + ".dsk"
			continue
		}

		Arch, err := os.Create(NombreArchivo) // valor ya existente en el main
		if err != nil {
			fmt.Println(err)
			global.Salida = "No se pudo crear"
			return Arch
		}

		/* Llenar el archivo con la cantidad de 0s especificada
		for i := 0; i < TamanhoArchivo; i++ {
			err := Escribir(Arch, byte(0), int64(i))
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}*/

		arreglo := make([]byte, 1024)
		// create array of byte(0)
		for i := 0; i <= TamanhoArchivo/1024; i++ {
			err := Escribir(Arch, arreglo, int64(i*1024))
			if err != nil {
				fmt.Println("Error: ", err)
				global.Salida = "Error"
			}
		}

		//ceros := make([]byte, TamanhoArchivo)
		//arch.Write(ceros)

		fmt.Printf("Archivo creado con %d bytes\n", int(TamanhoArchivo))
		global.Salida += "Archivo creado con " + strconv.Itoa(TamanhoArchivo) + " bytes"
		return Arch

	}
}

/*
func Escribir(nuevo interface{}) { // Interface permite qu ecepte cualquier objeto

	// abrir archivo
	file, err := os.OpenFile(NombreArchivo, os.O_RDWR, 0644)

	if err != nil {
		fmt.Println(err)
		return
	}

	file.Seek(0, 0)

	binary.Write(file, binary.LittleEndian, &nuevo)

	defer file.Close()

	fmt.Println("Se ha escrito en el archivo")

}
*/
func Escribir(file *os.File, data interface{}, position int64) error {

	//file.Seek(0, io.SeekCurrent) // Obtener la posición inicial
	//imprimirPos(*file)
	file.Seek(position, 0)

	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err objeto", err)
		global.Salida = "Err objeto"
		return err
	}

	//imprimirPos(*file)
	//file.Seek(0, io.SeekCurrent) // Obtener la posición final

	return nil
}
