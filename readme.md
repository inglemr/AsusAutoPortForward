# Kubernetes Auto Port Forwarder for ASUS Routers

This GoLang project is designed to monitor Kubernetes services with the label `service.kubernetes.io/autoportforward: "true"` and automatically create port forward rules on an ASUS router based on the configuration provided.

## Overview

The Kubernetes Auto Port Forwarder for ASUS Routers is a tool that simplifies the process of managing port forwarding rules on your ASUS router for Kubernetes services. It watches for Kubernetes services with a specific label and configures port forwarding rules accordingly.

## Prerequisites

Before using this tool, ensure you have the following prerequisites:

1. An ASUS router with firmware that supports the ASUSWRT API.
2. Kubernetes cluster configured and running.
3. Proper access and permissions to interact with your Kubernetes cluster.

## Installation

1. Clone this repository to your local machine:
  ```shell
  git clone https://github.com/your-username/kubernetes-auto-port-forwarder.git
  ```
2. Build the project
  ```shell
  go build
  ```

3. Set the required environment variables:
- `ROUTER_ADDRESS`: The IP address or hostname of your ASUS router.
- `ROUTER_USERNAME`: The username to access your ASUS router.
- `ROUTER_PASSWORD`: The password to access your ASUS router.
- `DEFAULT_TARGET_ADDRESS`: The default target address for port forwarding (e.g., your Kubernetes service).
- `KUBECONFIG`: Location of the kubeconfig file.


## Usage

Deploy your Kubernetes services with the label ```service.kubernetes.io/autoportforward: "true"``` to enable automatic port forwarding.

The Kubernetes Auto Port Forwarder will watch for services with this label and create port forwarding rules on your ASUS router automatically.

You can customize the target address for port forwarding by setting the DEFAULT_TARGET_ADDRESS environment variable.