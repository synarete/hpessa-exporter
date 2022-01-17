// SPDX-License-Identifier: Apache-2.0
package devmon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type client struct {
	ClientSet *kubernetes.Clientset
	Config    *rest.Config
}

func newClient() (*client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return newExternalClient()
	}
	cset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &client{}, err
	}

	return &client{
		ClientSet: cset,
		Config:    config,
	}, nil
}

func newExternalClient() (*client, error) {
	config, err := buildOutOfClusterConfig()
	if err != nil {
		return &client{}, err
	}
	cset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &client{}, err
	}

	return &client{
		ClientSet: cset,
		Config:    config,
	}, nil
}

func buildOutOfClusterConfig() (*rest.Config, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube/config")
	}

	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func GetRunningPod(ctx context.Context, clnt *client,
	nname types.NamespacedName) (*corev1.Pod, error) {
	pod, err := clnt.ClientSet.CoreV1().
		Pods(nname.Namespace).Get(ctx, nname.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if pod.Status.Phase != corev1.PodRunning {
		return pod, fmt.Errorf("pod %+v no running", nname)
	}
	return pod, nil
}

func ListAllRunningPods(ctx context.Context, clnt *client) ([]*corev1.Pod, error) {
	ret := []*corev1.Pod{}
	pods, err := clnt.ClientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return ret, err
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			dup := corev1.Pod{}
			pod.DeepCopyInto(&dup)
			ret = append(ret, &dup)
		}
	}
	return ret, nil
}

func ListRunningPodsWith(ctx context.Context, clnt *client,
	namePrefix, hostIP string) ([]*corev1.Pod, error) {
	ret := []*corev1.Pod{}
	pods, err := ListAllRunningPods(context.TODO(), clnt)
	if err != nil {
		return ret, err
	}
	for _, pod := range pods {
		if strings.HasPrefix(pod.Name, namePrefix) && pod.Status.HostIP == hostIP {
			ret = append(ret, pod)
		}
	}
	return ret, nil
}
