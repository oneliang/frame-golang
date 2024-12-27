package parallel

type Job struct {
	name               string
	jobConfiguration   *JobConfiguration
	sourceProcessorSet []SourceProcessor[any]
	firstJobStepList   []*JobStep
}

func NewJob(name string, jobConfiguration *JobConfiguration) *Job {
	if jobConfiguration == nil {
		jobConfiguration = DefaultJobConfiguration
	}
	return &Job{
		name:               name,
		jobConfiguration:   jobConfiguration,
		sourceProcessorSet: []SourceProcessor[any]{},
		firstJobStepList:   []*JobStep{},
	}
}

func (this *Job) AddSourceProcessor(sourceProcessor SourceProcessor[any]) {
	this.sourceProcessorSet = append(this.sourceProcessorSet, sourceProcessor)
}

func (this *Job) GenerateFirstJobStep() *JobStep {
	//parallel job step is not include the source processor
	firstJobStep := NewJobStep()
	this.firstJobStepList = append(this.firstJobStepList, firstJobStep)
	return firstJobStep
}
func (this *Job) Execute() {
	for _, sourceProcessor := range this.sourceProcessorSet {
		defaultSourceContext := NewDefaultSourceContext(sourceProcessor, this.firstJobStepList, this)
		sourceProcessor.Process(defaultSourceContext)
	}
	//for (parallelSourceProcessor in this.parallelSourceProcessorSet) {
	//	parallelSourceProcessor as ParallelSourceProcessor<Any?>
	//	val sourceCacheKey = parallelSourceProcessor.cacheKey
	//	if (sourceCacheKeySet.contains(parallelSourceProcessor.cacheKey)) {
	//		"duplicate cache key for source processor, source cache key:%s, source processor:%s".format(sourceCacheKey, parallelSourceProcessor).also {
	//			logError(it)
	//			error(it)
	//		}
	//	} else {
	//		sourceCacheKeySet += sourceCacheKey
	//	}
	//	val sourceData = this.cacheData?.getSourceData(sourceCacheKey)
	//	parallelSourceProcessor.initialize(sourceData)
	//	this.mainCoroutine.launch {
	//		val parallelJob = this as ParallelJob<Any?>
	//		val parallelSourceContext = DefaultParallelSourceContext(this.processCoroutine, parallelSourceProcessor, this.firstParallelJobStepList, parallelJob)
	//		parallelSourceProcessor.process(parallelSourceContext)
	//	}
	//}
}

func (this *Job) Finish() {

}

func CollectForProcessor(job *Job, jobStep *JobStep, value any, contextAction string, transformContext TransformContext[any]) {
	//logger.debug("parallelTransformContext:%s", parallelTransformContext)
	if jobStep.IsTransformProcessor() && transformContext != nil {
		if job.jobConfiguration.async {
			go func() {
				jobStep.transformProcessor.Process(value, transformContext)
			}()
		} else {
			jobStep.transformProcessor.Process(value, transformContext)
		}
	} else if jobStep.IsSinkProcessor() {
		for _, sinkProcessor := range jobStep.sinkProcessorList {
			if job.jobConfiguration.async {
				go func() {
					//logger.debug("sink processor, value:%s", value)
					sinkProcessor.Sink(value)
				}()
			} else {
				//logger.debug("sink processor, value:%s", value)
				sinkProcessor.Sink(value)
			}
			if (contextAction == CONTEXT_ACTION_SAVEPOINT || contextAction == CONTEXT_ACTION_FINISHED) && job.jobConfiguration.useCache {
				//val sinkKey = sinkProcessor.cacheKey
				//val sinkData = job.getSinkData(sinkKey) ?: CacheData.Data()
				//sinkProcessor.savepoint(sinkData)
				//job.updateSinkData(sinkKey, sinkData)
			}
		}
		//if contextAction == CONTEXT_ACTION_SAVEPOINT {
		//	job.saveCache()
		//} else if contextAction == CONTEXT_ACTION_FINISHED {
		//	job.saveCache()
		//	job.finish()
		//}
		if contextAction == CONTEXT_ACTION_FINISHED {
			job.Finish()
		}
	}
}
