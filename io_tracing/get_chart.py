import pandas as pd
import matplotlib.pyplot as plt

# 데이터 로드
df = pd.read_csv('strace_timeline.csv')

df = df[df['File Path'].str.endswith('.JPEG')]

# 시간(Time)을 가로축으로, 파일 경로(File Path)를 세로축으로 설정
plt.figure(figsize=(df['Time'].nunique()*0.8,
           df['File Path'].nunique()*0.3))  # 차트 크기 설정
plt.xlabel('Time (s)')  # x축 레이블 설정
plt.ylabel('File Path')  # y축 레이블 설정

# Syscall 종류(openat, lseek, read)에 따라 색상 지정
colors = {'openat': 'blue', 'lseek': 'red', 'read': 'green'}

# 파일 경로(File Path), Syscall 종류(Syscall), Offset를 이용하여 점 찍기
for index, row in df.iterrows():
    # 점의 크기 조정
    if row['Offset'] == 0:
        dot_size = 100
    else:
        dot_size = row['Offset'] / 1000  # Offset 값에 따라 크기 조정
    plt.scatter(row['Time'], row['File Path'],
                color=colors[row['Syscall']], s=dot_size)

plt.savefig('all_chart.png')
