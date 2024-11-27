# TODO

Looking to contribute to the project but you don't know how to start? The todo's list is an excelent starting point! Pick any of the actions and then create an issue with the label "TODO" so other maintainers know you will work on it.

---

- Setup screenshot tool and integrate it in CI
- Require mobile users to change to landscape mode before viewing big tables
- Close Pub/Sub event channels after node finishes a task
- Fix different node events being published to the same subscriber (UUID validation)
- Organize common template code (e.g., `error.gohtml`, websocket initialization, meta tags)
- Optimize network view template render
- Speedtests is not saved if client breaks Websocket connection (unbuffered channel not being consumed)
- Remove depending on environment variables (100% portability)
- Add initialization scripts for Windows
- Make Linux/macOS initialization scripts POSIX compliant (remove bash, use shell) 
- Cover comment //TODOs
- Fix init not printing build version

implement node-master handshake
implement node-master handshake from stdin
implement init node bootstrap from stdin invite url