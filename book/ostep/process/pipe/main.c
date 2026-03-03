// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Pipe shows how to change the environment for the child process after fork(2)
// but before running exec(3). This example uses pipe(2) to create IPC between
// date(1) and tr(1) command utilities using two processes:
//
// - the child process runs date(1) to output the result to STDOUT
//
// - the parent process runs tr(1) to translate lower-case alphabet to the
//   upper case from STDIN.
//
// Output:
// TUE MAR  3 11:40:58 CST 2026

#include <err.h>
#include <fcntl.h>
#include <stddef.h>
#include <string.h>
#include <sys/types.h>
#include <sysexits.h>
#include <unistd.h>

int IDX_PIPE_R = 0;
int IDX_PIPE_W = 1;

const char* PIPE_W_CMD = "/bin/date";
const char* PIPE_R_CMD = "/usr/bin/tr";

void pipe_write(int pipefd[2]) {
  char* argv[2];
  if (dup2(pipefd[IDX_PIPE_W], STDOUT_FILENO) == -1) {
    err(EX_OSERR, "dup2\n");
  }
  argv[0] = strdup(PIPE_W_CMD);
  argv[1] = NULL;
  if (execv(argv[0], argv) == -1) {
    err(EX_OSERR, "exec\n");
  }
}

void pipe_read(int pipefd[2]) {
  char* argv[4];
  if (dup2(pipefd[IDX_PIPE_R], STDIN_FILENO) == -1) {
    err(EX_OSERR, "dup pipe-r to STDIN\n");
  }
  argv[0] = strdup(PIPE_R_CMD);
  argv[1] = strdup("[:lower:]");
  argv[2] = strdup("[:upper:]");
  argv[3] = NULL;
  if (execv(argv[0], argv) == -1) {
    err(EX_OSERR, "exec\n");
  }
}

int main() {
  int pipefd[2];
  if (pipe2(pipefd, O_CLOEXEC) == -1) {
    err(EX_OSERR, "pipe\n");
  }
  switch (fork()) {
    case -1:
      err(EX_OSERR, "fork %d\n", getpid());
    case 0:  // child
      // write from child process to let it exit upon completion and close the
      // write end of the pipe.
      pipe_write(pipefd);
    default:  // parent
      pipe_read(pipefd);
  };
  return 0;
}
