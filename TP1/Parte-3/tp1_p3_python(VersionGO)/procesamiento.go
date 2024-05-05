package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var chP1ok chan bool = make(chan bool)

func ProcesarJson(rutaArchivo string, nombreArchivo string) {
	archivoJson, err := os.Open(rutaArchivo + nombreArchivo)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer archivoJson.Close()

	decoder := json.NewDecoder(archivoJson)

	_, err = decoder.Token()
	if err != nil {
		fmt.Println("Error al leer el inicio del array JSON:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	chSuc1 := make(chan Producto, 5)
	chSuc2 := make(chan Producto, 5)

	go ProcesarSucursal(chSuc1, 1, &wg)
	go ProcesarSucursal(chSuc2, 2, &wg)

	i := 0
	for decoder.More() {
		var producto Producto
		err := decoder.Decode(&producto)
		if err != nil {
			fmt.Println("Error al decodificar JSON:", err)
		}
		producto.pos = i
		i++
		switch producto.Sucursal {
		case 1:
			chSuc1 <- producto
		case 2:
			chSuc2 <- producto
		default:
			fmt.Println("Error: sucursal de producto no reconocida")
		}
	}

	close(chSuc1)
	close(chSuc2)

	_, err = decoder.Token()
	if err != nil {
		fmt.Println("Error al leer el fin del array JSON:", err)
		return
	}

	wg.Wait()
}

func ProcesarSucursal(chSuc chan Producto, suc int, wg *sync.WaitGroup) {
	defer wg.Done()

	montoTotal := 0
	maxValue := 0
	maxKey := ""
	ventasPorProd := make(map[string]int)
	ventasPorDia := make(map[string]CumuloPorFecha)

	for producto := range chSuc {
		montoTotal += producto.Monto
		ventasPorProd[producto.Producto]++

		fechaStr := producto.Fecha.String()
		ventasPorDia[fechaStr] = CumuloPorFecha{
			NroVentas:  ventasPorDia[fechaStr].NroVentas + 1,
			MontoTotal: ventasPorDia[fechaStr].MontoTotal + producto.Monto,
		}

	}

	for key, value := range ventasPorProd {
		if value > maxValue {
			maxValue = value
			maxKey = key
		}
	}

	switch suc {
	case 1:
		MostrarResultados(suc, montoTotal, maxKey, maxValue, ventasPorDia)
		chP1ok <- true
	case 2:
		if <-chP1ok {
			MostrarResultados(suc, montoTotal, maxKey, maxValue, ventasPorDia)
		}
	}
}
