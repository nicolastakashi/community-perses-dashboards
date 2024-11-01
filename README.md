# Perses Community Dashboards

This repository hosts the code for Perses Community Dashboards, designed to serve as the Prometheus mixins for the Perses project. Built with the Perses Go SDK, these dashboards are modular, reusable, and easy to integrate.

# Available Dashboards

- Prometheus
- Node Exporter
- Alert Manager

# Library Panels

In addition to community dashboards, this repository also contains a library of reusable panels used to construct these dashboards. If you’re building custom dashboards and need a specific panel, you can leverage these library panels to create your own tailored setups.

# Local Development

To start local development, use the following command:

```bash
make start-demo
```

This will launch a local Perses instance with predefined resources like Project and DataSources. You can access the Perses UI at `http://localhost:8080`.