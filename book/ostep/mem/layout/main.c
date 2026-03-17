// Copyright 2026 Samvel Khalatyan. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Layout demonstrates memory layout of core blocks: code, stack, heap.
// Traditionally, stack goes first. The stack and heap grow from the opposite
// ends to guarantee dynamic nature of the two (it is impossible to predict how
// large each is going to be).
//
// SYNOPSIS
//  layout
//
// OUTPUT
//  code:   main() at 0x201670
//  stack:  argc at 0x82080efb8
//  heap:   malloc() at 0x361343812000

#include <stdio.h>
#include <stdlib.h>

int main(int argc, char* _[]) {
  printf("code:\tmain() at %p\n", main);
  printf("stack:\targc at %p\n", &argc);
  void* p = malloc(1);
  printf("heap:\tmalloc() at %p\n", p);
  free(p);
  return 0;
}
