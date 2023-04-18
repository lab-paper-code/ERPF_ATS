# Analyzing The Data I/O Patterns of ImageNet

* ImageNet 모델 fine-tuning 시, 데이터셋에서 데이터를 읽는 패턴을 확인하고자 함.

## Prerequisite
* Python 3.X
  * 사용한 파이썬 라이브러리는 requirements.txt에 정리
* Strace
* fine-tuning에 사용할 데이터
  * 본 실험에서는 ImageNet이 제공하는 Validation set으로 fine-tuning을 수행함.

## Usage
### 1. Strace 로그 생성하기
다음과 같은 커맨드로 strace 로그를 생성할 수 있다. 
```shell
sh start_strace.sh <number_of_image> <batch_size> <epoch> <output_file_path>
```
### 2. 로그의 File I/O 연산을 정리한 csv 파일 생성하기
다음의 스크립트를 사용하여 로그에서 File I/O 연산만 모아서 정리할 수 있다.
```shell
python3 strace_log_trimmer.py <strace_log_path>
```
생성된 csv 파일의 내용은 다음과 같다. (예시)

| Time | File Path  | Syscall | Offset |
| ---- | ---------- | ------- | ------ |
| 0.01 | Data1.JPEG | openat  | 0      |
| 0.20 | Data1.JPEG | lseek   | 0      |
| 0.88 | Data1.JPEG | read    | 384759 |
| 1.22 | Data2.JPEG | openat  | 0      |

### 3. 그래프 생성하기
다음의 스크립트를 사용하여 앞서 생성한 형식의 csv 파일을 시각화한다.
```shell
python3 get_chart.py <csv_file_path> <output_file_path>
```
결과 예시는 다음과 같다. `openat()`은 파란색, `lseek()`는 빨간색, `read()`는 녹색으로 표시한다.
![example.png](example.png)