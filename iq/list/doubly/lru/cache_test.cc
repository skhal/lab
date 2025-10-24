// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
// clang-format off-next-line
//go:build ignore

#include "iq/list/doubly/lru/cache.h"

#include <algorithm>
#include <cctype>
#include <cstddef>
#include <optional>
#include <ostream>
#include <string>
#include <utility>
#include <vector>

#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"

namespace iq::list::doubly::lru {
namespace {

using ::testing::ContainerEq;
using ::testing::Eq;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

TEST(CacheDeathTest, FailsOnZeroCapacity) { EXPECT_DEATH(Cache cache(0), ""); }

struct CachePutTestParam {
  std::string name;
  std::size_t cap;
  std::vector<std::pair<Cache::Key, Cache::Value>> items;
  std::vector<Cache::Key> want;

  friend std::ostream& operator<<(std::ostream& os,
                                  const CachePutTestParam& tp) {
    return os << "cap: " << tp.cap << " items: {"
              << absl::StrJoin(tp.items, ", ", absl::PairFormatter(":")) << "}";
  }
};

class CachePutTest : public TestWithParam<CachePutTestParam> {};

TEST_P(CachePutTest, Put) {
  const CachePutTestParam& tp = GetParam();
  Cache cache(tp.cap);

  for (const auto& [key, value] : tp.items) {
    cache.Put(key, value);
  }

  EXPECT_THAT(cache.Keys(), ContainerEq(tp.want));
}

const CachePutTestParam kCachePutTestParams[]{
    {.name = "cap one fill", .cap = 1, .items = {{1, 10}}, .want = {1}},
    {
        .name = "cap one evict least recent",
        .cap = 1,
        .items = {{1, 10}, {2, 20}},
        .want = {2},
    },
    {.name = "cap two put one", .cap = 2, .items = {{1, 10}}, .want = {1}},
    {
        .name = "cap two fill",
        .cap = 2,
        .items = {{1, 10}, {2, 20}},
        .want = {2, 1},
    },
    {
        .name = "cap two evict least recent",
        .cap = 2,
        .items = {{1, 10}, {2, 20}, {3, 30}},
        .want = {3, 2},
    },
    {.name = "cap three put one", .cap = 3, .items = {{1, 10}}, .want = {1}},
    {
        .name = "cap three put two",
        .cap = 3,
        .items = {{1, 10}, {2, 20}},
        .want = {2, 1},
    },
    {
        .name = "cap three fill",
        .cap = 3,
        .items = {{1, 10}, {2, 20}, {3, 30}},
        .want = {3, 2, 1},
    },
    {
        .name = "cap three evict least recent",
        .cap = 3,
        .items = {{1, 10}, {2, 20}, {3, 30}, {4, 40}},
        .want = {4, 3, 2},
    },
};

INSTANTIATE_TEST_SUITE_P(CacheTest, CachePutTest, ValuesIn(kCachePutTestParams),
                         [](const TestParamInfo<CachePutTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); }, '_');
                           return name;
                         });

struct CacheGetTestParam {
  std::string name;
  std::size_t cap;
  std::vector<std::pair<Cache::Key, Cache::Value>> items;
  Cache::Key key;
  std::optional<Cache::Key> wantVal;
  std::vector<Cache::Key> wantKeys;

  friend std::ostream& operator<<(std::ostream& os,
                                  const CacheGetTestParam& tp) {
    return os << "cap: " << tp.cap << " items: {"
              << absl::StrJoin(tp.items, ", ", absl::PairFormatter(":")) << "}";
  }
};

class CacheGetTest : public TestWithParam<CacheGetTestParam> {};

TEST_P(CacheGetTest, Get) {
  const CacheGetTestParam& tp = GetParam();
  Cache cache(tp.cap);
  for (const auto& [key, val] : tp.items) {
    cache.Put(key, val);
  }

  const std::optional<Cache::Value> got = cache.Get(tp.key);

  EXPECT_THAT(got, Eq(tp.wantVal));
  EXPECT_THAT(cache.Keys(), ContainerEq(tp.wantKeys));
}

