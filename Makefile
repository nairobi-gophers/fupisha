SHELL:=/bin/bash                  
include .env

# ==============================================================================
# Running tests together with the staticcheck within the local computer
integration-test:
				@echo "++++ Run integration tests ++++"
				@CGO_ENABLED=0 go test -v ./api/v1/tests/ -count=1 
				@CGO_ENABLED=0 staticcheck ./...
				
unit-test:
		@echo "++++ Run unit tests ++++"
		@CGO_ENABLED=0 go test -v ./encoding/ -count=1 
		@CGO_ENABLED=0 staticcheck ./encoding/
		@CGO_ENABLED=0 go test -v ./provider/ -count=1
		@CGO_ENABLED=0 staticcheck ./provider/
		@CGO_ENABLED=0 go test -v ./store/postgres/ -count=1 
		@CGO_ENABLED=0 staticcheck ./store/postgres/
		


test:unit-test integration-test 

# ==============================================================================
	
up:
	@docker-compose up --detach --remove-orphans

down:
	@docker-compose down --remove-orphans
	
logs:
	@docker-compose logs -f 




