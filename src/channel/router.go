package channel

type Router struct {
	channels map[string]Channel
}

func NewRouter(chs map[string]Channel) *Router {
	return &Router{channels: chs}
}

func (r *Router) Route(name string) Channel {
	return r.channels[name]
}
