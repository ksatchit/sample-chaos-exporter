package main

import (
	"flag"
	"fmt"
	"log"
        "os"
        //"strings"
        "path/filepath"
        "k8s.io/client-go/util/homedir"
        "k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
        _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        v1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis"
        clientV1alpha1 "github.com/ksatchit/sample-chaos-exporter/pkg/clientset/v1alpha1"
)

var kubeconfig *string
var chaosengine *string
var chaosexperimentlist []string
var chaosresultmap map[string]string

func init() {
        if home := homedir.HomeDir(); home != ""  {
                kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
        } else {
               kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
        }

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

        //TODO: Handle cases where the resources are not present
        //TODO: Obtain/derive namespace of resources. Hardcoded to "default"

	//engines, err := clientSet.ChaosEngines("default").List(metav1.ListOptions{})
	engine, err := clientSet.ChaosEngines("default").Get(*chaosengine, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

        //result, err:= clientSet.ChaosResults("default").List(metav1.ListOptions{})
        //if err != nil {
        //        panic(err)
        //}

        /*METRIC*/
	fmt.Printf("No of chaos exp: %+v\n", len(engine.Spec.Experiments))
        fmt.Printf("------------------------\n")

        for _, element:= range engine.Spec.Experiments{
            // fmt.Println(index, "=>", element.Name)
            chaosexperimentlist = append(chaosexperimentlist, element.Name)
        }
        fmt.Printf("List of chaos experiments: %+v\n", chaosexperimentlist)
        fmt.Printf("------------------------\n")

        for _, test:= range chaosexperimentlist{
            chaosresultname := fmt.Sprintf("%s-%s", *chaosengine, test)
            testresult, err:= clientSet.ChaosResults("default").Get(chaosresultname, metav1.GetOptions{})
            if err != nil {
                   //err.Error() 
                   //testresult = "0" //result not available, mapped to "not-started"
                   //fmt.Println(err.Error())
                   fmt.Printf("experiment %s has not started", test)
            }

            fmt.Printf("result for %+v is %+v\n", chaosresultname, testresult.Spec.ExperimentStatus.Verdict)
            chaosresultmap["chaosresultname"] = testresult.Spec.ExperimentStatus.Verdict
        }

        //fmt.Printf("List of chaos results: %v\n", result)



        /*
	store := WatchResources(clientSet)

	for {
		projectsFromStore := store.List()
		fmt.Printf("project in store: %d\n", len(projectsFromStore))

		time.Sleep(2 * time.Second)
	}
        */
}
