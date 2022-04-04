run :
	go run cmd/main.go

docker :
	docker image build -f dockerfile . -t imagename
	docker container run -p 9000:8000 -d --name forum imagename

clean :
	docker system prune