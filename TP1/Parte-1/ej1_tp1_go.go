package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

const RAIZ = 'A'
const ARGUMENTOS = 1
const CHILDLESS = 0
const SECONDS_WAIT = 10

type ProcessPrinter struct {
	mu sync.Mutex
}

var myDictionary = map[rune][]rune{
	'A': {'B', 'C', 'D'},
	'B': {'E', 'F'},
	'C': {},
	'D': {'G'},
	'E': {},
	'F': {'H', 'I'},
	'G': {},
	'H': {},
	'I': {'J'},
	'J': {},
}

var printer = ProcessPrinter{}

func main() {
	var wg sync.WaitGroup
	var process rune

	if len(os.Args) == ARGUMENTOS {
		process = rune(RAIZ)
	} else {
		process = rune(os.Args[1][0])
	}

	printer.showProcess(string(process))
	children, ok := myDictionary[process]

	if !ok {
		fmt.Println("Clave no encontrada en el diccionario")
		return
	}

	if len(children) > CHILDLESS {
		wg.Add(len(children))
		for _, letter := range children {
		   createChildProcess(letter, &wg)
		}
	}
	time.Sleep(SECONDS_WAIT * time.Second)
	wg.Wait()
}
func createChildProcess(name rune, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command(os.Args[0], string(name))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("Error al ejecutar el proceso hijo:", err)
	}
}
func (pp *ProcessPrinter) showProcess(process string) {
	pp.mu.Lock()
	defer pp.mu.Unlock()
	fmt.Println("Proceso", process, "pid:", os.Getpid(), "ppid:", os.Getppid())
}