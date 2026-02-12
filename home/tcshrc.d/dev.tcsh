# Copyright 2026 Samvel Khalatyan. All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

if ( ! $?prompt ) exit

# dev alias starts a tmux session with shell history scoped to the session.
alias dev '\\
  tmux has -t "=$PWD:t:s/./_/" >& /dev/null || \\
        tmux new -d -c "$PWD" -s "$PWD:t:s/./_/" && \\
        tmux setenv -t "=$PWD:t" HISTFILE "$PWD"/.history && \\
        tmux respawnw -k -t "=$PWD:t":0 ; \\
  if ($?TMUX) \\
    tmux switchc -t "=$PWD:t:s/./_/" ; \\
  if ! ($?TMUX) \\
    tmux attach -t "=$PWD:t:s/./_/" \\
    '
