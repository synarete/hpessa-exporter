// SPDX-License-Identifier: Apache-2.0
package devmon_test

import (
	"testing"

	"github.com/red-hat-storage/hpessa-exporter/internal/devmon"
	"github.com/stretchr/testify/assert"
)

var (
	ssacliVersion = `

	SSACLI Version: 4.21.7.0 2020-07-15
	SOULAPI Version: 4.21.7.0 2020-07-15

 `
	ssacliCtrlAllShowConfigDetail = `

Smart HBA H240 in Slot 1 (RAID Mode)
   Bus Interface: PCI
   Slot: 1
   Serial Number: PDNNK0BRH571XZ
   Cache Serial Number: PDNNK0BRH571XZ
   Controller Status: OK
   Hardware Revision: B
   Firmware Version: 4.52-0
   Firmware Supports Online Firmware Activation: False
   Rebuild Priority: High
   Expand Priority: Medium
   Surface Scan Delay: 3 secs
   Surface Scan Mode: Idle
   Parallel Surface Scan Supported: Yes
   Current Parallel Surface Scan Count: 1
   Max Parallel Surface Scan Count: 16
   Queue Depth: Automatic
   Monitor and Performance Delay: 60  min
   Elevator Sort: Enabled
   Degraded Performance Optimization: Disabled
   Wait for Cache Room: Disabled
   Surface Analysis Inconsistency Notification: Disabled
   Post Prompt Timeout: 15 secs
   Cache Board Present: False
   Drive Write Cache: Disabled
   Controller Memory Size: 0.2
   SATA NCQ Supported: True
   Spare Activation Mode: Activate on physical drive failure (default)
   Controller Temperature (C): 41
   Number of Ports: 2 Internal only
   Encryption: Not Set
   Express Local Encryption: False
   Driver Name: hpsa
   Driver Version: 3.4.20
   Driver Supports SSD Smart Path: True
   PCI Address (Domain:Bus:Device.Function): 0000:03:00.0
   Negotiated PCIe Data Rate: PCIe 3.0 x8 (7880 MB/s)
   Controller Mode: RAID Mode
   Pending Controller Mode: RAID
   Port Max Phy Rate Limiting Supported: False
   Latency Scheduler Setting: Disabled
   Current Power Mode: MaxPerformance
   Survival Mode: Enabled
   Host Serial Number: SGH706XME6
   Sanitize Erase Supported: True
   Primary Boot Volume: logicaldrive 1 (600508B1001C90DB4A1FDCBCCB744F14)
   Secondary Boot Volume: None



   Internal Drive Cage at Port 1I, Box 1, OK

      Drive Bays: 4
      Port: 1I
      Box: 1
      Location: Internal

   Physical Drives
      physicaldrive 1I:1:1 (port 1I:box 1:bay 1, SATA HDD, 1 TB, OK)
      physicaldrive 1I:1:2 (port 1I:box 1:bay 2, SATA SSD, 400 GB, OK)


   Port Name: 1I
         Port ID: 1
         Port Connection Number: 1
         SAS Address: 50014380408E6FE4
         Port Location: Internal

   Port Name: 2I
         Port ID: 0
         Port Connection Number: 0
         SAS Address: 50014380408E6FE0
         Port Location: Internal

   Array: A
      Interface Type: Solid State SATA
      Unused Space: 0 MB (0.00%)
      Used Space: 372.58 GB (100.00%)
      Status: OK
      MultiDomain Status: OK
      Array Type: Data
      Smart Path: enable


      Logical Drive: 1
         Size: 372.58 GB
         Fault Tolerance: 0
         Heads: 255
         Sectors Per Track: 32
         Cylinders: 65535
         Strip Size: 256 KB
         Full Stripe Size: 256 KB
         Status: OK
         MultiDomain Status: OK
         Caching:  Disabled
         Unique Identifier: 600508B1001C90DB4A1FDCBCCB744F14
         Disk Name: /dev/sda
         Mount Points: 1024 MiB Partition   1 /boot
         Disk Partition Information
            Partition   1: Basic, 1024 MiB, /boot
         Boot Volume: Primary
         Drive Type: Data
         LD Acceleration Method: Smart Path


      physicaldrive 1I:1:2
         Port: 1I
         Box: 1
         Bay: 2
         Status: OK
         Drive Type: Data Drive
         Interface Type: Solid State SATA
         Size: 400 GB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/4096
         Firmware Revision: 4IWTHPG1
         Serial Number: BTHV603000TL400NGN
         WWID: 30014380408E6FC5
         Model: ATA     MK0400GEYKD
         SATA NCQ Capable: True
         SATA NCQ Enabled: True
         Current Temperature (C): 17
         Maximum Temperature (C): 35
         Usage remaining: 99.73%
         Power On Hours: 40698
         Estimated Life Remaining based on workload to date: 626359 days
         SSD Smart Trip Wearout: False
         PHY Count: 1
         PHY Transfer Rate: 6.0Gbps
         PHY Physical Link Rate: Unknown
         PHY Maximum Link Rate: Unknown
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 2 minute(s), 0 second(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None



   Array: B
      Interface Type: SATA
      Unused Space: 0 MB (0.00%)
      Used Space: 931.48 GB (100.00%)
      Status: OK
      MultiDomain Status: OK
      Array Type: Data
      Smart Path: disable


      Logical Drive: 2
         Size: 931.48 GB
         Fault Tolerance: 0
         Heads: 255
         Sectors Per Track: 32
         Cylinders: 65535
         Strip Size: 512 KB
         Full Stripe Size: 512 KB
         Status: OK
         MultiDomain Status: OK
         Caching:  Disabled
         Unique Identifier: 600508B1001CD499DA451C86BDC0A6CC
         Disk Name: /dev/sdb
         Mount Points: None
         Logical Drive Label: 06EE798FPDNNK0BRH571XZDADA
         Drive Type: Data
         LD Acceleration Method: All disabled


      physicaldrive 1I:1:1
         Port: 1I
         Box: 1
         Bay: 1
         Status: OK
         Drive Type: Data Drive
         Interface Type: SATA
         Size: 1 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPG4
         Serial Number: 17V3K9U7F1EA
         WWID: 30014380408E6FC4
         Model: ATA     MB1000GDUNU
         SATA NCQ Capable: True
         SATA NCQ Enabled: True
         Current Temperature (C): 21
         Maximum Temperature (C): 38
         PHY Count: 1
         PHY Transfer Rate: 6.0Gbps
         PHY Physical Link Rate: Unknown
         PHY Maximum Link Rate: Unknown
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 2 hour(s), 36 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None

`
	ssacliCtrlAllShowConfigDetail2 = `

HPE Smart Array P816i-a SR Gen10 in Slot 0 (Embedded)
   Bus Interface: PCI
   Slot: 0
   Serial Number: PEYHD0CRHB00IZ
   RAID 6 (ADG) Status: Enabled
   Controller Status: OK
   Hardware Revision: A
   Firmware Version: 1.34-0
   Firmware Supports Online Firmware Activation: False
   Rebuild Priority: High
   Expand Priority: Medium
   Surface Scan Delay: 3 secs
   Surface Scan Mode: Idle
   Parallel Surface Scan Supported: Yes
   Current Parallel Surface Scan Count: 1
   Max Parallel Surface Scan Count: 16
   Queue Depth: Automatic
   Monitor and Performance Delay: 60  min
   Elevator Sort: Enabled
   Degraded Performance Optimization: Disabled
   Inconsistency Repair Policy: Disabled
   Write Cache Bypass Threshold Size: 1040 KiB
   Wait for Cache Room: Disabled
   Surface Analysis Inconsistency Notification: Disabled
   Post Prompt Timeout: 15 secs
   Cache Board Present: True
   Cache Status: OK
   Cache Ratio: 10% Read / 90% Write
   Drive Write Cache: Disabled
   Total Cache Size: 4.0
   Total Cache Memory Available: 3.8
   No-Battery Write Cache: Disabled
   SSD Caching RAID5 WriteBack Enabled: True
   SSD Caching Version: 2
   Cache Backup Power Source: Batteries
   Battery/Capacitor Count: 1
   Battery/Capacitor Status: OK
   SATA NCQ Supported: True
   Spare Activation Mode: Activate on physical drive failure (default)
   Controller Temperature (C): 58
   Capacitor Temperature  (C): 46
   Number of Ports: 4 Internal only
   Encryption: Not Set
   Express Local Encryption: False
   Driver Name: smartpqi
   Driver Version: Linux 1.2.16-012
   PCI Address (Domain:Bus:Device.Function): 0000:5C:00.0
   Negotiated PCIe Data Rate: PCIe 3.0 x8 (7880 MB/s)
   Controller Mode: Mixed
   Port Max Phy Rate Limiting Supported: False
   Latency Scheduler Setting: Disabled
   Current Power Mode: MaxPerformance
   Survival Mode: Enabled
   Host Serial Number: SGH829WV4W
   Sanitize Erase Supported: True
   Sanitize Lock: None
   Sensor ID: 0
      Location: Capacitor
      Current Value (C): 46
      Max Value Since Power On: 55
   Sensor ID: 1
      Location: ASIC
      Current Value (C): 58
      Max Value Since Power On: 70
   Sensor ID: 2
      Location: Unknown
      Current Value (C): 49
      Max Value Since Power On: 59
   Primary Boot Volume: None
   Secondary Boot Volume: None



   Internal Drive Cage at Port 1I, Box 2, OK

      Drive Bays: 4
      Port: 1I
      Box: 2
      Location: Internal

   Physical Drives
      physicaldrive 1I:2:1 (port 1I:box 2:bay 1, SAS HDD, 6 TB, OK)
      physicaldrive 1I:2:2 (port 1I:box 2:bay 2, SAS HDD, 6 TB, OK)
      physicaldrive 1I:2:3 (port 1I:box 2:bay 3, SAS HDD, 6 TB, OK)
      physicaldrive 1I:2:4 (port 1I:box 2:bay 4, SAS HDD, 6 TB, OK)



   Internal Drive Cage at Port 2I, Box 3, OK

      Drive Bays: 4
      Port: 2I
      Box: 3
      Location: Internal

   Physical Drives
      physicaldrive 2I:3:1 (port 2I:box 3:bay 1, SAS HDD, 6 TB, OK)
      physicaldrive 2I:3:2 (port 2I:box 3:bay 2, SAS HDD, 6 TB, OK)
      physicaldrive 2I:3:3 (port 2I:box 3:bay 3, SAS HDD, 6 TB, OK)
      physicaldrive 2I:3:4 (port 2I:box 3:bay 4, SAS HDD, 6 TB, OK)



   Internal Drive Cage at Port 3I, Box 6, OK

      Drive Bays: 4
      Port: 3I
      Box: 6
      Location: Internal

   Physical Drives
      physicaldrive 3I:6:1 (port 3I:box 6:bay 1, SATA SSD, 960 GB, OK)
      physicaldrive 3I:6:2 (port 3I:box 6:bay 2, SATA SSD, 960 GB, OK)


   Port Name: 1I
         Port ID: 0
         Port Mode: Mixed
         Port Connection Number: 0
         SAS Address: 51402EC010337940
         Port Location: Internal

   Port Name: 2I
         Port ID: 1
         Port Mode: Mixed
         Port Connection Number: 1
         SAS Address: 51402EC010337944
         Port Location: Internal

   Port Name: 3I
         Port ID: 2
         Port Mode: Mixed
         Port Connection Number: 2
         SAS Address: 51402EC010337948
         Port Location: Internal

   Port Name: 4I
         Port ID: 3
         Port Mode: Mixed
         Port Connection Number: 3
         SAS Address: 51402EC01033794C
         Port Location: Internal

   Array: A
      Interface Type: SAS
      Unused Space: 1 MB (0.00%)
      Used Space: 21.83 TB (100.00%)
      Status: OK
      MultiDomain Status: OK
      Array Type: Data
      Smart Path: disable


      Logical Drive: 1
         Size: 10.92 TB
         Fault Tolerance: 6
         Heads: 255
         Sectors Per Track: 32
         Cylinders: 65535
         Strip Size: 256 KB
         Full Stripe Size: 512 KB
         Status: OK
         Unrecoverable Media Errors: None
         MultiDomain Status: OK
         Caching:  Enabled
         Parity Initialization Status: Initialization Completed
         Unique Identifier: 600508B1001CCD7B72DB95E459CDDCCC
         Disk Name: /dev/sdb
         Mount Points: 10.9 TiB Partition   4, 384 MiB Partition   3 /sysroot,/,
         Disk Partition Information
            Partition   4: Basic, 10.9 TiB, /sysroot,/,/etc,/usr,/var,
            Partition   3: Basic, 384 MiB, /boot
         Logical Drive Label: 00F87066PEYHD0CRHB00IZ E143
         Drive Type: Data
         LD Acceleration Method: Controller Cache


      physicaldrive 1I:2:1
         Port: 1I
         Box: 2
         Bay: 1
         Status: OK
         Drive Type: Data Drive
         Interface Type: SAS
         Size: 6 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPD2
         Serial Number: ZA19G381
         WWID: 5000C50094D7BEB1
         Model: HP      MB6000JVYYV
         Current Temperature (C): 35
         Maximum Temperature (C): 48
         PHY Count: 2
         PHY Transfer Rate: 12.0Gbps, Unknown
         PHY Physical Link Rate: 12.0Gbps, Unknown
         PHY Maximum Link Rate: 12.0Gbps, 12.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 11 hour(s), 40 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 5000C50094D7BEB3

      physicaldrive 1I:2:2
         Port: 1I
         Box: 2
         Bay: 2
         Status: OK
         Drive Type: Data Drive
         Interface Type: SAS
         Size: 6 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPD2
         Serial Number: ZA19ARG0
         WWID: 5000C50094DDBA4D
         Model: HP      MB6000JVYYV
         Current Temperature (C): 35
         Maximum Temperature (C): 50
         PHY Count: 2
         PHY Transfer Rate: 12.0Gbps, Unknown
         PHY Physical Link Rate: 12.0Gbps, Unknown
         PHY Maximum Link Rate: 12.0Gbps, 12.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 11 hour(s), 40 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 5000C50094DDBA4F

      physicaldrive 1I:2:3
         Port: 1I
         Box: 2
         Bay: 3
         Status: OK
         Drive Type: Data Drive
         Interface Type: SAS
         Size: 6 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPD2
         Serial Number: ZA19CRSY
         WWID: 5000C50094D31249
         Model: HP      MB6000JVYYV
         Current Temperature (C): 35
         Maximum Temperature (C): 50
         PHY Count: 2
         PHY Transfer Rate: 12.0Gbps, Unknown
         PHY Physical Link Rate: 12.0Gbps, Unknown
         PHY Maximum Link Rate: 12.0Gbps, 12.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 11 hour(s), 40 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 5000C50094D3124B

      physicaldrive 1I:2:4
         Port: 1I
         Box: 2
         Bay: 4
         Status: OK
         Drive Type: Data Drive
         Interface Type: SAS
         Size: 6 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPD2
         Serial Number: ZA19AMJJ
         WWID: 5000C50094DDBB81
         Model: HP      MB6000JVYYV
         Current Temperature (C): 37
         Maximum Temperature (C): 51
         PHY Count: 2
         PHY Transfer Rate: 12.0Gbps, Unknown
         PHY Physical Link Rate: 12.0Gbps, Unknown
         PHY Maximum Link Rate: 12.0Gbps, 12.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 11 hour(s), 40 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 5000C50094DDBB83



   Array: B
      Interface Type: SAS
      Unused Space: 1 MB (0.00%)
      Used Space: 21.83 TB (100.00%)
      Status: OK
      MultiDomain Status: OK
      Array Type: Data
      Smart Path: disable


      Logical Drive: 2
         Size: 10.92 TB
         Fault Tolerance: 6
         Heads: 255
         Sectors Per Track: 32
         Cylinders: 65535
         Strip Size: 256 KB
         Full Stripe Size: 512 KB
         Status: OK
         Unrecoverable Media Errors: None
         MultiDomain Status: OK
         Caching:  Enabled
         Parity Initialization Status: Initialization Completed
         Unique Identifier: 600508B1001C94F35699282EC75137E2
         Disk Name: /dev/sdc
         Mount Points: None
         Logical Drive Label: 04F87089PEYHD0CRHB00IZ BF80
         Drive Type: Data
         LD Acceleration Method: Controller Cache


      physicaldrive 2I:3:1
         Port: 2I
         Box: 3
         Bay: 1
         Status: OK
         Drive Type: Data Drive
         Interface Type: SAS
         Size: 6 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPD2
         Serial Number: ZA19FM7S
         WWID: 5000C50094D19209
         Model: HP      MB6000JVYYV
         Current Temperature (C): 38
         Maximum Temperature (C): 51
         PHY Count: 2
         PHY Transfer Rate: 12.0Gbps, Unknown
         PHY Physical Link Rate: 12.0Gbps, Unknown
         PHY Maximum Link Rate: 12.0Gbps, 12.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 11 hour(s), 40 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 5000C50094D1920B

      physicaldrive 2I:3:2
         Port: 2I
         Box: 3
         Bay: 2
         Status: OK
         Drive Type: Data Drive
         Interface Type: SAS
         Size: 6 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPD2
         Serial Number: ZA19ERG2
         WWID: 5000C50094CFF58D
         Model: HP      MB6000JVYYV
         Current Temperature (C): 38
         Maximum Temperature (C): 53
         PHY Count: 2
         PHY Transfer Rate: 12.0Gbps, Unknown
         PHY Physical Link Rate: 12.0Gbps, Unknown
         PHY Maximum Link Rate: 12.0Gbps, 12.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 11 hour(s), 40 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 5000C50094CFF58F

      physicaldrive 2I:3:3
         Port: 2I
         Box: 3
         Bay: 3
         Status: OK
         Drive Type: Data Drive
         Interface Type: SAS
         Size: 6 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPD2
         Serial Number: ZA19FNFX
         WWID: 5000C50094D125C5
         Model: HP      MB6000JVYYV
         Current Temperature (C): 36
         Maximum Temperature (C): 52
         PHY Count: 2
         PHY Transfer Rate: 12.0Gbps, Unknown
         PHY Physical Link Rate: 12.0Gbps, Unknown
         PHY Maximum Link Rate: 12.0Gbps, 12.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 11 hour(s), 40 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 5000C50094D125C7

      physicaldrive 2I:3:4
         Port: 2I
         Box: 3
         Bay: 4
         Status: OK
         Drive Type: Data Drive
         Interface Type: SAS
         Size: 6 TB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/512
         Rotational Speed: 7200
         Firmware Revision: HPD2
         Serial Number: ZA19FMFQ
         WWID: 5000C50094D18741
         Model: HP      MB6000JVYYV
         Current Temperature (C): 38
         Maximum Temperature (C): 52
         PHY Count: 2
         PHY Transfer Rate: 12.0Gbps, Unknown
         PHY Physical Link Rate: 12.0Gbps, Unknown
         PHY Maximum Link Rate: 12.0Gbps, 12.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Estimated Max Erase Time: 11 hour(s), 40 minute(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 5000C50094D18743



   Array: C
      Interface Type: Solid State SATA
      Unused Space: 0 MB (0.00%)
      Used Space: 1.75 TB (100.00%)
      Status: OK
      MultiDomain Status: OK
      Array Type: Data
      Smart Path: enable


      Logical Drive: 3
         Size: 894.22 GB
         Fault Tolerance: 1
         Heads: 255
         Sectors Per Track: 32
         Cylinders: 65535
         Strip Size: 256 KB
         Full Stripe Size: 256 KB
         Status: OK
         Unrecoverable Media Errors: None
         MultiDomain Status: OK
         Caching:  Disabled
         Unique Identifier: 600508B1001C32B269EB8948F3E5A8E4
         Disk Name: /dev/sdd
         Mount Points: None
         Logical Drive Label: 08F87097PEYHD0CRHB00IZ 67CC
         Mirror Group 1:
            physicaldrive 3I:6:1 (port 3I:box 6:bay 1, SATA SSD, 960 GB, OK)
         Mirror Group 2:
            physicaldrive 3I:6:2 (port 3I:box 6:bay 2, SATA SSD, 960 GB, OK)
         Drive Type: Data
         LD Acceleration Method: Smart Path


      physicaldrive 3I:6:1
         Port: 3I
         Box: 6
         Bay: 1
         Status: OK
         Drive Type: Data Drive
         Interface Type: Solid State SATA
         Size: 960 GB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/4096
         Firmware Revision: 4IYVHPG3
         Serial Number: BTYS825203G8960CGN
         WWID: 31402EC010337950
         Model: ATA     VK000960GWJPF
         SATA NCQ Capable: True
         SATA NCQ Enabled: True
         Current Temperature (C): 33
         Maximum Temperature (C): 52
         Usage remaining: 97.78%
         Power On Hours: 29441
         Estimated Life Remaining based on workload to date: 54030 days
         SSD Smart Trip Wearout: False
         PHY Count: 1
         PHY Transfer Rate: 6.0Gbps
         PHY Physical Link Rate: 6.0Gbps
         PHY Maximum Link Rate: 6.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Freeze Lock Supported: True
         Sanitize Anti-Freeze Lock Supported: True
         Sanitize Lock: None
         Sanitize Estimated Max Erase Time: 3 minute(s), 11 second(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: 1C98C6FF80760CD6

      physicaldrive 3I:6:2
         Port: 3I
         Box: 6
         Bay: 2
         Status: OK
         Drive Type: Data Drive
         Interface Type: Solid State SATA
         Size: 960 GB
         Drive exposed to OS: False
         Logical/Physical Block Size: 512/4096
         Firmware Revision: 4IYVHPG3
         Serial Number: BTYS8253008L960CGN
         WWID: 31402EC010337952
         Model: ATA     VK000960GWJPF
         SATA NCQ Capable: True
         SATA NCQ Enabled: True
         Current Temperature (C): 33
         Maximum Temperature (C): 57
         Usage remaining: 89.44%
         Power On Hours: 29420
         Estimated Life Remaining based on workload to date: 10382 days
         SSD Smart Trip Wearout: False
         PHY Count: 1
         PHY Transfer Rate: 6.0Gbps
         PHY Physical Link Rate: 6.0Gbps
         PHY Maximum Link Rate: 6.0Gbps
         Drive Authentication Status: OK
         Carrier Application Version: 11
         Carrier Bootloader Version: 6
         Sanitize Erase Supported: True
         Sanitize Freeze Lock Supported: True
         Sanitize Anti-Freeze Lock Supported: True
         Sanitize Lock: None
         Sanitize Estimated Max Erase Time: 3 minute(s), 11 second(s)
         Unrestricted Sanitize Supported: True
         Shingled Magnetic Recording Support: None
         Drive Unique ID: A4F747699CC33D09


   SEP (Vendor ID HPE, Model Smart Adapter) 379
      Device Number: 379
      Firmware Version: 1.34
      WWID: 51402EC010337940
      Port: Unknown
      Vendor ID: HPE
      Model: Smart Adapter

`
)

