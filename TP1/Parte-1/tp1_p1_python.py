import os
import time

TIME_SLEEP=20000


def imprimirProceso(proceso):
    pid = os.getpid()
    ppid = os.getppid()
    print(f"Proceso {proceso}: PID = {pid}, PID del Padre = {ppid}")

def crearProcesos(procesos, padre):
    hijos = procesos.get(padre, [])
    imprimirProceso(padre)
    for hijo in hijos:
        pid = os.fork()
        if pid < 0:
            print("Error pid del proceso")
        if pid == 0:  # es hijo
            crearProcesos(procesos, hijo)
            time.sleep(TIME_SLEEP)
            os._exit(0)
    for hijo in hijos:
        os.wait()


print("Comienzo a crear los procesos...\n")

procesos = {
    "A": ["B", "C", "D"],
    "B": ["E", "F"],
    "D": ["G"],
    "F": ["H", "I"],
    "I": ["J"]
}

crearProcesos(procesos, list(procesos.keys())[0])

print("\nCreación de procesos finalizada")