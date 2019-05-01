package main

/*
TODO: Implement the collector interface as per prometheus best practices
TODO: Freeze upon the default experiment list across chaos engines
TODO: Pass the application UUID as ENV to exporter container
TODO: Pass the chaosengine name as ENV to exporter container
TODO: Implement the logic to update status from controller on chaosengine CR 
TODO: Implement the client-go logic to extract the desired info from CR
*/

/* The chaos exporter collects and exposes the following metrics:

   # Total Number of ChaosExperiments
   # Total Number of PassedExperiments 
   # Total Number of FailedExperiments
   # Status of following experiments:
 
     a. pod_failure
     b. container_kill
     c. container_network_delay
     d. container_packet_loss
*/

import "github.com/prometheus/client_golang/prometheus"

const application_uuid = "3f2092f8-6400-11e9-905f-42010a800131"
const chaosEngineName = "engine-nginx"

var (
    experimentsTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "engine",
        Name:      "experiment_count",
        Help:      "Total number of experiments executed by the chaos engine",
    },
    []string{"app_uid"},
    )

    passedExperiments = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "engine",
        Name:      "passed_experiments",
        Help:      "Total number of passed experiments",
    },
    []string{"app_uid"},
    )

    failedExperiments = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "engine",
        Name:      "failed_experiments",
        Help:      "Total number of failed experiments",
    },
    []string{"app_uid"},
    )

    podFailureStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "experiment",
        Name:      "pod_failure_status",
        Help:      "Status of pod failure experiment",
    },
    []string{"app_uid"},
    )

   containerKillStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "experiment",
        Name:      "container_kill_status",
        Help:      "Status of container kill experiment",
    },
    []string{"app_uid"},
    )

   containerNetworkDelay = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "experiment",
        Name:      "container_network_delay_status",
        Help:      "Status of container network delay experiment",
    },
    []string{"app_uid"},
    )

   containerPacketLoss = prometheus.NewGaugeVec(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "experiment",
        Name:      "container_packet_loss_status",
        Help:      "Status of container packet loss experiment",
    },
    []string{"app_uid"},
    )
)

func init() {
	prometheus.MustRegister(experimentsTotal)
	prometheus.MustRegister(passedExperiments)
	prometheus.MustRegister(failedExperiments)
	prometheus.MustRegister(podFailureStatus)
	prometheus.MustRegister(containerKillStatus)
	prometheus.MustRegister(containerNetworkDelay)
	prometheus.MustRegister(containerPacketLoss)


        // Set default metrics for debug purposes
        experimentsTotal.WithLabelValues(application_uuid).Set(12)
        passedExperiments.WithLabelValues(application_uuid).Set(7)
        failedExperiments.WithLabelValues(application_uuid).Set(5)
        podFailureStatus.WithLabelValues(application_uuid).Set(3) //pass
        containerKillStatus.WithLabelValues(application_uuid).Set(2) //fail
        containerNetworkDelay.WithLabelValues(application_uuid).Set(1) //running 
        containerPacketLoss.WithLabelValues(application_uuid).Set(0) //not-started

}



