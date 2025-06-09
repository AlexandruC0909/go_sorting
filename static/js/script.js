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
      this.updateAlgorithmDetails(); // Initial call to set details for the default algorithm
      this.generateArray(); // Generates initial array
      this.watchAlgorithm(); // Add this line to set up the watcher
    },

    watchAlgorithm() {
      this.$watch('algorithm', (newValue, oldValue) => {
        console.log('Algorithm changed from', oldValue, 'to', newValue);
        this.updateAlgorithmDetails();
        this.resetVisualization(); // Reset visualization state
        // Optional: you might want to regenerate the array or clear it
        // this.generateArray(); // or this.currentArray = [];
        // For now, just updating details and resetting visualization is fine.
      });
    },

    updateAlgorithmDetails() {
      if (this.algorithm && algorithmData[this.algorithm]) {
        this.selectedAlgorithmDetails = algorithmData[this.algorithm];
      } else {
        this.selectedAlgorithmDetails = null; // Or some default state
      }
    },

    async generateArray() {
      // this.updateAlgorithmDetails(); // REMOVE this line, watcher handles it.
      try {
        const response = await fetch("/generate?size=" + this.arraySize);
        const data = await response.json();
        this.currentArray = data.array;
        this.resetVisualization();
      } catch (error) {
        console.error("Error generating array:", error);
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
      this.updateAlgorithmDetails(); // Add this line
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
