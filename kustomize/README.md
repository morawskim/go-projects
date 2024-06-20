# Kustomize KRM plugin

This repository contains a Kustomize plugin written in Go that transforms Kubernetes CronJob resources by modifying their suspend flag.

## Overview

The plugin reads a ResourceList containing Kubernetes resources, processes CronJob resources, and set value for field `spec.schedule`.
This is an example plugin demonstrating how to manipulate Kubernetes resources using Kustomize plugins and Go.

## Usage

To run unit tests - `make tests`

To build KRM plugin - `make build`

To dump Kustomize result - `make kustomize`

To apply changes in K8S cluster - `make kubectl-apply`.
The namespace "kustomize-plugin" has to exist.
To create this namespace you can use `kubectl apply -f kustomize/static/namespace.yml`
