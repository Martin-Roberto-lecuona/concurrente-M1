#include <iostream>
#include <fstream>
#include <string>
#include <cstdlib>
#include <ctime>

using namespace std;
int cant_lineas = 1000;   
int cant_palabras = 0;
int cant_caracteres = 0;
int cant_espacios = 0;
// Función para generar una palabra aleatoria de longitud dada
string generateRandomWord(int length) {
    static const char charset[] = "abcdefghijklmnopqrstuvwxyz";
    string word;

    for (int i = 0; i < length; ++i) {
        int random_index = rand() % (sizeof(charset) - 1);
        word += charset[random_index];
    }

    return word;
}

int main(int argc, char const *argv[]) {
    srand(time(nullptr));

    string file_name = "archivo.txt";
    fstream file;
    file.open(file_name, ios::out);
    if (!file) {
        cerr << "Error al crear el archivo." << endl;
        return -1;
    }
    for (int i = 0; i < cant_lineas; i++) {
        int words_per_line = rand() % 1000;
        for (int j = 0; j < words_per_line; j++) {
            int length = (rand() % 10) +1 ;
            string random_word = generateRandomWord(length);
            cant_palabras++;
            file << random_word << (length > 6 ? " " : "   ");
            cant_espacios+= length > 6 ? 1 : 3;
            cant_caracteres+=length;
        }
        file << endl; 
    }

    file.close();
    cout<<"ARCHIVO CREADO CON: "<<endl;
    cout<<"-----------------------------------------------"<<endl;
    cout<<"Cantidad de Lineas: "<<cant_lineas<<endl;
    cout<<"Cantidad de Caracteres: "<<cant_caracteres<<endl;
    cout<<"Cantidad de Palabras: "<<cant_palabras<<endl;
    cout<<"Cantidad de Espacios: "<<cant_espacios<<endl;
    cout<<"-----------------------------------------------"<<endl;

    return 0;
}
