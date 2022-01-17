# VERSION defines the project version for the bundle.
VERSION ?= 0.0.1

# Image URL to use all building/pushing image targets
IMAGE_REGISTRY ?= quay.io
REGISTRY_NAMESPACE ?= ssharon
IMAGE_NAME ?= hpessa-exporter
IMAGE_TAG ?= latest

# IMG defines the image used for the operator.
IMG ?= $(IMAGE_REGISTRY)/$(REGISTRY_NAMESPACE)/$(IMAGE_NAME):$(IMAGE_TAG)
