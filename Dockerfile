FROM golang:1.8

ENV SERVER_ADDR :80
ENV TEMP_DIR /data
EXPOSE 80

RUN apt-get update && apt-get install -y \
    libleptonica-dev \
    libtesseract-dev \
    tesseract-ocr

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"] # ["app"]