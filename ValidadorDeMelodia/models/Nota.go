package models

import (
	"math"
	"strconv"
	"strings"
)

type Nota struct {
	Nota       string  `json:"name"`
	Duracion   float64 `json:"duration"`
	Alteración string  `json:"alteration"`
	Octava     float64 `json:"3"`
	Error      int
	Frecuencia float64 `json:"frequency"`
}

func NewNota(octava, duracion float64, alteración, nota string) *Nota {
	return &Nota{
		Nota:       nota,
		Duracion:   duracion,
		Alteración: alteración,
		Octava:     octava,
	}
}

func (n *Nota) GetInfo() map[string]interface{} {
	var info = map[string]interface{}{}
	info["duration"] = n.Duracion
	if n.Nota == "Silencio" {
		info["type"] = "silence"
		return info
	}
	info["type"] = "note"
	info["name"] = n.Nota
	info["octave"] = n.Octava
	info["alteration"] = n.Alteración
	info["frequency"] = n.CalcularFrecuencia()
	return info
}

func (n *Nota) CalcularFrecuencia() float64 {
	var valorNota float64 = 0
	var valorAlteracion float64 = 0
	switch n.Nota {
	case "La":
		valorNota = 9
	case "Si":
		valorNota = 11
	case "Do":
		valorNota = 0
	case "Re":
		valorNota = 2
	case "Mi":
		valorNota = 4
	case "Fa":
		valorNota = 5
	case "Sol":
		valorNota = 7
	}

	switch n.Alteración {
	case "b":
		valorAlteracion = -1
	case "#":
		valorAlteracion = 1
	}
	n.Frecuencia = math.Round((440*math.Pow(2, ((valorNota+valorAlteracion+12*n.Octava)-57)/12))*100) / 100
	return n.Frecuencia

}

func (n *Nota) ModificarEspecificacion(clave, valor string) {
	switch clave {
	case "a":
		n.Alteración = valor
	case "d":
		if strings.Contains(valor, "/") {
			var numerosDeFraccion = strings.Split(valor, "/")
			var numerador, _ = strconv.ParseFloat(numerosDeFraccion[0], 64)
			var denominador, _ = strconv.ParseFloat(numerosDeFraccion[1], 64)
			n.Duracion = numerador / denominador
		} else {
			var valorDuracionFloat, _ = strconv.ParseFloat(valor, 64)
			n.Duracion = valorDuracionFloat
		}
	case "o":
		var valorOctavaFloat, _ = strconv.ParseFloat(valor, 64)
		n.Octava = valorOctavaFloat
	case "e":
		n.Error, _ = strconv.Atoi(valor)
	}
}
