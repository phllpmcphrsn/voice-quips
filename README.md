# Description
A marketplace for voice clips by me! The clips are stored in S3.

# Stack
Golang/Gin - REST API
MinIO - store actual audio file
Postgresql - store metadata and links to S3 files
(Maybe Redis - cache metadata for links previously used)

# MinIO setup
`docker run -d -p 9000:9000/tcp -p 9001:9001  minio/minio:latest server /data --console-address ":9001"`