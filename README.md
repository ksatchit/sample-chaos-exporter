# Sample-chaos-exporter

- This is a custom prometheus exporter to expose litmus chaos metrics
 
- The exporter is typically tied to a given chaosengine custom resource, 
  which, in-turn is associated with a given application deployment

- Two types of metrics are exposed: 

  - Fixed: TotalExperimentCount, TotalPassedTests, TotalFailedTests
  - Dymanic: Individual experiment status. The list of experiments may 
    vary across chaos engines. The exporter reports experiment status as
    per list in the chaosengine

- The metrics are of type Gauge, w/ each of the status metrics mapped to a 
  numeric value(not-executed:0, running:1, fail:2, pass:3)

- The metrics carry the application_uuid as label (this has to be passed as ENV)

## Steps: 

### Local Machine (test purposes) 

- Set the application deployment (assuming a live K8s cluster w/ app) UUID as ENV (APP_UUID)

- Set the ChaosEngine CR name as ENV (CHAOSENGINE) 
  - For CR spec, see: https://github.com/litmuschaos/chaos-operator/blob/master/deploy/crds/chaosengine.yaml

- If the experiments are not executed, apply the ChaosResult CRs manually 
  - For CR spec, see: https://github.com/litmuschaos/chaos-operator/blob/master/deploy/crds/chaosresult.yaml

- Run the exporter container (ksatchit/sample-chaos-exporter:ci) on host network. It is necessary to mount the kubeconfig
  & override entrypoint w/ `./exporter -kubeconfig <path>`

- Execute `curl 127.0.0.1:8080/metrics` to view metrics

- (Optional) Run the prometheus container (ksatchit/prometheus:trial) w/ config specified in the 
  prometheus/prometheus.yaml on host network

- (Optional) Run grafana container, add prometheus as data source

### On Kubernetes Cluster

- Install the RBAC (serviceaccount, role, rolebinding) as per deploy/rbac.md

- Deploy the chaos-exporter.yaml 

- From a cluster node, execute `curl <exporter-service-ip>:8080/metrics` 

- (Optional) Deploy prometheus with desired scrape interval & grafana deployments 

### Example Metrics

```
c_engine_experiment_count{app_uid="3f2092f8-6400-11e9-905f-42010a800131"} 2
# HELP c_engine_failed_experiments Total number of failed experiments
# TYPE c_engine_failed_experiments gauge
c_engine_failed_experiments{app_uid="3f2092f8-6400-11e9-905f-42010a800131"} 1
# HELP c_engine_passed_experiments Total number of passed experiments
# TYPE c_engine_passed_experiments gauge
c_engine_passed_experiments{app_uid="3f2092f8-6400-11e9-905f-42010a800131"} 1
# HELP c_exp_engine_nginx_container_kill 
# TYPE c_exp_engine_nginx_container_kill gauge
c_exp_engine_nginx_container_kill{app_uid="3f2092f8-6400-11e9-905f-42010a800131"} 2
# HELP c_exp_engine_nginx_pod_failure 
# TYPE c_exp_engine_nginx_pod_failure gauge
c_exp_engine_nginx_pod_failure{app_uid="3f2092f8-6400-11e9-905f-42010a800131"} 3
```


## Roadmap

(This will be updated/active) 

- Rethink metric types for status 
- Add more labels
