package knn

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/knn"
)

type Trainer struct {
	classifier *knn.KNNClassifier
}

func NewTrainer() *Trainer {
	return &Trainer{
		classifier: knn.NewKnnClassifier("euclidean", "linear", 2),
	}
}

func (t *Trainer) Train() {
	var err error
	// Load in a dataset, with headers. Header attributes will be stored.
	// Think of instances as a Data Frame structure in R or Pandas.
	// You can also create instances from scratch.
	rawData, err := base.ParseCSVToInstances("datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	// Print a pleasant summary of your data.
	fmt.Println(rawData)

	//Do a training-test split
	t.classifier.Fit(rawData)
}

func (t *Trainer) Predict(what base.FixedDataGrid) {
	// Calculates the Euclidean distance and returns the most popular label
	predictions, err := t.classifier.Predict(what)
	if err != nil {
		panic(err)
	}

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(what, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))
}
