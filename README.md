# PyLint
[![Build Status](http://ci.patwie.com/api/badges/PatWie/pylint/status.svg)](http://ci.patwie.com/PatWie/pylint)

A small webservice written in Go to lint python projects hosted at GitHub using pyflakes.
The main advantages of this service are:
- very lightweight
- very easy to setup in GitHub
- minimal permissions are required in GitHub with per-repository permission (all advantages from Github Apps)
- can be self-hosted

See the [badge-branch](https://github.com/PatWie/pylint/tree/badge) for the refactoring effort. This is my first GOlang project, it is a safe bet that I break many Go-idioms. Don't start to count them.

## Install

The essential steps are:

1. Register GitHub Application (**not** an OAuth App)
2. prepare pylint
3. start application using docker-compose


### Register GitHub Application

I assume you will host this service at `http://pylint.domain.com`.

Go to [https://github.com/settings/apps](https://github.com/settings/apps) and add a `New GitHub App`. 
You need to fill out:
- GitHub App Name: arbitrary name
- Homepage URL: http://pylint.domain.com
- User authorization callback URL: http://pylint.domain.com/auth-callback
- Webhook URL: http://pylint.domain.com/hook
- Webhook secret: some_random_strings

The Webhook secret is **not** optional for this webservice. You will need to note the `ID` (INTEGRATIONID) in the `About`-Section. Generate a `Private key` and download the key. Name it like `my-pylint-key.pem`. We do not need OAuth credentials.

### Prepare pylint

I propose to clone this repository

    git clone https://github.com/PatWie/pylint.git
    cd pylint
    cp .env.example .env
    edit .env

A few notes on `.env`:
- `PYLINTGO_PUBLICPORT` is the exposed port which can act as an endpoint in an NGINX reverse proxy
- `PYLINTGO_GITHUB_ADMINID` should remain 0 if you want to not restrict your github app, otherwise just use your github user-id. In this case, only your user account is allowed to install the applications (linting will work for any commit in that git-repo)
- `PYLINTGO_GITHUB_INTEGRATIONID`, see `ID` in the `About`-Section on the GitHub page
- `PYLINTGO_GITHUB_SECRET` is your Webhook secret
- `PYLINTGO_GITHUB_KEYPATH` should be `/keys/my-pylint-key.pem`
- `PYLINTGO_STATEURL` should be something like `http://pylint.domain.com/`. The trailing `/` is important.
- `PYLINTGO_NAME` any name

Now, copy the downloaded `my-pylint-key.pem` to `pylint/keys`. Make sure the `keys` directory can be accessed by the docker user and the `my-pylint-key.pem` is readable as well for the docker-user.

### start applications

    docker-compose up

That's all! This might take a while. Point your browser to `http://pylint.domain.com/home` to verify your installation. It should show you

    {"Msg":"PyLint Go is successfully running"}
