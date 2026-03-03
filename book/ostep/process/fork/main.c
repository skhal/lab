// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fork shows how to start a new process.
//
// Output (non-deterministic order after the first line):
// 10292 prepare for fork
// 10292 forked a child 10432
// 10292 done
// 10432 in child
// 10432 done

#include <err.h>
#include <stdio.h>
#include <sys/types.h>
#include <sysexits.h>
#include <unistd.h>

int main() {
  pid_t pid;
  printf("%d prepare for fork\n", getpid());
  switch (pid = fork()) {
    case -1:  // error
      err(EX_OSERR, "fork %d", getpid());
      break;
    case 0:  // child
      printf("%d in child\n", getpid());
      break;
    default:  // parent
      printf("%d forked a child %d\n", getpid(), pid);
      break;
  }
  printf("%d done\n", getpid());
  return 0;
}
