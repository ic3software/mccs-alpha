package dailyemail

import "sync"

type Pool struct {
	worker chan func()
	wg     sync.WaitGroup
}

func NewPool(maxGoroutines int) *Pool {
	p := Pool{
		worker: make(chan func()),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.worker {
				w()
			}
			p.wg.Done()
		}()
	}

	return &p
}

func (p *Pool) Run(w func()) {
	p.worker <- w
}

func (p *Pool) Shutdown() {
	close(p.worker)
	p.wg.Wait()
}
