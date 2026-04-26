# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/proto/generator.proto](#api_proto_generator-proto)
    - [Field](#generator-v1-Field)
    - [GenerateRequest](#generator-v1-GenerateRequest)
    - [GeneratedRecord](#generator-v1-GeneratedRecord)
    - [Option](#generator-v1-Option)
    - [Schema](#generator-v1-Schema)
  
    - [GeneratorService](#generator-v1-GeneratorService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_proto_generator-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/proto/generator.proto



<a name="generator-v1-Field"></a>

### Field
Field represents a single column in the dataset.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the column (e.g., &#34;id&#34;, &#34;customer_name&#34;, &#34;transaction_amount&#34;). |
| provider | [string](#string) |  | Identifier for the data provider/generator (e.g., &#34;uuid&#34;, &#34;int_range&#34;, &#34;email&#34;). |
| options | [Option](#generator-v1-Option) | repeated | List of configuration parameters for the provider. |






<a name="generator-v1-GenerateRequest"></a>

### GenerateRequest
GenerateRequest defines the configuration for a generation job.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| schema_id | [string](#string) |  | Reference to a previously persisted schema ID in the database. |
| inline_schema | [Schema](#generator-v1-Schema) |  | Ad-hoc schema definition for immediate use without prior persistence. |
| count | [uint64](#uint64) |  | Total number of records to be generated. Supports massive volumes. |
| seed | [uint32](#uint32) |  | Seed for the pseudo-random number generator. Use the same seed with the same schema to achieve deterministic output. |






<a name="generator-v1-GeneratedRecord"></a>

### GeneratedRecord
GeneratedRecord represents a single row of generated data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [string](#string) | repeated | Ordered list of values corresponding to the schema fields. Values are converted to string for transport but follow the schema order. |






<a name="generator-v1-Option"></a>

### Option
Option provides type-safe configuration for field providers.
This design avoids runtime string parsing, optimizing CPU cycles.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | Parameter name (e.g., &#34;min&#34;, &#34;max&#34;, &#34;format&#34;, &#34;distribution&#34;). |
| string_val | [string](#string) |  | For patterns, formats, or categories. |
| int_val | [int64](#int64) |  | For discrete bounds and counts. |
| float_val | [double](#double) |  | For continuous bounds and weights. |
| bool_val | [bool](#bool) |  | For toggleable features (e.g., &#34;nullable&#34;). |






<a name="generator-v1-Schema"></a>

### Schema
Schema represents the blueprint of the dataset to be generated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Human-readable name for the schema. |
| fields | [Field](#generator-v1-Field) | repeated | Ordered collection of fields that define the dataset structure. |





 

 

 


<a name="generator-v1-GeneratorService"></a>

### GeneratorService
GeneratorService defines the capabilities for synthetic data orchestration.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StreamGenerate | [GenerateRequest](#generator-v1-GenerateRequest) | [GeneratedRecord](#generator-v1-GeneratedRecord) stream | StreamGenerate initiates a server-side streaming connection. The server will push generated records as they are produced until the requested count is reached or the context is cancelled. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

