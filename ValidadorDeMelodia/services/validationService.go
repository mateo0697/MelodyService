package validationService

import (
	Nota "ValidadorDeMelodia/models"
	"regexp"
	"strconv"
	"strings"
)

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (s *ValidationService) ValidateMelody(melody string) map[string]interface{} {
	var notas = map[string]string{"A": "La", "B": "Si", "C": "Do", "D": "Re", "E": "Mi", "F": "Fa", "G": "Sol", "S": "Silencio"}
	var response = map[string]interface{}{}
	var notasValidas = []map[string]interface{}{}
	var separacion []string = strings.Split(melody, " ")
	var tempo, err = strconv.Atoi(separacion[0])
	if err != nil {
		var tempoErroneo []string = strings.Split(separacion[0], "")
		for i := 0; i <= len(tempoErroneo); i++ {
			if _, err := strconv.Atoi(tempoErroneo[i]); err != nil {
				response["cause"] = "error at position " + strconv.Itoa(i+1)
				return response
			}
		}
	} else {
		//Inicializo la distancia recorrida teniendo en cuenta el largo del tempo y el espacio que separa el tempo de las especificaciones
		var distanciaRecorridaHastaElMomento = len(separacion[0]) + 1
		for i := 1; i < len(separacion); i++ {
			var infoDeNota = separacion[i]
			var nota string = notas[strings.ToUpper(string(infoDeNota[0]))]
			if nota == "" {
				//Le sumo un lugar a la distancia para mostrar que el error esta en la notas
				response["cause"] = "error at position " + strconv.Itoa(distanciaRecorridaHastaElMomento+1)
				return response
			} else {
				// Le sumo uno a la distancia recorrida ya que la nota esta bien por lo cual se suma un lugar valido
				distanciaRecorridaHastaElMomento += 1
				infoDeNota = infoDeNota[1:]
				var Nota = validateNotes(infoDeNota, nota)
				if Nota.Error != 0 {
					response["cause"] = "error at position " + strconv.Itoa(distanciaRecorridaHastaElMomento+Nota.Error)
					return response
				}
				notasValidas = append(notasValidas, Nota.GetInfo())
				//Le sumo a la distancia recorrida todas las especificaciones de la nota mas el espacio
				distanciaRecorridaHastaElMomento += len(infoDeNota) + 1
			}
		}
	}
	response["tempo"] = map[string]interface{}{"value": tempo, "unit": "bpm"}
	response["notes"] = notasValidas
	return response
}

