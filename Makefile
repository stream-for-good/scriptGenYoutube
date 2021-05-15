build:
	docker build -f ./docker/Dockerfile -t stream4good/scriptgenyoutube .
push:
	docker push stream4good/scriptGenYoutube
run:	
	docker run -p 8080:10001 stream4good/scriptgenyoutube
	@echo "litening on localhost:8080"
