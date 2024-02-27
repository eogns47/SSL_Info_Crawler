# Golang 이미지를 기반으로 함
FROM golang:latest

# 작업 디렉토리 설정
WORKDIR /app

# 로컬의 소스 코드를 컨테이너의 작업 디렉토리로 복사
COPY . .

# 소스 코드 빌드
RUN go build -o ./app/main ./src

# 실행 파일 지정
ENTRYPOINT ["/app/main"]
