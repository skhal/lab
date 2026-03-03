// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Wait shows how to wait for a forked process to complete.
//
// Output:
// 33231 start
// 33383 child
// 33383 done
// 33231 child done 33383
// 33231 done

#include <err.h>
#include <stddef.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sysexits.h>
#include <unistd.h>

int main() {
  pid_t pid, wpid;
  printf("%d start\n", getpid());
  switch (pid = fork()) {
    case -1:  // error
      err(EX_OSERR, "fork %d\n", getpid());
      break;
    case 0:  // child
      printf("%d child\n", getpid());
      break;
    default:              // parent
                          // Wait for child to change state, typically done.
      wpid = wait(NULL);  // NULL - do not want child status info.
      if (wpid == -1) {
        err(EX_OSERR, "wait %d\n", getpid());
      }
      printf("%d child done %d [wait() = %d]\n", getpid(), pid, wpid);
  }
  printf("%d done\n", getpid());
  return 0;
}
