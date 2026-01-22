// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Wait shows how to wait for a forked child process.

#include <err.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

int main() {
  pid_t pid;
  printf("[%d]: start\n", getpid());
  switch (pid = fork()) {
    case -1:  // error
      err(1, "in pid %d\n", getpid());
      break;
    case 0:  // child
      printf("[%d]: child\n", getpid());
      break;
    default:       // parent
      wait(NULL);  // NULL for status - don't care
      printf("[%d]: child %d completed\n", getpid(), pid);
  }
  return 0;
}
