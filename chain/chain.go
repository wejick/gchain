package chain

import (
	"context"

	"github.com/wejick/gochain/model"
)

//go:generate mockery --name BaseChain
type BaseChain interface {
	// Run does prediction of input prompt of <string,string> and produce output of <string,string>
	// map of <string,string> of output and input prompt to accomodate many possible usecases
	Run(ctx context.Context, prompt map[string]string, options ...func(*model.Option)) (output map[string]string, err error)

	// SimpleRun does prediction of prompt of string and produce output of string
	// this is to accomodate simple prompt / output usage
	SimpleRun(ctx context.Context, prompt string, options ...func(*model.Option)) (output string, err error)
}
