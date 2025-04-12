package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github/codingMaster/internal/app/model"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/xuri/excelize/v2"
)

func (t *Application) tst(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./leabgin/index.html"))

	if err := temp.Execute(w, nil); err != nil {
		log.Panic(err)
	}
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	var login model.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if _, err := app.Store.Admin.Login(r.Context(), &login); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func (app *Application) Singup(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := app.Store.Admin.Signup(r.Context(), &user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

func (app *Application) AddProject(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.AddProject(r.Context(), &project); err != nil {
		http.Error(w, "Failed to add project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Project added successfully"))
}

func (app *Application) DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.DeleteProject(r.Context(), id); err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project deleted successfully"))
}

func (app *Application) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.UpdateProject(r.Context(), &project); err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project updated successfully"))
}

func (app *Application) Upprove(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.Upprove(r.Context(), id, payload.Status); err != nil {
		http.Error(w, "Failed to update project status", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project status updated successfully"))
}

func (app *Application) GetProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := app.Store.Role.GetProject(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve project", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (app *Application) Unpprove(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	if err := app.Store.Role.Unpprove(r.Context(), id); err != nil {
		http.Error(w, "Failed to reject project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Project rejected successfully"))
}

func (app *Application) UploadPDF(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Header.Get("Content-Type") != "application/pdf" {
		http.Error(w, "Only PDF files are allowed", http.StatusUnsupportedMediaType)
		return
	}

	filePath := "./uploads/" + handler.Filename

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func (app *Application) DownloadPDF(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := "./uploads/" + fileName

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/pdf")

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
}

func (app *Application) ConvertJSON(w http.ResponseWriter, r *http.Request) {
	queryValue := fmt.Sprintf("%x", time.Now().UnixNano())
	if queryValue == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	var data []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON file", http.StatusBadRequest)
		return
	}

	csvFilePath := "./uploads/" + queryValue + ".csv"
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		http.Error(w, "Error creating CSV file", http.StatusInternalServerError)
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	if len(data) > 0 {
		headers := make([]string, 0, len(data[0]))
		for key := range data[0] {
			headers = append(headers, key)
		}
		writer.Write(headers)

		for _, row := range data {
			record := make([]string, 0, len(row))
			for _, key := range headers {
				value := ""
				if row[key] != nil {
					value = row[key].(string)
				}
				record = append(record, value)
			}
			writer.Write(record)
		}
	}

	excelFilePath := "./uploads/" + queryValue + ".xlsx"
	excelFile := excelize.NewFile()
	sheetName := "Sheet1"
	excelFile.NewSheet(sheetName)

	if len(data) > 0 {
		headers := make([]string, 0, len(data[0]))
		for key := range data[0] {
			headers = append(headers, key)
		}
		for i, header := range headers {
			cell := string(rune('A'+i)) + "1"
			excelFile.SetCellValue(sheetName, cell, header)
		}

		for rowIndex, row := range data {
			for colIndex, key := range headers {
				cell := string(rune('A'+colIndex)) + strconv.Itoa(rowIndex+2)
				value := ""
				if row[key] != nil {
					value = row[key].(string)
				}
				excelFile.SetCellValue(sheetName, cell, value)
			}
		}
	}

	if err := excelFile.SaveAs(excelFilePath); err != nil {
		http.Error(w, "Error creating Excel file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Files created successfully: output.csv and output.xlsx"))
}

func (app *Application) DownloadCSV(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := "./uploads/" + fileName

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/csv")

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
}

func (app *Application) DownloadExcel(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	filePath := "./uploads/" + fileName

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
}

