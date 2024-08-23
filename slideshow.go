package main

import (
    "encoding/json"
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

const DELAY_IN_SECONDS = 4
const DELAY_IN_MS = DELAY_IN_SECONDS * 1000

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
    </head>
    <body>
        <img id="slideshow" src="" />
        <script>
            const DELAY_IN_MS = {{ .DelayInMs }};

            async function fetchImages() {
                const response = await fetch('/images');
                const images = await response.json();
                return images;
            }

            let index = 0;
            async function showImage() {
                const images = await fetchImages();
                if (images.length === 0) return;
                document.getElementById('slideshow').src = images[index];
                index = (index + 1) % images.length;
            }

            setInterval(showImage, DELAY_IN_MS);
            showImage();
        </script>
    </body>
    </html>`

    t := template.Must(template.New("slideshow").Parse(tmpl))
    if err := t.Execute(w, map[string]interface{}{
        "DelayInMs": DELAY_IN_MS,
    }); err != nil {
        http.Error(w, "Unable to load template", http.StatusInternalServerError)
        return
    }
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
    images, err := loadImages()
    if err != nil {
        http.Error(w, "Unable to load images", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(images)
}

func main() {
    http.Handle("/photos/", http.StripPrefix("/photos/", http.FileServer(http.Dir("photos"))))
    http.HandleFunc("/", slideshowHandler)
    http.HandleFunc("/images", imagesHandler)

    log.Println("Local server established at http://localhost:3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}