func TestParseSsaVersion(t *testing.T) {
	vers, err := devmon.ParseSsaVersion(ssacliVersion)
	assert.NoError(t, err)
	assert.NotNil(t, vers)
	assert.Equal(t, vers, "4.21.7.0 2020-07-15")
}

func TestParseSsaShowConfig(t *testing.T) {
	cfg, err := devmon.ParseSsaShowConfig(ssacliCtrlAllShowConfigDetail)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, len(cfg.Slots), 1)
	slot := cfg.Slots[0]
	assert.Equal(t, len(slot.Values), 46)
	assert.Equal(t, len(slot.Array), 2)

	arr := slot.Array[0]
	assert.Equal(t, len(arr.Values), 7)
	assert.Equal(t, len(arr.LogicalDrive), 1)
	ld := arr.LogicalDrive[0]
	assert.Equal(t, len(ld.Values), 17)
	assert.Equal(t, len(arr.PhysicalDrive), 1)
	pd := arr.PhysicalDrive[0]
	assert.Equal(t, len(pd.Values), 32)

	arr = slot.Array[1]
	assert.Equal(t, len(arr.Values), 7)
	assert.Equal(t, len(arr.LogicalDrive), 1)
	ld = arr.LogicalDrive[0]
	assert.Equal(t, len(ld.Values), 16)
	assert.Equal(t, len(arr.PhysicalDrive), 1)
	pd = arr.PhysicalDrive[0]
	assert.Equal(t, len(pd.Values), 29)
}

