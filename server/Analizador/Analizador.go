package Analizador

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"main.go/Comandos"
)

var re = regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

func getCommandAndParams(input string) (string, string) {
	parts := strings.Fields(input)
	if len(parts) > 0 {
		command := strings.ToLower(parts[0])
		params := strings.Join(parts[1:], " ")
		return command, params
	}
	return "", input
}

func Analizar(input string) {

	command, params := getCommandAndParams(input)

	//fmt.Println("Command: ", command, "Params: ", params)

	ReadComands(command, params)

	//mkdisk -size=3000 -unit=K -fit=BF
	//fdisk -size=300 -driveletter=A -name=Particion1
	//mount -driveletter=A -name=Particion1
	//mkfs -type=full -id=A119
	//login -user=root -pass=123 -id=A119

}

func ReadComands(command string, params string) {
	//fmt.Println(os.Args)

	switch { // uso de comandos
	case strings.EqualFold(command, "mkdisk"): // se usará
		Find_Mkdisk(params)
		break
	case strings.EqualFold(command, "rmdisk"):

		break
	case strings.EqualFold(command, "fdisk"): // Se usará
		Find_Fdisk(params)
		break
	case strings.EqualFold(command, "Mount"): // Se usará
		Find_Mount(params)
		break
	case strings.EqualFold(command, "Mkfs"):

		break
	case strings.EqualFold(command, "Execute"): // Se usará

		break
	case strings.EqualFold(command, "Rep"): // Se usará
		find_Rep(params)

		break
	case strings.EqualFold(command, "Login"):

		break
	}

}

func Find_Mkdisk(params string) {
	// Define flags
	fs := flag.NewFlagSet("mkdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	fit := fs.String("fit", "FF", "Ajuste")
	unit := fs.String("unit", "M", "Unidad")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	Comandos.MKDisk(size, fit, unit)

}

func Find_Fdisk(params string) {
	// Define flags
	fs := flag.NewFlagSet("fdisk", flag.ExitOnError)
	size := fs.Int("size", 0, "Tamaño")
	driveletter := fs.String("driveletter", "", "Letra")
	name := fs.String("name", "", "Nombre")
	unit := fs.String("unit", "m", "Unidad")
	types := fs.String("type", "p", "Tipo")
	fit := fs.String("fit", "f", "Ajuste")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size", "fit", "unit", "driveletter", "name", "type":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	// Call the function
	Comandos.FDisk(size, driveletter, name, unit, types, fit)
}

func Find_Mount(input string) {
	fs := flag.NewFlagSet("mount", flag.ExitOnError)

	driveletter := fs.String("driveletter", "", "Letra")
	name := fs.String("name", "", "Nombre")

	// Parse the flags
	fs.Parse(os.Args[1:])

	// find the flags in the input
	matches := re.FindAllStringSubmatch(input, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "driveletter", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}
	Comandos.Mount(driveletter, name)

}

func find_Rep(params string) {
	fs := flag.NewFlagSet("rep", flag.ExitOnError)
	name := fs.String("name", "0", "nombre del reporte")
	path := fs.String("path", "", "ubicación a generar el reporte")
	id := fs.String("id", "", "id de la partición")

	fs.Parse(os.Args[1:])

	matches := re.FindAllStringSubmatch(params, -1)

	// Process the input
	for _, match := range matches {
		flagName := match[1]
		flagValue := strings.ToLower(match[2])

		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "id", "path", "name":
			fs.Set(flagName, flagValue)
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	//fmt.Print("--->" + string(*id) + "<---")
	Comandos.Rep(name, path, id)
}
