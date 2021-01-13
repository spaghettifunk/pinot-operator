package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spaghettifunk/pinot-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("CronJob controller", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		PinotName      = "test-pinot"
		PinotNamespace = "test-pinot-namespace"
		ClusterName    = "test-pinot-cluster"

		timeout  = time.Second * 120
		duration = time.Second * 120
		interval = time.Millisecond * 250
	)

	Context("When updating PinotCluster Status", func() {
		It("Should create all the necessary components", func() {
			By("By creating a new Controller, Broker, Server and Zookeeper")
			ctx := context.Background()
			pinotCluster := &v1alpha1.Pinot{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apache.io/v1alpha1",
					Kind:       "Pinot",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      PinotName,
					Namespace: PinotNamespace,
				},
				Spec: v1alpha1.PinotSpec{
					ClusterName: ClusterName,
				},
			}

			// We'll need to retry getting this newly created Pinot, given that creation may not immediately happen.
			pinotLookupKey := types.NamespacedName{Name: PinotName, Namespace: PinotNamespace}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, pinotLookupKey, pinotCluster)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(pinotCluster.Spec.ClusterName).Should(Equal(ClusterName))
		})
	})
})
