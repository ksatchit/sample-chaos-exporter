package main

import (
	"flag"
	"fmt"
	"log"
        "os"
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

/* Define Variables */

var kubeconfig *string
var chaosengine *string
var chaosexperimentlist []string
var chaosresultmap map[string]string
var totalExpCount, totalPassedExp, totalFailedExp int

/* Init function to obtain:
   (a) Cluster Config Path as Arg
   (b) ChaosEngine Name as Arg
*/

func init() {
        if home := homedir.HomeDir(); home != ""  {
                kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
        } else {
               kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
        }

        chaosengine = flag.String("chaosengine", "", "name of the chaosengine")
        flag.Parse()
}

/* Main function to derive chaos metrics for a given engine:
     (a) totalExpCount: Total Number of Experiments 
     (b) passedExpCount: Number of Passed Experiments
     (c) failedExpCount: Number of Failed Experiments
     (d) <exp>Verdict: Individual Experiment Status {0-notstarted, 1-running, 2-fail, 3-pass}
*/

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
		err.Error()
	}

        //TODO: Handle cases where the resources are not present. Currently, it panics 
        //TODO: Obtain/derive namespace of resources. Hardcoded to "default"
        //TODO: Update the chaosresult to carry verdict alone. Status & Verdict are redundant
        //TODO: Ensure there are chaosresult CRs applied w/ charts (verdict: notstarted) 

	//engines, err := clientSet.ChaosEngines("default").List(metav1.ListOptions{})
	engine, err := clientSet.ChaosEngines("default").Get(*chaosengine, metav1.GetOptions{})
	if err != nil {
		err.Error()
	}

        /////////////////////////////////////////////////
        /*METRIC*/
        totalExpCount = len(engine.Spec.Experiments)  //
        /////////////////////////////////////////////////

        for _, element:= range engine.Spec.Experiments{
            // fmt.Println(index, "=>", element.Name)
            chaosexperimentlist = append(chaosexperimentlist, element.Name)
        }

        // Initialize the chaosresult map before entering loop
        chaosresultmap := make(map[string]string)

        for _, test:= range chaosexperimentlist{
            chaosresultname := fmt.Sprintf("%s-%s", *chaosengine, test)
            testresultdump, err:= clientSet.ChaosResults("default").Get(chaosresultname, metav1.GetOptions{})
            result := testresultdump.Spec.ExperimentStatus.Verdict
            if err != nil {
                   err.Error()
            }

            chaosresultmap[chaosresultname] = result
        }

        pcount, fcount := 0, 0
        for _, verdict := range chaosresultmap{
            if verdict == "pass" {
                   pcount++
            } else if verdict == "fail" {
                   fcount++
            }
        }

        /////////////////////////////////////////////////
        /*METRIC*/                                     //
        totalPassedExp = pcount                        //
        totalFailedExp = fcount                        //
        /////////////////////////////////////////////////

        fmt.Printf("%d %d %d\n", totalExpCount, totalPassedExp, totalFailedExp)
}
