package core

// Pipeline 负责按序执行 Plugin 链
type Pipeline struct {
	plugins []Plugin
}

func NewPipeline(ps ...Plugin) *Pipeline {
	return &Pipeline{plugins: ps}
}

func (p *Pipeline) Execute(r *Rocket) (Result, error) {
	for _, pl := range p.plugins {
		if err := pl.Web(r); err != nil {
			return Result{}, err
		}
	}
	return Result{}, nil
}
