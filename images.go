package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
)

var root = flag.String("root", ".", "file system path")

func main() {
	// http.HandleFunc("/blue/", blueHandler)
	// http.HandleFunc("/red/", redHandler)
	// http.Handle("/", redHandler)
	// http.Handle("/", http.FileServer(http.Dir(*root)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		redHandler(w, r)
	})
	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func redHandler(w http.ResponseWriter, r *http.Request) {

	width := r.URL.Query().Get("width")
	if width == "" {
		width = "500"
	}
	width1, err := strconv.Atoi(width)
	height := r.URL.Query().Get("height")
	height1, err := strconv.Atoi(height)
	quality := r.URL.Query().Get("quality")
	mode := r.URL.Query().Get("mode")

	src, err := imaging.Open("/home/houzzcart/houzzcart/backend/finalImage" + r.URL.Path)
	// src = imaging.Sharpen(src, 2)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	if mode == "fill" {
		src = imaging.Fill(src, width1, height1, imaging.Center, imaging.Lanczos)
	}
	if mode == "fit" {
		src = imaging.Fit(src, width1, height1, imaging.Lanczos)
	}
	if mode == "" {
		src = imaging.Fill(src, width1, height1, imaging.Center, imaging.Lanczos)
	}
	if quality != "" {
		// quality = "80"
		src = imaging.Resize(src, width1, height1, imaging.Lanczos)
	}
	// log.Printf(param1)
	// log.Printf(height)
	// log.Printf(quality)

	// src = imaging.Blur(src, 5)
	// src = imaging.CropAnchor(src, 300, 300, imaging.Center)

	var img image.Image = src
	writeImageWithTemplate(w, &img)
}

// ImageTemplate mC h

var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body style='margin:0'><img src="data:image/jpg;base64,{{.Image}}"></body>`

// Writeimagewithtemplate encodes an image 'img' in jpeg format and writes it into ResponseWriter using a template.
func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}
