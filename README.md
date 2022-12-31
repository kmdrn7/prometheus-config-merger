# prometheus-config-merger

Utility to merge rule_files and scrape_configs config from multiple prometheus config files into single config files.

Configurable through `prometheus-config-merger.yaml`

Example usage:
```
go run main.go merge
```
After running the command, observe the `target_prometheus_config` file for changes