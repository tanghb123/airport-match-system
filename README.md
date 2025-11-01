### airport-match-system



### 1.airport-match-system project

#### 1.1 config file:  /config/app.yaml

```note
# mysql

initial account: root
initial password: 123456
port: 3306

# reids
initial account: root
port: 6379

You can modify the configuration according to your own requirements
```

#### 1.2 start project

```note
cd airport-match-system
go mod tidy
go run main.go
```

#### 1.3 Interface Introduction

##### 1.3.1 driver create order

**Explanation**

```go
1、Generally, users need to register and log in. After successful login, a token will be returned. When accessing the interface subsequently, the token can be passed in the header for permission verification. Token verification can be executed through middleware. For simplicity, the driver's name can be directly passed in.

2.The generated driver order information is stored in Redis with orderID as the key. When querying order data through orderID later, it is first retrieved from Redis. If the data does not exist in Redis, it is then retrieved from the database.
3.Generate a driver rating randomly from 1 to 5, with 1 being the lowest and 5 being the highest.
```



Interface url

```
# POST    port：8081
127.0.0.1:8081/v1/driver/create_route
```

Input parameters **form-data** format

| Key              | Value               | Describe                                  |
| ---------------- | ------------------- | ----------------------------------------- |
| `user_name`      | dirver1             | driver user name, name it as you like     |
| `airport`        | newyourk airport    | name it as you like                       |
| `available_time` | 2025-11-10 15:00:00 | scheduled departure time                  |
| `vehicle_type`   | 1                   | you can fill in 1, 2, 3,4,5 ...(int type) |
| `price_per_km`   | 1.1                 | need greater than 0 (float64 type)        |

