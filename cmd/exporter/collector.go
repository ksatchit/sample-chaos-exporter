package main

/*
TODO: Implement the collector interface as per prometheus best practices
TODO: Freeze upon the default experiment list across chaos engines
TODO: Pass the application UUID as ENV to exporter container
TODO: Pass the chaosengine name as ENV to exporter container
TODO: Implement the logic to update status from controller on chaosengine CR 
*/

/* The chaos exporter collects and exposes the following metrics:

   # Total Number of ChaosExperiments
   # Total Number of PassedExperiments 
   # Total Number of FailedExperiments
   # Status of engine's experiments, such as:
 
     a. pod_failure
     b. container_kill
     c. container_network_delay
     d. container_packet_loss
*/

import (
        "os"
        "flag"
        // "fmt"
        "log"
        "strings"
        "path/filepath"
        "k8s.io/client-go/util/homedir"
        "k8s.io/client-go/tools/clientcmd"
        "k8s.io/client-go/rest"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/ksatchit/sample-chaos-exporter/pkg/util"
)

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

func init() {

        // Get cluster configuration
        var kubeconfig *string 

        // Get app details & chaoengine name from ENV
        app_uuid := os.Getenv("APP_UUID")
        chaosengine := os.Getenv("CHAOSENGINE")

        if home := homedir.HomeDir(); home != ""  {
                kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) path to the kubeconfig file")
        } else {
               kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
        }

        flag.Parse()

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

        if chaosengine == "" || app_uuid == "" {
               log.Printf("ERROR: please specify correct APP_UUID & CHAOSENGINE ENVs")
               os.Exit(1)
        }

        expTotal, passTotal, failTotal, expMap, err := util.GetChaosMetrics(config, chaosengine)
        if err != nil {
                panic(err.Error())
        }

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

	prometheus.MustRegister(experimentsTotal)
	prometheus.MustRegister(passedExperiments)
	prometheus.MustRegister(failedExperiments)

        experimentsTotal.WithLabelValues(app_uuid).Set(expTotal)
        passedExperiments.WithLabelValues(app_uuid).Set(passTotal)
        failedExperiments.WithLabelValues(app_uuid).Set(failTotal)
}



