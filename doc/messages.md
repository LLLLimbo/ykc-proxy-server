## Messages



### Login Verification (01)

| Field           | Type   | Descripton |
| --------------- | ------ | ---------- |
| header          | Header |            |
| id              | string | device id  |
| elcType         | int    |            |
| guns            | int    |            |
| protocolVersion | int    |            |
| softwareVersion | string |            |
| network         | int    |            |
| sim             | string |            |
| operator        | int    |            |



### Heartbeat (03)

| Field     | Type   | Description |
| --------- | ------ | ----------- |
| header    | Header |             |
| id        | string |             |
| gun       | string |             |
| gunStatus | int    |             |





### Billing model verification (05)

| Field            | Type   | Description |
| ---------------- | ------ | ----------- |
| header           | Header |             |
| Id               | string |             |
| billingModelCode | string |             |



### Offline data report (13)

| Field                   | Type   | Description |
| ----------------------- | ------ | ----------- |
| header                  | Header |             |
| tradeSeq                | string |             |
| Id                      | string |             |
| gunId                   | string |             |
| status                  | int    |             |
| reset                   | int    |             |
| plugged                 | int    |             |
| ov                      | int    |             |
| oc                      | int    |             |
| lineTemp                | int    |             |
| lineCode                | string |             |
| soc                     | int    |             |
| bpTopTemp               | int    |             |
| accumulatedChargingTime | int    |             |
| remainingTime           | int    |             |
| chargingDegrees         | int    |             |
| lossyChargingDegrees    | int    |             |
| chargedAmount           | int    |             |
| hardwareFailure         | int    |             |





### Remote bootstrap response (33)

| Field    | Type   | Description |
| -------- | ------ | ----------- |
| header   | Header |             |
| tradeSeq | string |             |
| Id       | string |             |
| gunId    | string |             |
| result   | bool   |             |
| reason   | int    |             |





### Remote shutdown response (35)

| Field  | Type   | Description |
| ------ | ------ | ----------- |
| header | Header |             |
| Id     | string |             |
| gunId  | string |             |
| result | bool   |             |
| reason | int    |             |





### Remote reboot response (91)

| Field  | Type   | Description |
| ------ | ------ | ----------- |
| header | Header |             |
| Id     | string |             |
| result | int    |             |





### Transaction record (3B)

| Field                     | Type   | Description |
| ------------------------- | ------ | ----------- |
| header                    | Header |             |
| tradeSeq                  | string |             |
| Id                        | string |             |
| gunId                     | string |             |
| startAt                   | int64  |             |
| endAt                     | int64  |             |
| sharpUnitPrice            | int64  |             |
| sharpElectricCharge       | int64  |             |
| lossySharpElectricCharge  | int64  |             |
| sharpPrice                | int64  |             |
| peakUnitPrice             | int64  |             |
| peakElectricCharge        | int64  |             |
| lossyPeakElectricCharge   | int64  |             |
| peakPrice                 | int64  |             |
| flatUnitPrice             | int64  |             |
| flatElectricCharge        | int64  |             |
| lossyFlatElectricCharge   | int64  |             |
| flatPrice                 | int64  |             |
| valleyUnitPrice           | int64  |             |
| valleyElectricCharge      | int64  |             |
| lossyValleyElectricCharge | int64  |             |
| valleyPrice               | int64  |             |
| initialMeterReading       | int64  |             |
| finalMeterReading         | int64  |             |
| totalElectricCharge       | int64  |             |
| lossyTotalElectricCharge  | int64  |             |
| consumptionAmount         | int64  |             |
| vin                       | string |             |
| startType                 | int    |             |
| transactionDateTime       | int64  |             |
| stopReason                | int    |             |
| physicalCardNumber        | string |             |



### Appendix

#### Header

| Field     | Type   | Description |
| --------- | ------ | ----------- |
| length    | int    |             |
| seq       | int    |             |
| encrypted | bool   |             |
| frameId   | string |             |