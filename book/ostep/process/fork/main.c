// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fork shows how to fork a child process.

#include <err.h>
#include <stdio.h>
#include <sys/types.h>
#include <unistd.h>

int main() {
  pid_t pid;
  printf("[%d]: prepare for fork\n", getpid());
  switch (pid = fork()) {
    case -1:  // error
      err(1, "parent pid %d", getpid());
      break;
    case 0:  // child
      printf("[%d]: in child\n", getpid());
      break;
    default:  // parent
      printf("[%d]: forked a child [%d]\n", getpid(), pid);
      break;
  }
  return 0;
}
