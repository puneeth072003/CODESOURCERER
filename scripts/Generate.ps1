$PROTO_DIR = "proto\protobuf"
$GENERATED_DIR = "proto\generated"

New-Item -ItemType Directory -Force -Path $GENERATED_DIR > $null 2>&1

protoc `
  --proto_path=$PROTO_DIR `
  --go_out=$GENERATED_DIR --go_opt=paths=source_relative `
  --go-grpc_out=$GENERATED_DIR --go-grpc_opt=paths=source_relative `
  $PROTO_DIR\*.proto