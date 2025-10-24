// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/singly/cycle/is_happy_number.h"

namespace iq::list::singly::cycle {
namespace {

int getNextNumber(int n) {
  int next = 0;
  for (; n != 0; n /= 10) {
    const int k = n % 10;
    if (k == 0) {
      continue;
    }
    next += k * k;
  }
  return next;
}

}  // namespace

bool IsHappyNumber(int n) {
  int slow = n;
  int fast = getNextNumber(n);
  while (fast != 1) {
    if (slow == fast) {
      return false;
    }
    slow = getNextNumber(slow);
    fast = getNextNumber(fast);
    fast = getNextNumber(fast);
  }
  return true;
}

}  // namespace iq::list::singly::cycle
