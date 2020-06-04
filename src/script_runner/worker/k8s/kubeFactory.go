package worker

import (
	"bufio"
	"context"
	"fmt"
	"gitlab.com/grchive/grchive/core"
	"io"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	watchtools "k8s.io/client-go/tools/watch"
	"os"
	"strconv"
	"strings"
)

type KubeWorker struct {
	client         *kubernetes.Clientset
	spec           apiv1.PodSpec
	pod            *apiv1.Pod
	runUniqueLabel string
	labels         map[string]string

	log strings.Builder
}

func (w KubeWorker) wait(condition func(e watch.Event) (bool, error)) error {
	selector := labels.Set(w.labels).String()

	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return w.client.CoreV1().Pods("backend").List(metav1.ListOptions{
				LabelSelector: selector,
			})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return w.client.CoreV1().Pods("backend").Watch(metav1.ListOptions{
				LabelSelector: selector,
			})
		},
	}

	precondition := func(store cache.Store) (bool, error) {
		_, exists, err := store.Get(&metav1.ObjectMeta{
			Name: w.pod.ObjectMeta.Name,
		})
		if err != nil {
			return true, err
		}

		return exists, nil
	}

	_, err := watchtools.UntilWithSync(
		context.Background(),
		lw,
		&apiv1.Pod{},
		precondition,
		condition,
	)
	return err
}

func (w KubeWorker) waitReady() error {
	return w.wait(func(e watch.Event) (bool, error) {
		switch e.Type {
		case watch.Added, watch.Modified:
			pod, ok := e.Object.(*apiv1.Pod)
			if !ok || pod.ObjectMeta.Name != w.pod.ObjectMeta.Name {
				return false, nil
			}

			for _, status := range pod.Status.ContainerStatuses {
				if status.State.Running != nil || status.State.Terminated != nil {
					return true, nil
				}
			}
		}
		return false, nil
	})
}

func (w KubeWorker) waitTerminated() (int, error) {
	retCode := 0
	err := w.wait(func(e watch.Event) (bool, error) {
		switch e.Type {
		case watch.Added, watch.Modified:
			pod, ok := e.Object.(*apiv1.Pod)
			if !ok || pod.ObjectMeta.Name != w.pod.ObjectMeta.Name {
				return false, nil
			}

			for _, status := range pod.Status.ContainerStatuses {
				if status.State.Terminated != nil {
					retCode = int(status.State.Terminated.ExitCode)
					return true, nil
				}
			}
		}
		return false, nil
	})
	return retCode, err
}

func (w *KubeWorker) Run() (int, error) {
	w.labels = map[string]string{
		"grchive.script": w.runUniqueLabel,
	}
	w.log = strings.Builder{}

	var err error
	core.Info("Create Pods...")
	w.pod, err = w.client.CoreV1().Pods("backend").Create(&apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "script-runner-worker",
			Labels:       w.labels,
		},
		Spec: w.spec,
	})

	if err != nil {
		return -1, err
	}

	core.Info("Wait For Pod Ready...")
	err = w.waitReady()
	if err != nil {
		return -1, err
	}

	core.Info("Wait For Pod Finish...")
	return w.waitTerminated()
}

func (w KubeWorker) Logs() (string, error) {
	core.Info("Log Pods...")
	logReq := w.client.CoreV1().Pods("backend").GetLogs(w.pod.ObjectMeta.Name, &apiv1.PodLogOptions{})

	reqStream, err := logReq.Stream()
	if err != nil {
		return "", err
	}
	defer reqStream.Close()

	r := bufio.NewReader(reqStream)
	for {
		bytes, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		_, err = w.log.Write(bytes)
		if err != nil {
			break
		}
	}

	return w.log.String(), nil
}

func (w KubeWorker) Cleanup() {
	w.client.CoreV1().Pods("backend").Delete(w.pod.ObjectMeta.Name, nil)
}

type KubeFactory struct {
}

func (f KubeFactory) CreateWorker(opts WorkerOptions) (Worker, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	worker := KubeWorker{}
	worker.client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	worker.runUniqueLabel = fmt.Sprintf("script-runner-%d", opts.RunId)

	diskLimit, err := resource.ParseQuantity("1.0Gi")
	if err != nil {
		return nil, err
	}

	pullSecrets := []apiv1.LocalObjectReference{}
	if val, ok := os.LookupEnv("IMAGE_PULL_SECRET"); ok {
		pullSecrets = append(pullSecrets, apiv1.LocalObjectReference{
			Name: val,
		})
	}

	worker.spec = apiv1.PodSpec{
		RestartPolicy: apiv1.RestartPolicyNever,
		Volumes: []apiv1.Volume{
			apiv1.Volume{
				Name: "work-dir",
				VolumeSource: apiv1.VolumeSource{
					EmptyDir: &apiv1.EmptyDirVolumeSource{
						Medium:    apiv1.StorageMediumDefault,
						SizeLimit: &diskLimit,
					},
				},
			},
			apiv1.Volume{
				Name: "maven-dir",
				VolumeSource: apiv1.VolumeSource{
					EmptyDir: &apiv1.EmptyDirVolumeSource{
						Medium:    apiv1.StorageMediumDefault,
						SizeLimit: &diskLimit,
					},
				},
			},
		},
		Containers: []apiv1.Container{
			apiv1.Container{
				Name:            "script-runner-worker-container",
				Image:           core.EnvConfig.ScriptRunner.RunnerImage,
				ImagePullPolicy: apiv1.PullAlways,
				Args: []string{
					"-group",
					opts.ClientLibGroupId,
					"-artifact",
					opts.ClientLibArtifactId,
					"-vers",
					opts.ClientLibVersion,
					"-class",
					opts.RunClassName,
					"-fn",
					opts.RunFunctionName,
					"-meta",
					opts.RunMetadataName,
					"-runId",
					strconv.FormatInt(opts.RunId, 10),
				},
				VolumeMounts: []apiv1.VolumeMount{
					apiv1.VolumeMount{
						Name:      "work-dir",
						MountPath: "/data",
					},
					apiv1.VolumeMount{
						Name:      "maven-dir",
						MountPath: "/root/.m2/repository",
					},
				},
			},
		},
		ImagePullSecrets: pullSecrets,
	}

	return &worker, nil
}
