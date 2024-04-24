package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const SIZE int = 9000
const THREADS int = 100
const ESCALAR int = 7

var matrixA [SIZE][SIZE]int
var matrixResSequ [SIZE][SIZE]int
var matrixResConc [SIZE][SIZE]int

func printMatrix(matrix *[SIZE][SIZE]int) {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			fmt.Print("|")
			if matrix[i][j] < 10 {
				fmt.Print(" ")
			}
			if matrix[i][j] < 100 {
				fmt.Print(" ")
			}
			fmt.Print(matrix[i][j])
			fmt.Print("|")
		}
		fmt.Println()
	}
}

func initMatrix(m *[SIZE][SIZE]int) {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			m[i][j] = rand.Intn(100)
		}
	}
}

func multSeqMatrix() {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			matrixResSequ[i][j] = matrixA[i][j] * ESCALAR
		}
	}
}
func multConcMatrix() {
	var wg sync.WaitGroup
	wg.Add(THREADS)

	blockSize := SIZE / THREADS

	for i := 0; i < THREADS; i++ {
		go multRows(i*blockSize, (i+1)*blockSize, &wg)
	}

	wg.Wait()
}

func multRows(startRow int, endRow int, wg *sync.WaitGroup) {
	defer wg.Done()
	for r := startRow; r < endRow; r++ {
		for c := 0; c < SIZE; c++ {
			matrixResConc[r][c] = matrixA[r][c] * ESCALAR
		}
	}
}

func compareRes() bool {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			if matrixResConc[i][j] != matrixResSequ[i][j] {
				fmt.Println(matrixResConc[i][j], "!=", matrixResSequ[i][j], "i", i, "j", j)
				return false
			}
		}
	}
	return true
}
func printAllMatrix() {
	fmt.Println("Matrix A:")
	printMatrix(&matrixA)

	fmt.Println("\nRes Seq:")
	printMatrix(&matrixResSequ)

	fmt.Println("\nRes Conc:")
	printMatrix(&matrixResConc)
}

func main() {
	fmt.Println("matrix size: ", SIZE)
	fmt.Println("Threads: ", THREADS)
	initMatrix(&matrixA)

	startSeq := time.Now()
	multSeqMatrix()
	endSeq := time.Since(startSeq).Milliseconds()

	startConc := time.Now()
	multConcMatrix()
	endConc := time.Since(startConc).Milliseconds()

	// printAllMatrix()

	fmt.Println("\nSECUENCIAL == CONCURRENTE:", compareRes())

	fmt.Println("SECUENCIAL: ", endSeq)
	fmt.Println("CONCURRENTE: ", endConc)
}