const CacheGetTestParam kCacheGetTestParams[]{
    // capacity=1
    {
        .name = "cap one get hit",
        .cap = 1,
        .items = {{1, 10}},
        .key = 1,
        .wantVal = 10,
        .wantKeys = {1},
    },
    {
        .name = "cap one get miss",
        .cap = 1,
        .items = {{1, 10}},
        .key = 2,
        .wantVal = std::nullopt,
        .wantKeys = {1},
    },
    // capacity=2
    {
        .name = "cap two one item get hit",
        .cap = 2,
        .items = {{1, 10}},
        .key = 1,
        .wantVal = 10,
        .wantKeys = {1},
    },
    {
        .name = "cap two one item get miss",
        .cap = 2,
        .items = {{1, 10}},
        .key = 2,
        .wantVal = std::nullopt,
        .wantKeys = {1},
    },
    {
        .name = "cap two two items get least recent makes it most recent",
        .cap = 2,
        .items = {{1, 10}, {2, 20}},
        .key = 1,
        .wantVal = 10,
        .wantKeys = {1, 2},
    },
    {
        .name = "cap two two items get most recent",
        .cap = 2,
        .items = {{1, 10}, {2, 20}},
        .key = 2,
        .wantVal = 20,
        .wantKeys = {2, 1},
    },
    {
        .name = "cap two two items get miss",
        .cap = 2,
        .items = {{1, 10}, {2, 20}},
        .key = 3,
        .wantVal = std::nullopt,
        .wantKeys = {2, 1},
    },
    // capacity=3
    {
        .name = "cap three one item get hit",
        .cap = 3,
        .items = {{1, 10}},
        .key = 1,
        .wantVal = 10,
        .wantKeys = {1},
    },
    {
        .name = "cap three one item get miss",
        .cap = 3,
        .items = {{1, 10}},
        .key = 2,
        .wantVal = std::nullopt,
        .wantKeys = {1},
    },
    {
        .name = "cap three two items get least recent makes it most recent",
        .cap = 3,
        .items = {{1, 10}, {2, 20}},
        .key = 1,
        .wantVal = 10,
        .wantKeys = {1, 2},
    },
    {
        .name = "cap three two items get most recent",
        .cap = 3,
        .items = {{1, 10}, {2, 20}},
        .key = 2,
        .wantVal = 20,
        .wantKeys = {2, 1},
    },
    {
        .name = "cap three two items get miss",
        .cap = 3,
        .items = {{1, 10}, {2, 20}},
        .key = 3,
        .wantVal = std::nullopt,
        .wantKeys = {2, 1},
    },
    {
        .name = "cap three three items get least recent makes it most recent",
        .cap = 3,
        .items = {{1, 10}, {2, 20}, {3, 30}},
        .key = 1,
        .wantVal = 10,
        .wantKeys = {1, 3, 2},
    },
    {
        .name = "cap three three items get middle item",
        .cap = 3,
        .items = {{1, 10}, {2, 20}, {3, 30}},
        .key = 2,
        .wantVal = 20,
        .wantKeys = {2, 3, 1},
    },
    {
        .name = "cap three three items get most recent",
        .cap = 3,
        .items = {{1, 10}, {2, 20}, {3, 30}},
        .key = 3,
        .wantVal = 30,
        .wantKeys = {3, 2, 1},
    },
    {
        .name = "cap three three items get miss",
        .cap = 3,
        .items = {{1, 10}, {2, 20}, {3, 30}},
        .key = 4,
        .wantVal = std::nullopt,
        .wantKeys = {3, 2, 1},
    },
};

INSTANTIATE_TEST_SUITE_P(CacheTest, CacheGetTest, ValuesIn(kCacheGetTestParams),
                         [](const TestParamInfo<CacheGetTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); }, '_');
                           return name;
                         });

}  // namespace
}  // namespace iq::list::doubly::lru
