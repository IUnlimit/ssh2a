package cache

type Pair struct {
	k, v interface{}
}

func NewPair(k, v interface{}) *Pair {
	return &Pair{
		k,
		v,
	}
}

func (p *Pair) getK() interface{} {
	return p.k
}

func (p *Pair) getV() interface{} {
	return p.v
}
