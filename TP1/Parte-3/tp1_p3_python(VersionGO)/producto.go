package main

import (
	"fmt"
	"strconv"
	"strings"
)

type CustomFecha struct {
	Dia  int `json:"day"`
	Mes  int `json:"month"`
	Anio int `json:"year"`
}

type Producto struct {
	Sucursal int         `json:"SUCURSAL"`
	Fecha    CustomFecha `json:"FECHA"`
	Producto string      `json:"PRODUCTO"`
	Monto    int         `json:"MONTO"`
	pos      int
}

func (f CustomFecha) String() string {
	return fmt.Sprintf("%v-%v-%v", f.Dia, f.Mes, f.Anio)
}

func (cd *CustomFecha) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	partes := strings.Split(s, "-")
	if len(partes) != 3 {
		return fmt.Errorf("formato de fecha invalido: %s", s)
	}

	dia, err := strconv.Atoi(partes[0])
	if err != nil {
		return err
	}
	mes, err := strconv.Atoi(partes[1])
	if err != nil {
		return err
	}
	anio, err := strconv.Atoi(partes[2])
	if err != nil {
		return err
	}

	cd.Dia = dia
	cd.Mes = mes
	cd.Anio = anio

	return nil
}
