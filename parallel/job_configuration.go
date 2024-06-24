package parallel

var (
	DefaultJobConfiguration = NewJobConfiguration(false, false)
)

type JobConfiguration struct {
	async    bool
	useCache bool
}

func NewJobConfiguration(async bool, useCache bool) *JobConfiguration {
	return &JobConfiguration{
		async:    async,
		useCache: useCache,
	}
}
