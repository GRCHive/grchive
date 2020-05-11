package worker

type Worker interface {
	Run() (int, error)
	Logs() (string, error)
	Cleanup()
}

type WorkerOptions struct {
	ClientLibGroupId    string
	ClientLibArtifactId string
	ClientLibVersion    string
	RunClassName        string
	RunFunctionName     string
	RunMetadataName     string
	RunId               int64
}

type WorkerFactory interface {
	CreateWorker(opts WorkerOptions) (Worker, error)
}
