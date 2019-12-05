# api

```sh
# register
curl -X POST http://127.0.0.1:7070/user/register -d '{ "userName":"ljs", "passWord":"ljs" }' -H 'content-type: application/json'

curl -X POST http://127.0.0.1:7070/user/register -d '{ "userName":"ljs1", "passWord":"ljs1" }' -H 'content-type: application/json'

# login
curl -X POST http://127.0.0.1:7070/user/login -d '{ "userName":"ljs", "passWord":"ljs" }' -H 'content-type: application/json'

curl -X POST http://127.0.0.1:7070/user/login -d '{ "userName":"ljs1", "passWord":"ljs1" }' -H 'content-type: application/json'

# checkAuth
curl -X POST http://127.0.0.1:7070/user/checkAuth -d '{ "authToken":"jwBfUHIhuM7S1KNNnsM_S9Sz3_KLlM9g3jqPYGGVVdo=" }' -H 'content-type: application/json'

curl -X POST http://127.0.0.1:7070/user/checkAuth -d '{ "authToken":"7F8tqhMRHcBGyVKzKePUzW1lE_A_efMMaepmdo7nZYE=" }' -H 'content-type: application/json'

# logout
curl -X POST http://127.0.0.1:7070/user/logout -d '{ "authToken":"fwCV1YozmvUOI1Fwr0k-ge-cWu0m6o_k5lZtlpxRUBk=" }' -H 'content-type: application/json'

# push
curl -X POST http://127.0.0.1:7070/push/push -d '{ "msg":"message1", "toUserId":"6", "roomId":1, "authToken":"jwBfUHIhuM7S1KNNnsM_S9Sz3_KLlM9g3jqPYGGVVdo=" }' -H 'content-type: application/json'
```
