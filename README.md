# MXTrack: Security Observability Framework for ML/AI Model File Loading

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="brand/logo-dark.png">
    <source media="(prefers-color-scheme: light)" srcset="brand/logo-light.png">
    <img alt="MXTrack Logo" src="brand/logo.png" width="320">
  </picture>
</p>

<p align="center">
  <a href="README.md"><img alt="README in English" src="https://img.shields.io/badge/English-d9d9d9"></a>
  <a href="docs/README_CN.md"><img alt="简体中文版自述文件" src="https://img.shields.io/badge/简体中文-d9d9d9"></a>
</p>


[![Go Report Card](https://goreportcard.com/badge/github.com/mxcrafts/mxtrack)](https://goreportcard.com/report/github.com/mxcrafts/mxtrack)
[![GoDoc](https://godoc.org/github.com/mxcrafts/mxtrack?status.svg)](https://godoc.org/github.com/mxcrafts/mxtrack)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Overview

> MXTrack is an open-source security observability tool designed to monitor and analyze potential risks during the loading and execution of machine learning (ML) and artificial intelligence (AI) model files. Built with Golang and eBPF (Extended Berkeley Packet Filter), MXTrack combines the efficiency of low-level kernel tracing with the robustness of modern systems programming to deliver high-performance, low-overhead monitoring. By focusing on critical system behaviors and configurations, MXTrack helps developers, MLOps engineers, and security researchers identify vulnerabilities, unauthorized access, and anomalous activities in ML/AI workflows.

## Technical Highlights

- eBPF-Powered Efficiency
Leverages eBPF to perform lightweight, kernel-level event tracing without requiring kernel modifications. This minimizes runtime overhead (<3% CPU in most cases) while enabling real-time observation of system calls, network traffic, and file operations.

- Golang Performance & Portability
Utilizes Golang's concurrency model and cross-platform capabilities to ensure high-throughput event processing and seamless deployment across Linux distributions.

- Zero-Dependency Monitoring
Avoids reliance on external kernel modules or agents, reducing attack surfaces and operational complexity.

## Features

- 🔍 **File Monitoring**: Monitor file operations (create, delete, modify, etc.) in specified directories
- 🚀 **Process Monitoring**: Track execution of specified commands
- 🌐 **Network Monitoring**: Monitor network activity on specific ports and protocols
- 📝 **Log Management**: Support log rotation, compression, and retention policies
- ⚡ **High Performance**: Low-overhead system monitoring based on eBPF technology
- 🔧 **Configurable**: Flexible monitoring policy configuration via TOML files

## Why MXTrack?

- Low Overhead, High Fidelity
eBPF's kernel-space execution eliminates costly context switches, enabling precise tracking of system events without degrading model inference or training performance.

- Real-Time Alerts
Integrates with logging systems (e.g., Elasticsearch, Prometheus) for proactive threat response.

- Extensible Architecture
Supports plugins for custom detectors and integrations, with Golang's static binary packaging simplifying deployment.

## Use Cases

- MLOps Pipelines: Enhance security in CI/CD workflows by auditing model deployment processes.

- Research Environments: Safeguard experimental models and datasets from unintended access or tampering.

- Compliance: Meet regulatory requirements (e.g., GDPR, HIPAA) by enforcing strict access controls and audit trails.

## Quick Start

### Prerequisites

- Linux kernel version >= 4.18
- Go version >= 1.21
- LLVM/Clang 11+

### Installation

```bash
# build from source
git clone https://github.com/mxcrafts/mxtrack.git
cd mxtrack
make && ./bin/mxtrack
```

### Configuration

Create policy.toml configuration file:

```toml

[file_monitor]
enabled = true
directories = ["/path/to/monitor"]

[exec_monitor]
enabled = true
watch_commands = ["bash", "python"]

[network_monitor]
enabled = true
ports = [80, 443]
protocols = ["tcp", "udp"]

[log]
level = "info"
format = "json"
output_path = "/var/log/mxtrack/app.log"
max_size = 100  # MB
max_age = 7     # Days
max_backups = 5
compress = true

```

### Running

```bash
sudo mxtrack -config policy.toml
```


## Development

### Build Dependencies

```bash

# Build dependencies (Ubuntu)
sudo apt-get install -y clang llvm libelf-dev

# Common commands
make test       # Run unit tests
make generate   # Generate eBPF code
make package    # Create release package

```

### Performance Metrics

### Generate eBPF Code

```bash
make generate
```

### Contributing

Pull Requests and Issues are welcome! Please check our Contributing Guide for details.

### Performance Benchmarks

| Monitor Type | Event Latency | CPU Usage | Memory Usage |
|-------------|---------------|------------|--------------|
| File Monitor| < 1ms | < 1% | ~10MB |
| Process Monitor| < 0.5ms | < 0.5% | ~5MB |
| Network Monitor| < 1ms | < 1% | ~15MB |

### License
This project is licensed under the [MIT License](LICENSE).

### Contact

- Issues: GitHub Issues
- Email: support@mxcrafts.com
- Community: Discussions

### Acknowledgments

- eBPF
- Cilium
- Go
