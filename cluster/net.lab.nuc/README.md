NAME
====

**nuc.lab.net** -- Intel NUC server with FreeBSD

DESCRIPTION
===========

`nuc.lab.net` is a FreeBSD server on [Intel NUC Kit NUC6i7KYK](https://www.intel.com/content/www/us/en/products/sku/89187/intel-nuc-kit-nuc6i7kyk/specifications.html).

Hardware
--------

-	Processor: [Intel Core i7-677HQ](https://www.intel.com/content/www/us/en/products/sku/93341/intel-core-i76770hq-processor-6m-cache-up-to-3-50-ghz/specifications.html)

-	RAM: 2x 16 GB DDR4-2400 16GB HyperX HX424S14IBK2/32 kit

-	SSD: 2x 1TB NVMe WD WDS100T3X0C Black SN 750 Gen3 PCIe M.2 2280

Sync
----

```console
% rsync -arvz --files-from=./rsync.files-from op@nuc.lab.net:/ ./
```
