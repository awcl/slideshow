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

const delayInSeconds = 4

var (
    images     []string
    currentIdx int
    mu         sync.Mutex
)

func loadImages() ([]string, error) {
    var imgList []string

    err := filepath.WalkDir("photos", func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            ext := strings.ToLower(filepath.Ext(d.Name()))
            if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
                relativePath, _ := filepath.Rel("photos", path)
                imgList = append(imgList, "/photos/"+relativePath)
            }
        }
        return nil
    })
    if err != nil {
        return nil, err
    }

    sort.Strings(imgList)
    return imgList, nil
}

func formatTime(t time.Time) string {
    hour := t.Format("15")
    minute := t.Format("04")
    day := t.Format("02")
    month := t.Format("Jan")
    year := t.Format("2006")
    weekday := t.Format("Mon")
    timezone := t.Format("MST")
    return hour + ":" + minute + " " + timezone + " // " + weekday + " " + day + " " + month + " " + year
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

    now := time.Now()
    timestamp := now.UnixNano() / int64(time.Millisecond)

    currentTime := formatTime(now)

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
                position: relative;
            }
            #slideshow {
                position: absolute;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                object-fit: contain;
            }
            .overlay {
                position: absolute;
                top: 5px;
                left: 50%;
                transform: translateX(-50%);
                background-color: rgba(0, 0, 0, 0.6);
                color: #fff;
                padding: 4px;
                border-radius: 5px;
                font-size: 16px;
                z-index: 10;
                white-space: nowrap;
            }
        </style>
        <meta http-equiv="refresh" content="{{ .RefreshInterval }}" />
    </head>
    <body>
        <div class="overlay">
            {{ .CurrentTime }} // {{ add .CurrentIndex 1 }} of {{ .TotalImages }} Images
        </div>
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

    t := template.Must(template.New("slideshow").Funcs(template.FuncMap{
        "add": func(a, b int) int { return a + b },
    }).Parse(tmpl))

    if err := t.Execute(w, map[string]interface{}{
        "Images":           images,
        "CurrentIndex":     idx,
        "TotalImages":      totalImages,
        "RefreshInterval":  delayInSeconds,
        "Timestamp":        timestamp,
        "CurrentTime":      currentTime,
    }); err != nil {
        http.Error(w, "Unable to load template :(", http.StatusInternalServerError)
        log.Printf("Template execution error: %v", err)
    }
}

func updateIndex() {
    mu.Lock()
    defer mu.Unlock()

    newImages, err := loadImages()
    if err != nil {
        log.Printf("Error loading images during index update: %v", err)
        return
    }

    if len(newImages) == 0 {
        images = nil
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
        log.Fatalf("Error loading images: %v", err)
    }

    if len(images) > 0 {
        currentIdx = 0
    }

    http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("photos"))))
    http.HandleFunc("/", slideshowHandler)

    log.Println("Local server established at http://localhost:3000/")

    go func() {
        for {
            time.Sleep(delayInSeconds * time.Second)
            updateIndex()
        }
    }()

    if err := http.ListenAndServe(":3000", nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