func validateNotes(infoDeNota string, nota string) *Nota.Nota {
	var isSilence = nota == "Silencio"
	//octava, duracion, alteración, nota
	var Nota = Nota.NewNota(4, 1, "n", nota)
	//Tipos de especificaciones validas
	var clavesEspecificacionValidas = map[string]string{"a": "^[nb#]$", "d": "^([0-4]||[0-9]+\\/[0-9]+)$", "o": "^[0-8]$"}
	//Valores que pueden tomar el ultimo caracter de un valor de especificacion, un numero para o y d, y n/b/# para la alteracion
	var valoresDeEspecificacionValidos = map[string]string{"0": "numero", "1": "numero", "2": "numero", "3": "numero", "4": "numero", "5": "numero", "6": "numero", "7": "numero", "8": "numero", "9": "numero", "#": "alteracion", "n": "alteracion", "b": "alteracion"}

	//Funcion que dependiendo de la clave valida con distintos regex
	validarEspecificacion := func(clave string, valor string) int {
		var regex = clavesEspecificacionValidas[clave]
		if !regexp.MustCompile(regex).MatchString(valor) {
			return 1
		} else if regex == "^([0-4]||[0-9]+\\/[0-9]+)$" && strings.Contains(valor, "/") {
			var numerosDeFraccion = strings.Split(valor, "/")
			var numerador, _ = strconv.ParseFloat(numerosDeFraccion[0], 64)
			var denominador, _ = strconv.ParseFloat(numerosDeFraccion[1], 64)
			if numerador/denominador > 4 {
				return 1
			}
		}
		return 0
	}

	//Si no hay informacion se pasa a la siguiente
	lenInforDeNota := len(infoDeNota)
	if lenInforDeNota == 0 {
		Nota.ModificarEspecificacion("e", "0")
		return Nota
	} else {
		//Primero valido que el primer caracter sea { si no ya falla
		if string(infoDeNota[0]) != "{" {
			Nota.ModificarEspecificacion("e", "1")
			return Nota
		}
		var recorridoHastaElMomento = 1
		//Variable que guarda la clave de especificacion actual
		var claveEspecificacion = ""
		//Variable que determina si esta el igual
		var equalAproved = false
		//Variable que guarda el valor de especificacion actual
		var valorEspecificacion = ""

		//Recorro cada uno de los caracteres de las especificaciones
		for i := 1; i < lenInforDeNota; i++ {
			var caracterActual = string(infoDeNota[i])
			//Primero valido si el caracter es una clave de especificacion
			if clavesEspecificacionValidas[caracterActual] != "" {
				if claveEspecificacion == "" {
					//Si es una clave y todavia no tengo guardada una:
					if isSilence && caracterActual != "d" {
						//Si es un silencio y es una clave distinta a duracion falla
						Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+1))
						return Nota
					}
					//Si no guardo la clave actual y continuo
					claveEspecificacion = caracterActual
					recorridoHastaElMomento++
					continue
				} else {
					//Si es una clave pero ya tengo guardada una verifico si es que falta el separador o hay uno invalido y luego falla
					if valoresDeEspecificacionValidos[string(infoDeNota[i-1])] != "" {
						recorridoHastaElMomento++
					}
					Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+len(valorEspecificacion)))
					return Nota
				}
			} else if claveEspecificacion == "" {
				//Si no es una clave y no tengo una en memoria significa que hay una clave invalida
				Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+1))
				return Nota
			}
			if !equalAproved {
				//Si llego hasta aca significa que no es una clave valida peor ya existe una en memoria por lo cual el siguiente caracter tiene que ser =, si no lo es falla
				if caracterActual != "=" {
					Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+1))
					return Nota
				} else {
					equalAproved = true
					recorridoHastaElMomento++
				}
			} else {
				//Si llega hasta aca significa que el caracter no es una clave valida, tengo clave actual y ya tengo el igual
				if caracterActual == ";" {
					//Si el caracter es el de separacion, valido el valor de especificacion
					var resultado = validarEspecificacion(claveEspecificacion, valorEspecificacion)
					if resultado != 0 {
						Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+resultado))
						return Nota
					} else {
						//Si el valor es valido agrego la distancia del valor al recorrido y reseteo variables
						Nota.ModificarEspecificacion(claveEspecificacion, valorEspecificacion)
						recorridoHastaElMomento = recorridoHastaElMomento + len(valorEspecificacion) + 1
						claveEspecificacion = ""
						equalAproved = false
						valorEspecificacion = ""
					}
					if i == lenInforDeNota-2 || i == lenInforDeNota-1 {
						Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento))
						return Nota
					}
				} else if i == lenInforDeNota-1 {
					//Si estoy en el ultimo lugar de la especificacion:
					if valoresDeEspecificacionValidos[caracterActual] != "" {
						//Y el caracter es un valor de especificacion valido, se lo sumo al valor actual y verifico
						valorEspecificacion += caracterActual
						var resultado = validarEspecificacion(claveEspecificacion, valorEspecificacion)
						if resultado != 0 {
							//Si el valor no es valido devuelvo el error
							Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+resultado))
							return Nota
						}
						//Si el valor es valido tengo que devolver el error porque el ultimo caracter de la especificacion no es un }
						Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+len(valorEspecificacion)+1))
						return Nota
					} else {
						//Si el ultimo lugar no es un valor valido significa que es o un finalizador o un valor invalido por lo que primero valido el valor
						var resultado = validarEspecificacion(claveEspecificacion, valorEspecificacion)
						if resultado != 0 {
							//Si el valor no es valido devuelvo el error
							Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+resultado))
							return Nota
						}
						Nota.ModificarEspecificacion(claveEspecificacion, valorEspecificacion)
						if caracterActual != "}" {
							//Y si el ultimo lugar no es } y el valor es valido fallo
							Nota.ModificarEspecificacion("e", strconv.Itoa(recorridoHastaElMomento+len(valorEspecificacion)+1))
							return Nota
						}
					}
				} else {
					//Si llego hasta aca significa que tengo que sumar un caracter al valor a validar
					valorEspecificacion += caracterActual
				}
			}
		}
	}
	//Si llegue hasta aca es valida la especifacion
	return Nota
}

