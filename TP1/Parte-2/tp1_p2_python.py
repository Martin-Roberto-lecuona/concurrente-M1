import numpy as np
import threading
import time

LIMITE_SUPERIOR = 9
LIMITE_INFERIOR = 0
TAM_MATRIZ = 9000
CANT_HILOS = 100
ESCALAR = 7
TIME_MILISECONDS = 1000

def create_matriz():
    return np.random.randint(LIMITE_INFERIOR,LIMITE_SUPERIOR, size=(TAM_MATRIZ, TAM_MATRIZ))

def compare_equals_matrices(matriz, matrizB):
    return np.array_equal(matriz,matrizB)

def inicializar_matriz():
    result_matriz = []
    for i in range(TAM_MATRIZ):
        row = []
        for j in range(TAM_MATRIZ):
            row.append(0)
        result_matriz.append(row)
    return result_matriz

def product_matrizXscalar_secuencial():
    for i in range(TAM_MATRIZ):
        for j in range(TAM_MATRIZ):
            matrizCS[i][j] = matriz_A[i][j] * ESCALAR

    return matrizCS

def product_matrizXscalar(start, end):
    for i in range(start,end):
        for j in range(TAM_MATRIZ):
            matrizCC[i][j] = matriz_A[i][j] * ESCALAR

def product_matrizXscalar_concurrente():
    threads = []
    blockSize = TAM_MATRIZ // CANT_HILOS
    for i in range(CANT_HILOS):
        thread = threading.Thread(target=product_matrizXscalar, args=(i*blockSize,(i+1)*blockSize,))
        thread.start()
        threads.append(thread)

    for j in threads:
        j.join()

matrizCC = inicializar_matriz()
matrizCS = inicializar_matriz()
matriz_A = create_matriz()

print(f"La matriz cuadrada es {TAM_MATRIZ} x {TAM_MATRIZ}")
print(f"Hilos ejecutandose concurrentemente: {CANT_HILOS}")

calculate_time_concurrente = time.time()
product_matrizXscalar_concurrente()
elapsed_time_2_ms = (time.time() - calculate_time_concurrente) * TIME_MILISECONDS

calculate_time_secuencial = time.time()
product_matrizXscalar_secuencial()
elapsed_time_1_ms = (time.time() - calculate_time_secuencial) * TIME_MILISECONDS

print(f"La aplicacion SECUENCIAL tardó {elapsed_time_1_ms} milisegundos en ejecutarse.")
print(f"La aplicacion CONCURRENTE tardó {elapsed_time_2_ms} milisegundos en ejecutarse.")

equals = compare_equals_matrices(matrizCC,matrizCS)

if equals:
    print("Las matrices son iguales")
else:
    print("Las matrices no son iguales")