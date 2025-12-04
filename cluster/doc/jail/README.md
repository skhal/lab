Name
====

**jail** - setup jails

Summary
-------

This section describes how to setup and manage thick jails using base tools, jail(8) and jexec(8).

First, [bootstrap](./bootstrap.md) the OS to host jails. It creates ZFS datasets, sets up networking stack, and configures jail service.

In order to simplify maintenance, jails split into a template and a running container. This setup makes it very quick to create a new jail using a ZFS clone of a template.

The instructions distinguish two types kinds of jails:

-	FreeBSD: hosts FreeBSD userland suitable to run FreeBSD binaries, create an isolated environment to mimic a virtual machine, etc.

	-	[FreeBSD template](./freebsd_template.md)
	-	[FreeBSD jail](./freebsd_jail.md)

-	Linux: a FreeBSD jail with Linux userland bootstrapped at `/compat/<distribution>`. Even though it combines FreeBSD and Linux environment, Linux jails are designed to be primarily used as a Linux environment chroot'ed at `/compat/<distribution>`.

	-	[Ubuntu template](./ubuntu_template.md)
	-	[Ubuntu jail](./ubuntu_jail.md)
