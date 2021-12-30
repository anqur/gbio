package registries

import (
	"math/rand"
	"time"
)

type LB interface {
	Pick(endpoints []string) string
}

type randLB struct {
	r *rand.Rand
}

func RandLB() LB {
	return &randLB{r: rand.New(rand.NewSource(time.Now().Unix()))}
}

func (r *randLB) Pick(endpoints []string) string {
	return endpoints[r.r.Intn(len(endpoints))]
}

type firstLB struct{}

func FirstLB() LB                               { return new(firstLB) }
func (*firstLB) Pick(endpoints []string) string { return endpoints[0] }
