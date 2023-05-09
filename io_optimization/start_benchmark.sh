#!/bin/sh

n=5
ssd_path=$1
hdd_path=$2
batch_size=${3:-32}
epochs=${4:-5}

ssd_checkpoint="./checkpoints"
hdd_checkpoint="/hdd/checkpoints"

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
    python3.8 target_scripts/image_fine-tuning.py --dataset_path=$ssd_path --batch_size=$batch_size --epochs=$epochs --checkpoint_path=$ssd_checkpoint 1> fine-tuning_log_ssd 2> error_log_ssd
    end=$(date +%s.%N)
    runtime=$(echo "$end - $start" | bc)
    sdd_execution_time=$(echo "$sdd_execution_time + $runtime" | bc)
done
sdd_execution_time_avg=$(echo "$sdd_execution_time / $n" | bc)
echo "Average execution time (SSD): $sdd_execution_time"

hdd_execution_time=0
sudo hdparm -W 0 /dev/sda
for i in $(seq 1 $n); do
    start=$(date +%s.%N)
    python3.8 target_scripts/image_fine-tuning.py --dataset_path=$hdd_path --batch_size=$batch_size --epochs=$epochs --checkpoint_path=$hdd_checkpoint 1> fine-tuning_log_hdd 2> error_log_hdd
    end=$(date +%s.%N)
    runtime=$(echo "$end - $start" | bc)
    hdd_execution_time=$(echo "$hdd_execution_time + $runtime" | bc)
done
sudo hdparm -W 1 /dev/sda
hdd_execution_time_avg=$(echo "$hdd_execution_time / $n" | bc)
echo "Average execution time (HDD): $hdd_execution_time"