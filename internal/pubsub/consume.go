package pubsub

type Acktype int

const (
	Ack Acktype = iota
	NackDiscard
	NackRequeue
)

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)
