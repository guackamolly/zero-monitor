# transport-security

This document covers the security measures applied to secure the transport layers used in the tool.

## Master HTTP Server

If you want to expose the master node http server behind a reverse proxy or directly to the internet, you should start a TLS server. This can be done by using a self certificate and updating the `.env` file.

```txt
...
server_tls_crt_fp=<self-certificate.crt>
server_tls_key_fp=<self-certificate.key>
...
```

(Alternatively you can use the `init` program to generate the certificate files in the configuration directory. More info [here](init.md)).

## Pub/Sub Server

Master and nodes exchange messages in a pub/sub manner through a [ZeroMQ TCP channel](http://api.zeromq.org/4-2:zmq-tcp). This
channel does not operate with TLS, but instead encrypts sensitive messages before being publishing them in the channel. The
reasoning behind this model is that ZeroMQ is not directly compatible with TLS and encrypting all messages leads to CPU exhaustion on low-end devices.

Sensitive messages are encrypted using an `AES-256 GCM` cipher block, which key is exchanged when a node tries to join the network. This key is encrypted using a `RSA 2048 PKCS 1` public key that must be derived from a private key known in the master server. The easiest way to generate these keys is using the `init` program, but alternatively you can use *openssl* or another trusted source to generate the private + public key and updating the `.env` file.

```txt
...
mq_transport_pem_key=<private-key.pem>
mq_transport_pub_key=<public-key.pub>
...
```

**Note**: Master only needs the **private** key file, where as **nodes** only the public key file.

---

Both symmetric key and nonce are generated from a secure random source, and are attached on the `Metadata` field of the encrypted message. The following diagram helps visualizing the key exchange protocol.

![sequence diagram that illustrates how the key that is used to encrypt sensitive messages is exchanged between nodes](static/key-exchange.svg)
