#!/bin/bash
git ls-remote https://github.com/actions/checkout.git refs/tags/v2 | awk '{print $1}'
git ls-remote https://github.com/crate-ci/typos.git HEAD | awk '{print $1}'
git ls-remote https://github.com/fit2cloud/LLM-CodeReview-Action.git HEAD | awk '{print $1}'
