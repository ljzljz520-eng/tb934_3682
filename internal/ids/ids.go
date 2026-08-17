package ids

import (
	"fmt"
	"sync/atomic"
)

type Generator struct {
	prefix string
	seq    uint64
}

func New(prefix string) *Generator {
	return &Generator{prefix: prefix}
}

func (g *Generator) Next() string {
	n := atomic.AddUint64(&g.seq, 1)
	return fmt.Sprintf("%s-%04d", g.prefix, n)
}

func (g *Generator) NextFor(kind string) string {
	n := atomic.AddUint64(&g.seq, 1)
	return fmt.Sprintf("%s-%s-%04d", g.prefix, kind, n)
}
