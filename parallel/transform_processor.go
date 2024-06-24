package parallel

type TransformProcessor[IN any, OUT any] interface {
	Process(value IN, transformContext TransformContext[OUT])
}
