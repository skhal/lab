// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <err.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

int main() {
  pid_t pid;
  printf("[%d] start\n", getpid());
  switch (pid = fork()) {
    case -1:  // fail
      err(1, "parent %d", getpid());
      break;
    case 0: {  // child
      char* argv[2];
      argv[0] = strdup("/bin/date");
      argv[1] = NULL;
      if (execv(argv[0], argv) == -1) {
        err(1, "[%d]: child: exec %s", getpid(), argv[0]);
        exit(1);
      }
      break;
    }
    default:  // parent
      printf("[%d] forked child [%d]\n", getpid(), pid);
      pid_t wpid = wait(NULL);
      printf("[%d] wait: rc %d\n", getpid(), wpid);
      break;
  }
  return 0;
}
