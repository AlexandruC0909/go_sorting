package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
    "path/filepath"
    "strings"
    "mime"
    "embed"
    "github.com/tdewolff/minify"
    "github.com/tdewolff/minify/css"
    "github.com/tdewolff/minify/js"
)
//go:embed static/**/*
var staticFiles embed.FS

//go:embed templates/*.html
var htmlFiles embed.FS
type SortStep struct {
	Array     []int `json:"array"`
	Comparing []int `json:"comparing"`
	Swapping  []int `json:"swapping"`
	Sorted    []int `json:"sorted"`
}

type SortResult struct {
	Steps []SortStep `json:"steps"`
	Name  string     `json:"name"`
}

func main() {
    mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/sort", sortHandler)
	mux.HandleFunc("/generate", generateHandler)


    mux.Handle("/static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        filePath := r.URL.Path

		isJS := strings.HasSuffix(filePath, ".js")
		isCSS := strings.HasSuffix(filePath, ".css")

        data, err := staticFiles.ReadFile("static/" + filePath)
        if err != nil {
            http.NotFound(w, r)
            return
        }

        if isJS || isCSS {
            mediaType := "application/javascript"
			if isCSS {
				mediaType = "text/css"
			}

            minifiedData, err := minifyContent(data, mediaType)
            if err != nil {
                http.Error(w, "Failed to minify CSS", http.StatusInternalServerError)
                return
            }
            data = minifiedData
        }

        ext := strings.ToLower(filepath.Ext(filePath))
        mimeType := mime.TypeByExtension(ext)
        if mimeType != "" {
            w.Header().Set("Content-Type", mimeType)
        } else {
            w.Header().Set("Content-Type", "application/octet-stream")
        }

        w.Write(data)
    })))
	
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	sizeStr := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 5 || size > 50 {
		size = 20
	}
	
	rand.Seed(time.Now().UnixNano())
	array := make([]int, size)
	for i := range array {
		array[i] = rand.Intn(100) + 1
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]int{"array": array})
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Array     []int  `json:"array"`
		Algorithm string `json:"algorithm"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	var result SortResult
	
	switch request.Algorithm {
	case "bubble":
		result = bubbleSort(request.Array)
	case "selection":
		result = selectionSort(request.Array)
	case "insertion":
		result = insertionSort(request.Array)
	case "quick":
		result = quickSort(request.Array)
	default:
		result = bubbleSort(request.Array)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func bubbleSort(arr []int) SortResult {
	steps := []SortStep{}
	array := make([]int, len(arr))
	copy(array, arr)
	n := len(array)
	
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			steps = append(steps, SortStep{
				Array:     copyArray(array),
				Comparing: []int{j, j + 1},
			})
			
			if array[j] > array[j+1] {
				steps = append(steps, SortStep{
					Array:    copyArray(array),
					Swapping: []int{j, j + 1},
				})
				
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
		
		sorted := make([]int, i+1)
		for k := 0; k <= i; k++ {
			sorted[k] = n - 1 - k
		}
		steps = append(steps, SortStep{
			Array:  copyArray(array),
			Sorted: sorted,
		})
	}
	
	return SortResult{Steps: steps, Name: "Bubble Sort"}
}

func selectionSort(arr []int) SortResult {
	steps := []SortStep{}
	array := make([]int, len(arr))
	copy(array, arr)
	n := len(array)
	
	for i := 0; i < n-1; i++ {
		minIdx := i
		
		for j := i + 1; j < n; j++ {
			steps = append(steps, SortStep{
				Array:     copyArray(array),
				Comparing: []int{minIdx, j},
			})
			
			if array[j] < array[minIdx] {
				minIdx = j
			}
		}
		
		if minIdx != i {
			steps = append(steps, SortStep{
				Array:    copyArray(array),
				Swapping: []int{i, minIdx},
			})
			
			array[i], array[minIdx] = array[minIdx], array[i]
		}
		
		sorted := make([]int, i+1)
		for k := 0; k <= i; k++ {
			sorted[k] = k
		}
		steps = append(steps, SortStep{
			Array:  copyArray(array),
			Sorted: sorted,
		})
	}
	
	return SortResult{Steps: steps, Name: "Selection Sort"}
}

func insertionSort(arr []int) SortResult {
	steps := []SortStep{}
	array := make([]int, len(arr))
	copy(array, arr)
	n := len(array)
	
	for i := 1; i < n; i++ {
		key := array[i]
		j := i - 1
		
		steps = append(steps, SortStep{
			Array:     copyArray(array),
			Comparing: []int{i},
		})
		
		for j >= 0 && array[j] > key {
			steps = append(steps, SortStep{
				Array:     copyArray(array),
				Comparing: []int{j, j + 1},
			})
			
			array[j+1] = array[j]
			j--
		}
		
		array[j+1] = key
		
		sorted := make([]int, i+1)
		for k := 0; k <= i; k++ {
			sorted[k] = k
		}
		steps = append(steps, SortStep{
			Array:  copyArray(array),
			Sorted: sorted,
		})
	}
	
	return SortResult{Steps: steps, Name: "Insertion Sort"}
}

func quickSort(arr []int) SortResult {
	steps := []SortStep{}
	array := make([]int, len(arr))
	copy(array, arr)
	
	quickSortHelper(array, 0, len(array)-1, &steps)
	
	return SortResult{Steps: steps, Name: "Quick Sort"}
}

func quickSortHelper(array []int, low, high int, steps *[]SortStep) {
	if low < high {
		pi := partition(array, low, high, steps)
		quickSortHelper(array, low, pi-1, steps)
		quickSortHelper(array, pi+1, high, steps)
	}
}

func partition(array []int, low, high int, steps *[]SortStep) int {
	pivot := array[high]
	i := low - 1
	
	for j := low; j < high; j++ {
		*steps = append(*steps, SortStep{
			Array:     copyArray(array),
			Comparing: []int{j, high},
		})
		
		if array[j] < pivot {
			i++
			if i != j {
				*steps = append(*steps, SortStep{
					Array:    copyArray(array),
					Swapping: []int{i, j},
				})
				array[i], array[j] = array[j], array[i]
			}
		}
	}
	
	*steps = append(*steps, SortStep{
		Array:    copyArray(array),
		Swapping: []int{i + 1, high},
	})
	
	array[i+1], array[high] = array[high], array[i+1]
	return i + 1
}

func copyArray(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	return result
}

func minifyContent(content []byte, mediaType string) ([]byte, error) {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)

	return m.Bytes(mediaType, content)
}

