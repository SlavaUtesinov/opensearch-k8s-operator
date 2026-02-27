package helmtests

import (
	"context"
	"time"

	opsterv1 "github.com/Opster/opensearch-k8s-operator/opensearch-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The cluster has been created using Helm outside of this test. This test verifies the presence of the resources after the cluster is created.
var _ = Describe("DeployWithHelm", Ordered, func() {
	name := "opensearch-cluster"
	namespace := "default"

	When("cluster is created using helm", Ordered, func() {
		It("should have 3 ready master pods", func() {
			sts := appsv1.StatefulSet{}
			Eventually(func() int32 {
				err := k8sClient.Get(context.Background(), client.ObjectKey{Name: name + "-masters", Namespace: namespace}, &sts)
				if err == nil {
					GinkgoWriter.Printf("%+v\n", sts.Status)
					pods := &corev1.PodList{}
					err := k8sClient.List(context.Background(), pods, client.InNamespace(namespace))
					if err == nil {
						for _, pod := range pods.Items {
							revision, ok := pod.Labels["controller-revision-hash"]
							GinkgoWriter.Printf("Pod: %s\tPhase: %s", pod.Name, pod.Status.Phase)
							if ok {
								GinkgoWriter.Printf("\tRevision: %s\t Image: %s", revision, pod.Spec.Containers[0].Image)
							}
							GinkgoWriter.Println()
						}
					} else {
						GinkgoWriter.Println(err)
					}
					cluster := &opsterv1.OpenSearchCluster{}
					if err := k8sClient.Get(context.Background(), client.ObjectKey{Name: name, Namespace: namespace}, cluster); err != nil {
						GinkgoWriter.Printf("err: %s\n", err.Error())
					} else {
						GinkgoWriter.Printf("Cluster: %+v\n", cluster.Status)
					}

					return sts.Status.ReadyReplicas
				}
				GinkgoWriter.Println(err)
				return 0
			}, time.Minute*15, time.Second*5).Should(Equal(int32(3)))
		})

		It("should have a ready dashboards pod", func() {
			deployment := appsv1.Deployment{}
			Eventually(func() int32 {
				err := k8sClient.Get(context.Background(), client.ObjectKey{Name: name + "-dashboards", Namespace: namespace}, &deployment)
				if err == nil {
					return deployment.Status.ReadyReplicas
				}
				return 0
			}, time.Minute*5, time.Second*5).Should(Equal(int32(1)))
		})
	})
})
