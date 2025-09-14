// smackapi/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Task represents a single task
type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"` // e.g., "pending", "in-progress", "done"
}

var (
	tasks   = make(map[int]Task)
	nextID  = 1
	tasksMu sync.Mutex
)

// helper to respond with JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// GET /tasks - list all tasks
func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	taskList := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		taskList = append(taskList, t)
	}
	respondJSON(w, http.StatusOK, taskList)
}

// POST /tasks - create a new task
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tasksMu.Lock()
	t.ID = nextID
	nextID++
	t.Status = "pending"
	tasks[t.ID] = t
	tasksMu.Unlock()

	respondJSON(w, http.StatusCreated, t)
}

// GET /tasks/{id} - get a task by ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	tasksMu.Lock()
	t, ok := tasks[id]
	tasksMu.Unlock()
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, t)
}

// PUT /tasks/{id} - update a task status
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var update struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tasksMu.Lock()
	task, ok := tasks[id]
	if !ok {
		tasksMu.Unlock()
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	task.Status = update.Status
	tasks[id] = task
	tasksMu.Unlock()

	respondJSON(w, http.StatusOK, task)
}

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listTasksHandler(w, r)
		case http.MethodPost:
			createTaskHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTaskHandler(w, r)
		case http.MethodPut:
			updateTaskHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := "8080"
	fmt.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
