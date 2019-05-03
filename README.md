# sample-chaos-exporter
custom prometheus exporter to expose litmus chaos metrics 

## Steps: 

- Run the exporter container (ksatchit/sample-chaos-exporter:trial) on host network
- Run the prometheus container (ksatchit/prometheus:trial) w/ config specified in the 
  prometheus/prometheus.yaml on host network
- Run grafana container and add prometheus as data source
