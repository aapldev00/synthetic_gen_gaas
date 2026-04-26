# Variables
PROTO_DIR=api/proto
DOC_DIR=docs
PKG_DIR=pkg/genproto
# This must match the name in your 'go mod init'
MODULE_NAME=github.com/aapldev00/synthetic_gen_gaas

# Main command to generate code and documentation
proto:
	@echo "Generating Go code and documentation..."
	@mkdir -p $(PKG_DIR) $(DOC_DIR)
	protoc --go_out=. --go_opt=module=$(MODULE_NAME) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE_NAME) \
		--doc_out=./$(DOC_DIR) --doc_opt=markdown,PROTO_DOCUMENTATION.md \
		$(PROTO_DIR)/*.proto
	@echo "Done! Files generated in $(PKG_DIR) and documentation in $(DOC_DIR)"

# Clean up generated files
clean:
	@echo "Cleaning up generated files..."
	rm -rf $(PKG_DIR)/* $(DOC_DIR)/*
	@echo "Cleanup complete."

.PHONY: proto clean