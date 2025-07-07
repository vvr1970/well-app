package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"wells-app/models"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var tmpl *template.Template

func init() {
	// Подключение к PostgreSQL
	var err error
	db, err = models.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	// Загрузка шаблонов
	//tmpl = template.Must(template.ParseGlob("templates/*.html"))
	// Создаем FuncMap с нашими функциями
	funcMap := template.FuncMap{
		"formatNumber": func(n float64) string {
			return fmt.Sprintf("%.0f", n)
		},
		"title": strings.Title,
	}

	// Загружаем шаблоны с функциями
	tmpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
}

func main() {

	r := mux.NewRouter()

	// Статические файлы (графики)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Роуты
	r.HandleFunc("/", listWells).Methods("GET")
	r.HandleFunc("/create", createWellForm).Methods("GET")
	r.HandleFunc("/create", createWell).Methods("POST")
	r.HandleFunc("/edit/{id}", editWellForm).Methods("GET")
	r.HandleFunc("/edit/{id}", editWell).Methods("POST")
	r.HandleFunc("/delete/{id}", deleteWell).Methods("GET")
	r.HandleFunc("/view/{id}", viewWell).Methods("GET") // Важно!

	log.Println("Server started on :8081")
	http.ListenAndServe(":8081", r)
}

// Список всех скважин
func listWells(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, location, Pbuf, status FROM wells")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var wells []models.GasWell
	for rows.Next() {
		var well models.GasWell
		err := rows.Scan(&well.ID, &well.Name, &well.Location, &well.Pbuf, &well.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		wells = append(wells, well)
	}

	tmpl.ExecuteTemplate(w, "base", wells)
}

// Форма добавления
func createWellForm(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "basecreate", nil)
	if err != nil {
		log.Printf("Ошибка рендеринга: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

// Добавление скважины
func createWell(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	location := r.FormValue("location")
	gammag, _ := strconv.ParseFloat(r.FormValue("gammag"), 64)
	temp, _ := strconv.ParseFloat(r.FormValue("temp"), 64)
	tempust, _ := strconv.ParseFloat(r.FormValue("tempust"), 64)
	depth, _ := strconv.ParseFloat(r.FormValue("depth"), 64)
	pbuf, _ := strconv.ParseFloat(r.FormValue("pbuf"), 64)
	ptb, _ := strconv.ParseFloat(r.FormValue("ptb"), 64)
	ppl, _ := strconv.ParseFloat(r.FormValue("ppl"), 64)
	pz, _ := strconv.ParseFloat(r.FormValue("pz"), 64)
	q, _ := strconv.ParseFloat(r.FormValue("q"), 64)
	roughness, _ := strconv.ParseFloat(r.FormValue("roughness"), 64)
	diametr, _ := strconv.ParseFloat(r.FormValue("diametr"), 64)
	a, _ := strconv.ParseFloat(r.FormValue("a"), 64)
	b, _ := strconv.ParseFloat(r.FormValue("b"), 64)
	mu, _ := strconv.ParseFloat(r.FormValue("mu"), 64)
	wgf, _ := strconv.ParseFloat(r.FormValue("wgf"), 64)
	rog, _ := strconv.ParseFloat(r.FormValue("rog"), 64)
	hw, _ := strconv.ParseFloat(r.FormValue("hw"), 64)
	qmin, _ := strconv.ParseFloat(r.FormValue("qmin"), 64)
	pmax, _ := strconv.ParseFloat(r.FormValue("pmax"), 64)
	status := r.FormValue("status")

	sqlCrt := "INSERT INTO wells (name, location, gammag, temp, tempust, depth, pbuf,ptb, ppl, pz,q,roughness,"
	sqlCrt = sqlCrt + "diametr,a,b,mu,wgf,rog,hw,qmin,pmax, status)"
	sqlCrt = sqlCrt + " VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)"
	_, err := db.Exec(
		sqlCrt,
		name, location, gammag, temp, tempust, depth, pbuf, ptb, ppl, pz, q, roughness, diametr, a, b, mu, wgf,
		rog, hw, qmin, pmax, status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Форма редактирования
func editWellForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var well models.GasWell
	sqlEdit := "SELECT id, name, location, gammag, temp, tempust, depth, pbuf,ptb, ppl, pz,q,roughness,"
	sqlEdit = sqlEdit + "diametr,a,b,mu,wgf,rog,hw,qmin,pmax, status FROM wells WHERE id = $1"
	err := db.QueryRow(sqlEdit, id).Scan(
		&well.ID, &well.Name, &well.Location, &well.GammaG, &well.Temp, &well.TempUst, &well.Depth, &well.Pbuf,
		&well.Ptb, &well.Ppl, &well.Pz, &well.Q, &well.Roughness, &well.Diameter, &well.A, &well.B, &well.Mu,
		&well.WGF, &well.Rog, &well.Hw, &well.Qmin, &well.Pmax, &well.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "baseedit", well)
	if err != nil {
		log.Printf("Ошибка рендеринга: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}

// Обновление данных
func editWell(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	name := r.FormValue("name")
	location := r.FormValue("location")
	gammag, _ := strconv.ParseFloat(r.FormValue("gammag"), 64)
	temp, _ := strconv.ParseFloat(r.FormValue("temp"), 64)
	tempust, _ := strconv.ParseFloat(r.FormValue("tempust"), 64)
	depth, _ := strconv.ParseFloat(r.FormValue("depth"), 64)
	pbuf, _ := strconv.ParseFloat(r.FormValue("pbuf"), 64)
	ptb, _ := strconv.ParseFloat(r.FormValue("ptb"), 64)
	ppl, _ := strconv.ParseFloat(r.FormValue("ppl"), 64)
	pz, _ := strconv.ParseFloat(r.FormValue("pz"), 64)
	q, _ := strconv.ParseFloat(r.FormValue("q"), 64)
	roughness, _ := strconv.ParseFloat(r.FormValue("roughness"), 64)
	diametr, _ := strconv.ParseFloat(r.FormValue("diametr"), 64)
	a, _ := strconv.ParseFloat(r.FormValue("a"), 64)
	b, _ := strconv.ParseFloat(r.FormValue("b"), 64)
	mu, _ := strconv.ParseFloat(r.FormValue("mu"), 64)
	wgf, _ := strconv.ParseFloat(r.FormValue("wgf"), 64)
	rog, _ := strconv.ParseFloat(r.FormValue("rog"), 64)
	hw, _ := strconv.ParseFloat(r.FormValue("hw"), 64)
	qmin, _ := strconv.ParseFloat(r.FormValue("qmin"), 64)
	pmax, _ := strconv.ParseFloat(r.FormValue("pmax"), 64)
	status := r.FormValue("status")
	sqlEdit := "UPDATE wells SET name=$1, location=$2, gammag=$3, temp=$4, tempust=$5, depth=$6, pbuf=$7,"
	sqlEdit = sqlEdit + "ptb=$8, ppl=$9, pz=$10, q=$11, roughness=$12, diametr=$13, a=$14, b=$15, mu=$16, wgf=$17,"
	sqlEdit = sqlEdit + "rog=$18, hw=$19, qmin=$20, pmax=$21, status=$22 WHERE id=$23"
	_, err := db.Exec(
		sqlEdit,
		name, location, gammag, temp, tempust, depth, pbuf, ptb, ppl, pz, q, roughness, diametr, a, b, mu, wgf,
		rog, hw, qmin, pmax, status, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Удаление скважины
func deleteWell(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	_, err := db.Exec("DELETE FROM wells WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Просмотр деталей скважины
func viewWell(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Ошибка парсинга ID:", err)
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	sqlView := `
        SELECT id, name, location, gammag, temp, tempust, depth, pbuf, ptb, ppl, pz, q, 
               roughness, diametr, a, b, mu, wgf, rog, hw, qmin, pmax, status 
        FROM wells 
        WHERE id = $1`

	var well models.GasWell
	err = db.QueryRow(sqlView, id).Scan(
		&well.ID, &well.Name, &well.Location, &well.GammaG, &well.Temp, &well.TempUst, &well.Depth, &well.Pbuf,
		&well.Ptb, &well.Ppl, &well.Pz, &well.Q, &well.Roughness, &well.Diameter, &well.A, &well.B, &well.Mu,
		&well.WGF, &well.Rog, &well.Hw, &well.Qmin, &well.Pmax, &well.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Скважина не найдена", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = tmpl.ExecuteTemplate(w, "baseview", well)
	if err != nil {
		log.Printf("Ошибка рендеринга: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
	}
}
