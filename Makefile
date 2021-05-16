SHELL:=/bin/bash                    
include .env

# ==============================================================================
# Running tests together with the staticcheck within the local computer

test:
	CGO_ENABLED=0 go test -v ./... -count=1
	CGO_ENABLED=0 staticcheck ./...

# ==============================================================================
	
up:
	@docker-compose up --detach --remove-orphans

down:
	@docker-compose down --remove-orphans
	
logs:
	@docker-compose logs -f 




