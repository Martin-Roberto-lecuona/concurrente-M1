#include <chrono>
#include <iostream>
#include <thread>
#include <vector>

#define TAM_MATRIZ 9000
#define CANT_HILOS 100
#define ESCALAR 7

using namespace std;
using namespace std::chrono;

int matrizMult[TAM_MATRIZ][TAM_MATRIZ] = {0};
int matrizRS[TAM_MATRIZ][TAM_MATRIZ] = {0};
int matrizRC[TAM_MATRIZ][TAM_MATRIZ] = {0};

void cargarMatrizMult()
{
  for (int i = 0; i < TAM_MATRIZ; i++)
  {
    for (int j = 0; j < TAM_MATRIZ; j++)
    {
      matrizMult[i][j] = rand() % 10;
    }
  }
}

void cargarMatrizRS()
{
  for (int i = 0; i < TAM_MATRIZ; i++)
  {
    for (int j = 0; j < TAM_MATRIZ; j++)
    {
      matrizRS[i][j] = matrizMult[i][j] * ESCALAR;
    }
  }
}

void cargarFilasRC(int inicio, int fin)
{
  for (int i = inicio; i < fin; i++)
  {
    for (int j = 0; j < TAM_MATRIZ; j++)
    {
      matrizRC[i][j] = matrizMult[i][j] * ESCALAR;
    }
  }
}

void cargarMatrizRC()
{
  vector<thread> hilos;
  int filasPorHilo = TAM_MATRIZ / CANT_HILOS;
  int cantHilosAmpliados =
      TAM_MATRIZ % CANT_HILOS; // Si existe resto, lo repartire entre hilos.
  int inicio;
  int fin = 0;
  // Puede que cantHilosAmpliados sea 0, por lo que no entraria al primer for y
  // fin no tendria valor, y en el segundo for inicio=fin, entonces fin=0 para
  // que inicio empiece desde el principio.

  for (int i = 0; i < cantHilosAmpliados; i++)
  {
    int filasFinales = filasPorHilo + 1;
    inicio = i * filasFinales;
    fin = inicio + filasFinales;

    hilos.push_back(thread(cargarFilasRC, inicio, fin));
  }

  for (int i = cantHilosAmpliados; i < CANT_HILOS; i++)
  {
    inicio = fin;
    fin = inicio + filasPorHilo;

    hilos.push_back(thread(cargarFilasRC, inicio, fin));
  }

  for (int i = 0; i < CANT_HILOS; i++)
  {
    hilos[i].join();
  }

  hilos.clear();
}

void compararMatrices()
{
  bool hayDiferencia = false;

  for (int i = 0; i < TAM_MATRIZ; i++)
  {
    for (int j = 0; j < TAM_MATRIZ; j++)
    {
      int valorRS = matrizRS[i][j];
      int valorRC = matrizRC[i][j];
      if (valorRS != valorRC)
      {
        hayDiferencia = true;
        cout << "Diferencia: RS[" << i << "][" << j << "]=" << valorRS
             << " | RS[" << i << "][" << j << "]=" << valorRC << endl;
      }
    }
  }

  if (!hayDiferencia)
  {
    cout << "Las matrices RS y RC son iguales." << endl;
  }
}

int main(int argc, char const *argv[])
{
  srand(time(NULL));

  cargarMatrizMult();

  cout << "Tamaño de matriz: " << TAM_MATRIZ << " x " << TAM_MATRIZ << endl;

  // Iniciar el reloj
  auto inicio = high_resolution_clock::now();
  // auto le dice al compilador que determine el tipo de la variable en funcion
  // del tipo de expresion con la que se inicializa. Util cuando el tipo del
  // objeto devuelto por una funcion es complejo o dificil de escribir
  // explicitamente. La variable sigue siendo estaticamente tipada, no dinamica.

  // Llamar a la función que deseo medir
  cargarMatrizRS();

  // Detener el reloj
  auto fin = high_resolution_clock::now();

  // Calcular la duración en milisegundos
  auto duracion = duration_cast<milliseconds>(fin - inicio);

  // Imprimir el tiempo transcurrido
  cout << "La función cargarMatrizRS() tardó " << duracion.count()
       << " milisegundos." << endl;

  // Iniciar el reloj
  inicio = high_resolution_clock::now();

  // Llamar a la función que deseo medir
  cargarMatrizRC();

  // Detener el reloj
  fin = high_resolution_clock::now();

  // Calcular la duración en milisegundos
  duracion = duration_cast<milliseconds>(fin - inicio);

  // Imprimir el tiempo transcurrido
  cout << "La función cargarMatrizRC() utilizó " << CANT_HILOS
       << " hilo/s y tardó " << duracion.count() << " milisegundos." << endl;

  compararMatrices();

  return 0;
}