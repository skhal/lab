# NAME

**padding** - the number of bites to be added to the block size


# DESCRIPTION

Consider the following code from Go archive/tar:

https://github.com/golang/go/blob/177cd8d7633843c3eabf8f887381cadbeed3514b/src/archive/tar/format.go#L152-L156

## Problem

It states that if the block size is 2^N and offset is any number between zero
(inclusive) and the block size, then the expression `-offset & (blockSize - 1)`
gives the number of bytes between the block size and the offset, that is:

```
blockSize - offset = -offset & (blockSize - 1)
```

## Proof

The blockSize is:

```math
blockSize = 2^N
```

Then:

```math
blockSize - 1 = \sum_{i=0}^{N-1}{2^i}
```

Offset is some number that is less than the block size. When represented as a
polynomial of the powers of 2, it will use only N-1 powers.

Let's denote $`\nu_i`$ the i-th power of 2 factor, which is either 0 or 1. The
offset then can be represented as:

```math
offset = \sum_{i=0}^{N-1}{\nu_i 2^i}
```

There are different ways to represent negative numbers in computers
([Wikipedia](https://en.wikipedia.org/wiki/Signed_number_representations) with
two'2 complement being dominant these days. Two's complement of a number is
given by the complement of the number plus one:

```math
-N => ~N + 1
```

Applying this idea to the offset, it translates into:

```math
\begin{align}
-offset &= ~\left( \sum_{i=0}^{N-1}{\nu_i 2^i} \right) + 1 \\
        &= \sum_{i=0}^{N-1}{\bar{\nu}_i 2^i} + 2^0 \\
\end{align}
```

where $`\bar{\nu}_i`$ is the complement of $`\nu_i`$.

Final expression:

```math
\begin{align}
-offset \& (blockSize - 1) &= \left( \sum_{i=0}^{N-1}{\bar{\nu}_i 2^i} + 2^0 \right) & \sum_{j=0}^{N-1}{2^j} \\
  &= \sum_{i=0}^{N-1}{\bar{\nu}_i 2^i} + 2^0
\end{align}
```

Verify the result by adding it to the offset - we expect to get the block size:

```math
\begin{align}
offset + [ -offset & ( blockSize - 1 ) ] &= \sum_{i=0}^{N-1}{\nu_i 2^i} + \sum_{i=0}^{N-1}{\bar{\nu}_i 2^i} + 2^0 \\
  &= \sum_{i=0}^{N-1}{2^i} + 2^0 \\
  &= 2^N
\end{align}
```
