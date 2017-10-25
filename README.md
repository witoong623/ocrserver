# ocrserver
## Docker
Need to mount volume to /data.  This image exposed port 80.
`sudo docker run -it --rm -v ocr-server-data:/data -p 127.0.0.1:8080:80 --name ocr-server-testing ocr-server`
