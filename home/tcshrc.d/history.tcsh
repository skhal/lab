# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

set history = 10000
set savehist = (10000 merge)
# Isolate Linux VMs history from BSD (default)
if ( "$OSTYPE" == "linux" ) then
  set histfile = "$HOME"/.history."$OSTYPE"
endif
