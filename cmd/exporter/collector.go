package main

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

var (
    experimentsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "engine",
        Name:      "experiment_count",
        Help:      "Total number of experiments executed by the chaos engine",
    })

    passedExperiments = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "engine",
        Name:      "passed_experiments",
        Help:      "Total number of passed experiments",
    })

    failedExperiments = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "engine",
        Name:      "failed_experiments",
        Help:      "Total number of failed experiments",
    })

    podFailureStatus = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "experiment",
        Name:      "pod_failure_status",
        Help:      "Status of pod failure experiment",
    })

   containerKillStatus = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "experiment",
        Name:      "container_kill_status",
        Help:      "Status of container kill experiment",
    })

   containerNetworkDelay = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "experiment",
        Name:      "container_network_delay_status",
        Help:      "Status of container network delay experiment",
    })

   containerPacketLoss = prometheus.NewGauge(prometheus.GaugeOpts{
        Namespace: "chaos",
        Subsystem: "experiment",
        Name:      "container_packet_loss_status",
        Help:      "Status of container packet loss experiment",
    })
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

        experimentsTotal.Set(12)
        passedExperiments.Set(7)
        failedExperiments.Set(5)
        podFailureStatus.Set(3) //pass
        containerKillStatus.Set(2) //fail
        containerNetworkDelay.Set(1) //running 
        containerPacketLoss.Set(0) //not-started

}



