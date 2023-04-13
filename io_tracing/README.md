strace를 이용한 이미지 분류 프로그램 system-call 로그분석 

 - config.py: 이미지 분류와 분석 파이썬 코드의 config 값 저장  
 - image_classification.py: 이미지 분류 작업을 실행하는 파이썬 스크립트 (Usage: python3 image_classification.py -n=image_count)
 - strace_analyzer.py: strace log 파일을 통해 summary와 plot을 만드는 파이썬 스크립트
 - strace_log_trimmer.py: strace log에서 확인된 파일 읽기 작업을 (파일 경로, 해당 파일로부터 읽은 바이트 수, 해당 파일의 물리적 오프셋) 형식으로 정리하여 파일로 출력하는 파이썬 스크립트. 결과 파일은 ./results에 저장됨 (Usage: python3 strace_log_trimmer.py <log file path>)
 - ./plots: strace log를 분석해 나온 plot 이미지 파일들을 저장한 디렉토리  
 - ./results: 이미지 분류 결과 텍스트 파일를 저장한 디렉토리 
 - requirements.txt: 필요한 파이썬 모듈을 나열한 텍스트 파일 (Usage: pip3 install -r requirements.txt)
 - start_strace.sh: strace로 image_classification.py를 실행하며 log를 저장하는 shell script (Usage: sh start_strace.sh n / n=image_count)
 - start_strace_visualization.sh: strace log로부터 그래프를 만드는 shell script (Usage: sh start_strace_visualization.sh n / n=image_count)
 - ./strace_outputs: strace log를 저장하는 디렉토리 
 - ./summary : strace_analyzer.py를 실행한 후 텍스트 파일 형식의 요약 정보들을 저장하는 디렉토리 