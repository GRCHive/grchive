package worker

type KubeWorker struct {
}

type KubeFactory struct {
}

func (f KubeFactory) CreateWorker(opts WorkerOptions) (Worker, error) {
	return nil, nil
}
