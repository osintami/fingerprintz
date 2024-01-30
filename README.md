# [![OSINTAMI Logo](./images/logo.png)OSINTAMI](https://www.osintami.com)

OSINTAMI is putting together data that can be used to track onliine entities.  Thousands of hours have gone into 
the making of this product and curating open source intelligence data and normalizing it into an easy to query
set of data items and rules.  Use cases include demandgen pipeline, cookie pool data enrichment, identity graphs, 
fraud detection, real time forum bot detection and prefiltering web traffic.  Any data can be imported or accessed 
remotely once data items are defined in the data dictionary with an appropriate connector.  CSV, JSON, TXT data can 
typicallbe imported with zero code changes.

## Architecture
There are three Go services that make up the OSINTAMI data server.  The OSINTAMI API Gateway, the ETLr,
the Normalized OSINT Data Server.

### API Gateway
Written in Go, the API Gateway is a reverse HTTPS proxy  It is configuration driven and integrates with Stripe.  User
ACLs are role based and once authenticated via API key the role is pass to the sub-service along with the request.

* go-chi/chi/v5 server
* configuration driven, specify URL mappings in/out with user role required, see config.json
* create users and roles using the REST interface, stored in Postgres
* audit users by API usage
* token counter and enforcer, returns 412 when well over the limit
* pixel fire API for website cookies and demandgen funnel tracking
* throttling using go-chi/httprate on a per API key basis
* Stripe payment integration (new account, monthly subscription paid, cancelled)
* integration with NODS possible to pre-filter out non-desireables (tor, cloud nodes, proxy, vpn, bot, blacklisted, etc.)


### ETLr
Written in Go, the ETLr manages scheduled data pulls from dozens of sources and normalizes the data 
into a consistent format for use in NODS for data items and rules.

* go-chi/chi/v5 server
* configuration driven, specify remote data sources to import and transform and publish
* internal cron for data refresh cycles, i.e. collect Tor data every hour, others daily or weekly
* mmdb, csv, yaml, txt input and output support
* builds data dictionary for each source as part of a ETL cycle
* publishes data and dictionary

#### Add New ETLr Job
1) See config.json for examples of supported input types
2) Look in associated source_<name>.go to see how they are parsed
3) Copy closest match source_<name>.go to a new name and edit accordingly
4) Define data collection items and code up parser
5) Disable all entries in config.json except the new entry (or create a new file)
6) Add new source to etl_manager.go under createInstance(source Source)
7) Fire up the ETLr in the debugger
8) Force a data collection run with http://127.0.0.1:8081/etlr/v1/refresh/{name}
9) cd etlr/etl
10) ./test.sh
11) check for failing tests and look for test coverage of your new source
12) create a test data set in etlr/etl/test/source/{name}.{type}
13) edit sources_all_test.go and add your source to TestSources(t *testing.T)
14) add your source to the etlr/etl/test/config.json file
15) run the tests again, debug, fix up and ensure 100% coverage of your new source
16) configure nods/config.json to include the new source
17) test your new data items http://localhost:8082/nods/v1/data/ip/{name}/{item}?ip=187.190.197.253
as found in the data/{name}.json item schema file

### Normalized OSINT Data Server
Written in Go, the NODS uses a URI to access data items, which is all documented in a data dictionary.  Rules
can be written using any of the data items and rules themselves are data items.

* go-chi/chi/v5 server
* configuration driven, specify data sources and access method (API, local mmdb file, code, etc.)
* hundreds of data items (374) from dozens data sources (94) currently supported
* pull data items by category key (IP address, email, domain, phone, browser)
* pull data items by data vendor for a single category
* pull data items from all vendors for a single category
* custom rules can be defined and accessed via a data items call and can be nested
* the data dictionary defines all data items and rules that can be accessed

