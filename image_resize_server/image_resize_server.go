package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func ResizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resized := resize.Resize(480, 300, img, resize.NearestNeighbor)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "image/jpeg")
		jpeg.Encode(w, resized, nil)
	} else {
		IndexHandler(w, r)
	}
}

func ResizedSampleHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "test_resized.jpg")
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/resize", ResizeHandler)
	http.HandleFunc("/sample", ResizedSampleHandler)

	http.ListenAndServe(":3000", nil)
}
