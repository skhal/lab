// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Exec shows how to run a program from a process.
//
// Output:
// 81980 start
// 81980 forked a child 82167
// 82167 child
// Tue Mar  3 10:10:31 CST 2026
// 81980 child done 82167 [wait() = 82167]
// 81980 done

#include <err.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sysexits.h>
#include <unistd.h>

char COMMAND[] = "/bin/date";

int main() {
  pid_t pid, wpid;
  printf("%d start\n", getpid());
  switch (pid = fork()) {
    case -1:  // fail
      err(EX_OSERR, "fork %d", getpid());
      break;
    case 0: {  // child
      printf("%d child\n", getpid());
      char* argv[2];
      argv[0] = strdup(COMMAND);
      argv[1] = NULL;
      if (execv(argv[0], argv) == -1) {
        err(EX_OSERR, "%d exec %s", getpid(), argv[0]);
      }
      break;
    }
    default:  // parent
      printf("%d forked a child %d\n", getpid(), pid);
      if ((wpid = wait(NULL)) == -1) {
        err(EX_OSERR, "wait %d\n", getpid());
      }
      printf("%d child done %d [wait() = %d]\n", getpid(), pid, wpid);
      break;
  }
  printf("%d done\n", getpid());
  return 0;
}
