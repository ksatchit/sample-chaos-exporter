/* The chaos exporter collects and exposes the following type of metrics:

   Fixed (always captured):
     - Total number of chaos experiments 
     - Total number of passed experiments 
     - Total Number of failed experiments
 
   Dynamic (experiment list may vary based on c.engine):
     - States of individual chaos experiments
     - {not-executed:0, running:1, fail:2, pass:3}
       TODO: Improve representaion of test state

   Common experiments include:
 
     - pod_failure
     - container_kill
     - container_network_delay
     - container_packet_loss
*/

package main

import (
  "os"
  "flag"
  "net/http"
  "strings"
  "github.com/ksatchit/sample-chaos-exporter/pkg/util"
  log "github.com/Sirupsen/logrus"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/client-go/rest" 
)

// Declare general variables (cluster ops, error handling, misc)
var kubeconfig string
var config *rest.Config
var err error
var app_uuid, chaosengine string 

func main(){

    // Get app details & chaoengine name from ENV 
    app_uuid := os.Getenv("APP_UUID")
    chaosengine := os.Getenv("CHAOSENGINE")

    flag.StringVar(&kubeconfig, "kubeconfig", "", "path to the kubeconfig file")
    flag.Parse()

    // Use in-cluster config if kubeconfig file not available
    if kubeconfig == "" {
        log.Info("using the in-cluster config")
        config, err = rest.InClusterConfig()
    } else {
        log.Info("using configuration from '%s'", kubeconfig)
        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
    }

    if err != nil {
        panic(err.Error())
    }

    // Validate availability of mandatory ENV
    if chaosengine == "" || app_uuid == "" {
        log.Fatal("ERROR: please specify correct APP_UUID & CHAOSENGINE ENVs")
        os.Exit(1)
    }

    // Declare the fixed chaos metrics 
    var (
        experimentsTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "c",
            Subsystem: "engine",
            Name:      "experiment_count",
            Help:      "Total number of experiments executed by the chaos engine",
        },
        []string{"app_uid"},
        )

        passedExperiments = prometheus.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "c",
            Subsystem: "engine",
            Name:      "passed_experiments",
            Help:      "Total number of passed experiments",
        },
        []string{"app_uid"},
        )

        failedExperiments = prometheus.NewGaugeVec(prometheus.GaugeOpts{
            Namespace: "c",
            Subsystem: "engine",
            Name:      "failed_experiments",
            Help:      "Total number of failed experiments",
        },
        []string{"app_uid"},
        )

    )

    // Get the chaos metrics for the specified chaosengine 
    expTotal, passTotal, failTotal, expMap, err := util.GetChaosMetrics(config, chaosengine)
    if err != nil {
        //panic(err.Error())
        log.Fatal("Unable to get metrics: ", err.Error())
    }

    // Define, register & set the dynamically obtained chaos metrics (experiment state)
    for index, verdict := range expMap{
        sanitized_exp_name := strings.Replace(index, "-", "_", -1)

        var (
            tmpExp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
                Namespace: "c",
                Subsystem: "exp",
                Name:      sanitized_exp_name,
                Help: "",
                },
                []string{"app_uid"},
                )
         )

         prometheus.MustRegister(tmpExp)
         tmpExp.WithLabelValues(app_uuid).Set(verdict)
    }

    // Register the fixed chaos metrics
    prometheus.MustRegister(experimentsTotal)
    prometheus.MustRegister(passedExperiments)
    prometheus.MustRegister(failedExperiments)

    // Set the fixed chaos metrics
    experimentsTotal.WithLabelValues(app_uuid).Set(expTotal)
    passedExperiments.WithLabelValues(app_uuid).Set(passTotal)
    failedExperiments.WithLabelValues(app_uuid).Set(failTotal)

    //This section will start the HTTP server and expose
    //any metrics on the /metrics endpoint.
    http.Handle("/metrics", promhttp.Handler())
    log.Info("Beginning to serve on port :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
