## Install

The essential steps are:

1. Register your GitHub Application (**not** an OAuth App)
2. edit `pylint-configuration.yml`
3. start application

### Register GitHub Application

I assume you will host this service at `http://pylint.domain.com`.

Go to [https://github.com/settings/apps](https://github.com/settings/apps) and add a `New GitHub App`.
You need to fill out:
- GitHub App Name: arbitrary name
- Homepage URL: http://pylint.domain.com
- Webhook URL: http://pylint.domain.com/hook
- Webhook secret: <some-random-strings-as-a-serect>

The web-hook secret `<some-random-strings-as-a-serect>` is **not** optional for this web-service. You will need to note the `ID` (INTEGRATIONID) in the `About`-Section. Generate a `Private key` and download the key. We do not need OAuth credentials.

### Start Application: Pre-Build binaries

Download the latest release and edit `pylint-configuration.yml` according to your setup.

### Start Application: Docker-Compose

I propose to clone this repository

    git clone https://github.com/PatWie/pylint.git
    cd pylint
    cp pylint-configuration.example.yml pylint-configuration.yml
    edit pylint-configuration.yml

After editing `pylint-configuration.yml`, just run

    docker-compose up

That's all! This might take a while. Point your browser to `http://localhost:8080` to verify your installation. It should show you

    active

For docker-compose setups I suggest to put the key-file into `keys` and use the path "/keys/<keyfilename>"
Make sure, the `my-pylint-key.pem` is readable as well for the docker-user.

## Configuration

A few notes on `pylint-configuration.yml`:

```yaml
github:
  integration_id: 00000                        # see `ID` in the `About`-Section on the GitHub page
  secret: "<some-random-strings-as-a-serect>"  # is your webhook secret
pylint:
  name: PyLint                                 # any name (will be displayed next to the commits)
  port: 8080                                   # is the exposed port which can act as an endpoint in an NGINX reverse proxy
  url: http://pylint.domain.com                # should be something like `http://pylint.domain.com` without trailing `/`
  reports_path: /data/reports                  # path where the text-files of the reports should be saved
  key_file: /keys/private-key.pem              # path to your private key of the GitHub Application
  database_file: ./pylint.db                   # path to sqlite database file
redis:
  host: redis                                  # hostname containing the redis service (commonly "localhost" or "redis")
  port: 6379                                   # redis port (default is 6379)

```
