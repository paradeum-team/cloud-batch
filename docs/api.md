


# Cloud Batch API
  

## Informations

### Version

v0.0.3

### Contact

  

### Terms Of Service

https://gitlab.paradeum.com/pld/cloud-batch

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## Access control

### Security Schemes

#### ApiKeyAuth (header: Authorization)



> **Type**: apikey

## All endpoints

###  auth

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /auth/login | [post auth login](#post-auth-login) | post Auth |
| PUT | /auth/passwd | [put auth passwd](#put-auth-passwd) | Update password |
  


###  servers

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /batch/servers | [delete batch servers](#delete-batch-servers) | Batch Get Servers |
| GET | /servers | [get servers](#get-servers) | Batch Get Servers |
| POST | /batch/servers | [post batch servers](#post-batch-servers) | Batch Create Servers |
  


###  version

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /version | [get version](#get-version) | Get Version |
  


## Paths

### <span id="delete-batch-servers"></span> Batch Get Servers (*DeleteBatchServers*)

```
DELETE /batch/servers
```

#### Produces
  * application/json

#### Security Requirements
  * ApiKeyAuth

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| batchDeleteServersForm | `body` | [ModelsBatchDeleteServersForm](#models-batch-delete-servers-form) | `models.ModelsBatchDeleteServersForm` | | ✓ | | batchDeleteServersForm |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-batch-servers-200) | OK | OK |  | [schema](#delete-batch-servers-200-schema) |
| [500](#delete-batch-servers-500) | Internal Server Error | Internal Server Error |  | [schema](#delete-batch-servers-500-schema) |
| [default](#delete-batch-servers-default) | |  |  | [schema](#delete-batch-servers-default-schema) |

#### Responses


##### <span id="delete-batch-servers-200"></span> 200 - OK
Status: OK

###### <span id="delete-batch-servers-200-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="delete-batch-servers-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="delete-batch-servers-500-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="delete-batch-servers-default"></span> Default Response


###### <span id="delete-batch-servers-default-schema"></span> Schema

  

[AppResponseString](#app-response-string)

### <span id="get-servers"></span> Batch Get Servers (*GetServers*)

```
GET /servers
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| batch_number | `query` | string | `string` |  |  |  | batch_number |
| project | `query` | string | `string` |  |  |  | project |
| provider | `query` | string | `string` |  |  |  | provider |
| status | `query` | []string | `[]string` | `multi` |  |  | status |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-servers-200) | OK | OK |  | [schema](#get-servers-200-schema) |
| [500](#get-servers-500) | Internal Server Error | Internal Server Error |  | [schema](#get-servers-500-schema) |
| [default](#get-servers-default) | |  |  | [schema](#get-servers-default-schema) |

#### Responses


##### <span id="get-servers-200"></span> 200 - OK
Status: OK

###### <span id="get-servers-200-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="get-servers-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-servers-500-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="get-servers-default"></span> Default Response


###### <span id="get-servers-default-schema"></span> Schema

  

[AppResponseString](#app-response-string)

### <span id="get-version"></span> Get Version (*GetVersion*)

```
GET /version
```

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-version-200) | OK | OK |  | [schema](#get-version-200-schema) |
| [500](#get-version-500) | Internal Server Error | Internal Server Error |  | [schema](#get-version-500-schema) |

#### Responses


##### <span id="get-version-200"></span> 200 - OK
Status: OK

###### <span id="get-version-200-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="get-version-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="get-version-500-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

### <span id="post-auth-login"></span> post Auth (*PostAuthLogin*)

```
POST /auth/login
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| auth | `body` | [ModelsAuth](#models-auth) | `models.ModelsAuth` | | ✓ | | auth |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-auth-login-200) | OK | OK |  | [schema](#post-auth-login-200-schema) |
| [400](#post-auth-login-400) | Bad Request | Bad Request |  | [schema](#post-auth-login-400-schema) |
| [500](#post-auth-login-500) | Internal Server Error | Internal Server Error |  | [schema](#post-auth-login-500-schema) |

#### Responses


##### <span id="post-auth-login-200"></span> 200 - OK
Status: OK

###### <span id="post-auth-login-200-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="post-auth-login-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-auth-login-400-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="post-auth-login-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-auth-login-500-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

### <span id="post-batch-servers"></span> Batch Create Servers (*PostBatchServers*)

```
POST /batch/servers
```

#### Produces
  * application/json

#### Security Requirements
  * ApiKeyAuth

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| batchCreateServers | `body` | [ModelsBatchCreateServersForm](#models-batch-create-servers-form) | `models.ModelsBatchCreateServersForm` | | ✓ | | batchCreateServers |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-batch-servers-200) | OK | OK |  | [schema](#post-batch-servers-200-schema) |
| [500](#post-batch-servers-500) | Internal Server Error | Internal Server Error |  | [schema](#post-batch-servers-500-schema) |
| [default](#post-batch-servers-default) | |  |  | [schema](#post-batch-servers-default-schema) |

#### Responses


##### <span id="post-batch-servers-200"></span> 200 - OK
Status: OK

###### <span id="post-batch-servers-200-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="post-batch-servers-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-batch-servers-500-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="post-batch-servers-default"></span> Default Response


###### <span id="post-batch-servers-default-schema"></span> Schema

  

[AppResponseString](#app-response-string)

### <span id="put-auth-passwd"></span> Update password (*PutAuthPasswd*)

```
PUT /auth/passwd
```

#### Produces
  * application/json

#### Security Requirements
  * ApiKeyAuth

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| updateAuth | `body` | [ModelsUpdateAuth](#models-update-auth) | `models.ModelsUpdateAuth` | | ✓ | | updateAuth |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#put-auth-passwd-200) | OK | OK |  | [schema](#put-auth-passwd-200-schema) |
| [400](#put-auth-passwd-400) | Bad Request | Bad Request |  | [schema](#put-auth-passwd-400-schema) |
| [401](#put-auth-passwd-401) | Unauthorized | Bad Request |  | [schema](#put-auth-passwd-401-schema) |
| [500](#put-auth-passwd-500) | Internal Server Error | Internal Server Error |  | [schema](#put-auth-passwd-500-schema) |

#### Responses


##### <span id="put-auth-passwd-200"></span> 200 - OK
Status: OK

###### <span id="put-auth-passwd-200-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="put-auth-passwd-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="put-auth-passwd-400-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="put-auth-passwd-401"></span> 401 - Bad Request
Status: Unauthorized

###### <span id="put-auth-passwd-401-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

##### <span id="put-auth-passwd-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="put-auth-passwd-500-schema"></span> Schema
   
  

[AppResponseString](#app-response-string)

## Models

### <span id="app-response-string"></span> app.ResponseString


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| code | string| `string` |  | |  |  |
| data | [interface{}](#interface)| `interface{}` |  | |  |  |
| errDetail | [interface{}](#interface)| `interface{}` |  | |  |  |
| msg | string| `string` |  | |  |  |
| requestId | string| `string` |  | |  |  |



### <span id="models-auth"></span> models.Auth


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| password | string| `string` | ✓ | |  |  |
| username | string| `string` | ✓ | |  |  |



### <span id="models-batch-create-servers-form"></span> models.BatchCreateServersForm


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| count | integer| `int64` | ✓ | `1`| 创建数量 |  |
| cpu_core_count | integer| `int64` | ✓ | `4`| cpu 核数 |  |
| memory_size_mb | integer| `int64` | ✓ | `8192`| 内存（单位 M) |  |
| project | string| `string` | ✓ | `"bfs"`| 项目标签 |  |
| provider | string| `string` | ✓ | `"aliyun"`| 云提供者 |  |
| region_scope | string| `string` | ✓ | `"cnml"`| 区域前缀: all 全部, cnml 中国大陆, cn 中国， ap 亚太， us 美洲， eu 欧洲, me 中东 |  |



### <span id="models-batch-delete-servers-form"></span> models.BatchDeleteServersForm


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| batch_number | string| `string` |  | | 创建主机的批号 |  |
| ids | []string| `[]string` |  | | 主机 ids, 如果设置了 ids, 忽略其它参数 |  |
| project | string| `string` |  | `"test"`| 项目标签 |  |
| provider | string| `string` |  | `"aliyun"`| 云提供者 |  |
| region_scope | string| `string` |  | `"cnml"`| 区域前缀: all 全部, cnml 中国大陆, cn 中国， ap 亚太， us 美洲， eu 欧洲, me 中东 |  |



### <span id="models-update-auth"></span> models.UpdateAuth


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| oldPassword | string| `string` | ✓ | |  |  |
| password | string| `string` | ✓ | |  |  |
| username | string| `string` | ✓ | |  |  |


