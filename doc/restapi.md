## API List



### Verification Response(02)

Path: `/proxy/02`

Request body:

| Field  | Type   | Description              |
| ------ | ------ | ------------------------ |
| header | Header |                          |
| id     | string | device id                |
| result | bool   | true-pass   false-reject |



Example request:

```json
{
    "header":{
        "encrypted": false,
        "seq": 72
    },
    "id": "12345620230378",
    "result": true
}
```





Response body:

| Field   | Type   | Description   |
| ------- | ------ | ------------- |
| message | string | error message |





### Billing model verification(06)

Path: `/proxy/06`

Request body:

| Field            | Type   | Description              |
| ---------------- | ------ | ------------------------ |
| header           | Header |                          |
| id               | string | device id                |
| billingModelCode | string | code of billing model    |
| result           | bool   | true-pass   false-reject |



Example request:

```json
{
    "header":{
        "encrypted": false,
        "seq": 72
    },
    "id": "12345620230378",
    "billingModelCode": "0000",
    "result": true
}
```





Response body:

| Field   | Type   | Description   |
| ------- | ------ | ------------- |
| message | string | error message |





### Remote bootstrap(34)

Path: `/proxy/34`

Request body:

| Field            | Type   | Description            |
| ---------------- | ------ | ---------------------- |
| header           | Header |                        |
| tradeSeq         | string | trade sequence number  |
| id               | string | device id              |
| gunId            | string | gun id                 |
| logicCard        | string | number of logic card   |
| physicalCard     | string | number of physic card  |
| billingModelCode | string | code of billing model  |
| balance          | int    | account balance (x100) |



Example request:

```json
{
    "header":{
        "encrypted": false,
        "seq": 73
    },
    "id": "12345620230378",
    "tradeSeq": "55031412782305012018061914444680",
    "gunId": "01",
    "logicCard": "0000001000000573",
    "physicalCard": "00000000D14B0A54",
    "balance": 1000000
}
```





Response body:

| Field   | Type   | Description   |
| ------- | ------ | ------------- |
| message | string | error message |



### Remote shutdown(36)

Path: `/proxy/36`

Request body:

| Field  | Type   | Description |
| ------ | ------ | ----------- |
| header | Header |             |
| id     | string | device id   |
| gunId  | string | gun id      |



Example request:

```json
{
    "header":{
        "encrypted": false,
        "seq": 73
    },
    "id": "12345620230378",
    "gunId": "01",
}
```





Response body:

| Field   | Type   | Description   |
| ------- | ------ | ------------- |
| message | string | error message |





### Confirme transaction record(40)

Path: `/proxy/40`

Request body:

| Field    | Type   | Description           |
| -------- | ------ | --------------------- |
| header   | Header |                       |
| id       | string | device id             |
| tradeSeq | string | trade sequence number |
| result   | int    | 0-pass 1-reject       |



Example request:

```json
{
    "header":{
        "encrypted": false,
        "seq": 73
    },
    "id": "12345620230378",
    "tradeSeq": "55031412782305012018061914444680",
    "result": 0
}
```





Response body:

| Field   | Type   | Description   |
| ------- | ------ | ------------- |
| message | string | error message |





### Remote reboot(92)

Path: `/proxy/92`

Request body:

| Field   | Type   | Description                             |
| ------- | ------ | --------------------------------------- |
| header  | Header |                                         |
| id      | string | device id                               |
| control | string | 1- reboot instantly  2-reboot when free |



Example request:

```json
{
    "header":{
        "encrypted": false,
        "seq": 73
    },
    "id": "12345620230378",
    "control": 1
}
```





Response body:

| Field   | Type   | Description   |
| ------- | ------ | ------------- |
| message | string | error message |

