// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"reflect"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"github.com/vdaas/vald/internal/errors"
)

type LLM interface {
	Embed(ctx context.Context, doc string) ([]float32, error)
}

type OpenAI interface {
	LLM
}

type openAI struct {
	model   openai.EmbeddingModel
	token   string
	baseURL string
	client  *openai.Client
}

func NewOpenAI(opts ...OpenAIOption) (OpenAI, error) {
	o := &openAI{}
	for _, opt := range append(defaultOpenAIOptions, opts...) {
		if err := opt(o); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if o.token == "" {
		return nil, errors.NewErrInvalidOption("token", o.token, errors.New("token must not be empty"))
	}
	cfg := openai.DefaultConfig(o.token)
	if o.baseURL != "" {
		cfg.BaseURL = o.baseURL
	}
	o.client = openai.NewClientWithConfig(cfg)
	return o, nil
}

func (o *openAI) Embed(ctx context.Context, doc string) ([]float32, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}
	embeddings, err := o.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: o.model,
		Input: doc,
	})
	if err != nil {
		return nil, err
	}
	for _, embedding := range embeddings.Data {
		if embedding.Embedding != nil {
			return embedding.Embedding, nil
		}
	}
	return nil, errors.New("openai: no embedding returned for input")
}
