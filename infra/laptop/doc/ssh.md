<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

NAME
====

**ssh** -- setup SSH client

DESCRIPTION
===========

Outline:

-	Generate SSH key for `lab.net` cluster
-	Configure SSH client to auto-add keys to SSH agent and forward the agent.
-	Copy the key to the remote hosts

	```console
	% ssh-keygen -C "sn.khalatyan@gmail.com" -f ~/.ssh/id_lab -t ed25519
	% cat <<EOF >~/.ssh/config
	Host *.lab.net
	AddKeysToAgent yes
	ForwardAgent yes
	IdentityFile ~/.ssh/id_lab
	EOF
	```

Distribute the public SSH key to the remote hosts:

```console
% ssh-copy-id -i ~/.ssh/id_lab.pub dev.nuc.lab.net
/usr/bin/ssh-copy-id: INFO: Source of key(s) to be installed: "/Users/skhalatyan/.ssh/id_dev.nuc.lab.net.pub"
/usr/bin/ssh-copy-id: INFO: attempting to log in with the new key(s), to filter out any that are already installed
/usr/bin/ssh-copy-id: INFO: 1 key(s) remain to be installed -- if you are prompted now it is to install the new keys
(skhalatyan@dev.nuc.lab.net) Password for skhalatyan@dev.nuc.lab.net:

Number of key(s) added:        1

Now try logging into the machine, with: "ssh -i /Users/skhalatyan/.ssh/id_dev.nuc.lab.net 'dev.nuc.lab.net'"
and check to make sure that only the key(s) you wanted were added.
```

Verity it works:

```console
% ssh dev.nuc.lab.net hostname
Enter passphrase for key '/Users/skhalatyan/.ssh/id_lab':
dev.nuc.lab.net
```
