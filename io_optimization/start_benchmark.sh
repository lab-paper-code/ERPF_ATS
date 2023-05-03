#!/bin/sh

n=100
ssd_path=$1
hdd_path=$2
batch_size=${3:-32}
epochs=${4:-10}

. ./env/bin/activate

dataset_size=$(find $ssd_path -type f | wc -l)
echo "Current setting ========================="
echo "\tDataset path in SSD: $ssd_path"
echo "\tDataset path in HDD: $hdd_path"
echo "\tTotal dataset size: $dataset_size"
echo "\tBatch size: $batch_size"
echo "\tEpochs: $epochs"
echo "========================================="
echo "Total I/O: $(expr $dataset_size \* $epochs)"
echo "========================================="

sdd_execution_time=0
for i in $(seq 1 $n); do
    start=$(date +%s.%N)
    python3.8 target_scripts/image_fine-tuning.py $ssd_path $dataset_size $batch_size $epochs > /dev/null
    end=$(date +%s.%N)
    runtime=$(echo "$end - $start" | bc)
    sdd_execution_time=$(echo "$sdd_execution_time + $runtime" | bc)
done
sdd_execution_time_avg=$(echo "$sdd_execution_time / $n" | bc)
echo "Average execution time (SSD): $sdd_execution_time"

hdd_execution_time=0
hdparm -A0 /dev/sda
for i in $(seq 1 $n); do
    start=$(date +%s.%N)
    python3.8 target_scripts/image_fine-tuning.py $hdd_path $dataset_size $batch_size $epochs  > /dev/null
    end=$(date +%s.%N)
    runtime=$(echo "$end - $start" | bc)
    hdd_execution_time=$(echo "$hdd_execution_time + $runtime" | bc)
done
hdparm -A1 /dev/sda
hdd_execution_time_avg=$(echo "$hdd_execution_time / $n" | bc)
echo "Average execution time (HDD): $hdd_execution_time"