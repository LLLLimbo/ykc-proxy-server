## Messages



### Login Verification (01)

| Field           | Type   | Descripton                             |
| --------------- | ------ | -------------------------------------- |
| header          | Header |                                        |
| id              | string | device id                              |
| elcType         | int    | 0-direct 1-cross                       |
| guns            | int    | number of gun                          |
| protocolVersion | int    | protocol version                       |
| softwareVersion | string | software version                       |
| network         | int    | network type 0-SIM 1-LAN 2-WAN 3-OTHER |
| sim             | string | SIM card number                        |
| operator        | int    | network operator 0-CMCC 1-CTCC 2-CUCC  |



### Heartbeat (03)

| Field     | Type   | Description      |
| --------- | ------ | ---------------- |
| header    | Header |                  |
| id        | string | device id        |
| gun       | string | gun id           |
| gunStatus | int    | 0-normal 1-error |





### Billing model verification (05)

| Field            | Type   | Description               |
| ---------------- | ------ | ------------------------- |
| header           | Header |                           |
| id               | string | device id                 |
| billingModelCode | string | billing model's unique id |



### Offline data report (13)

| Field                   | Type   | Description                                          |
| ----------------------- | ------ | ---------------------------------------------------- |
| header                  | Header |                                                      |
| tradeSeq                | string | trade sequence number                                |
| id                      | string | device id                                            |
| gunId                   | string | gun id                                               |
| status                  | int    | gun's status 0-offline 1-error 2-free 3-charging     |
| reset                   | int    | is gun reset 0-false 1-true 2-unknown                |
| plugged                 | int    | is gun plugged 0-false 1-true                        |
| ov                      | int    | output voltage (X10)                                 |
| oc                      | int    | output current (X10)                                 |
| lineTemp                | int    | the wire temperature of the gun (Offset -50)         |
| lineCode                | string | the wire number of gun                               |
| soc                     | int    |                                                      |
| bpTopTemp               | int    | the highest temperature of battery pack (Offset -50) |
| accumulatedChargingTime | int    | accumulated charging duration (in minutes)           |
| remainingTime           | int    | remaining time (in minutes)                          |
| chargingDegrees         | int    | charging degrees (X10000)                            |
| lossyChargingDegrees    | int    | lossy charging degrees (X10000)                      |
| chargedAmount           | int    | charged amount (X10000)                              |
| hardwareFailure         | int    | hardware failure code                                |



### Charging finished (19)

| Field                            | Type   | Description           |
| -------------------------------- | ------ | --------------------- |
| header                           | Header |                       |
| tradeSeq                         | string | trade sequence number |
| id                               | string | device id             |
| gunId                            | string | gun id                |
| bmsSoc                           | int    |                       |
| bmsBatteryPackLowestVoltage      | int    |                       |
| bmsBatteryPackHighestVoltage     | int    |                       |
| bmsBatteryPackLowestTemperature  | int    |                       |
| bmsBatteryPackHighestTemperature | int    |                       |
| cumulativeChargingDuration       | int    |                       |
| outputPower                      | int    |                       |
| chargingUnitId                   | int    |                       |



### Remote bootstrap response (33)

| Field    | Type   | Description                                                  |
| -------- | ------ | ------------------------------------------------------------ |
| header   | Header |                                                              |
| tradeSeq | string | trade sequence number                                        |
| id       | string | device id                                                    |
| gunId    | string | gun id                                                       |
| result   | bool   | bootstrap result 0-fail 1-success                            |
| reason   | int    | fail reason  0-none 1-device id not match 2-gun is already in charging 3-device on failure 4-device offline 5-gun is not plugged |





### Remote shutdown response (35)

| Field  | Type   | Description                                                  |
| ------ | ------ | ------------------------------------------------------------ |
| header | Header |                                                              |
| Id     | string | device id                                                    |
| gunId  | string | gun id                                                       |
| result | bool   | shutdown result 0-fail 1-success                             |
| reason | int    | fail reason  0-none 1-device id not match 2-gun is not in charging 3-other |





### Set billing model response (57)

| Field  | Type   | Description      |
| ------ | ------ | ---------------- |
| header | Header |                  |
| id     | string | device id        |
| result | int    | 0-fail 1-success |





### Remote reboot response (91)

| Field  | Type   | Description      |
| ------ | ------ | ---------------- |
| header | Header |                  |
| id     | string | device id        |
| result | int    | 0-fail 1-success |





### Transaction record (3B)

| Field                     | Type   | Description                           |
| ------------------------- | ------ | ------------------------------------- |
| header                    | Header |                                       |
| tradeSeq                  | string | trade sequence number                 |
| id                        | string | device id                             |
| gunId                     | string | gun id                                |
| startAt                   | int64  | charging start time                   |
| endAt                     | int64  | charging end time                     |
| sharpUnitPrice            | int64  | sharp unit price (X100000)            |
| sharpElectricCharge       | int64  | sharp electric charge (X10000)        |
| lossySharpElectricCharge  | int64  | lossy sharp electric charge (X10000)  |
| sharpPrice                | int64  | sharp price (X10000)                  |
| peakUnitPrice             | int64  | peak unit price (X100000)             |
| peakElectricCharge        | int64  | peak electric charge (X10000)         |
| lossyPeakElectricCharge   | int64  | lossy peak electric charge (X10000)   |
| peakPrice                 | int64  | peak price (X10000)                   |
| flatUnitPrice             | int64  | flat unit price (X100000)             |
| flatElectricCharge        | int64  | flat electric charge (X10000)         |
| lossyFlatElectricCharge   | int64  | lossy flat electric charge (X10000)   |
| flatPrice                 | int64  | flat price (X10000)                   |
| valleyUnitPrice           | int64  | valley unit price (X100000)           |
| valleyElectricCharge      | int64  | valley electric charge (X10000)       |
| lossyValleyElectricCharge | int64  | lossy valley electric charge (X10000) |
| valleyPrice               | int64  | valley price (X10000)                 |
| initialMeterReading       | int64  | initial meter reading (X10000)        |
| finalMeterReading         | int64  | final meter reading (X10000)          |
| totalElectricCharge       | int64  | total electric charge (X10000)        |
| lossyTotalElectricCharge  | int64  | lossy total electric charge (X10000)  |
| consumptionAmount         | int64  | consumption amount (X10000)           |
| vin                       | string | VIN code                              |
| startType                 | int    | 1-APP 2-card 4-offline card 5-vin     |
| transactionDateTime       | int64  | transaction time                      |
| stopReason                | int    | stop reason                           |
| physicalCardNumber        | string | physical card number                  |



### Appendix

#### Header

| Field     | Type   | Description |
| --------- | ------ | ----------- |
| length    | int    |             |
| seq       | int    |             |
| encrypted | bool   |             |
| frameId   | string |             |