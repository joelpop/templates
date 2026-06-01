# Coverage Gaps Surface ACs Without Tests

At the pre-PR gate, ask "which ACs have a passing test?" rather than "what's the line-coverage percentage?" so gaps are visible — coverage of a line under an unrelated test does not fill a missing AC test.

At the pre-PR gate, the Unit Tester asks "which ACs have a passing test?" — not
"what's the line-coverage percentage?". A line of code without a test for the AC
it implements is a hole; coverage of that line under an unrelated test does not
fill it.

If an AC has no test, the Unit Tester reports the gap to the Coder and
Architect. Closing it may mean (a) writing the missing test, (b) recognizing the
AC was never implemented and implementing it, or (c) recognizing the AC was
misstated and revising it. All three are legitimate; shipping without one of them
is not.
