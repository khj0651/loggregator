# FlashFlood

Discover your installation's Loggregator capacity by flooding the system with log messages until it breaks.

## How it works

This app deploys to Cloud Foundry. It emits log messages – and reads its own log stream – to detect message loss in the Loggregator system.

It begins in an "idle" state. When started, it emits logs at progessively faster rates until stopped. At each stage, records are kept of the number of logs sent and received.

## Deploying
1. Create a user (or reuse an existing one) in your Cloud Foundry installation.
1. `cf push flashflood --no-start` to push code without starting (as environment variables must be set first)
1. Set your environment variables: 
```
cf set-env flashflood USERNAME [username]
cf set-env flashflood PASSWORD [password]
cf set-env flashflood UAA_URL [url of UAA]
cf set-env flashflood LOGGREGATOR_URL [url of loggregator]
```
1. `cf start flashflood` to stage and run the app with the given variables. Note: do **not** scale up, as the log streams from all instances would be comingled and produce erroneous results.

## Using

1. To begin flooding, issue a GET request (or visit in a web browser) `[flashflood-url]/start`.
1. To see results (updated as the test progresses), visit `[flashflood-url]/results`.
1. To stop the test at any time, GET `[flashflood-url]/stop`. The results will be kept until the next `/start` request.

## Analyzing results

