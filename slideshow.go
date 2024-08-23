package main

import (
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const DELAY_IN_SECONDS = 4

func loadImages() ([]string, error) {
    files, err := os.ReadDir("photos")
    if err != nil {
        return nil, err
    }

    var images []string
    for _, file := range files {
        if !file.IsDir() {
            ext := strings.ToLower(filepath.Ext(file.Name()))
            if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
                images = append(images, "/photos/"+file.Name())
            }
        }
    }
    return images, nil
}

func slideshowHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Image Slideshow</title>
        <style>
            body {
                text-align: center;
                background: #000;
                margin: 0;
                overflow: hidden;
                height: 100vh;
                width: 100vw;
            }
            #slideshow {
                max-width: 100%;
                max-height: 100%;
                object-fit: contain;
            }
        </style>
        <meta http-equiv="refresh" content="{{ .RefreshInterval }}" />
    </head>
    <body>
        {{ if .Images }}
            <img id="slideshow" src="{{ index .Images .CurrentIndex }}?{{ .Timestamp }}" />
        {{ else }}
            <p>No images available.</p>
        {{ end }}
    </body>
    </html>`

    images, err := loadImages()
    if err != nil {
        http.Error(w, "Unable to load images", http.StatusInternalServerError)
        log.Println("Error loading images:", err)
        return
    }

    currentIndex := (time.Now().Second() / DELAY_IN_SECONDS) % len(images)

    t := template.Must(template.New("slideshow").Parse(tmpl))
    if err := t.Execute(w, map[string]interface{}{
        "Images":           images,
        "CurrentIndex":     currentIndex,
        "RefreshInterval":  DELAY_IN_SECONDS,
        "Timestamp":        time.Now().Unix(),
    }); err != nil {
        http.Error(w, "Unable to load template", http.StatusInternalServerError)
        log.Println("Template execution error:", err)
        return
    }
}

func main() {
    http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("photos"))))
    http.HandleFunc("/", slideshowHandler)

    log.Println("Local server established at http://localhost:3000/")
    if err := http.ListenAndServe(":3000", nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}
