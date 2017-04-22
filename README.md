# DASH
A circle ci dashboard written in Go/React that displays your followed builds with metadata.

## Getting Started
1. Set an environment variable:
    ```bash
    export CIRCLE_CI_AUTH_TOKEN={my_ci_token}
    ```
    You can find this under your profile settings/api tokens.

2. Build binary:
    ```bash
    go build
    ```

3. Run:
    ```bash
    ./GoDash
    ```
    If you do not set the environment variable, the program will fail

4. Navigate to `localhost:8080/dash`

    This will list all your followed projects and current build info
