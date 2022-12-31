# prometheus-config-merger

Utility to merge rule_files and scrape_configs config from multiple prometheus config files into single config files.

Let's say we have 2 prometheus config files

*prometheus_current.yaml*
```
global:
  evaluation_interval: 1m
  scrape_interval: 15s
rule_files:
  - /etc/config/recording_rules.yml
  - /etc/config/alerting_rules.yml
scrape_configs:
  - job_name: prometheus
    ...
...
```

*prometheus_operator.yaml*
```
global:
  evaluation_interval: 30s
  scrape_interval: 30s
rule_files:
  - /etc/prometheus/rules/prometheus-k8s-rulefiles-0/*.yaml
scrape_configs:
  - job_name: serviceMonitor/default/fluentd/0
    ...
...
```

For some reason, we need the `prometheus-server` to be able to read the config through those 2 config files (for **rule_files** and **scrape_configs** fields only). So, basically we need to have those multiple prometheus config files merged, to let prometheus server use it as it's main config file.

To be able to merge those config files, run this command
```
prometheus-config-merger merge

or

go run main.go merge
```

After running the command, it will create a new `prometheus.yaml` file where the merged configs is stored

*prometheus.yaml*
```
global:
  evaluation_interval: 1m
  scrape_interval: 15s
rule_files:
  - /etc/config/recording_rules.yml
  - /etc/config/alerting_rules.yml
  - /etc/prometheus/rules/prometheus-k8s-rulefiles-0/*.yaml
scrape_configs:
  - job_name: prometheus
    ...
  - job_name: serviceMonitor/default/fluentd/0
    ...
...
```

You can observe the values of **rule_files** and **scrape_configs** in new prometheus config file is taken from those 2 segregated prometheus config `prometheus_current.yaml` and `prometheus_operator.yaml`

PS: It only merge the **rule_files** and **scrape_configs** fields from multiple prometheus config files, and using one of them as base config. Behavior is configurable through `prometheus-config-merger.yaml`
