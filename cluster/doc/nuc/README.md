Name
====

**nuc** - ASUS NUC

Hardware
========

-	ASUS NUC 15 Pro+ NUC15CRSU7, [Support](https://www.asus.com/us/supportonly/nuc15crsu9/helpdesk_bios/)
	-	Intel Core Ultra 7 [255H](https://www.intel.com/content/www/us/en/products/sku/241751/intel-core-ultra-7-processor-255h-24m-cache-up-to-5-10-ghz/specifications.html)
	-	G.SKILL Ripjaws SO-DIMM 32GB (16GB x2) DDR5 5600
	-	WD Black SN8100 2TB PCIe 5.0x4 M.2 2280 NVMe

BIOS
====

-	Set date and time in UTC
-	Disable "Secure Boot" to boot non-Windows OS, [ref](https://www.asus.com/global/support/faq/1044664/).

Update
------

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
