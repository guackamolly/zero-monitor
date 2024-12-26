# 0.0.0-3ceb9b7d95e9b7974a376d170e06433897b40a51

- (feat): automatically reconnect nodes if master restarts
- (fix): not being able to view detailed information of nodes on webkit browsers

# 0.0.0-79735a0f2d2817f9c0dbf974b5acd100c2c5c92d

- (chore): disable packages route for guest viewers

# 0.0.0-34a44199aa82296d3c3347fe4b1431bcb7657ee8

- (feat): autostart agents after server restarts
- (fix): redirect to last visited page after sign-in

# 0.0.0-9dd1f1907eec81bb3251a2f7001fa5c2f0fee754

- (BREAKING) nodes are now identifies using an UUID generated once during the first startup (saved in the config folder)
- add support for multiple VPS
- display master agent version in the home page
- add support for removing nodes from the network
- correct sign in button wording (before was "Create" now is "Sign-in")