success result

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "driver_order_id": "driver-60020899-91c9-4588-a1b3-90e21405b99d",
        "status": 0  // 0: not matched, 1: matched
    }
}
```

fail result

```json
{
    "code": 10004,
    "message": "Invalid vehicle type",
    "data": null
}
```

error code

| Code  | Message                  |
| ----- | ------------------------ |
| 10001 | exist blank parameter    |
| 10002 | invalid format           |
| 10003 | Invalid available time   |
| 10004 | Invalid vehicle type     |
| 10005 | Invalid per km parameter |
| 10006 | Failed to create order   |
| 10007 | Duplicated order         |
| 10008 | Failed to marshal order  |



##### 1.3.2 passenger create order

**Explanation**

```go
1.maxPricePerKm represents the maximum price offered by passengers, and there is no passenger rating.
```



Interface url

```
# POST    port：8081
127.0.0.1:8081/v1/passenger/create_route
```

Input parameters **form-data** format

| Key                | Value               | Describe                                                       |
| ------------------ | ------------------- | -------------------------------------------------------------- |
| `user_name`        | passenger1          | passenger user name, name it as you like                       |
| `airport`          | newyourk airport    | name it as you like                                            |
| `available_time`   | 2025-11-10 15:00:00 | scheduled departure time                                       |
| `vehicle_type`     | 1                   | you can fill in 1, 2, 3,4,5 ...(int type)                      |
| `max_price_per_km` | 1.1                 | maximum price per kilometer,need greater than 0 (float64 type) |

success result



```json
{
    "code": 200,
    "message": "success",
    "data": {
        "passenger_order_id": "passenger-4e5640d9-0d6e-4a3e-b9e8-b47fb7565430",
        "status": 0 // 0: not matched, 1: matched
    }
}
```

fail result

```json
{
    "code": 10004,
    "message": "Invalid vehicle type",
    "data": null
}
```



error code

| Code  | Message                  |
| ----- | ------------------------ |
| 10001 | exist blank parameter    |
| 10002 | invalid format           |
| 10003 | Invalid available time   |
| 10004 | Invalid vehicle type     |
| 10005 | Invalid per km parameter |
| 10006 | Failed to create order   |
| 10007 | Duplicated order         |
| 10008 | Failed to marshal order  |



##### 1.3.3 passenger find match order

**Explanation**

```go
1.When the user passes in the "passenger_order_id", first verify whether the order exists.
2.Retrieve the fee charged per kilometer by the platform from the FeeConfig table ("feePerKm"). This is a simplified version. In actual projects, you can cache these configuration items in Redis, and set up a scheduled task to periodically update the data in the FeeConfig table to Redis, reducing database queries.
3.To ensure platform revenue, the user's "MaxPricePerKm-feePerKm" must be greater than the value of the driver's order "price_per_km", while Airport, AvailableTime, and VehicleType must be equal. Among these orders, the one with the minimum "price_per_km", the maximum "rating", and the earliest "created_at"" will be selected.
4.It should be noted that the "price_per_km" in "respons" represents the amount of money that the user needs to pay. The driver will charge a fee that is this amount minus the platform fee("feePerKm").
```

Interface url


    # POST    port：8081
    127.0.0.1:8081/v1/passenger/find_match

Input parameters form-data format

| Key                  | Value                                          | Describe             |
| -------------------- | ---------------------------------------------- | -------------------- |
| `passenger_order_id` | passenger-90151c82-b436-4aad-b91f-f3c38c1dfe0d | "passenger-" + uuid, |

success result



```json
{
    "code": 200,
    "message": "success",
    "data": {
        "airport": "newyourk airport",
        "available_time": "2025-11-10T15:00:00+08:00",
        "driver_name": "dirver3",
        "driver_order_id": "driver-f525dfac-926e-4166-b9f0-23e46e95d70e",
        "price_per_km": 1.2000000000000002,
        "rating": 5,
        "vehicle_type": 3
    }
}
```

fail result



```json
{
    "code": 10009,
    "message": "Failed to check order",
    "data": null
}
```

error code

| Code  | Message                         |
| ----- | ------------------------------- |
| 10009 | Failed to check order           |
| 10010 | Failed to get fee               |
| 10011 | Failed to match order           |
| 10012 | Failed to create match relation |
| 10013 | Order has matched               |

#### 

##### 1.3.4 execute match order

**Explanation**

```go
1.Update PassengerPublishOrder, DriverPublishOrder, MatchRelation, and MatchOrderResult table through transaction
```



Interface url


    # POST    port：8081
    127.0.0.1:8081/v1/passenger/execute_match



Input parameters form-data format

| Key                  | Value                                          | Describe             |
| -------------------- | ---------------------------------------------- | -------------------- |
| `passenger_order_id` | passenger-90151c82-b436-4aad-b91f-f3c38c1dfe0d | "passenger-" + uuid, |
| `driver_order_id`    | driver-f525dfac-926e-4166-b9f0-23e46e95d70e    | "driver-" + uuid,    |

success result



```json
{
    "code": 200,
    "message": "success",
    "data": {
        "ID": 2,
        "MatchOrderID": "5cdbe9e5-91b2-4f4d-9abd-d6a7779027ac",
        "PassengerID": "passenger-90151c82-b436-4aad-b91f-f3c38c1dfe0d",
        "DriverID": "driver-f525dfac-926e-4166-b9f0-23e46e95d70e",
        "DriverRating": 5,
        "VehicleType": 3,
        "DriverPricePerKm": 1.1,
        "PassengerPricePerKm": 1.2000000000000002,
        "FeePerKm": 0.1,
        "Airport": "newyourk airport",
        "AvailableTime": "2025-11-10T15:00:00+08:00",
        "CreatedAt": "2025-11-01T12:26:57.421+08:00",
        "UpdatedAt": "2025-11-01T12:26:57.421+08:00"
    }
}
```

fail result

```json
{
    "code": 10014,
    "message": "Failed to check relation",
    "data": null
}
```

error code

| Code  | Message                        |
| ----- | ------------------------------ |
| 10010 | Failed to get fee              |
| 10013 | Order has matched              |
| 10014 | Failed to check relation       |
| 10015 | Invalid passenger order        |
| 10016 | Invalid driver order           |
| 10017 | Failed to execute match orders |

#### 
