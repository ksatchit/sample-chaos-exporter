package main

import (
	"flag"
	"fmt"
	"log"
        "os"
	// "time"
        "path/filepath"
        "k8s.io/client-go/util/homedir"
        "k8s.io/client-go/tools/clientcmd"
        _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	//"github.com/martin-helmich/kubernetes-crd-example/api/types/v1alpha1"
	//clientV1alpha1 "github.com/martin-helmich/kubernetes-crd-example/clientset/v1alpha1"

        v1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis"
        clientV1alpha1 "github.com/ksatchit/sample-chaos-exporter/pkg/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

var kubeconfig *string
var chaosengine *string

func init() {
        if home := homedir.HomeDir(); home != ""  {
                kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
        } else {
               kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
        }
        // flag.Parse()

        chaosengine = flag.String("chaosengine", "", "name of the chaosengine")
        flag.Parse()
}

func main() {

        var config *rest.Config
        var err error

        if *kubeconfig == "" {
               log.Printf("using the in-cluster config")
               config, err = rest.InClusterConfig()
        } else {
               log.Printf("using configuration from '%s'", *kubeconfig)
               config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
        }

        if err != nil {
                panic(err.Error())
        }

        if *chaosengine == "" {
               log.Printf("ERROR: please specify name of chaosengine, exiting")
               os.Exit(1)
        }

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//engines, err := clientSet.ChaosEngines("default").List(metav1.ListOptions{})
	engine, err := clientSet.ChaosEngines("default").Get(*chaosengine, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

        
	fmt.Printf("No of chaos exp: %+v\n", len(engine.Spec.Experiments))
        fmt.Printf("No of passed chaos exp: %+v\n", engine.Status)


        /*
	store := WatchResources(clientSet)

	for {
		projectsFromStore := store.List()
		fmt.Printf("project in store: %d\n", len(projectsFromStore))

		time.Sleep(2 * time.Second)
	}
        */
}    
