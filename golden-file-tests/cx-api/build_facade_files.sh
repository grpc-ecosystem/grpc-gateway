# Directory containing the protobuf files
mkdir -p proto
cp -a src/** proto/
cp -a deps/** proto/
proto_dir="proto"
go_out_dir="internal"
mod_prefix="github.com/coralogix"
mod_name="$mod_prefix/openapi-facade/go"
proto_files=($(find "$proto_dir" -name "*.proto" -print))
openapi_args=""

# Build arguments for import paths of all modules
for proto_file in "${proto_files[@]}" 
do
    out_module=$(dirname $proto_file)
  
    if [[ $out_module == *"coralogix"* ]]; then
        mod_path="${out_module##*/com/}"
	# For all other protos, the package path is the same as the directory path
	openapi_args+="--openapiv3_opt=M${proto_file##*$proto_dir/}=${mod_name}/${go_out_dir}/${mod_path} "
    fi
done

protofile_list=""

for proto_file in "${proto_files[@]}" 
do
    protofile_list+="${proto_file} "
done

protoc --proto_path=$proto_dir --openapiv3_out=.. --openapiv3_opt=allow_merge=true,ignore_additional_bindings=true,openapi_naming_strategy=simple $openapi_args $protofile_list
