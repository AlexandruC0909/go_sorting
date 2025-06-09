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
	
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sorting Algorithm Visualizer</title>
    <script defer src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .controls {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
            flex-wrap: wrap;
            align-items: center;
        }
        .btn {
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            transition: background-color 0.2s;
        }
        .btn-primary {
            background: #007bff;
            color: white;
        }
        .btn-primary:hover {
            background: #0056b3;
        }
        .btn-secondary {
            background: #6c757d;
            color: white;
        }
        .btn-secondary:hover {
            background: #545b62;
        }
        .btn-success {
            background: #28a745;
            color: white;
        }
        .btn-success:hover {
            background: #1e7e34;
        }
        select, input {
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
        }
        .visualization {
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 20px;
            min-height: 300px;
            background: #fafafa;
            display: flex;
            align-items: end;
            justify-content: center;
            gap: 2px;
        }
        .bar {
            background: #007bff;
            transition: all 0.3s ease;
            border-radius: 2px 2px 0 0;
            min-width: 8px;
            position: relative;
        }
        .bar.comparing {
            background: #ffc107;
        }
        .bar.swapping {
            background: #dc3545;
        }
        .bar.sorted {
            background: #28a745;
        }
        .bar-value {
            position: absolute;
            top: -20px;
            left: 50%;
            transform: translateX(-50%);
            font-size: 10px;
            font-weight: bold;
        }
        .step-info {
            margin: 20px 0;
            padding: 10px;
            background: #e9ecef;
            border-radius: 4px;
            font-family: monospace;
        }
        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
        }
        .speed-control {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .array-display {
            margin: 10px 0;
            font-family: monospace;
            font-size: 14px;
            background: #f8f9fa;
            padding: 10px;
            border-radius: 4px;
            word-break: break-all;
        }
    </style>
</head>
<body>
    <div class="container" x-data="sortingApp()">
        <h1>Sorting Algorithm Visualizer</h1>
        
        <div class="controls">
            <select x-model="algorithm">
                <option value="bubble">Bubble Sort</option>
                <option value="selection">Selection Sort</option>
                <option value="insertion">Insertion Sort</option>
                <option value="quick">Quick Sort</option>
            </select>
            
            <div class="speed-control"> <!-- Reusing class for layout, can be renamed -->
                <label for="arraySizeSlider">Array Size:</label>
                <input type="range" id="arraySizeSlider" x-model="arraySize" min="5" max="50" step="1">
                <span x-text="arraySize" style="min-width: 25px; text-align: right;"></span>
            </div>
            
            <button class="btn btn-secondary" @click="generateArray()" :disabled="isLoading || isPlaying">Generate Array</button>
            <button class="btn" 
                    :class="{'btn-primary': !isPlaying, 'btn-warning': isPlaying}" 
                    @click="isPlaying ? pauseAnimation() : startOrResumeAnimation()"
                    :disabled="isLoading || (currentArray.length === 0 && !isPlaying)">
                <span x-text="isPlaying ? 'Pause' : 'Play'"></span>
            </button>
            
            <div class="speed-control">
                <label>Speed:</label>
                <input type="range" x-model="speed" min="1" max="10" step="1">
                <span x-text="speed"></span>
            </div>
        </div>

        <div class="array-display" x-show="currentArray.length > 0">
            Array: [<span x-text="currentArray.join(', ')"></span>]
        </div>

        <div class="visualization">
            <template x-for="(value, index) in currentArray" :key="index">
                <div class="bar" 
                     :style="'height: ' + (value * 4) + 'px; width: ' + Math.max(8, 800/currentArray.length) + 'px'"
                     :class="{
                         'comparing': comparingIndices.includes(index),
                         'swapping': swappingIndices.includes(index),
                         'sorted': sortedIndices.includes(index)
                     }">
                    <div class="bar-value" x-text="value"></div>
                </div>
            </template>
        </div>

        <div class="step-info" x-show="currentStep >= 0 && sortResult && sortResult.steps">
            <div>Algorithm: <span x-text="sortResult?.name || 'Unknown'"></span></div>
            <div>Step: <span x-text="currentStep + 1"></span> / <span x-text="sortResult?.steps?.length || 0"></span></div>
            <div x-show="comparingIndices.length > 0">Comparing indices: <span x-text="comparingIndices.join(', ')"></span></div>
            <div x-show="swappingIndices.length > 0">Swapping indices: <span x-text="swappingIndices.join(', ')"></span></div>
        </div>

        <div class="loading" x-show="isLoading">
            Generating sorting steps...
        </div>
    </div>

    <script>
        function sortingApp() {
            return {
                algorithm: 'bubble',
                arraySize: 20,
                currentArray: [],
                sortResult: null,
                currentStep: -1,
                isPlaying: false,
                isLoading: false,
                speed: 5,
                comparingIndices: [],
                swappingIndices: [],
                sortedIndices: [],

                init() {
                    this.generateArray();
                },

                async generateArray() {
                    try {
                        const response = await fetch('/generate?size=' + this.arraySize);
                        const data = await response.json();
                        this.currentArray = data.array;
                        this.resetVisualization();
                    } catch (error) {
                        console.error('Error generating array:', error);
                    }
                },

                pauseAnimation() {
                    this.isPlaying = false;
                    console.log('Animation paused');
                },

                async startOrResumeAnimation() {
                    console.log('startOrResumeAnimation called');
                    if (!this.currentArray || this.currentArray.length === 0) {
                        console.log('Cannot start: No array generated.');
                        return;
                    }

                    if (!this.sortResult || this.sortResult.steps.length === 0) {
                        console.log('No sort result found, calling startSort...');
                        await this.startSort(); 
                        if (!this.sortResult || this.sortResult.steps.length === 0) {
                            console.error('Failed to get sort steps.');
                            this.isPlaying = false; 
                            return;
                        }
                    }
                    
                    console.log('Setting isPlaying to true and starting animation loop.');
                    this.isPlaying = true;
                    this.playAnimation(); 
                },

                async startSort() {
                    this.isLoading = true;
                    // Reset state for a new sort operation
                    this.sortResult = null; 
                    this.currentStep = -1;
                    this.comparingIndices = [];
                    this.swappingIndices = [];
                    this.sortedIndices = [];
                    this.isPlaying = false; // Ensure animation is stopped

                    try {
                        const response = await fetch('/sort', {
                            method: 'POST',
                            headers: { 'Content-Type': 'application/json', },
                            body: JSON.stringify({ array: this.currentArray, algorithm: this.algorithm })
                        });
                        if (!response.ok) { throw new Error("HTTP error! status: " + response.status); }
                        const result = await response.json();
                        if (result && result.steps && result.steps.length > 0) {
                            this.sortResult = result;
                            console.log('Sort completed with', result.steps.length, 'steps');
                        } else {
                            console.error('Invalid sort result received');
                            this.sortResult = null;
                        }
                        this.isLoading = false;
                    } catch (error) {
                        console.error('Error sorting:', error);
                        this.isLoading = false;
                        this.sortResult = null;
                    }
                },

                async playAnimation() {
                    // This function will be more thoroughly refactored in the next step.
                    // For now, ensure it respects this.isPlaying and can be called by startOrResumeAnimation.
                    if (!this.sortResult || !this.sortResult.steps || this.sortResult.steps.length === 0) {
                        console.log('No sort result for playAnimation');
                        this.isPlaying = false; // Ensure consistency
                        return;
                    }

                    // If currentStep indicates a completed animation, reset to play again.
                    if (this.currentStep >= this.sortResult.steps.length - 1) {
                        this.currentStep = -1;
                        this.resetVisualizationStateForAnimation(); // Prepare for a new run
                    }
                    
                    let startFrom = this.currentStep > -1 ? this.currentStep : 0;
                    if (startFrom === 0 && this.currentStep === -1) { // if truly starting from the beginning
                         this.resetVisualizationStateForAnimation();
                    }

                    for (let i = startFrom; i < this.sortResult.steps.length; i++) {
                        if (!this.isPlaying) {
                            console.log('Animation paused at step', i);
                            this.currentStep = i; // Save progress
                            return; 
                        }
                        this.currentStep = i;
                        const step = this.sortResult.steps[i];
                        this.currentArray = [...step.array];
                        this.comparingIndices = step.comparing || [];
                        this.swappingIndices = step.swapping || [];
                        this.sortedIndices = step.sorted || [];
                        await this.sleep(1100 - (this.speed * 100));
                    }

                    if (this.isPlaying) { // Animation completed
                        this.isPlaying = false;
                        this.currentStep = -1; // Reset for next time
                        this.sortedIndices = this.currentArray.map((_, idx) => idx); // Mark all as sorted
                    }
                },

                resetVisualizationStateForAnimation() {
                    // Does not reset currentArray or sortResult
                    this.comparingIndices = [];
                    this.swappingIndices = [];
                    this.sortedIndices = [];
                    // currentStep reset is handled by the caller or within playAnimation logic
                },

                resetVisualization() { // Called by Generate Array or when new sort is initiated
                    this.currentStep = -1;
                    this.comparingIndices = [];
                    this.swappingIndices = [];
                    this.sortedIndices = [];
                    this.sortResult = null; // This is key, new sort needed
                    this.isPlaying = false; // Stop any ongoing animation
                },

                sleep(ms) {
                    return new Promise(resolve => setTimeout(resolve, ms));
                }
            }
        }
    </script>
</body>
</html>
`
	
	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	t.Execute(w, nil)
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
