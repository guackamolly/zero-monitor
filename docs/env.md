# env

This document aims at explaining the meaning behind the environment values required to be set before running master and node binaries.

## Master

Master expects the .env file to be present in the working directory under the name `.env` or under the configuration directory under the name `master.env`.

|Environment Variable|Description|Example Value|Default Value|
|--------------------|-----------|-------------|-------------|
|`server_host`|Specifies the IP/Host name to bind the web-server listener connection|`0.0.0.0`|-|
|`server_port`|Specifies the port to bind the web-server listener connection|`8080`|-|
|`mq_sub_host`|Specifies the IP/host name to bind the ZeroMQ listener connection|`0.0.0.0`|-|
|`mq_sub_port`|Specifies the port to bind the ZeroMQ listener connection|`36113`|-|
|`mq_transport_pem_key`|Specifies the path of the `RSA 2048 PKCS 1` private key used for decrypting the communication during the key-exchange of node-master|`~/.config/zero-monitor/mq.pem`|`${CFG_DIR}/mq.pem`|
|`mq_transport_pub_key`|Specifies the path of the `RSA 2048 PKCS 1` public key used for encrypting the communication during the key-exchange of node-master|`~/.config/zero-monitor/mq.pub`|`${CFG_DIR}/mq.pub`|

---

Additionally you can set these variables to customize your experience while using the tool:

|Environment Variable|Description|Example Value|Default Value|
|--------------------|-----------|-------------|-------------|
|`bolt_db_path`|Specifies the path of the Bolt in-memory database used to store speedtests results|`~/.config/zero-monitor/master.db`|`${WORKING_DIR}/master.db`|
|`server_tls_crt_fp`|Specifies the path of the signed certificate to encrypt web-server communication (HTTPS)|`~/.config/zero-monitor/master.crt`|-|
|`server_tls_crt_key`|Specifies the path of the signed certificate private key to encrypt web-server communication (HTTPS)|`~/.config/zero-monitor/master.crt.key`|-|
|`server_virtual_host`|If you want to deploy the web server as a virtual path|`/zero-monitor`|-|

## Node

Node expects the .env file to be present in the working directory under the name `.env` or under the configuration directory under the name `node.env`.

|Environment Variable|Description|Example Value|Default Value|
|--------------------|-----------|-------------|-------------|
|`mq_sub_host`|Specifies the IP/host name of **master** ZeroMQ connection|`azure-proxy`|-|
|`mq_sub_port`|Specifies the port of **master** ZeroMQ connection |`36113`|-|
|`mq_transport_pub_key`|Specifies the path of the `RSA 2048 PKCS 1` public key used for encrypting the communication during the key-exchange of node-master|`~/.config/zero-monitor/mq.pub`|`${CFG_DIR}/mq.pub`|
|`mq_invite_code`|Specifies the invite code generated on the master dashboard, to join the network|`8aea958d-22da-4f8f-9521-2645f4ec497a`|-|