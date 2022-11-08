.PHONY: decrypt_conf decrypt_realm encrypt_conf encrypt_realm
KEY_FILE=$(HOME)/key.bin

decrypt_conf:
	openssl enc -aes-256-cbc -d -md sha512 -pbkdf2 -iter 10000 -in resources/.gophoto-prod.enc -out resources/.gophoto-prod.yaml -pass file:$(KEY_FILE)

decrypt_realm:
	openssl enc -aes-256-cbc -d -md sha512 -pbkdf2 -iter 10000 -in ./resources/keycloak/gophoto-realm-prod.enc -out ./resources/keycloak/gophoto-realm-prod.json -pass file:$(KEY_FILE)

encrypt_conf:
	openssl enc -aes-256-cbc -md sha512 -pbkdf2 -iter 10000 -in resources/.gophoto-prod.yaml -out resources/.gophoto-prod.enc -pass file:$(KEY_FILE) && rm resources/.gophoto-prod.yaml

encrypt_realm:
	openssl enc -aes-256-cbc -md sha512 -pbkdf2 -iter 10000 -in ./resources/keycloak/gophoto-realm-prod.json -out ./resources/keycloak/gophoto-realm-prod.enc -pass file:$(KEY_FILE) && rm resources/keycloak/gophoto-realm-prod.json

.PHONY: fake.session
fake.session: SESSION='{"user":{"id":"userid","username":"toto","first_name":"Toto","last_name":"Toto","role":"admin","can_share":true},"session_id":"session_id"}'
fake.session:
	@echo -n $(SESSION) | base64 -w0

