---
category: documentation
---

# AWS

## Import swagger documentation into AWS API Gateway
AWS Api gateway service allows to import swagger documentation to create REST API. The process is very straigh foward and can be found [here](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-import-api.html).
These are some considerations when importing the documentation:

1. Remove model circular dependencies(Circular dependencies aren't supported by the parser).
2. Remove security-related annotations (These annotations aren't well supported by the parser).
3. Max length of fields is reviewed by the parser but the errors aren't self-explanatory. Review the [documentation](https://swagger.io/specification/v2/) to verify that the requirements are met.
4. API gateway errors aren't good enough for troubleshooting so I recommend using this [page](https://apidevtools.org/swagger-parser/online/) for structure validation.
