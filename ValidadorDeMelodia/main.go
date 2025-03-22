package main

import (
	Nota "ValidadorDeMelodia/models"
	validationService "ValidadorDeMelodia/services"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/speaker"
)

type Melodia struct {
	Melodia string      `json:"melody"`
	Tempo   Tempo       `json:"tempo"`
	Notes   []Nota.Nota `json:"notes"`
}

type Tempo struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"` // En este caso, podemos suponer que siempre será "bpm"
}

func validarMelodia(w http.ResponseWriter, r *http.Request) {
	var melody Melodia
	err := json.NewDecoder(r.Body).Decode(&melody)
	if err != nil {
		http.Error(w, "Error parseando la melodia", http.StatusBadRequest)
	}
	validationSvc := validationService.NewValidationService()

	var result = validationSvc.ValidateMelody(melody.Melodia)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func playMelody(w http.ResponseWriter, r *http.Request) {
	var melody Melodia
	err := json.NewDecoder(r.Body).Decode(&melody)
	if err != nil {
		http.Error(w, "Error parseando la melodia", http.StatusBadRequest)
	}
	sr := beep.SampleRate(44100)
	sounds := []beep.Streamer{}
	beatsPerSecond := float64(melody.Tempo.Value) / 60.0
	for _, note := range melody.Notes {
		noteDuration := 1.0 / beatsPerSecond * float64(note.Duracion)
		noteDurationSamples := int(sr.N(time.Duration(noteDuration) * time.Second))
		stream, _ := generators.SineTone(sr, note.Frecuencia)
		sounds = append(sounds, beep.Take(noteDurationSamples, stream))
	}
	speaker.Init(sr, 4410)
	speaker.Play(beep.Seq(sounds...))
}

func main() {
	http.HandleFunc("/melody/validate", validarMelodia)
	http.HandleFunc("/melody/play", playMelody)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}

}
