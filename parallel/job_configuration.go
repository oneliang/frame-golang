package parallel

var (
	DefaultJobConfiguration = NewJobConfiguration(false, false, 4)
)

type JobConfiguration struct {
	async    bool
	useCache bool
	poolSize uint
}

func NewJobConfiguration(async bool, useCache bool, poolSize uint) *JobConfiguration {
	return &JobConfiguration{
		async:    async,
		useCache: useCache,
		poolSize: poolSize,
	}
}
