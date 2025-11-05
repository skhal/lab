Name
====

**36** - evaluate calculator expression

Description
===========

Problem
-------

Evaluate a calculator expression.

Assumptions:

-	the numbers are non-negative integers
-	operators are plus and minus
-	parentheses are allowed

Example
-------

Expression "1 + (3 - 2)" results in 2.

Solution
--------

<details>
<summary>Details<summary>

- Parse the expression by extract tokens: numbers, operators, parentheses, etc.
- Use stack to keep track for numbers and operators.
- Evaluate operators as soon as the right operand becomes available.
- Be careful with handing parenthesis and error checking to validate the
  expression.

Complexity:

- _time_: O(n)
- _space_: O(n)

</details>