### Whoami Server
Written in Go, the Whoami server generates JWT tokens for an entity based on several pieces of data.  A users
SHA256 lowercase email identifier, IP address, User-Agent, GEO, device identifier, partner identifier, and
fifty plus OSINT data signals make up a unique digital fingerprint for a user.  A user can have multiple 
fingerprints per application, typically two or three, one for each device.  This is the beginning of being able
to monitor identities across companies and products for the purpose of stopping fraud.  More to come on the 
identity co-op.

* go-chi/chi/v5 server
* API to fingerprint an identity
* API for re-checking the fingerprint and reporting risk factors
* API for feedback (login success/failure) -- COMING SOON-ISH

## Server Setup
### Ubuntu Development Server

The following steps will get an Ubuntu server up and running pretty quickly but this
is by no means a proper bullet proof production environment.

```
    ssh {admin}@api.osintami.com

    sudo apt-get update
    sudo apt install git

    # create osintami service user
    sudo groupadd --system osintami
    sudo useradd --system -g osintami osintami
    sudo usermod -s /bin/bash osintami
    sudo mkdir -p /home/osintami/{gateway etlr nods whoami monster data logs}
    sudo chown -R osintami:osintami /home/osintami

    # create SSL cert/key for your domain
    sudo snap install --classic certbot
    sudo ln -s /snap/bin/certbot /usr/bin/certbot
    sudo certbot certonly --standalone

    # change the owner of the SSL certs to the service user
    sudo chown -R osintami:osintami /etc/letsencrypt/live/api.osintami.com

    # set server timezone
    sudo dpkg-reconfigure tzdata

    # install and configure postgres
    sudo apt install postgresql postgresql-contrib
    service postgresql status
    # NOTE:  you will need this database and password to configure the API Gateway and Whoami services
    sudo -u postgres psql
        \password
        CREATE DATABASE osintami;
        \q

    # install go
    sudo apt install golang-go
```

### Configure github.com
```
    # create an approved cert for github.com and install it on github.com
    cd ~/.ssh
    ssh-keygen

    vi ~/.gitconfig

        [user]
            email = {email}@gmail.com
            name = {name}
        [init]
            defaultBranch = main
        [url "git@github.com"]
            insteadof = https://github.com
        [push]
            autoSetupRemote = true
```

### Clone Server Components
```
    mkdir -p go/src && cd go/src
    git clone git@github.com:osintami/fingerprintz.git
    cd fingerprintz
```

### API Gateway Setup
```
    # add postgres database and password and SSL key/cert files
    vi /home/osintami/gateway/.env
    POSTGRES_DB=osintami
    POSTGRES_HOST=127.0.0.1
    POSTGRES_HOST_AUTH_METHOD=trust
    POSTGRES_PASSWORD=xxx
    POSTGRES_PORT=5432
    POSTGRES_USER=postgres
    CERT_FILE=/etc/letsencrypt/live/api.osintami.com/fullchain.pem
    KEY_FILE=/etc/letsencrypt/live/api.osintami.com/privkey.pem
    LISTEN_ADDR=0.0.0.0:443
    LOG_LEVEL=INFO
    PATH_PREFIX=/
    EMAIL_ALERT_FROM=OSINTAMI@osintami.com
    EMAIL_ALERT_PASSWORD=xxx
    EMAIL_ALERT_SMTP_SERVER=smtp.gmail.com
    EMAIL_ALERT_SMTP_PORT=587
    OSINTAMI=http://127.0.0.1:8082/nods/v1/data/
    
    cd gateway
    ./deploy.sh

    sudo -u postgres psql
        \c osintami
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        INSERT INTO accounts VALUES (0, NOW(), NOW(), null, 'Administrator', 'admin@osintami.com', uuid_generate_v4(), 0, 'admin', null, true);
        SELECT * FROM accounts;
        \q

    # NOTE:  record api_key for admin account testing

    # start service and check for startup errors in the log
    sudo service osintami-gateway restart
    sudo journalctl -u osintami-gateway -f
    sudo cat /home/osintami/logs/gateway.log
    sudo systemctl enable osintami-gateway 

    # TEST:  should get an access denied error, not page not found
    curl https://api.osintami.com/data/items
    {"error":"not authorized"}

```

