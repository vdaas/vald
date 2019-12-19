package tensorflow

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/vdaas/vald/internal/errors"
)

type TFModel struct {
	exportDir string
	tags      []string
	options   *tf.SessionOptions
	model     *tf.SavedModel
}

type Options func(*tf.SessionOptions)

type Feed struct {
	InputBytes    []byte
	OperationName string
	OutputIndex   int
}

type Fetch struct {
	OperationName string
	OutputIndex   int
}

func NewTFModel(exportDir string, tags []string, options ...Options) *TFModel {
	model := TFModel{exportDir, tags, nil, nil}
	for _, o := range options {
		o(model.options)
	}
	return &model
}

func WithOptions(options *tf.SessionOptions) Options {
	return func(ops *tf.SessionOptions) {
		options = ops
	}
}

func (tfModel *TFModel) LoadModel() error {
	model, err := tf.LoadSavedModel(tfModel.exportDir, tfModel.tags, tfModel.options)
	if err != nil {
		return err
	}
	tfModel.model = model
	return nil
}

func (tfModel *TFModel) Close() {
	tfModel.model.Session.Close()
}

func (tfModel TFModel) GetVector(feeds []Feed, fetches []Fetch, targets []*tf.Operation) ([][][]float64, error) {
	getFeeds := func(feeds []Feed) (map[tf.Output]*tf.Tensor, error) {
		result := map[tf.Output]*tf.Tensor{}

		for _, feed := range feeds {
			inputTensor, err := tf.NewTensor([]string{string(feed.InputBytes)})
			if err != nil {
				return nil, err
			}
			result[tfModel.model.Graph.Operation(feed.OperationName).Output(feed.OutputIndex)] = inputTensor
		}

		return result, nil
	}

	input, err := getFeeds(feeds)
	if err != nil {
		return nil, err
	}

	output := []tf.Output{}
	for _, fetch := range fetches {
		output = append(output, tfModel.model.Graph.Operation(fetch.OperationName).Output(fetch.OutputIndex))
	}

	result, err := tfModel.model.Session.Run(input, output, targets)
	if err != nil {
		return nil, err
	}

	values := [][][]float64{}
	for i := range result {
		value, ok := result[i].Value().([][]float64)
		if ok {
			values = append(values, value)
		} else {
			return nil, errors.New("cast error")
		}
	}
	return values, nil
}
