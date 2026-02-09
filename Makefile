.PHONY: start

start:
	docker start pgadmin
	docker start postgres-db
	docker start localstack