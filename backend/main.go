package main


import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Lugar struct {
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Ubicacion   string `json:"ubicacion"`
	Imagen      string `json:"imagen"`
	Horario     string `json:"horario"`
}

type Usuario struct {
	Nombre    string `json:"nombre"`
	Usuario   string `json:"usuario"`
	Password  string `json:"password"`
	Imagen    string `json:"imagen"`
	Ubicacion string `json:"ubicacion"`
}

var lugares []Lugar
var usuarios = map[string]Usuario{
	"juanp": {
		Nombre:    "Alejandro Rivera",
		Usuario:   "juanp",
		Password:  "123456",
		Imagen:    "ruta/a/la/imagen.jpg",
		Ubicacion: "Ubicacion de Juan",
	},
}

func init() {
	lugares = []Lugar{
		{
			"Mombasa Restaurante",
			"Supper Club Platos de autor y de origen",
			"https://maps.app.goo.gl/qcszb5tn8m4D42ct8",
			"imgs/mombasa.png",
			"9 AM - 10 PM",
		},
		{
			"The Rooftop Garden",
			"Resto-bar en la ciudad de Medellín, la mejor música, comida, parche futbolero, fiesta ¡Y mucho más!",
			"https://maps.app.goo.gl/sWGmWEQyCdLzeDU29",
			"imgs/rooftop.png",
			"8 AM - 8 PM",
		},
		{
			"La República Restaurante Bar",
			"Somos un restaurante bar muy Colombiano, donde encuentras comida típica, música crossover, DJ y artistas en vivo los fines de semana, contamos con 10 pantallas gigantes para disfrutar los mejores partidos y demas deportes de todos los gustos.",
			"https://maps.app.goo.gl/cpzVpHvU5MsAYTfJ8",
			"imgs/republica.png",
			"10 AM - 8 PM",
		},
		{
			"37 PARK BISTRÓ BAR",
			"Almuerzo, Cena, Brunch, Abierto hasta tarde, Bebidas",
			"https://maps.app.goo.gl/keK58AaMjPgGHpt39",
			"imgs/37Park.png",
			"6 PM - 2 AM",
		},
	}
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func getLugares(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	if len(lugares) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]Lugar{})
		return
	}

	err := json.NewEncoder(w).Encode(lugares)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al codificar los lugares"})
		return
	}
}

func getUsuario(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios["juanp"])
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Datos de inicio de sesión inválidos"})
		return
	}

	user, exists := usuarios[loginData.Username]
	if !exists || user.Password != loginData.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Usuario o contraseña incorrectos"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Inicio de sesión exitoso"})
}

func main() {
	println("Iniciando servidor...")

	http.HandleFunc("/api/lugares", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}
		getLugares(w)
	})

	http.HandleFunc("/api/usuario", func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		getUsuario(w)
	})

	http.HandleFunc("/api/login", loginHandler)

	// Servir archivos estáticos del frontend
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	// Leer puerto de variable de entorno (Render la define)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Servidor corriendo en http://0.0.0.0:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
