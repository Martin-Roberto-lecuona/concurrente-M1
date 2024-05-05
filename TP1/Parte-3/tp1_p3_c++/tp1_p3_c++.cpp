#include <iostream>
#include <thread>
#include <vector>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <map>
#include <vector>
#include <fstream>
#include <atomic>
#include <sstream>
#include <chrono>
#include <mutex> A
#include <unordered_map>

// Lo que arruina mucho el tiempo es el mapa debido a que hay que usar mutex. Hay bibliotecas como tbb o Folly que tienen 
// mapas con una sicronizacion interna, parecido a un atomic pero se complica mucho pasarlo de maquina a maquina por que hay que 
// instalar muchas dependencias

using namespace std;
using namespace std::chrono;
// MUTEX
mutex mtx_palabras; 
mutex mtx_caracteres;
mutex mtx_espacios;

mutex mtx_mapa; 

// Global Vars
int cant_lineas = 0; 
atomic<int> cant_palabras(0);
atomic<int> cant_caracteres(0);
atomic<int> cant_espacios(0);

int cant_lineas_mtx = 0; 
int cant_palabras_mtx = 0;
int cant_caracteres_mtx = 0;
int cant_espacios_mtx = 0;

unordered_map<string, int> word_count;
unordered_map<string, int> word_count_mtx;

// Funciones
void resolucion_atomic(string file_name);
void creacion_hilos_atomic(fstream *file);
void mostrar_atomic();
void calculo_atomic(string line);
void resolucion_mutex(string file_name);
void creacion_hilos_mutex(fstream *file);
void mostrar_mutex();
void calculo_mutex(string line);

// MAIN
int main(int argc, char const *argv[]){

    string file_name;
    
    if(argc == 2)
    {
        file_name=argv[1];
        cout <<"RTA :"<< endl;
    }
    else
    {
        cerr << "Falta nombre de archivo"<<endl;
        return -1;
    }

    try 
    {
        resolucion_atomic(file_name);
        resolucion_mutex(file_name);
        
    } catch (const exception &e) 
    {
        cerr << "Error durante el cálculo: " << e.what() << endl;
        return -1;
    }

    return 0;
}

void resolucion_atomic(string file_name)
{
    fstream file;
    file.open(file_name, ios::in);
    if (!file) 
		cerr << "Error al abrir el archivo";
    cout<<"================= C++ ========================="<<endl;
    auto inicio = high_resolution_clock::now();
    
    creacion_hilos_atomic(&file);
    
    auto fin = high_resolution_clock::now();
    auto duracion = duration_cast<milliseconds>(fin - inicio);
    
    mostrar_atomic();
    cout<<"=========== calculo_con_atomic " << duracion.count() << " ms." "================"<<endl;
    file.close();
}

void creacion_hilos_atomic(fstream *file)
{
    string  line;
    vector<thread> hilos;

    int treads_not_detached = 0;
    while(getline(*file,line))
    {
        hilos.push_back(thread(calculo_atomic, line));
        treads_not_detached++;
        if (treads_not_detached > 10000) 
        {
            for (auto &t : hilos) 
                if (t.joinable()) 
                    t.join();
            treads_not_detached = 0;
        } 
        
        cant_lineas++;
    }
    
    for (auto &t : hilos) 
        if (t.joinable()) 
            t.join();

    file->close();
}

void mostrar_atomic()
{
    cout << "Cantidad de Lineas: " << cant_lineas << endl;
    cout << "Cantidad de Caracteres: " << cant_caracteres << endl;
    cout << "Cantidad de Palabras: " << cant_palabras << endl;
    cout << "Cantidad de Espacios: " << cant_espacios << endl;
    int max_val = 0;
    string max_key;
    for (const auto &p : word_count)
    {
        if (p.second > max_val)
        {
            max_val = p.second;
            max_key = p.first;
        }
    }
    cout << "Palabra con mas repeticion '" << max_key << "', con " << max_val << " repeticiones " << endl;
}

void calculo_atomic(string line)
{
    istringstream iss(line);
    string word;
    while (iss >> word) 
    {
         if (word != " ") 
         {
            atomic_fetch_add(&cant_palabras, 1);
            mtx_mapa.lock();
            word_count[word]++; 
            mtx_mapa.unlock();
         }
    }
    for (char c : line) 
        if (c == ' ') 
            atomic_fetch_add(&cant_espacios, 1);
        else
            atomic_fetch_add(&cant_caracteres, 1);
}

void resolucion_mutex(string file_name)
{
    fstream file;
    file.open(file_name, ios::in);
    if (!file) 
		cerr << "Error al abrir el archivo";
    cout<<"================= C++ ========================="<<endl;
    auto inicio = high_resolution_clock::now();
    
    creacion_hilos_mutex(&file);
    
    auto fin = high_resolution_clock::now();
    auto duracion = duration_cast<milliseconds>(fin - inicio);

    mostrar_mutex();
    cout<<"=========== calculo_con_mutex " << duracion.count() << " ms." "================"<<endl;
    file.close();
}

void creacion_hilos_mutex(fstream *file)
{
    string  line;
    vector<thread> hilos;
    int treads_not_detached = 0;
    while(getline(*file,line))
    {
        hilos.push_back(thread(calculo_mutex, line));
        treads_not_detached++;
        if (treads_not_detached > 10000) 
        {
            for (auto &t : hilos) 
                if (t.joinable()) 
                    t.join();
            treads_not_detached = 0;
        } 
        cant_lineas_mtx++;
    }
    
    for (auto &t : hilos) 
        if (t.joinable()) 
            t.join();

    file->close();
}

void mostrar_mutex()
{
    cout << "Cantidad de Lineas: " << cant_lineas_mtx << endl;
    cout << "Cantidad de Caracteres: " << cant_caracteres_mtx << endl;
    cout << "Cantidad de Palabras: " << cant_palabras_mtx << endl;
    cout << "Cantidad de Espacios: " << cant_espacios_mtx << endl;
    int max_val = 0;
    string max_key;
    for (const auto &p : word_count_mtx)
    {
        if (p.second > max_val)
        {
            max_val = p.second;
            max_key = p.first;
        }
    }
    cout << "Palabra con mas repeticion '" << max_key << "', con " << max_val << " repeticiones " << endl;
}

void calculo_mutex(string line) {
    istringstream iss(line);
    string word;

    while (iss >> word) {
        if (word != " ") {
            mtx_palabras.lock();
            cant_palabras_mtx++;
            mtx_palabras.unlock();
            mtx_mapa.lock();
            word_count_mtx[word]++;
            mtx_mapa.unlock();
        }
    }

    for (char c : line) {
        if (c == ' ') {
            mtx_espacios.lock(); 
            cant_espacios_mtx++;
            mtx_espacios.unlock();  
        } else {
            mtx_caracteres.lock();
            cant_caracteres_mtx++;
            mtx_caracteres.unlock();
        }
    }
}

