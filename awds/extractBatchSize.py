import re
import os


def extract_batch_size(results):
    batch_sizes = []
    for line in results:
        columns = line.split()
        if len(columns) >= 10:
            batch_size = int(columns[9])
            batch_sizes.append(batch_size)
    return batch_sizes


f = open(
    "./input.txt",
)
results = f.readlines()

batch_sizes = extract_batch_size(results)
print("Batch Sizes:", batch_sizes)
print("Sum of Batch Sizes:", sum(batch_sizes) + 30)
