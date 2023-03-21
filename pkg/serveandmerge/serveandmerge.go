package serveandmerge

import (
	"log"

	"github.com/oklog/run"
)

type Runner interface {
	Start() error
}

type ServeAndMerge struct {
}

func NewServeAndMerge() *ServeAndMerge {
	return &ServeAndMerge{}
}

func (snm *ServeAndMerge) Start(runners []Runner) error {

	var g run.Group

	// run runners
	{
		g.Add(func() error {
			return runners[0].Start()
		}, func(err error) {
			if err != nil {
				log.Fatal(err.Error())
			}
		})
	}

	{
		g.Add(func() error {
			return runners[1].Start()
		}, func(err error) {
			if err != nil {
				log.Fatal(err.Error())
			}
		})
	}

	if err := g.Run(); err != nil {
		return err
	}

	return nil
}
