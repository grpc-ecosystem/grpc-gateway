| Title                                  | Category      |
| -------------------------------------- | ------------- |
| Use go templates in protofile comments | Documentation |

# Use go templates in protofile comments

Use [Go templates](https://golang.org/pkg/text/template/ "Package template") in your protofile comments to allow more advanced documentation such as:  
* Documentation about fields in the proto objects.  
* Import the content of external files (such as [Markdown](https://en.wikipedia.org/wiki/Markdown "Markdown Github")). 

## How to use it

By default this function is turned off, so if you want to use it you have to set the ```use_go_templates``` flag to true inside of the ```swagger_out``` flag.
```bash
--swagger_out=use_go_templates=true:.
```

## Example

Example of a bash script with the ```use_go_templates``` flag set to true:

```bash
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
  --go_out=plugins=grpc:. \
  --grpc-gateway_out=logtostderr=true:. \
  --swagger_out=logtostderr=true,use_go_templates=true:. \
  *.proto 
```

For an example of a protofile which has Go templates enabled, [click here](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/examples/proto/examplepb/use_go_template.proto "Example protofile with Go template").
