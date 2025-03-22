Desafio tecnico de Melodias:

- Requisito necesarios para utilizarlo:
  - Go

- Forma de utilizarlo:
  - Ingresar a la carpeta "ValidadorDeMelodia" y ejecutar el comando go run .\main.go
  - Se habilitaran 2 endpoits:
    - POST /melody/validate: Validador de melodias
      - Espera siguiente formato de body:
         { "melody": "60 A{d=7/4;o=3;a=#} B{o=2;d=1/4} S G{d=2}" }
    - POST /melody/play: Reproductor de melodias
      - Espera siguiente formato de body
        { 
    "tempo": { 
        "value": 60, 
        "unit": "bpm" 
    }, 
    "notes": [ 
        { 
            "type": "note", 
            "name": "la", 
            "octave": 3, 
            "alteration": "#", 
            "duration": 1.75, 
            "frequency": 233.08 
        }
    ] 
}
