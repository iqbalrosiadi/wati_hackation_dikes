package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"

// 	"github.com/sjwhitworth/golearn/base"
// 	"github.com/sjwhitworth/golearn/evaluation"
// 	"github.com/sjwhitworth/golearn/knn"
// )

// func main() {
// 	// Initiate user input reader
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		// Print the instruction to the reader in the console
// 		fmt.Println("Pleasse select the following option:")
// 		fmt.Println("A. Train the model")
// 		fmt.Println("B. Load the model")
// 		fmt.Println("C. Predict the model")
// 		fmt.Println("D. Exit")

// 		// Read the user's input
// 		char, _, err := reader.ReadRune()
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		switch char {
// 		case 'A':
// 			filePath, err := reader.ReadString('\n')
// 			if err != nil {
// 				panic(err)
// 			}
// 			TrainModel(filePath)
// 			break
// 		case 'B':
// 			fmt.Println("B Key Pressed")
// 			break
// 		case 'C':
// 			fmt.Println("C Key Pressed")
// 			break
// 		case 'D':
// 			os.Exit(0)
// 		}
// 	}

// 	// Load trained model
// 	classifier, _ := LoadModel2("knn_model.bin")

// 	// Load data need to predict
// 	rawData, err := base.ParseCSVToInstances("datasets/iris_headers.csv", true)
// 	if err != nil {
// 		panic(err)
// 	}
// 	_, testData := base.InstancesTrainTestSplit(rawData, 0.50)

// 	predictions, err := classifier.Predict(testData)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(predictions)

// 	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
// 	if err != nil {
// 		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
// 	}
// 	fmt.Println(evaluation.GetSummary(confusionMat))

// }

// func TrainModel(datasetPath string) {
// 	trainData, err := base.ParseCSVToInstances("datasetPath", true)
// 	if err != nil {
// 		panic(err)
// 	}

// 	//Do a training-test split
// 	classifier := knn.NewKnnClassifier("euclidean", "linear", 2)
// 	classifier.Fit(trainData)
// 	classifier.Save("knn_model.bin")
// }

// func LoadModel2(filename string) (*knn.KNNClassifier, error) {
// 	classifier := &knn.KNNClassifier{}
// 	classifier.Load("knn_model.bin")
// 	return classifier, nil
// }
