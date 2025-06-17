package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Employee struct {
	ID         int     `json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	Department string  `json:"department"`
	JobTitle   string  `json:"job_title"`
	HireDate   string  `json:"hire_date"`
	Salary     float64 `json:"salary"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

var (
	db *sql.DB

	// Prometheus metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
		},
		[]string{"method", "endpoint"},
	)

	employeesCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "employees_created_total",
			Help: "Total number of employees created",
		},
	)

	employeesTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "employees_total",
			Help: "Total number of employees in database",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(employeesCreated)
	prometheus.MustRegister(employeesTotal)
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")

	// Update employee count on startup
	updateEmployeeCount()

	// Setup routes
	r := mux.NewRouter()
	r.Use(metricsMiddleware)

	r.HandleFunc("/", serveIndex).Methods("GET")
	r.HandleFunc("/api/employees", createEmployee).Methods("POST")
	r.HandleFunc("/api/employees", getEmployees).Methods("GET")
	r.HandleFunc("/api/employees/{id}", deleteEmployee).Methods("DELETE")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture status code
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		endpoint := r.URL.Path
		method := r.Method
		status := strconv.Itoa(ww.statusCode)

		httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO employees (first_name, last_name, email, phone, department, job_title, hire_date, salary)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := db.QueryRow(query, emp.FirstName, emp.LastName, emp.Email, emp.Phone,
		emp.Department, emp.JobTitle, emp.HireDate, emp.Salary).Scan(&emp.ID, &emp.CreatedAt, &emp.UpdatedAt)

	if err != nil {
		log.Printf("Error creating employee: %v", err)
		http.Error(w, "Failed to create employee", http.StatusInternalServerError)
		return
	}

	employeesCreated.Inc()
	updateEmployeeCount()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, first_name, last_name, email, phone, department, job_title, hire_date, salary, created_at, updated_at FROM employees ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying employees: %v", err)
		http.Error(w, "Failed to fetch employees", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.ID, &emp.FirstName, &emp.LastName, &emp.Email, &emp.Phone,
			&emp.Department, &emp.JobTitle, &emp.HireDate, &emp.Salary, &emp.CreatedAt, &emp.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning employee: %v", err)
			continue
		}
		employees = append(employees, emp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM employees WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete employee", http.StatusInternalServerError)
		return
	}

	updateEmployeeCount()
	w.WriteHeader(http.StatusNoContent)
}

func updateEmployeeCount() {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM employees").Scan(&count)
	if err != nil {
		log.Printf("Error counting employees: %v", err)
		return
	}
	employeesTotal.Set(float64(count))
}
