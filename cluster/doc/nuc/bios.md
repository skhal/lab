Name
====

**bios** - Basic Input/Output System setup

Configuration
=============

-	Set date and time in UTC
-	Disable "Secure Boot" to boot non-Windows OS, [ref](https://www.asus.com/global/support/faq/1044664/).
-	Disable PCIe ASPM Support, [ref](https://superuser.com/questions/1822809/why-does-disabling-active-power-management-in-bios-double-nvme-speed)

Update
======

The update runs for about 5 minutes.

1.	Format a USB drive as FAT16
2.	Download "BIOS Full Package" archive (`.zip`) from the [Support](https://www.asus.com/us/supportonly/nuc15crsu7/helpdesk_bios/) page. Verify the checksum:

	```console
	% # generate a checksum file
	% sha256 ./CRARL579.0027.zip > bios.checksum
	% # update the checksum in bios.checksum
	% # verify the checksum
	% shasum -c bios.checksum --ignore-missing
	./CRARL579.0027.zip: OK
	```

3.	Copy capsule file (`.cap`) from "Capsule File for BIOS Flash through F7" folder to the USB drive.

4.	Start ASUS NUC with the USB drive in. Keep F7 pressed to start BIOS update.
