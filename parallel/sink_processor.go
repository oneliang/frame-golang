package parallel

type SinkProcessor[IN any] interface {
	Sink(value IN)
}
