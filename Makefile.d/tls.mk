#
# Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

CERT_DIR ?= $(ROOTDIR)/internal/test/data/tls
DOMAIN ?= $(NAME).$(ORG).org
EMAIL ?= $(NAME)@$(ORG).org
STREET_ADDR ?= 1-3 Kioicho, Tokyo Garden Terrace Kioicho Tower, Chiyoda-ku, Tokyo 102-8282, Japan

CA_CN ?= $(DOMAIN) Root CA
CA_DAYS ?= 3650
END_ENTITY_DAYS ?= 825

CA_KEY := $(CERT_DIR)/ca.key
CA_CRT := $(CERT_DIR)/ca.crt
CA_PEM := $(CERT_DIR)/ca.pem
CA_SRL := $(CERT_DIR)/ca.srl
SERVER_KEY := $(CERT_DIR)/server.key
SERVER_CSR := $(CERT_DIR)/server.csr
SERVER_CRT := $(CERT_DIR)/server.crt
CLIENT_KEY := $(CERT_DIR)/client.key
CLIENT_CSR := $(CERT_DIR)/client.csr
CLIENT_CRT := $(CERT_DIR)/client.crt

INVALID_CA := $(CERT_DIR)/invalid-ca.pem
INVALID_SERVER_CRT := $(CERT_DIR)/invalid-server.crt

.PHONY: certs/gen
## generates SSL Certrificates for Testing
certs/gen: certs/clean $(CERT_DIR) $(SERVER_CRT) $(CLIENT_CRT) $(CA_PEM) $(INVALID_CA) $(INVALID_SERVER_CRT) certs/tidy certs/verify
	@echo "✅ Certificates generated in '$(CERT_DIR)'"

$(CERT_DIR):
	mkdir -p $@

$(CA_KEY):
	openssl genpkey -algorithm ed25519 -out $(CA_KEY)

$(CA_CRT): $(CA_KEY)
	openssl req -x509 -new -nodes \
	-key	$(CA_KEY) \
	-sha256 \
	-days	$(CA_DAYS) \
	-subj	"/C=JP/ST=Tokyo/L=Chiyoda-ku/streetAddress=$(STREET_ADDR)/O=$(ORG)/CN=$(CA_CN)/emailAddress=$(EMAIL)" \
	-addext "basicConstraints = critical,CA:TRUE,pathlen:1" \
	-addext "keyUsage = critical,keyCertSign,cRLSign" \
	-addext "subjectKeyIdentifier = hash" \
	-addext "authorityKeyIdentifier = keyid:always" \
	-addext "authorityInfoAccess = OCSP;URI:http://ocsp.$(DOMAIN)" \
	-addext "crlDistributionPoints = URI:http://crl.$(DOMAIN)/ca.crl" \
	-out	$(CA_CRT)

$(CA_PEM): $(CA_CRT) $(CA_KEY)
	cat $^ > $@
	chmod 600 $@

$(SERVER_KEY):
	openssl genpkey -algorithm ed25519 -out $(SERVER_KEY)

$(SERVER_CSR): $(SERVER_KEY)
	openssl req -new -batch \
	-key	$(SERVER_KEY) \
	-subj "/C=JP/ST=Tokyo/L=Chiyoda-ku/streetAddress=$(STREET_ADDR)/O=$(ORG)/OU=Backend/CN=$(DOMAIN)/emailAddress=$(EMAIL)" \
	-addext "subjectAltName = DNS:$(DOMAIN),IP:127.0.0.1" \
	-addext "basicConstraints = critical,CA:FALSE" \
	-addext "keyUsage = critical,digitalSignature,keyEncipherment" \
	-addext "extendedKeyUsage = serverAuth" \
	-addext "subjectKeyIdentifier = hash" \
	-out $(SERVER_CSR)

