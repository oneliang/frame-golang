package parallel

import "fmt"

type TransformContext[IN any] interface {
	Collect(value IN)
}

type DefaultTransformContext[IN any] struct {
	job                 *Job
	jobStep             *JobStep
	parentContextAction string
	transformContextMap map[string]TransformContext[any]
}

func NewDefaultTransformContext[IN any](job *Job, jobStep *JobStep, parentContextAction string) TransformContext[IN] {
	return &DefaultTransformContext[IN]{
		job:                 job,
		jobStep:             jobStep,
		parentContextAction: parentContextAction,
		transformContextMap: make(map[string]TransformContext[any]),
	}
}

func (this *DefaultTransformContext[IN]) Collect(value IN) {
	//logger.debug("transform processor, value:%s", value)
	var nextJobStep *JobStep = nil
	if this.jobStep.HasNextJobStep() {
		nextJobStep = this.jobStep.nextJobStep
	}
	if nextJobStep == nil {
		//logger.error(
		//	"this parallel job step is used for a transform processor, but next parallel job step is null, you may be need to add a transform processor or a sink processor for next job, this:%s",
		//	this.parallelJobStep
		//)
	} else {
		var fixNextJobStep = nextJobStep
		var transformContext TransformContext[any] = nil
		if fixNextJobStep.IsTransformProcessor() {
			key := fmt.Sprintf("%p_%s", fixNextJobStep, this.parentContextAction)
			mapValue, ok := this.transformContextMap[key]
			if !ok {
				transformContext = NewDefaultTransformContext[any](this.job, fixNextJobStep, this.parentContextAction)
				this.transformContextMap[key] = transformContext
			} else {
				transformContext = mapValue
			}
		}
		CollectForProcessor(this.job, fixNextJobStep, value, this.parentContextAction, transformContext)
	}
}
