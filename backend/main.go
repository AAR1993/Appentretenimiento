package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq" // ← nuevo driver
	"log"
	"net/http"
	"os"
)

// DB es la conexión global que usarán los repositorios.
var DB *sql.DB

// InitDB abre la conexión y realiza un ping para validar que esté viva.
func InitDB() {
	// Obtén DSN de variables de entorno:  user:pass@tcp(host:port)/dbname?parseTime=true
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		// Ejemplo local por defecto
		dsn = "postgresql://alejandro:W29nP10f5i1m7lwmiODFHKW6k8ItZrxt@dpg-d1910s15pdvs73dqvn80-a.oregon-postgres.render.com/entretenimiento"
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("abriendo BD: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("ping BD: %v", err)
	}
	fmt.Println("Conexión a PostgreSQL exitosa")
}

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

//var lugares []Lugar
//var usuarios = map[string]Usuario{
//	"juanp": {
//		Nombre:    "Alejandro Rivera",
//		Usuario:   "juanp",
//		Password:  "123456",
//		Imagen:    "ruta/a/la/imagen.jpg",
//		Ubicacion: "Ubicacion de Juan",
//	},
//}

//func init() {
//	lugares = []Lugar{
//		{
//			"Mombasa Restaurante",
//			"Supper Club Platos de autor y de origen",
//			"https://maps.app.goo.gl/qcszb5tn8m4D42ct8",
//			"imgs/mombasa.png",
//			"9 AM - 10 PM",
//		},
//		{
//			"The Rooftop Garden",
//			"Resto-bar en la ciudad de Medellín, la mejor música, comida, parche futbolero, fiesta ¡Y mucho más!",
//			"https://maps.app.goo.gl/sWGmWEQyCdLzeDU29",
//			"imgs/rooftop.png",
//			"8 AM - 8 PM",
//		},
//		{
//			"La República Restaurante Bar",
//			"Somos un restaurante bar muy Colombiano, donde encuentras comida típica, música crossover, DJ y artistas en vivo los fines de semana, contamos con 10 pantallas gigantes para disfrutar los mejores partidos y demas deportes de todos los gustos.",
//			"https://maps.app.goo.gl/cpzVpHvU5MsAYTfJ8",
//			"imgs/republica.png",
//			"10 AM - 8 PM",
//		},
//		{
//			"37 PARK BISTRÓ BAR",
//			"Almuerzo, Cena, Brunch, Abierto hasta tarde, Bebidas",
//			"https://maps.app.goo.gl/keK58AaMjPgGHpt39",
//			"imgs/37Park.png",
//			"6 PM - 2 AM",
//		},
//	}
//}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Cache-Control, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func getLugares(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	lugares, err := ListLugares() // ← leer de la BD
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error al obtener los lugares",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lugares)
}

//func getUsuario(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(usuarios["juanp"])
//}

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

	// Get user from database
	user, err := GetUsuarioByUsername(loginData.Username)
	if err != nil {
		// User not found or other error
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Usuario o contraseña incorrectos"})
		return
	}

	// Verify password (plain text comparison - in production, use bcrypt)

	if user.Password != loginData.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Usuario o contraseña incorrectos"})
		return
	}

	// Login successful
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Inicio de sesión exitoso",
		"user": map[string]interface{}{
			"nombre":    user.Nombre,
			"usuario":   user.Usuario,
			"imagen":    user.Imagen,
			"ubicacion": user.Ubicacion,
		},
	})
}

func createLugarHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var lugar Lugar
	if err := json.NewDecoder(r.Body).Decode(&lugar); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validaciones mínimas
	if lugar.Nombre == "" || lugar.Descripcion == "" {
		http.Error(w, "Nombre y descripción son obligatorios", http.StatusBadRequest)
		return
	}

	id, err := CreateLugar(lugar)
	if err != nil {
		log.Printf("Error creando lugar: %v", err)
		http.Error(w, "Error creando lugar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          id,
		"nombre":      lugar.Nombre,
		"descripcion": lugar.Descripcion,
		"ubicacion":   lugar.Ubicacion,
		"imagen":      lugar.Imagen, // ← nuevo
		"horario":     lugar.Horario,
	})
}

// Register a new user
func registerHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	fmt.Println("registrando usuario")
	if r.Method == "OPTIONS" {
		fmt.Println("Entró a options")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		fmt.Println("Entró a post")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser Usuario
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		fmt.Printf("Error al decodificar el body %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if newUser.Usuario == "" || newUser.Password == "" || newUser.Nombre == "" {
		fmt.Println("Error registrando por parametros")
		http.Error(w, "Username, password, and name are required", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	if _, err := GetUsuarioByUsername(newUser.Usuario); err == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Create the user
	id, err := CreateUsuario(newUser)
	if err != nil {
		fmt.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Return the created user (without password)
	newUser.Password = "" // Don't return the password
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        id,
		"usuario":   newUser.Usuario,
		"nombre":    newUser.Nombre,
		"imagen":    newUser.Imagen,
		"ubicacion": newUser.Ubicacion,
	})
}

// withCORS envuelve un handler agregando cabeceras CORS y manejando OPTIONS.
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func lugaresHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getLugares(w, r)
	case http.MethodPost:
		createLugarHandler(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	log.Println("Iniciando servidor...")
	InitDB()
	http.HandleFunc("/api/lugares", withCORS(lugaresHandler))
	//http.HandleFunc("/api/usuario", withCORS(getUsuario))
	http.HandleFunc("/api/login", withCORS(loginHandler))
	http.HandleFunc("/api/usuarios", withCORS(registerHandler))

	// Servir archivos estáticos del frontend
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	// Leer puerto de variable de entorno (Render la define)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Servidor corriendo en http://0.0.0.0:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
