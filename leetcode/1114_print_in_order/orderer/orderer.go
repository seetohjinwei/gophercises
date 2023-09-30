package orderer

type Orderer struct {
	order      []int
	secondWait chan struct{}
	thirdWait  chan struct{}
	done       chan struct{}
}

func New() *Orderer {
	return &Orderer{
		order:      []int{},
		secondWait: make(chan struct{}),
		thirdWait:  make(chan struct{}),
		done:       make(chan struct{}),
	}
}

func (o *Orderer) first() {
	o.order = append(o.order, 1)
	o.secondWait <- struct{}{}
}

func (o *Orderer) second() {
	<-o.secondWait
	o.order = append(o.order, 2)
	o.thirdWait <- struct{}{}
}

func (o *Orderer) third() {
	<-o.thirdWait
	o.order = append(o.order, 3)
	o.done <- struct{}{}
}

func SameOrder(o *Orderer, orders [3]int) {
	mapping := []func(){o.first, o.second, o.third}

	for _, order := range orders {
		go mapping[order-1]()
	}

	<-o.done
}
