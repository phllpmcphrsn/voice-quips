# Description
A marketplace for voice clips by me! The clips are stored in S3.

# Stack
Golang/Gin - REST API
MinIO - store actual audio file
Postgresql - store metadata and links to S3 files
(Maybe Redis - cache metadata for links previously used)