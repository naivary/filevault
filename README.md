# Filevault

Filevault is a simple fileserver which is providing an http interface to save, delete and retrieve files.

## Provided API-Interface

### `POST /api/v1/`
This path is storing the given file which is provided under the multipart form key `file` in the directory defined by the query parameter `dir`.
It will automatically append the filename to the `dir` so it is expected that `dir` is only containing the path of teh directory.

### `GET /api/v1/`
Using the query parameter `path` you can retrieve a saved file if the file exists. The `Content-Type` will automatically be 
set to the right mime-type if an official mime-type exists.

### `DELETE /api/v1/` 
Using the query parameter `path` you can delete a saved file if the file exists. If the file does not exist the error will be ignored. 
If after the deletion of the file the directory is empty, the directory will also be deleted. A non-empty directory cannot be deleted.

### `GET /api/v1/health`
This endpoint allows for an automatic healthcheck of the application which is useful in the context of Docker and Kubernetes.

