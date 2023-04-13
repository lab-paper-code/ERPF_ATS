#!/bin/sh
img_num=$1

strace -o strace_outputs/strace${img_num}.log -r -C -e trace=file,read,write python3 image_classification.py -n ${img_num}
