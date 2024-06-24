package parallel

import "fmt"

const (
	CONTEXT_ACTION_NONE      = "NONE"
	CONTEXT_ACTION_FINISHED  = "FINISHED"
	CONTEXT_ACTION_SAVEPOINT = "SAVEPOINT"
)

type SourceContext[IN any] interface {
	Collect(value IN, contextAction string)
}

type DefaultSourceContext[IN any] struct {
	sourceProcessor     SourceProcessor[IN]
	jobStepList         []*JobStep
	job                 *Job
	transformContextMap map[string]TransformContext[any]
}

func NewDefaultSourceContext[IN any](sourceProcessor SourceProcessor[IN], jobStepList []*JobStep, job *Job) SourceContext[IN] {
	return &DefaultSourceContext[IN]{
		sourceProcessor:     sourceProcessor,
		jobStepList:         jobStepList,
		job:                 job,
		transformContextMap: make(map[string]TransformContext[any]),
	}
}
func (this *DefaultSourceContext[IN]) Collect(value IN, contextAction string) {
	if len(this.jobStepList) == 0 {
		//this.job.finish()
	}
	for _, jobStep := range this.jobStepList {
		var transformContext TransformContext[any] = nil
		if jobStep.IsTransformProcessor() {
			key := fmt.Sprintf("%p_%s", jobStep, contextAction)
			mapValue, ok := this.transformContextMap[key]
			if !ok {
				transformContext = NewDefaultTransformContext[any](this.job, jobStep, contextAction)
				this.transformContextMap[key] = transformContext
			} else {
				transformContext = mapValue
			}
		}

		CollectForProcessor(this.job, jobStep, value, contextAction, transformContext)
	}

	// for parallel source  processor
	if !this.job.jobConfiguration.async {
		if (contextAction == CONTEXT_ACTION_SAVEPOINT || contextAction == CONTEXT_ACTION_FINISHED) && this.job.jobConfiguration.useCache {
			//val sourceKey = this.sourceProcessor.cacheKey
			//val sourceData = this.parallelJob.getSourceData(sourceKey) ?: CacheData.Data()
			//this.parallelSourceProcessor.savepoint(sourceData)
			//this.parallelJob.updateSourceData(sourceKey, sourceData)
		}
	}
}
