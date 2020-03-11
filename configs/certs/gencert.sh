#!/usr/bin/env sh

# Put the value you want here.
export ORG=madlab
# Put the value you want here.
export NAME=madprobe

cfssl gencert -initca ca-csr.json | cfssljson -bare $ORG-ca

cfssl gencert \
  -ca=$ORG-ca.pem \
  -ca-key=$ORG-ca-key.pem \
  -config=ca-config.json \
  -profile=server \
  server-csr.json | cfssljson -bare ${NAME}-server

cfssl gencert \
  -ca=$ORG-ca.pem \
  -ca-key=$ORG-ca-key.pem \
  -config=ca-config.json \
  -profile=client \
  client-csr.json | cfssljson -bare ${NAME}-client
