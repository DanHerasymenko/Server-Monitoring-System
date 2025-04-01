# Author: Herasymenko Daniil

---

# Server Monitoring System

This project implements a scalable **Server Monitoring System** designed to collect, store, and visualize system metrics (CPU, RAM, Disk) from multiple agents.

Agents are installed on remote machines and periodically collect metrics, which are sent via **gRPC** to a central server. The server caches the latest metrics for real-time access and stores historical data for further analysis.

---

## Key Features

- **Metric Collection** from multiple agents
- **Real-Time gRPC Streaming**
- **Centralized Caching (Redis)**
- **Historical Data Storage (PostgreSQL)**
- **Prometheus + Grafana Dashboards**
- **Cross-platform Agent (Windows/Linux)**
- **Graceful Shutdown & Logging**

---

## Project Components

### 1. **Agent**
- Collects system metrics using `gopsutil`
- Sends metrics to the server via gRPC stream
- Implemented as a Windows/Linux service
- Configurable via environment variables

### 2. **Server**
- Listens for agent gRPC streams
- Validates and logs received metrics
- Caches latest metrics in **Redis**
- Stores historical metrics in **PostgreSQL**
- Exposes Prometheus-compatible endpoint

### 3. **Redis (Cache)**
- Stores the **latest metrics** per agent
- Ensures **fast access** for dashboard queries
- Implements TTL to detect inactive agents

### 4. **PostgreSQL (Storage)**
- Long-term storage of all metric submissions
- Enables querying for reports, trends, and analytics

### 5. **Prometheus & Grafana**
- Prometheus scrapes metric endpoint
- Grafana visualizes agent health, usage stats, etc.
- Dashboards include CPU, RAM, Disk history and uptime

### 6. **Logging System**
- Uses `slog` for structured JSON logging
- Logs are rotated and written to disk
- Logs include AgentIP and ServerIP context

---

## Technology Stack

- **Language**: Go
- **Protocol**: gRPC
- **Metrics Collection**: gopsutil
- **Database**: PostgreSQL
- **Cache**: Redis
- **Monitoring**: Prometheus + Grafana
- **Service Management**: systemd / Windows Services
- **Containerization**: Docker & Docker Compose

---

## Usage

- Run server via `go run cmd/server/main.go`
- Install and start agent via `monitoring-agent.exe install && monitoring-agent.exe start`
- View logs in `/logs/monitoring_agent.log`
- View dashboards on `localhost:3000` (Grafana)

---

## Documentation

- **Proto schema** in `proto/monitoring.proto`
- **gRPC UI Testing** with [grpcui](https://github.com/fullstorydev/grpcui) command `grpcui.exe -plaintext localhost:50051`
- Prometheus and Grafana configs in `deployments/`

---

This project provides a modern, extensible framework for tracking the health and performance of distributed systems.

