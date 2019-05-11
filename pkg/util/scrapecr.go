package util 

import (
	"fmt"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
        _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        v1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis"
        clientV1alpha1 "github.com/ksatchit/sample-chaos-exporter/pkg/clientset/v1alpha1"
)

/* Function to derive chaos metrics for a given engine:
     (a) totalExpCount: Total Number of Experiments 
     (b) passedExpCount: Number of Passed Experiments
     (c) failedExpCount: Number of Failed Experiments
     (d) <exp>Verdict: Individual Experiment Status {0-notstarted, 1-running, 2-fail, 3-pass}
*/

// Define Variables 
var chaosexperimentlist []string
var chaosresultmap map[string]string
var statusmap map[string]float64

var numericstatus = map[string]float64{
         "not-executed": 0,
         "running":      1,
         "fail":         2,
         "pass":         3,
}

func statusConv (expstatus string)(numeric float64){
         if numeric, ok := numericstatus[expstatus]; ok {
            return numeric
            fmt.Printf("%v", numeric)
         }
         return 127
}

func GetChaosMetrics(cfg *rest.Config, cEngine string)(totalExpCount, totalPassedExp, totalFailedExp float64, rMap map[string]float64, err error){

	v1alpha1.AddToScheme(scheme.Scheme)

	clientSet, err := clientV1alpha1.NewForConfig(cfg)
	if err != nil {
		//err.Error()
                return 0, 0, 0, nil, err
	}

        //TODO: Handle cases where the resources are not present. Currently, it panics 
        //TODO: Obtain/derive namespace of resources. Hardcoded to "default"
        //TODO: Update the chaosresult to carry verdict alone. Status & Verdict are redundant
        //TODO: Ensure there are chaosresult CRs applied w/ charts (verdict: notstarted) 

	//engines, err := clientSet.ChaosEngines("default").List(metav1.ListOptions{})
	engine, err := clientSet.ChaosEngines("default").Get(cEngine, metav1.GetOptions{})
	if err != nil {
		//err.Error()
                return 0, 0, 0, nil, err
	}

        /////////////////////////////////////////////////////////
        /*METRIC*/
        totalExpCount = float64(len(engine.Spec.Experiments))  //
        /////////////////////////////////////////////////////////

        for _, element:= range engine.Spec.Experiments{
            // fmt.Println(index, "=>", element.Name)
            chaosexperimentlist = append(chaosexperimentlist, element.Name)
        }

        // Initialize the chaosresult map before entering loop
        chaosresultmap := make(map[string]string)

        for _, test:= range chaosexperimentlist{
            chaosresultname := fmt.Sprintf("%s-%s", cEngine, test)
            testresultdump, err:= clientSet.ChaosResults("default").Get(chaosresultname, metav1.GetOptions{})
            result := testresultdump.Spec.ExperimentStatus.Verdict
            if err != nil {
                   //err.Error()
                   return 0, 0, 0, nil, err
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
        totalPassedExp = float64(pcount)               //
        totalFailedExp = float64(fcount)               //
        /////////////////////////////////////////////////
        fmt.Printf("%+v %+v %+v\n", totalExpCount, totalPassedExp, totalFailedExp)

        //Map verdict to numerical values {0-notstarted, 1-running, 2-fail, 3-pass}
        statusmap := make(map[string]float64)
        for index, status := range chaosresultmap{
            val := statusConv(status)
            statusmap[index] = val
        }
        fmt.Printf("%+v\n", statusmap)

        return totalExpCount, totalPassedExp, totalFailedExp, statusmap, nil
}
