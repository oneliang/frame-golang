package parallel

type SourceProcessor[OUT any] interface {
	Process(sourceContext SourceContext[OUT])
}
