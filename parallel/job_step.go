package parallel

import (
	"log"
)

type JobStep struct {
	nextJobStep        *JobStep
	transformProcessor TransformProcessor[any, any]
	sinkProcessorList  []SinkProcessor[any]
}

func NewJobStep() *JobStep {
	return &JobStep{
		nextJobStep:        nil,
		transformProcessor: nil,
		sinkProcessorList:  []SinkProcessor[any]{},
	}
}

func (this *JobStep) AddTransformProcessor(transformProcessor TransformProcessor[any, any]) *JobStep {
	if this.transformProcessor != nil {
		log.Fatalf("parallel transform processor has been initialized, only can initialize one time")
	}
	this.transformProcessor = transformProcessor
	jobStep := NewJobStep()
	this.nextJobStep = jobStep
	return jobStep
}

func (this *JobStep) AddSinkProcessor(sinkProcessor SinkProcessor[any]) {
	if this.sinkProcessorList == nil {
		this.sinkProcessorList = []SinkProcessor[any]{}
	}
	this.sinkProcessorList = append(this.sinkProcessorList, sinkProcessor)
}

func (this *JobStep) IsTransformProcessor() bool {
	return len(this.sinkProcessorList) == 0 && this.nextJobStep != nil && this.transformProcessor != nil
}

func (this *JobStep) IsSinkProcessor() bool {
	return len(this.sinkProcessorList) > 0 && this.nextJobStep == nil && this.transformProcessor == nil
}

func (this *JobStep) HasNextJobStep() bool {
	return this.nextJobStep != nil
}
