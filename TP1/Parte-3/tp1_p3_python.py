from collections import Counter
from multiprocessing import Pipe
import os
import sys

SUCURSAL_B = 1
SUCURSAL_C = 2
CHILD = 0
CANT_EMPATE = 3


def escribir_registros(archivo):
    datos = """SUCURSAL,FECHA,PRODUCTO,MONTO
1,01-05-2024,Producto A,100.00
1,01-05-2024,Producto B,150.00
1,01-05-2024,Producto D,60.00
1,02-05-2024,Producto C,80.00
2,02-05-2024,Producto D,130.00
1,03-05-2024,Producto D,75.00
2,04-05-2024,Producto A,85.00
1,05-05-2024,Producto C,160.00
2,06-05-2024,Producto C,95.00
1,07-05-2024,Producto D,180.00
1,08-05-2024,Producto B,120.00
2,09-05-2024,Producto C,110.00
1,10-05-2024,Producto B,200.00
1,10-05-2024,Producto E,150.00
1,11-05-2024,Producto A,140.00
2,12-05-2024,Producto B,90.00
1,13-05-2024,Producto C,170.00"""

    with open(archivo, 'w') as f:
        f.write(datos)


def leer_registros(archivo):
    registros = []
    with open(archivo, 'r') as file:
        next(file)  # Ignorar la primera línea (encabezados)
        for linea in file:
            sucursal, fecha, producto, monto = linea.strip().split(',')
            registros.append({
                'SUCURSAL': int(sucursal),
                'FECHA': fecha,
                'PRODUCTO': producto,
                'MONTO': float(monto)
            })
    return registros


def procesar_datos(registros):
    sucursal = registros[0]['SUCURSAL']
    total_monto = 0
    productos_cant_ventas = Counter()
    fechas_cant_ventas = Counter()
    montos_por_fecha = Counter()

    for registro in registros:
        monto = registro['MONTO']
        producto = registro['PRODUCTO']
        fecha = registro['FECHA']

        total_monto += monto
        productos_cant_ventas[producto] += 1
        fechas_cant_ventas[fecha] += 1
        montos_por_fecha[fecha] += monto

    producto_mas_vendido = max(productos_cant_ventas,
                               key=productos_cant_ventas.get)

    max_cant_ventas = max(fechas_cant_ventas.values())
    fechas_max_cant_ventas = [
        fecha for fecha, cantidad in fechas_cant_ventas.items()
        if cantidad == max_cant_ventas
    ]
    # Limitar las fechas a las primeras 3 en caso de empate
    fechas_max_cant_ventas = fechas_max_cant_ventas[:CANT_EMPATE]

    max_monto_fecha = max(montos_por_fecha.values())
    fechas_max_monto = [
        fecha for fecha, monto in montos_por_fecha.items()
        if monto == max_monto_fecha
    ]
    fechas_max_monto = fechas_max_monto[:CANT_EMPATE]

    resultados = {
        'Número de sucursal': sucursal,
        'Monto total': total_monto,
        'Producto más vendido': producto_mas_vendido,
        'Fecha/s con mayor cantidad de ventas': fechas_max_cant_ventas,
        'Mayor cantidad de ventas en una fecha': max_cant_ventas,
        'Fecha/s con mayor monto vendido': fechas_max_monto,
        'Mayor monto vendido en una fecha': max_monto_fecha
    }
    return resultados


def imprimir_resultados(resultados):
    for clave, valor in resultados.items():
        print(f"{clave}: {valor}")
    print("\n")


if __name__ == '__main__':
    archivo = 'registros_ventas.txt'
    pipe_ab, pipe_ba = Pipe()  # Para la comunicación entre A y B
    pipe_ac, pipe_ca = Pipe()  # Para la comunicación entre A y C

    escribir_registros(archivo)
    registros = leer_registros(archivo)
    registrosB = []
    registrosC = []

    for registro in registros:
        if registro['SUCURSAL'] == SUCURSAL_B:
            registrosB.append(registro)
        elif registro['SUCURSAL'] == SUCURSAL_C:
            registrosC.append(registro)

    pid_b = os.fork()
    if pid_b < 0:
        sys.exit("Error al crear el proceso hijo B.")

    if pid_b == CHILD:  # Proceso hijo B
        pipe_ab.close()
        registrosB = pipe_ba.recv()
        resultadosB = procesar_datos(registrosB)
        pipe_ba.send(resultadosB)
        pipe_ba.close()
        os._exit(0)
    else:  # Proceso padre A
        pipe_ba.close()
        pipe_ab.send(registrosB)
        resultadosB = pipe_ab.recv()
        pipe_ab.close()
        os.wait()

    pid_c = os.fork()
    if pid_c < 0:
        sys.exit("Error al crear el proceso hijo C.")

    if pid_c == CHILD:  # Proceso hijo C
        pipe_ac.close()
        registrosC = pipe_ca.recv()
        resultadosC = procesar_datos(registrosC)
        pipe_ca.send(resultadosC)
        pipe_ca.close()
        os._exit(0)
    else:  # Proceso padre A
        pipe_ca.close()
        pipe_ac.send(registrosC)
        resultadosC = pipe_ac.recv()
        pipe_ac.close()
        os.wait()

imprimir_resultados(resultadosB)
imprimir_resultados(resultadosC)
