package main

import (
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
    "sort"
    "sync"
)

const DELAY_IN_SECONDS = 4

var (
    images     []string
    currentIdx int
    mu         sync.Mutex
)

func loadImages() ([]string, error) {
    files, err := os.ReadDir("photos")
    if err != nil {
        return nil, err
    }

    var imgList []string
    for _, file := range files {
        if !file.IsDir() {
            ext := strings.ToLower(filepath.Ext(file.Name()))
            if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
                imgList = append(imgList, "/photos/"+file.Name())
            }
        }
    }

    sort.Strings(imgList)
    return imgList, nil
}

func slideshowHandler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    idx := currentIdx
    totalImages := len(images)
    mu.Unlock()

    if totalImages == 0 {
        http.Error(w, "No images available :(", http.StatusNotFound)
        return
    }

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
            .overlay {
                position: absolute;
                top: 10px;
                right: 10px;
                background-color: rgba(0, 0, 0, 0.25);
                color: #fff;
                padding: 5px 10px;
                border-radius: 3px;
                font-size: 14px;
            }
        </style>
        <meta http-equiv="refresh" content="{{ .RefreshInterval }}" />
    </head>
    <body>
        <div class="overlay">{{ add .CurrentIndex 1 }} of {{ .TotalImages }}</div>
        {{ if .Images }}
            {{ with index .Images .CurrentIndex }}
                <img id="slideshow" src="{{ . }}?{{ $.Timestamp }}" />
            {{ else }}
                <p>No image available :(</p>
            {{ end }}
        {{ else }}
            <p>No images available :(</p>
        {{ end }}
    </body>
    </html>`

    now := time.Now()
    timestamp := now.UnixNano() / int64(time.Millisecond)

    t := template.Must(template.New("slideshow").Funcs(template.FuncMap{
        "add": func(a, b int) int { return a + b },
    }).Parse(tmpl))
    err := t.Execute(w, map[string]interface{}{
        "Images":           images,
        "CurrentIndex":     idx,
        "TotalImages":      totalImages,
        "RefreshInterval":  DELAY_IN_SECONDS,
        "Timestamp":        timestamp,
    })
    if err != nil {
        http.Error(w, "Unable to load template :(", http.StatusInternalServerError)
        log.Println("Template execution error:", err)
    }
}

func updateIndex() {
    mu.Lock()
    defer mu.Unlock()

    newImages, err := loadImages()
    if err != nil {
        log.Println("Error loading images during index update:", err)
        return
    }

    if len(newImages) == 0 {
        images = []string{}
        currentIdx = 0
        return
    }

    if len(images) != len(newImages) || !equal(images, newImages) {
        images = newImages
        if len(images) > 0 {
            currentIdx = currentIdx % len(images)
        }
    } else {
        if len(images) > 0 {
            currentIdx = (currentIdx + 1) % len(images)
        }
    }
}

func equal(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func main() {
    var err error
    images, err = loadImages()
    if err != nil {
        log.Fatal("Error loading images:", err)
    }

    if len(images) > 0 {
        currentIdx = 0
    }

    http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("photos"))))
    http.HandleFunc("/", slideshowHandler)

    log.Println("Local server established at http://localhost:3000/")

    go func() {
        for {
            time.Sleep(DELAY_IN_SECONDS * time.Second)
            updateIndex()
        }
    }()

    if err := http.ListenAndServe(":3000", nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}
