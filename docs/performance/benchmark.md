# Benchmark

This document explain how to perform a benchmark of approximate nearest neighbor on the Vald cluster, such as [ann-benchmarks](https://github.com/erikbern/ann-benchmarks).

If you want to perform load testing on the Vald cluster, please refer to [this document](loadtest.md).
We also created a guideline of the unit benchmark testing, please refer to [this document](unit_benchmark.md).

## Overview

In this document, we will perform the [ann-benchmarks](https://github.com/erikbern/ann-benchmarks) on Vald cluster.

The test will perform on the whole Vald cluster, which means the request will boardcast to multiple Vald agents instead of single Vald agent, and return the aggreated result.

The result may different from you environment due to the network configuration and overhead of the Vald cluster. In this document we will also explain how to perform the ann-benchmark testing on your environment.

## Benchmark matrix

### Result

Explain the result.

## Benchmarking tools

