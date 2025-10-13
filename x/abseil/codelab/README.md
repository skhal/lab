# NAME

**abseil-codelab** - tutorial on Abseil integration with Bazel


# DESCRIPTION

This is pretty much reproduced [tutorial](https://abseil.io/docs/cpp/quickstart)
from Abseil.

```console
% bazel run //:hello_world
INFO: Analyzed target //:hello_world (0 packages loaded, 0 targets configured).
INFO: Found 1 target...
Target //:hello_world up-to-date:
  bazel-bin/hello_world
INFO: Elapsed time: 0.085s, Critical Path: 0.00s
INFO: 1 process: 1 action cache hit, 1 internal.
INFO: Build completed successfully, 1 total action
INFO: Running command line: bazel-bin/hello_world
Joined string: foo-bar-baz
```


# SEE ALSO

* https://abseil.io/docs/cpp/quickstart
