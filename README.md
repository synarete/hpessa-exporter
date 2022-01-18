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
	- `/usr/sbin/ssacli`
	- `/opt/smartstorageadmin/ssacli/bin/ssacli`
	- `/opt/hp/ssacli/bld/ssacli`


## Deployment 
Use deployment yaml from this repository:

```oc apply -f deply/hpessa-exporter.yaml```

Verify that exporter daemonset pods are up and running:

```
$ oc get pods -n openshift-storage-hpessa
NAME                    READY   STATUS    RESTARTS   AGE
hpessa-exporter-2cbvl   1/1     Running   0          35s
hpessa-exporter-sk2cn   1/1     Running   0          35s
hpessa-exporter-v699z   1/1     Running   0          35s
```

Verify that each of the running pods exports local metrics

```
$ oc curl http://hpessa-exporter-2cbvl:8080/metrics -n openshift-storage-hpessa
# HELP hpessa_blkdev_size_bytes Block device size in bytes.
# TYPE hpessa_blkdev_size_bytes gauge
hpessa_blkdev_size_bytes{major="8",minor="16",model="LOGICAL VOLUME",name="sdb",vendor="HPE"} 2.3441958064e+10
# HELP hpessa_ssa_physical_device_power_hours Power on in hours
# TYPE hpessa_ssa_physical_device_power_hours gauge
hpessa_ssa_physical_device_power_hours{bay="1",box="2",dev="/dev/sdb",id="physicaldrive 1I:2:1",uniqueid="5000C50094D7BEB3"} -1
# HELP hpessa_ssa_physical_device_size Size in bytes of physical device
# TYPE hpessa_ssa_physical_device_size gauge
hpessa_ssa_physical_device_size{bay="1",box="2",dev="/dev/sdb",id="physicaldrive 1I:2:1",uniqueid="5000C50094D7BEB3"} 6.597069766656e+12
# HELP hpessa_ssa_physical_device_status Status of physical device
# TYPE hpessa_ssa_physical_device_status gauge
hpessa_ssa_physical_device_status{bay="1",box="2",dev="/dev/sdb",id="physicaldrive 1I:2:1",uniqueid="5000C50094D7BEB3"} 0
# HELP hpessa_ssa_physical_device_temp_curr Current temperature of physical device
# TYPE hpessa_ssa_physical_device_temp_curr gauge
hpessa_ssa_physical_device_temp_curr{bay="1",box="2",dev="/dev/sdb",id="physicaldrive 1I:2:1",uniqueid="5000C50094D7BEB3"} 34
# HELP hpessa_ssa_physical_device_temp_maxi Maximal temperature of physical device
# TYPE hpessa_ssa_physical_device_temp_maxi gauge
hpessa_ssa_physical_device_temp_maxi{bay="1",box="2",dev="/dev/sdb",id="physicaldrive 1I:2:1",uniqueid="5000C50094D7BEB3"} 48
...
```
