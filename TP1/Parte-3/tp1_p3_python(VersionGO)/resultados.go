package main

import (
	"fmt"
	"sort"
)

type CumuloPorFecha struct {
	NroVentas  int
	MontoTotal int
}

func MostrarResultados(suc int, montoTotal int, maxKey string, maxValue int, ventasPorDia map[string]CumuloPorFecha) {
	fmt.Printf("Sucursal: %v\n", suc)
	fmt.Printf("Monto total: $%v\n", montoTotal)
	fmt.Printf("Producto más vendido: %v con %v ventas\n", maxKey, maxValue)

	fechasPorNroVentas := ordenarFechas(ventasPorDia, func(i, j string) bool {
		return ventasPorDia[i].NroVentas > ventasPorDia[j].NroVentas
	})
	fechasPorMontoTotal := ordenarFechas(ventasPorDia, func(i, j string) bool {
		return ventasPorDia[i].MontoTotal > ventasPorDia[j].MontoTotal
	})

	MostrarMaximoVentasPorDia(ventasPorDia, fechasPorNroVentas)

	MostrarMaximoMontoPorDia(ventasPorDia, fechasPorMontoTotal)
}

func ordenarFechas(ventasPorDia map[string]CumuloPorFecha, less func(i string, j string) bool) []string {
	var fechas []string
	for fecha := range ventasPorDia {
		fechas = append(fechas, fecha)
	}

	sort.Slice(fechas, func(i, j int) bool {
		return less(fechas[i], fechas[j])
	})
	return fechas
}

func MostrarMaximoMontoPorDia(ventasPorDia map[string]CumuloPorFecha, fechasPorMontoTotal []string) {
	fmt.Println("Fecha con mayor monto vendido:")
	maxMontoTotal := ventasPorDia[(fechasPorMontoTotal)[0]].MontoTotal
	for i, fecha := range fechasPorMontoTotal {
		if i >= 3 || ventasPorDia[fecha].MontoTotal != maxMontoTotal {
			break
		}
		fmt.Printf("Fecha: %v, Monto vendido: $%v\n", fecha, ventasPorDia[fecha].MontoTotal)
	}
}

func MostrarMaximoVentasPorDia(ventasPorDia map[string]CumuloPorFecha, fechasPorNroVentas []string) {
	fmt.Println("Fecha con mayor cantidad de ventas:")
	maxNroVentas := ventasPorDia[(fechasPorNroVentas)[0]].NroVentas
	for i, fecha := range fechasPorNroVentas {
		if i >= 3 || ventasPorDia[fecha].NroVentas != maxNroVentas {
			break
		}
		fmt.Printf("Fecha: %v, Cantidad de ventas: %v\n", fecha, ventasPorDia[fecha].NroVentas)
	}
}