/*
func (s *Service) ValidateMelody(melody string) string {
	var notas = map[string]string{"A": "La", "B": "Si", "C": "Do", "D": "Re", "E": "Mi", "F": "Fa", "G": "Sol", "S": "Silencio"}
	var separacion []string = strings.Split(melody, " ")
	var _, err = strconv.Atoi(separacion[0])
	if err != nil {
		var tempoErroneo []string = strings.Split(separacion[0], "")
		for i := 0; i <= len(tempoErroneo); i++ {
			if _, err := strconv.Atoi(tempoErroneo[i]); err != nil {
				fmt.Println("Error en la melodia en la posicion ", i+1)
				return 0
			}
		}
	} else {
		//Inicializo la distancia recorrida teniendo en cuenta el largo del tempo y el espacio que separa el tempo de las especificaciones
		var distanciaRecorridaHastaElMomento = len(separacion[0]) + 1
		for i := 1; i < len(separacion); i++ {
			var infoDeNota = separacion[i]
			var nota string = notas[strings.ToUpper(string(infoDeNota[0]))]
			if nota == "" {
				//Le sumo un lugar a la distancia para mostrar que el error esta en la notass
				fmt.Println("Error en la melodia en la posicion ", distanciaRecorridaHastaElMomento+1)
			} else {
				// Le sumo uno a la distancia recorrida ya que la nota esta bien por lo cual se suma un lugar valido
				distanciaRecorridaHastaElMomento += 1
				infoDeNota = infoDeNota[1:]
				var posicionDeError = validateNotes(infoDeNota, nota)
				if posicionDeError != 0 {
					fmt.Println("Error en la melodia en la posicion ", distanciaRecorridaHastaElMomento+posicionDeError)
					return 0
				}
				//Le sumo a la distancia recorrida todas las especificaciones de la nota mas el espacio
				distanciaRecorridaHastaElMomento += len(infoDeNota) + 1
			}
		} //texto = "60 A{d=7/4;o=3;a=#} B{o=2;d=1/2} S A{d=2;a=n} G{a=b} B S{d=1/3}"
		fmt.Println("Melodia valida")
	}
	return 0
}

func validateNotes(infoDeNota string, nota string) int {
	var isSilence = nota == "Silencio"
	//Tipos de especificaciones validas
	var clavesEspecificacionValidas = map[string]string{"a": "^[nb#]$", "d": "^([0-4]||[0-9]+\\/[0-9]+)$", "o": "^[0-8]$"}
	//Valores que pueden tomar el ultimo caracter de un valor de especificacion, un numero para o y d, y n/b/# para la alteracion
	var valoresDeEspecificacionValidos = map[string]string{"0": "numero", "1": "numero", "2": "numero", "3": "numero", "4": "numero", "5": "numero", "6": "numero", "7": "numero", "8": "numero", "9": "numero", "#": "alteracion", "n": "alteracion", "b": "alteracion"}

	//Funcion que dependiendo de la clave valida con distintos regex
	validarEspecificacion := func(clave string, valor string) int {
		var regex = clavesEspecificacionValidas[clave]
		if !regexp.MustCompile(regex).MatchString(valor) {
			return 1
		} else if regex == "^([0-4]||[0-9]+\\/[0-9]+)$" && strings.Contains(valor, "/") {
			var numerosDeFraccion = strings.Split(valor, "/")
			var numerador, _ = strconv.ParseFloat(numerosDeFraccion[0], 64)
			var denominador, _ = strconv.ParseFloat(numerosDeFraccion[1], 64)
			if numerador/denominador > 4 {
				return 1
			}
		}
		return 0
	}

	//Si no hay informacion se pasa a la siguiente
	lenInforDeNota := len(infoDeNota)
	if lenInforDeNota == 0 {
		return 0
	} else {
		//Primero valido que el primer caracter sea { si no ya falla
		if string(infoDeNota[0]) != "{" {
			return 1
		}
		var recorridoHastaElMomento = 1
		//Variable que guarda la clave de especificacion actual
		var claveEspecificacion = ""
		//Variable que determina si esta el igual
		var equalAproved = false
		//Variable que guarda el valor de especificacion actual
		var valorEspecificacion = ""

		//Recorro cada uno de los caracteres de las especificaciones
		for i := 1; i < lenInforDeNota; i++ {
			var caracterActual = string(infoDeNota[i])
			//Primero valido si el caracter es una clave de especificacion
			if clavesEspecificacionValidas[caracterActual] != "" {
				if claveEspecificacion == "" {
					//Si es una clave y todavia no tengo guardada una:
					if isSilence && caracterActual != "d" {
						//Si es un silencio y es una clave distinta a duracion falla
						return recorridoHastaElMomento + 1
					}
					//Si no guardo la clave actual y continuo
					claveEspecificacion = caracterActual
					recorridoHastaElMomento++
					continue
				} else {
					//Si es una clave pero ya tengo guardada una verifico si es que falta el separador o hay uno invalido y luego falla
					if valoresDeEspecificacionValidos[string(infoDeNota[i-1])] != "" {
						recorridoHastaElMomento++
					}
					return recorridoHastaElMomento + len(valorEspecificacion)
				}
			} else if claveEspecificacion == "" {
				//Si no es una clave y no tengo una en memoria significa que hay una clave invalida
				return recorridoHastaElMomento + 1
			}
			if !equalAproved {
				//Si llego hasta aca significa que no es una clave valida peor ya existe una en memoria por lo cual el siguiente caracter tiene que ser =, si no lo es falla
				if caracterActual != "=" {
					return recorridoHastaElMomento + 1
				} else {
					equalAproved = true
					recorridoHastaElMomento++
				}
			} else {
				//Si llega hasta aca significa que el caracter no es una clave valida, tengo clave actual y ya tengo el igual
				if caracterActual == ";" {
					//Si el caracter es el de separacion, valido el valor de especificacion
					var resultado = validarEspecificacion(claveEspecificacion, valorEspecificacion)
					if resultado != 0 {
						return recorridoHastaElMomento + resultado
					} else {
						//Si el valor es valido agrego la distancia del valor al recorrido y reseteo variables
						recorridoHastaElMomento = recorridoHastaElMomento + len(valorEspecificacion) + 1
						claveEspecificacion = ""
						equalAproved = false
						valorEspecificacion = ""
					}
					if i == lenInforDeNota-2 || i == lenInforDeNota-1 {
						return recorridoHastaElMomento
					}
				} else if i == lenInforDeNota-1 {
					//Si estoy en el ultimo lugar de la especificacion:
					if valoresDeEspecificacionValidos[caracterActual] != "" {
						//Y el caracter es un valor de especificacion valido, se lo sumo al valor actual y verifico
						valorEspecificacion += caracterActual
						var resultado = validarEspecificacion(claveEspecificacion, valorEspecificacion)
						if resultado != 0 {
							//Si el valor no es valido devuelvo el error
							return recorridoHastaElMomento + resultado
						}
						//Si el valor es valido tengo que devolver el error porque el ultimo caracter de la especificacion no es un }
						return recorridoHastaElMomento + len(valorEspecificacion) + 1
					} else {
						//Si el ultimo lugar no es un valor valido significa que es o un finalizador o un valor invalido por lo que primero valido el valor
						var resultado = validarEspecificacion(claveEspecificacion, valorEspecificacion)
						if resultado != 0 {
							//Si el valor no es valido devuelvo el error
							return recorridoHastaElMomento + resultado
						}
						if caracterActual != "}" {
							//Y si el ultimo lugar no es } y el valor es valido fallo
							return recorridoHastaElMomento + len(valorEspecificacion) + 1
						}
					}
				} else {
					//Si llego hasta aca significa que tengo que sumar un caracter al valor a validar
					valorEspecificacion += caracterActual
				}
			}
		}
	}
	//Si llegue hasta aca es valida la especifacion

	return 0
}

*/