func TestParseConfigToLogical(t *testing.T) {
	cfg, err := devmon.ParseSsaShowConfig(ssacliCtrlAllShowConfigDetail)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	ldm, err := devmon.ParseConfigToLogical(cfg)
	assert.NoError(t, err)
	assert.Equal(t, len(ldm.DevMap), 2)
	ldia := ldm.DevMap["sda"]
	assert.Equal(t, ldia.DiskName, "/dev/sda")
	assert.Equal(t, ldia.Size, "372.58 GB")
	assert.Greater(t, ldia.SizeBytes, uint64(372)*devmon.Giga)
	assert.Greater(t, uint64(373)*devmon.Giga, ldia.SizeBytes)
	ldib := ldm.DevMap["sdb"]
	assert.Equal(t, ldib.DiskName, "/dev/sdb")
	assert.Equal(t, ldib.Size, "931.48 GB")
	assert.Greater(t, ldib.SizeBytes, uint64(931)*devmon.Giga)
	assert.Greater(t, uint64(932)*devmon.Giga, ldib.SizeBytes)

	for _, ldi := range ldm.DevMap {
		for _, pdi := range ldi.PhysicalDrives {
			assert.Greater(t, pdi.SizeBytes, uint64(0))
			assert.Greater(t, pdi.TempCurr, int64(0))
			assert.Greater(t, pdi.TempMaxi, int64(0))
			assert.True(t, (pdi.PowerHours > 0) || (pdi.PowerHours == -1))
		}
	}
}

func TestParseConfigToLogical2(t *testing.T) {
	cfg, err := devmon.ParseSsaShowConfig(ssacliCtrlAllShowConfigDetail2)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	ldm, err := devmon.ParseConfigToLogical(cfg)
	assert.NoError(t, err)
	assert.Equal(t, len(ldm.DevMap), 3)

	ldb := ldm.DevMap["sdb"]
	assert.Equal(t, ldb.UniqueID, "600508B1001CCD7B72DB95E459CDDCCC")
	assert.Equal(t, len(ldb.PhysicalDrives), 4)

	for _, ldi := range ldm.DevMap {
		for _, pdi := range ldi.PhysicalDrives {
			assert.Greater(t, pdi.SizeBytes, uint64(0))
			assert.Greater(t, pdi.TempCurr, int64(0))
			assert.Greater(t, pdi.TempMaxi, int64(0))
			assert.True(t, (pdi.PowerHours > 0) || (pdi.PowerHours == -1))
		}
	}
}
