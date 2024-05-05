package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// / otra vez el hecho de tener el mapa te arruina el tiempo de atomic
var (
	cant_lineas     int64 = 0
	cant_caracteres int64 = 0
	cant_palabras   int64 = 0
	cant_espacios   int64 = 0

	cant_lineas_mtx     int64 = 0
	cant_caracteres_mtx int64 = 0
	cant_palabras_mtx   int64 = 0
	cant_espacios_mtx   int64 = 0

	mtx_caracteres sync.Mutex
	mtx_palabras   sync.Mutex
	mtx_espacios   sync.Mutex

	// Mapas para contar palabras y su frecuencia
	palabra_count     = sync.Map{}             // Para la solución atómica
	palabra_count_mtx = make(map[string]int64) // Para la solución con mutex
	mtx_map           sync.Mutex               // Mutex para la protección del mapa con mutex
)

func main() {
	file_name := ""
	if len(os.Args) == 2 {
		file_name = os.Args[1]
	} else {
		fmt.Println("Falta nombre de archivo")
		return
	}

	// Solución atómica
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	startSeq := time.Now()
	solucion_atomic(scanner)
	endSeq := time.Since(startSeq).Milliseconds()
	file.Close()

	// Solución con mutex
	file, err = os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	scanner = bufio.NewScanner(file)
	startSeq2 := time.Now()
	solucion_mutex(scanner)
	endSeq2 := time.Since(startSeq2).Milliseconds()

	printRta(endSeq, endSeq2)
}

func printRta(endSeq int64, endSeq2 int64) {
	fmt.Println("=================== GO ========================")
	fmt.Println("Cantidad de Lineas: ", cant_lineas)
	fmt.Println("Cantidad de Caracteres: ", cant_caracteres)
	fmt.Println("Cantidad de Palabras: ", cant_palabras)
	fmt.Println("Cantidad de Espacios: ", cant_espacios)

	// Mostrar el mapa de palabras para la solución atómica
	var max_val int64 = 0
	max_key := ""
	palabra_count.Range(func(key, value any) bool {
		val := atomic.LoadInt64(value.(*int64)) // Desreferenciar el valor
		if val > max_val {
			max_val = val
			max_key = key.(string) // Convertir la clave a string
		}
		return true
	})
	fmt.Println("Palabra con mas repeticion '", max_key, "', con ", max_val, " repeticiones ")
	fmt.Println("=====Tiempo con GO y atomic:", endSeq, "ms=======")

	fmt.Println("=================== GO ========================")
	fmt.Println("Cantidad de Lineas: ", cant_lineas_mtx)
	fmt.Println("Cantidad de Caracteres: ", cant_caracteres_mtx)
	fmt.Println("Cantidad de Palabras: ", cant_palabras_mtx)
	fmt.Println("Cantidad de Espacios: ", cant_espacios_mtx)

	// Mostrar el mapa de palabras para la solución con mutex
	mtx_map.Lock() // Bloquear el acceso al mapa
	max_val = 0
	max_key = ""
	for key, value := range palabra_count_mtx {
		if value > max_val {
			max_val = value
			max_key = key
		}
	}
	fmt.Println("Palabra con mas repeticion '", max_key, "', con ", max_val, " repeticiones ")
	mtx_map.Unlock() // Desbloquear el acceso al mapa
	fmt.Println("=====Tiempo con GO y mutex:", endSeq2, "ms=======")
}

func solucion_atomic(scanner *bufio.Scanner) {
	var wg sync.WaitGroup
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go calculo_atomic(line, &wg)
		atomic.AddInt64(&cant_lineas, 1) // Operación atómica
	}
	wg.Wait()
}

func calculo_atomic(line string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Calcular frecuencia de palabras
	iss := bufio.NewScanner(strings.NewReader(line))
	iss.Split(bufio.ScanWords)

	for iss.Scan() {
		word := iss.Text()
		atomic.AddInt64(&cant_palabras, 1)

		// Actualizar el mapa atómico
		if val, ok := palabra_count.Load(word); ok {
			atomic.AddInt64(val.(*int64), 1)
		} else {
			initialValue := int64(1)
			palabra_count.Store(word, &initialValue)
		}
	}

	for _, c := range line {
		if c == ' ' {
			atomic.AddInt64(&cant_espacios, 1)
		} else {
			atomic.AddInt64(&cant_caracteres, 1)
		}
	}
}

func solucion_mutex(scanner *bufio.Scanner) {
	var wg sync.WaitGroup
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go calculo_mutex(line, &wg)
		atomic.AddInt64(&cant_lineas_mtx, 1) // Operación atómica
	}
	wg.Wait()
}

func calculo_mutex(line string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Calcular frecuencia de palabras
	iss := bufio.NewScanner(strings.NewReader(line))
	iss.Split(bufio.ScanWords) // Dividir por palabras

	for iss.Scan() {
		word := iss.Text()

		mtx_map.Lock()
		cant_palabras_mtx++
		if _, exists := palabra_count_mtx[word]; exists {
			palabra_count_mtx[word]++
		} else {
			palabra_count_mtx[word] = 1
		}
		mtx_map.Unlock()
	}

	// Contar caracteres y espacios
	for _, c := range line {
		if c == ' ' {
			mtx_espacios.Lock()
			cant_espacios_mtx++
			mtx_espacios.Unlock()
		} else {
			mtx_caracteres.Lock()
			cant_caracteres_mtx++
			mtx_caracteres.Unlock()
		}
	}
}
