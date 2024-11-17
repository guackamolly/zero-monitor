# Server configuration
server_host=0.0.0.0
server_port=9090

## If exposing the server directly to the internet or behind a reverse proxy
server_tls_crt_fp=
server_tls_key_fp=

## If deploying as a virtual host
server_virtual_host=/
## Relative or absolute path to the directory that contains the views to serve
server_public_root=web/

# Message queue configuration

## If these variables aren't filled, pem/pub files are lookup on configuration directory
mq_transport_pem_key=
mq_transport_pub_key=

# Database configuration
bolt_db_path=master.db
