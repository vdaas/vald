# tikv

This package is from [pingcap](https://github.com/pingcap/kvproto) repository and processed in the steps below:

1. Copy only required files
2. Remove unnecessary `import` and `option` statements
3. Remove redundant messages and services
4. Comment out unused message fields (Avoid deletion to reduce future merge conflicts)
5. Adjust package names and paths
6. Change metapb.Region to metapb.Region2 to avoid conflicts.
