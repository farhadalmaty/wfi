package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type PageData struct {
	Title string
}

func RegisterHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(w, nil)
	})
	http.HandleFunc("/register", registerHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func DbStart() {
	dbPath := "bd.db"
	// Проверяем, существует ли файл базы данных
    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        // Если файл не существует, создаем новую базу данных
        err := createDatabase(dbPath)
        if err != nil {
			fmt.Println("Failed to create database:", err)
            // log.Fatal("Failed to create database:", err)
        }
    }

    // Теперь база данных существует, вы можете открыть ее для дальнейшей работы
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
		fmt.Println("Failed to open database:", err)
    }
    defer db.Close()
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := PageData{
			Title: "Register",
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		
		if userExists(username) {
			fmt.Println("Username already exists:", username)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
		// сохранение в БД
		// Открываем соединение с базой данных
        db, err := sql.Open("sqlite3", "bd.db")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer db.Close()

        // Подготавливаем запрос для вставки данных в базу данных
        stmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
        if err != nil {
			fmt.Println("Failed to prepare statement:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer stmt.Close()

        // Выполняем запрос для вставки данных
        _, err = stmt.Exec(username, password)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

		fmt.Println("Username:", username, "Password:", password)
		// переход на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func userExists(username string) bool {
    db, err := sql.Open("sqlite3", "bd.db")
    if err != nil {
        log.Fatal(err)
		return false
    }
    defer db.Close()

	// Создаем необходимые таблицы в базе данных
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL
        )
    `)
    if err != nil {
        fmt.Println("Failed to create table:", err)
		return false
    }


    // Подготавливаем запрос для проверки наличия пользователя в базе данных
    stmt, err := db.Prepare("SELECT COUNT(*) FROM users WHERE username = ?")
    if err != nil {
        log.Fatal(err)
		return false
    }
    defer stmt.Close()

    var count int
    // Выполняем запрос
    err = stmt.QueryRow(username).Scan(&count)
    if err != nil {
        log.Fatal(err)
    }

    // Если количество пользователей с указанным именем больше нуля, значит, пользователь существует
    return count > 0
}

func createDatabase(dbFilePath string) error {
	fmt.Println("Creating database:", dbFilePath)
    // Создаем новый файл базы данных
    dbFile, err := os.Create(dbFilePath)
    if err != nil {
        return err
    }
    dbFile.Close()

    // Открываем новую базу данных
    db, err := sql.Open("sqlite3", dbFilePath)
    if err != nil {
        return err
    }
    defer db.Close()

    // Создаем необходимые таблицы в базе данных
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL
        )
    `)
    if err != nil {
        return err
    }

    return nil
}