### ETLr Setup
```  
    # setup environment
    vi /home/osintami/etlr/.env

    MAXMIND_API_KEY=
    MAXMIND_ACCOUNT=
    UDGER_API_KEY=
    ABUSEIPDB_API_KEY=
    DBIP_API_KEY=
    IP2LOCATION_API_KEY=
    LOG_LEVEL=INFO
    GOGC=200
    TIDY_UP=true
    LISTEN_ADDR=127.0.0.1:8081
    PATH_PREFIX=/etlr
    TEMP_PATH=/tmp/
    DATA_PATH=/home/osintami/data/

    cd etlr
    ./deploy.sh

    # start service and check for startup errors in the log
    sudo service osintami-etlr restart
    sudo journalctl -u osintami-etlr -f
    sudo cat /home/osintami/logs/etlr.log
    sudo systemctl enable osintami-etlr

    # TEST: tormetrics.mmdb and tormetrics.json updated in the DATA_PATH /home/osintami/data/.
    curl http://127.0.0.1:8081/etlr/v1/refresh/tormetrics
    {"message":"sucess"}

    # PROD:  override schedule and pull in latest OSINT worldwide data now
    curl http://127.0.0.1:8081/etlr/v1/refresh/all
```

### Normalize OSINT Data Server Setup
```
     # setup environment
    vi /home/osintami/nods/.env

    PWNED_API_KEY=
    IPINFO_API_KEY=
    IP2LOCATION_API_KEY=
    DBIP_API_KEY=
    APIVOID_API_KEY=
    IPQS_API_KEY=
    XCONNECT_USER=
    XCONNECT_PASS=

    LOCAL_DB_PATH=/home/osintami/data/
    LISTEN_ADDR=127.0.0.1:8082
    LOG_LEVEL=INFO
    PATH_PREFIX=/nods

    cd nods
    ./deploy.sh

    # start service and check for startup errors in the log
    sudo service osintami-nods restart
    sudo journalctl -u osintami-nods -f
    sudo cat /home/osintami/logs/nods.log
    sudo systemctl enable osintami-nods

    # test the dictionary
    http://127.0.0.1:8082/nods/v1/data/items?role=admin

    [
        {
            "Item": "browser/useragent/deviceBrand",
            "Enabled": true,
            "GJSON": "Device.Brand",
            "Description": "Device brand",
            "Type": "String"
        }
        ...
    ]
```

#### API - Data Dictionary
```
 https://api.osintami.com/data/items?key=<admin api key>&csv=true

```

#### API - Specific Item
```
https://api.osintami.com/data/ip/google/cloud.isGoogle?ip=34.173.187.95

{
    "Item": "google.cloud.isGoogle",
    "Result": {
        "Type": 0
    },
    "Keys": {
        "ip": "34.173.187.95"
    },
    "Error": "item does not exist"
}
```

#### API - Rule Based Item

Executes the following rule, made up of multiple cloud data sources to determine if
the specified IP address is a cloud node.  This is useful as part of a more sophiscated bot
check rule.  Rules can invoke other rules, making this API pretty flexible.  Existing rules
can be found in the data dictionary.

```
{
    "Item": "rule/osint/isCloudNode",
    "Enabled": true,
    "GJSON": "isCloudNode",
    "Query": "[ip/udger/cloud.isCloudNode] || [ip/google/cloud.isGoogle] || [ip/digitalocean/cloud.isDigitalOcean] || [ip/amazon/cloud.isAmazon] || [ip/oracle/cloud.isOracle] || [ip/azure/cloud.isAzure] || [ip/cloudflare/cloud.isCloudflare]",
    "Description": "This IP address belongs to a cloud provider.",
    "Type": "Boolean"
}
```

```
https://api.osintami.com/data/rule/osint/isCloudNode?ip=34.173.187.95

{
    "Item": "rule/osintami/isCloudNode",
    "Result": {
        "Type": 1,
        "Bool": true
    },
    "Keys": {
        "ip": "34.173.187.95"
    },
    "Error": ""
}

```
