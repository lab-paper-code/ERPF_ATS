#!/bin/sh
dataset_path=$2

# Download ImageNet Dataset
wget $dataset_path

mv ./ILSVRC2012_img_val.tar ./images
cd images

# Unzip
tar -xvf ILSVRC2012_img_val.tar
rm ILSVRC2012_img_val.tar

# 이미지를 카테고리 별로 분류
sh image_categorize.sh

# 카테고리 별 이미지 파일을 2개씩만 남김 (총 2,000장이 되게끔)
for dir in */
do
  cd "$dir"
  jpeg_files=$(ls *.JPEG 2>/dev/null)

  if [ $(echo "$jpeg_files" | wc -l) -gt 2 ]
  then
    excess=$(echo "$jpeg_files" | head -n -2)
    rm $excess
  fi
  cd ..
done

cd ..