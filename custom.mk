KEY_FILE=$(HOME)/key.bin

decrypt_conf:
	openssl enc -aes-256-cbc -d -md sha512 -pbkdf2 -iter 10000 -in resources/.gophoto-prod.enc -out resources/.gophoto-prod.yaml -pass file:$(KEY_FILE)

decrypt_realm:
	openssl enc -aes-256-cbc -d -md sha512 -pbkdf2 -iter 10000 -in ./resources/keycloak/gophoto-realm-prod.enc -out ./resources/keycloak/gophoto-realm-prod.json -pass file:$(KEY_FILE)

encrypt_conf:
	openssl enc -aes-256-cbc -md sha512 -pbkdf2 -iter 10000 -in resources/.gophoto-prod.yaml -out resources/.gophoto-prod.enc -pass file:$(KEY_FILE) && rm resources/.gophoto-prod.yaml

encrypt_realm:
	openssl enc -aes-256-cbc -md sha512 -pbkdf2 -iter 10000 -in ./resources/keycloak/gophoto-realm-prod.json -out ./resources/keycloak/gophoto-realm-prod.enc -pass file:$(KEY_FILE) && rm resources/keycloak/gophoto-realm-prod.json


