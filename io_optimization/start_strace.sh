#!/bin/sh
image_num=$1
batch_size=$2
epoch=$3
output_file_path=$4

strace -o ${output_file_path} -r -C -f -e trace=openat,lseek,read,write python3 image_training.py ${image_num} ${batch_size} ${epoch}
