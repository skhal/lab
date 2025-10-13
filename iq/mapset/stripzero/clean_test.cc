// Copyright 2025 Samvel Khalatyan. All rights reserved.
//
//go:build ignore

#include <cctype>
#include <ostream>
#include <string>
#include <string_view>

#include "absl/strings/str_format.h"
#include "absl/strings/str_join.h"
#include "gmock/gmock.h"
#include "gtest/gtest.h"
#include "iq/mapset/stripzero/clean.h"

namespace iq::mapset::stripzero {

class RowFormatter {
 public:
  RowFormatter(std::string_view prefix)
      : is_first_row_(true), prefix_(prefix) {}

  void operator()(std::string* s, const Row& r) {
    if (is_first_row_) {
      is_first_row_ = false;
      s->append(absl::StrFormat("%s", absl::StrJoin(r, " ")));
      return;
    }
    s->append(absl::StrFormat("%s%s", prefix_, absl::StrJoin(r, " ")));
  }

 private:
  bool is_first_row_;
  const std::string prefix_;
};

std::ostream& operator<<(std::ostream& os, [[maybe_unused]] const Matrix& m) {
  return os << "[" << absl::StrJoin(m, "\n", RowFormatter(" ")) << "]";
}

namespace {

using ::testing::Eq;
using ::testing::Pointwise;
using ::testing::TestParamInfo;
using ::testing::TestWithParam;
using ::testing::ValuesIn;

constexpr char kCharUnderscore = '_';

struct CleanTestParam {
  std::string name;
  Matrix matrix;
  Matrix want;

  friend std::ostream& operator<<(std::ostream& os, const CleanTestParam& tp) {
    return os << "matrix:\n" << tp.matrix;
  }
};

class CleanTest : public TestWithParam<CleanTestParam> {};

TEST_P(CleanTest, Test) {
  const CleanTestParam& tp = GetParam();
  Matrix m{tp.matrix};

  Clean(m);

  EXPECT_THAT(m, Pointwise(Eq(), tp.want));
}

const CleanTestParam kCleanTestParams[]{
    {.name = "empty", .matrix = {}, .want = {}},
    // 1x1
    {.name = "one by one non zero",
     .matrix =
         {
             {1},
         },
     .want =
         {
             {1},
         }},
    {.name = "one by one zero",
     .matrix =
         {
             {0},
         },
     .want =
         {
             {0},
         }},
    // row
    {
        .name = "row non zero",
        .matrix =
            {
                {1, 2, 3},
            },
        .want =
            {
                {1, 2, 3},
            },
    },
    {
        .name = "row with zero",
        .matrix =
            {
                {1, 2, 0},
            },
        .want =
            {
                {0, 0, 0},
            },
    },
    // column
    {
        .name = "col non zero",
        .matrix =
            {
                {1},
                {2},
                {3},
            },
        .want =
            {
                {1},
                {2},
                {3},
            },
    },
    {
        .name = "col with zero",
        .matrix =
            {
                {1},
                {0},
                {3},
            },
        .want =
            {
                {0},
                {0},
                {0},
            },
    },
    // m-by-n matrix
    {
        .name = "matrix non zero",
        .matrix =
            {
                {1, 1, 1, 1},
                {1, 1, 1, 1},
                {1, 1, 1, 1},
            },
        .want =
            {
                {1, 1, 1, 1},
                {1, 1, 1, 1},
                {1, 1, 1, 1},
            },
    },
    {
        .name = "matrix one zero",
        .matrix =
            {
                {1, 1, 1, 1},
                {1, 0, 1, 1},
                {1, 1, 1, 1},
            },
        .want =
            {
                {1, 0, 1, 1},
                {0, 0, 0, 0},
                {1, 0, 1, 1},
            },
    },
    {
        .name = "matrix two zero",
        .matrix =
            {
                {1, 0, 1, 1},
                {1, 1, 1, 1},
                {1, 1, 1, 0},
            },
        .want =
            {
                {0, 0, 0, 0},
                {1, 0, 1, 0},
                {0, 0, 0, 0},
            },
    },
    {
        .name = "matrix two zero in row",
        .matrix =
            {
                {1, 0, 1, 1},
                {1, 1, 1, 1},
                {1, 1, 1, 1},
            },
        .want =
            {
                {0, 0, 0, 0},
                {1, 0, 1, 1},
                {1, 0, 1, 1},
            },
    },
    {
        .name = "matrix two two zeros in row",
        .matrix =
            {
                {1, 0, 0, 1},
                {1, 1, 1, 1},
                {1, 1, 1, 1},
            },
        .want =
            {
                {0, 0, 0, 0},
                {1, 0, 0, 1},
                {1, 0, 0, 1},
            },
    },
    {
        .name = "matrix two zero in col",
        .matrix =
            {
                {1, 1, 1, 1},
                {1, 1, 1, 1},
                {0, 1, 1, 1},
            },
        .want =
            {
                {0, 1, 1, 1},
                {0, 1, 1, 1},
                {0, 0, 0, 0},
            },
    },
    {
        .name = "matrix two two zeros in col",
        .matrix =
            {
                {1, 1, 1, 1},
                {0, 1, 1, 1},
                {0, 1, 1, 1},
            },
        .want =
            {
                {0, 1, 1, 1},
                {0, 0, 0, 0},
                {0, 0, 0, 0},
            },
    },
};

INSTANTIATE_TEST_SUITE_P(CleanTest, CleanTest, ValuesIn(kCleanTestParams),
                         [](const TestParamInfo<CleanTestParam>& info) {
                           std::string name = info.param.name;
                           std::replace_if(
                               name.begin(), name.end(),
                               [](char c) { return !std::isalnum(c); },
                               kCharUnderscore);
                           return name;
                         });

}  // namespace
}  // namespace iq::mapset::stripzero
