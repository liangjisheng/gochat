# api

```sh
# register
curl -X POST http://127.0.0.1:7070/user/register -d '{ "userName":"ljs", "passWord":"ljs" }' -H 'content-type: application/json'

# login
curl -X POST http://127.0.0.1:7070/user/login -d '{ "userName":"ljs", "passWord":"ljs" }' -H 'content-type: application/json'

# checkAuth
curl -X POST http://127.0.0.1:7070/user/checkAuth -d '{ "authToken":"fwCV1YozmvUOI1Fwr0k-ge-cWu0m6o_k5lZtlpxRUBk=" }' -H 'content-type: application/json'

# logout
curl -X POST http://127.0.0.1:7070/user/logout -d '{ "authToken":"fwCV1YozmvUOI1Fwr0k-ge-cWu0m6o_k5lZtlpxRUBk=" }' -H 'content-type: application/json'
```
