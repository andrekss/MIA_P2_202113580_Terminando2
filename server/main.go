package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"main.go/Analizador"
	"main.go/global"
)

func listarArchivos(ruta string) ([]string, error) {
	// Lee el contenido de la carpeta especificada
	archivos, err := ioutil.ReadDir(ruta)
	if err != nil {
		return nil, err
	}

	// Arreglo para almacenar los nombres de los archivos
	nombres := make([]string, 0, len(archivos))

	// Recorre los archivos y agrega sus nombres al arreglo
	for _, archivo := range archivos {
		nombres = append(nombres, archivo.Name())
	}

	return nombres, nil
}

func main() {
	app := fiber.New()

	app.Use(cors.New())

	// Definir una ruta de ejemplo
	app.Post("/Execute", func(Req *fiber.Ctx) error {

		Analizador.Analizar(string(Req.Body()))
		resp := Req.SendString(global.Salida)
		global.Salida = ""
		return resp
	})

	app.Get("/Discos", func(c *fiber.Ctx) error {
		// Llama a la función listarArchivos para obtener los nombres de los archivos
		nombres, err := listarArchivos("./MIA/P1")
		if err != nil {
			log.Fatalf("Error al listar archivos: %v", err)
		}

		// Convierte el arreglo de strings a JSON
		respuesta, err := json.Marshal(nombres)
		if err != nil {
			log.Fatalf("Error al convertir a JSON: %v", err)
		}

		// Retorna el JSON como respuesta
		return c.Send(respuesta)
	})

	app.Get("/Particion", func(c *fiber.Ctx) error {

		// Convierte el arreglo de strings a JSON
		respuesta, err := json.Marshal(global.Particiones)
		if err != nil {
			log.Fatalf("Error al convertir a JSON: %v", err)
		}

		return c.Send(respuesta)

	})

	port := ":5000"

	// Iniciar el servidor
	log.Fatal(app.Listen(port))
}
