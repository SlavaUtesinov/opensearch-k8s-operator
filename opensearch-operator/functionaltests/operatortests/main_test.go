package operatortests

import (
	"testing"

	opsterv1 "github.com/Opster/opensearch-k8s-operator/opensearch-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/describe"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var k8sClient client.Client

var k8sClientSet *kubernetes.Clientset

var decriber describe.ResourceDescriber

func TestAPIs(t *testing.T) {
	config, err := clientcmd.BuildConfigFromFlags("", "../kubeconfig")
	if err != nil {
		panic(err.Error())
	}
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(opsterv1.AddToScheme(scheme))

	k8sClient, err = client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		panic(err.Error())
	}

	k8sClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var ok bool
	decriber, ok = describe.DescriberFor(schema.GroupKind{Group: appsv1.GroupName, Kind: "StatefulSet"}, config)
	if !ok {
		panic("can't create describer")
	}

	RegisterFailHandler(Fail)

	RunSpecs(t, "FunctionalTests")
}
