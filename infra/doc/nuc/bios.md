<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# Name

**bios** - Basic Input/Output System setup

# Status

Installed: v0028

# Configuration

- Set date and time in UTC
- Disable "Secure Boot" to boot non-Windows OS, [ref](https://www.asus.com/global/support/faq/1044664/).
- Disable PCIe ASPM Support, [ref](https://superuser.com/questions/1822809/why-does-disabling-active-power-management-in-bios-double-nvme-speed)
- Disable Bluetooth

# Upgrade

The upgrade runs for about 5 minutes. Some versions of BIOS (579.0030) may
reject to upgrade components like thunderbolt when the monitor is connected
via USB Type-C - connect the monitor on HDMI port.

1. Format a USB drive as FAT32
2. Download "BIOS Full Package" archive (`.zip`) from the [Support](https://www.asus.com/us/supportonly/nuc15crsu7/helpdesk_bios/) page.

Verify the checksum:

```console
$ sha256 ./CRARL579.0030.zip > /tmp/bios.checksum

# update the checksum in /tmp/bios.checksum

$ shasum -c /tmp/bios.checksum --ignore-missing
./CRARL579.0030.zip: OK
```

3. Copy capsule file (`.cap`) from "Capsule File for BIOS Flash through F7" folder to the USB drive.

4. Start ASUS NUC with the USB drive in. Keep F7 pressed to start BIOS update.
