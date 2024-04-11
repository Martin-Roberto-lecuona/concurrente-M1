#include <iostream>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <sys/types.h>
#include <map>
#include <vector>

#define SLEEP_TIME 20

using namespace std;

struct ProcesoInfo
{
    char proc;
    pid_t pid;
    pid_t ppid;
};

void imprimir_info(char proc)
{
    ProcesoInfo p = {proc, getpid(), getppid()};
    cout << "Proceso: " << p.proc << ", PID: " << p.pid << ", PPID: " << p.ppid << endl;
}

void imprimir_error(char proc)
{
    cerr << "Error al crear proceso " << proc << endl;
    exit(EXIT_FAILURE);
}

void crear_arbol(const std::pair<const char, std::vector<char>> &p, std::map<char, std::vector<char>> &procesoRelaciones)
{
    char padre = p.first;
    imprimir_info(padre);

    for (auto const &hijo : p.second)
    {
        if (procesoRelaciones.find(hijo) != procesoRelaciones.end())
        {
            // hijo que es padre
            pid_t pid_padre = fork();
            if (pid_padre < 0)
                imprimir_error(padre);
            else if (pid_padre == 0)
            {
                crear_arbol(*procesoRelaciones.find(hijo), procesoRelaciones);
                wait(NULL);
                sleep(SLEEP_TIME);
                exit(EXIT_SUCCESS);
            }
        } else
        {
            // hijo que es hoja
            pid_t pid_hoja = fork();
            if (pid_hoja < 0)
                imprimir_error(hijo);
            else if (pid_hoja == 0)
            {
                imprimir_info(hijo);
                sleep(SLEEP_TIME);
                exit(EXIT_SUCCESS);
            }
        }
    }
}

int main(int argc, char *argv[])
{
    if (argc == 2 && (strcmp(argv[1], "-h") == 0 || strcmp(argv[1], "--help") == 0))
    {
        cout << "Este proceso crea un árbol de procesos" << endl;
        return EXIT_SUCCESS;
    }

    char raiz = 'A';
    map<char, vector<char>> procesoRelaciones =
    {
        {'A', {'B', 'C', 'D'}},
        {'B', {'E', 'F'}},
        {'F', {'H', 'I'}},
        {'I', {'J'}},
        {'D', {'G'}}
    };

    crear_arbol(*procesoRelaciones.find(raiz), procesoRelaciones);

    while (wait(NULL) > 0);

    return EXIT_SUCCESS;
}

