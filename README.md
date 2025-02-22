# gokapiX

An implementation of BM25 variants surveyed in *"Improvements to BM25 and Language Models Examined"* by Trotman et al., 2014. <br/>
Paper Link: https://www.cs.otago.ac.nz/homepages/andrew/papers/2014-2.pdf

## Overview

**gokapiX** is a ranking algorithm designed to implement multiple BM25 variants, improving upon traditional BM25 scoring for information retrieval. This implementation is inspired by [rank_bm25](https://github.com/dorianbrown/rank_bm25).

Name Credits: [gokapi](https://github.com/raphaelsty/gokapi)

Algorithms implemented:
- Atire BM25
- BM25L
- BM25+
- BM25-Adpt
- BM25T

## Installation
To install **gokapiX**, use:

```sh
go get github.com/niranjanorkat/gokapiX
