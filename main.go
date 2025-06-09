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
)

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
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/sort", sortHandler)
	http.HandleFunc("/generate", generateHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
			// Show comparison
			steps = append(steps, SortStep{
				Array:     copyArray(array),
				Comparing: []int{j, j + 1},
			})
			
			if array[j] > array[j+1] {
				// Show swap
				steps = append(steps, SortStep{
					Array:    copyArray(array),
					Swapping: []int{j, j + 1},
				})
				
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
		
		// Mark as sorted
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