$(SERVER_CRT): $(SERVER_CSR) $(CA_CRT)
	printf '%s\n' \
	'[ v3_req ]' \
	'subjectAltName = DNS:$(DOMAIN),IP:127.0.0.1' \
	'basicConstraints = critical,CA:FALSE' \
	'keyUsage = critical,digitalSignature,keyEncipherment' \
	'extendedKeyUsage = serverAuth' \
	'subjectKeyIdentifier = hash' \
	'authorityKeyIdentifier = keyid,issuer' \
	'authorityInfoAccess = OCSP;URI:http://ocsp.$(DOMAIN)' \
	'crlDistributionPoints = URI:http://crl.$(DOMAIN)/ca.crl' \
	| openssl x509 -req \
	-in	$(SERVER_CSR) \
	-CA	$(CA_CRT) \
	-CAkey	$(CA_KEY) \
	-CAcreateserial \
	-days	$(END_ENTITY_DAYS) \
	-sha256 \
	-extfile	/dev/stdin \
	-extensions v3_req \
	-out	$(SERVER_CRT)

$(CLIENT_KEY):
	openssl genpkey -algorithm ed25519 -out $(CLIENT_KEY)

$(CLIENT_CSR): $(CLIENT_KEY)
	openssl req -new -batch \
	-key	$(CLIENT_KEY) \
	-subj "/C=JP/ST=Tokyo/L=Chiyoda-ku/streetAddress=$(STREET_ADDR)/O=$(ORG)/OU=Frontend/CN=client.$(DOMAIN)/emailAddress=$(EMAIL)" \
	-addext "subjectAltName = email:$(EMAIL)" \
	-addext "basicConstraints = critical,CA:FALSE" \
	-addext "keyUsage = critical,digitalSignature" \
	-addext "extendedKeyUsage = clientAuth" \
	-addext "subjectKeyIdentifier = hash" \
	-addext "authorityInfoAccess = OCSP;URI:http://ocsp.$(DOMAIN)" \
	-addext "crlDistributionPoints = URI:http://crl.$(DOMAIN)/ca.crl" \
	-out $(CLIENT_CSR)

$(CLIENT_CRT): $(CLIENT_CSR) $(CA_CRT)
	printf '%s\n' \
	'[ v3_req ]' \
	'subjectAltName	= email:$(EMAIL)' \
	'basicConstraints = critical,CA:FALSE' \
	'keyUsage = critical,digitalSignature' \
	'extendedKeyUsage = clientAuth' \
	'subjectKeyIdentifier = hash' \
	'authorityKeyIdentifier = keyid,issuer' \
	'authorityInfoAccess = OCSP;URI:http://ocsp.$(DOMAIN)' \
	'crlDistributionPoints = URI:http://crl.$(DOMAIN)/ca.crl' \
	| openssl x509 -req \
	-in	$(CLIENT_CSR) \
	-CA	$(CA_CRT) -CAkey $(CA_KEY) -CAcreateserial \
	-days	$(END_ENTITY_DAYS) \
	-sha256 \
	-extfile	/dev/stdin -extensions v3_req \
	-out	$(CLIENT_CRT)

$(INVALID_CA): $(CLIENT_CRT)
	@echo "→ Generating invalid CA bundle (client.crt as CA)"
	cp $< $@

$(INVALID_SERVER_CRT): $(SERVER_CSR)
	@echo "→ Generating completely invalid server.crt (copying CSR as-is)"
	cp $< $@

.PHONY: certs/verify
## verify SSL Certrificates for Testing
certs/verify: $(CA_PEM) $(SERVER_CRT) $(CLIENT_CRT)
	@echo "🔍 Verifying certificates..."
	@openssl verify -CAfile $(CA_PEM) \
	$(CA_PEM) \
	$(SERVER_CRT) \
	$(CLIENT_CRT) \
	&& echo "✅ All certificates are valid" \
	|| (echo "❌ Certificate verification failed" && false)

.PHONY: certs/tidy
## Tidy SSL Certrificates for Testing
certs/tidy:
	@echo "🧹 Removing intermediate files..."
	@rm -f $(CA_KEY) \
	$(CA_CRT) \
	$(CA_SRL) \
	$(SERVER_CSR) \
	$(CLIENT_CSR)

.PHONY: certs/clean
## Cleanup & Remove All SSL Certrificates for Testing from CERT Dir
certs/clean:
	@echo "🗑️ Removing entire '$(CERT_DIR)' directory"
	rm -rf $(CERT_DIR)