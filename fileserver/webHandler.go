package fileserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	fileServerAdmin = http.FileServer(http.Dir("./test/test1"))

	mimeTypes = map[string]string{
		".css":  "text/css",
		".js":   "application/javascript",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".html": "text/html",
		".json": "application/json",
		".xml":  "application/xml",
		".txt":  "text/plain",
		".pdf":  "application/pdf",
	}
)

//Handel the websits

func GetHomePage(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./web/dist/homePage/index.html"))

	if err := temp.Execute(w, nil); err != nil {
		log.Panic(err)
	}
}

func GetWebEvent(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./web/dist/event/Event.html"))

	if err := temp.Execute(w, nil); err != nil {
		log.Panic(err)
	}
}

func GetAddSection(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("./web/dist/Admin/index.html"))

	if err := temp.Execute(w, nil); err != nil {
		log.Panic(err)
	}
}

//Handel the Sever files

func HandelWeb(w http.ResponseWriter, r *http.Request) {
	newPath := strings.TrimPrefix(r.URL.Path, "/")
	filePath := filepath.Join("web", "dist", "homePage", newPath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
        http.NotFound(w, r)
        return
    }

	ext := filepath.Ext(filePath)
	fmt.Println(ext)
	if mimeType, exists := mimeTypes[ext]; exists {
        w.Header().Set("Content-Type", mimeType)
    }

	http.ServeFile(w, r, filePath)
}


func HandelAdmin(w http.ResponseWriter, r *http.Request) {
	newPath := strings.TrimPrefix(r.URL.Path, "/Admin/")
	filePath := filepath.Join("web", "dist", "Admin", newPath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
        http.NotFound(w, r)
        return
    }

	ext := filepath.Ext(filePath)
	fmt.Println(ext)
	if mimeType, exists := mimeTypes[ext]; exists {
        w.Header().Set("Content-Type", mimeType)
    }

	http.ServeFile(w, r, filePath)
}


func HandelOutPutCss(w http.ResponseWriter, r *http.Request) {
	newPath := strings.TrimPrefix(r.URL.Path, "/web/dist/")
	filePath := filepath.Join("web", "dist", newPath)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
        http.NotFound(w, r)
        return
    }

	ext := filepath.Ext(filePath)
	fmt.Println(ext)
	if mimeType, exists := mimeTypes[ext]; exists {
        w.Header().Set("Content-Type", mimeType)
    }

	http.ServeFile(w, r, filePath)
}

// This section will handel the upload of the files
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	const maxUploadSize = 10 * 1024 * 1024 // 10 mb
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		fmt.Println(err)
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid imge", http.StatusBadRequest)
		return
	}
	uploadDir := "./web/dist/Admin/imges_admin_section"
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), header.Filename)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error while creating the URL file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		fmt.Println(err)

		http.Error(w, "Error while saving the file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"fileURL": "/imges_admin_section/" + filename,
		"message": "File uploaded successfully",
	})
}
