#!/bin/bash
./typos > typos_output.txt
cat typos_output.txt | grep "error: " | wc -l
