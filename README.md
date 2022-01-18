# hpessa-exporter

## Overview
Openshift's **hpessa-exporter** allows users to export SMART information of 
local storage devices as [Prometheus](https://prometheus.io) metrics, by using
HPE Smart Storage Administrator tool. Those metrics (device status, temperature 
etc.) may hint administrators on potential up-coming failures.


## Prerequisites
- Install [RedHat's Openshift](https://www.redhat.com/en/openshift-4) (4.9+)
- Install [HPE Smart Storage Administrator](https://support.hpe.com/hpesc/public/docDisplay?docId=emr_na-c04455150) 
  on local host machines. Ensure that the **ssacli** is installed on one of the
  following locations:
	- /usr/sbin/ssacli
	- /opt/smartstorageadmin/ssacli/bin/ssacli
	- /opt/hp/ssacli/bld/ssacli


## Deployment 
Use deployment yaml from this repository:

```oc apply -f deply/hpessa-exporter.yaml```

Ensure that exporter daemonset pods are up and running:
