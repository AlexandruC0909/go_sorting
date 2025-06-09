function sortingApp() {
  const algorithmData = {
    bubble: {
      name: "Bubble Sort",
      description: "A simple comparison sort algorithm. It repeatedly steps through the list, compares adjacent elements and swaps them if they are in the wrong order.",
      bestCase: "O(n)",
      averageCase: "O(n^2)",
      worstCase: "O(n^2)",
      space: "O(1)"
    },
    selection: {
      name: "Selection Sort",
      description: "An in-place comparison sort. It divides the input list into two parts: a sorted sublist of items which is built up from left to right and a sublist of the remaining unsorted items.",
      bestCase: "O(n^2)",
      averageCase: "O(n^2)",
      worstCase: "O(n^2)",
      space: "O(1)"
    },
    insertion: {
      name: "Insertion Sort",
      description: "A simple sorting algorithm that builds the final sorted array one item at a time. It is much less efficient on large lists than more advanced algorithms such as quicksort, heapsort, or merge sort.",
      bestCase: "O(n)",
      averageCase: "O(n^2)",
      worstCase: "O(n^2)",
      space: "O(1)"
    },
    quick: {
      name: "Quick Sort",
      description: "An efficient, divide-and-conquer sorting algorithm. It works by selecting a 'pivot' element from the array and partitioning the other elements into two sub-arrays, according to whether they are less than or greater than the pivot.",
      bestCase: "O(n log n)",
      averageCase: "O(n log n)",
      worstCase: "O(n^2)",
      space: "O(log n)"
    },
    merge: {
      name: "Merge Sort",
      description: "A divide-and-conquer algorithm that divides the array into halves, sorts them recursively, and then merges them.",
      bestCase: "O(n log n)",
      averageCase: "O(n log n)",
      worstCase: "O(n log n)",
      space: "O(n)"
    },
    heap: {
      name: "Heap Sort",
      description: "A comparison-based algorithm using a binary heap. It builds a max heap, then repeatedly extracts the max element.",
      bestCase: "O(n log n)",
      averageCase: "O(n log n)",
      worstCase: "O(n log n)",
      space: "O(1)"
    },
    shell: {
      name: "Shell Sort",
      description: "An improvement over insertion sort that allows comparison and exchange of items that are far apart, using decreasing gap intervals.",
      bestCase: "O(n log n)",
      averageCase: "O(n (log n)^2)",
      worstCase: "O(n^2)",
      space: "O(1)"
    },
    cocktail: {
      name: "Cocktail Shaker Sort",
      description: "A bidirectional bubble sort that sorts in both directions on each pass through the list.",
      bestCase: "O(n)",
      averageCase: "O(n^2)",
      worstCase: "O(n^2)",
      space: "O(1)"
    }
  };

  return {
    algorithm: "bubble",
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
    selectedAlgorithmDetails: null,

    init() {
      this.updateAlgorithmDetails(); 
      this.generateArray(); 
      this.watchAlgorithm(); 
    },

    watchAlgorithm() {
      this.$watch('algorithm', (newValue, oldValue) => {
        console.log('Algorithm changed from', oldValue, 'to', newValue);
        this.updateAlgorithmDetails();
        this.resetVisualization(); 
      });
    },

    updateAlgorithmDetails() {
      if (this.algorithm && algorithmData[this.algorithm]) {
        this.selectedAlgorithmDetails = algorithmData[this.algorithm];
      } else {
        this.selectedAlgorithmDetails = null; 
      }
    },

    async generateArray() {
      try {
        const size = parseInt(this.arraySize); 
        if (isNaN(size) || size < 5 || size > 50) { 
            console.error("Invalid array size:", this.arraySize);
            this.currentArray = []; 
            this.resetVisualization();
            return;
        }

        const newArray = [];
        for (let i = 0; i < size; i++) {
          newArray.push(Math.floor(Math.random() * 100) + 1);
        }
        this.currentArray = newArray;
        this.resetVisualization();
      } catch (error) {
        console.error("Error generating array locally:", error);
        this.currentArray = [];
        this.resetVisualization();
      }
    },

    pauseAnimation() {
      this.isPlaying = false;
      console.log("Animation paused");
    },

    async startOrResumeAnimation() {
      console.log("startOrResumeAnimation called");
      if (!this.currentArray || this.currentArray.length === 0) {
        console.log("Cannot start: No array generated.");
        return;
      }

      if (!this.sortResult || this.sortResult.steps.length === 0) {
        console.log("No sort result found, calling startSort...");
        await this.startSort();
        if (!this.sortResult || this.sortResult.steps.length === 0) {
          console.error("Failed to get sort steps.");
          this.isPlaying = false;
          return;
        }
      }

      console.log("Setting isPlaying to true and starting animation loop.");
      this.isPlaying = true;
      this.playAnimation();
    },

    async startSort() {
      this.updateAlgorithmDetails(); 
      this.isLoading = true;
      this.sortResult = null;
      this.currentStep = -1;
      this.comparingIndices = [];
      this.swappingIndices = [];
      this.sortedIndices = [];
      this.isPlaying = false;

      try {
        const response = await fetch("/sort", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            array: this.currentArray,
            algorithm: this.algorithm,
          }),
        });
        if (!response.ok) {
          throw new Error("HTTP error! status: " + response.status);
        }
        const result = await response.json();
        if (result && result.steps && result.steps.length > 0) {
          this.sortResult = result;
          console.log("Sort completed with", result.steps.length, "steps");
        } else {
          console.error("Invalid sort result received");
          this.sortResult = null;
        }
        this.isLoading = false;
      } catch (error) {
        console.error("Error sorting:", error);
        this.isLoading = false;
        this.sortResult = null;
      }
    },

    async playAnimation() {
      if (
        !this.sortResult ||
        !this.sortResult.steps ||
        this.sortResult.steps.length === 0
      ) {
        console.log("No sort result for playAnimation");
        this.isPlaying = false; 
        return;
      }

      if (this.currentStep >= this.sortResult.steps.length - 1) {
        this.currentStep = -1;
        this.resetVisualizationStateForAnimation(); 
      }

      let startFrom = this.currentStep > -1 ? this.currentStep : 0;
      if (startFrom === 0 && this.currentStep === -1) {
        this.resetVisualizationStateForAnimation();
      }

      for (let i = startFrom; i < this.sortResult.steps.length; i++) {
        if (!this.isPlaying) {
          console.log("Animation paused at step", i);
          this.currentStep = i; 
          return;
        }
        this.currentStep = i;
        const step = this.sortResult.steps[i];
        this.currentArray = [...step.array];
        this.comparingIndices = step.comparing || [];
        this.swappingIndices = step.swapping || [];
        this.sortedIndices = step.sorted || [];
        await this.sleep(1100 - this.speed * 100);
      }

      if (this.isPlaying) {
        this.isPlaying = false;
        this.currentStep = -1; 
        this.sortedIndices = this.currentArray.map((_, idx) => idx); 
      }
    },

    resetVisualizationStateForAnimation() {
      this.comparingIndices = [];
      this.swappingIndices = [];
      this.sortedIndices = [];
    },

    resetVisualization() {
      this.currentStep = -1;
      this.comparingIndices = [];
      this.swappingIndices = [];
      this.sortedIndices = [];
      this.sortResult = null; 
      this.isPlaying = false;
    },

    sleep(ms) {
      return new Promise((resolve) => setTimeout(resolve, ms));
    },
  };
}
