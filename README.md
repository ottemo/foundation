Ottemo
=========

[![Go Report Status](https://goreportcard.com/badge/github.com/ottemo/commerce)](https://goreportcard.com/report/github.com/ottemo/commerce)

a small, wicked fast ecommerce platform built for gophers

## Install and Setup

Commerce supports golang versions 1.7+ and requires you to [install the lastest stable version of dep](https://golang.github.io/dep/docs/installation.html) to manage installation of dependencies. Before you run commerce for the first time, you need to create an empty database and create a user which will have permissions to create/update/insert/delete. 

Download commerce to your desired GOROOT 
```
# on my local GOROOT is found at ~/code/go
> cd ~/code/go/src/github.com

# create a directory for ottemo
> mkdir ottemo

# clone the repository
> git clone https://github.com/ottemo/commerce
```

Install and update dependencies for commerce
```
# change into the commerce directory
> cd commerce

# run dep ensure to install dependencies
> dep ensure
```

Build commerce with support for MySQL. There is a build script which adds in extra information, but this is not required.
```
# there is a build script in the bin directory
> bin/make.sh -tags mysql 
```

Note, you may also build commerce using the `go build` command, but remember to tell it which version of the database drivers to use. For instance - `go build -tags mysql`

The final step before running commerce for the first time is to configure a few basic settings. There exists a sample configuration file which follows the ini file format. 
```
# copy the sample file
> cp ottemo.sample.ini ottemo.ini
```

Now edit the file with your favorite editor. Here is a basic version of the file for MySQL. Change the values USER and PASSWORD according to the credentials you created when first setting up the empty database. 
```
; minimal settings for mysql driver
db.mysql.db=commerce
db.mysql.maxConnections=50
db.mysql.poolConnections=10
db.mysql.uri=USER:PASSWORD@/commerce

; let commerce know where to find/store images and videos
media.fsmedia.folder=./media/
media.resize.images.onfly=true

; if you are doing development and not using HTTPS to access commmerce set this to false 
secure_cookie=false

; to allow cross domain cookies set your desired domain (we do use cross domain cookies)
xdomain.master=http://*.local.dev/
```

Now you may run commerce from the commandline
```
> ./commerce
```

You should a message similiar to the following  printed to stdout:
```
2019-04-25T15:57:45Z Connecting to MySQL DB. Timeout: 10 seconds.
Ottemo v1.4.5-jwv_basic_build_run_instructions_497-b1489 (mysql) [2019-04-25T08:57:18-07:00]
REST API Service [HTTPRouter] starting to listen on :3000
2019-04-25T15:57:45Z DB connection established.
```

Commerce will create a `var` folder which will contain logs folder and a session folder, since we didn't compile it to use redis for sessions, (which is typical for local development).


## License

[Mozilla Publice License 2.0](LICENSE.md)
## Terms and Conditions

All Submissions you make to Ottemo, Inc. (“Ottemo”) through GitHub are subject to the following terms and conditions:

1. You grant Ottemo a perpetual, worldwide, non-exclusive, no charge, royalty free, irrevocable license under your applicable copyrights and patents to reproduce, prepare derivative works of, display, publicly perform, sublicense and distribute any feedback, ideas, code, or other information (“Submission”) you submit through GitHub.
2. Your Submission is an original work of authorship and you are the owner or are legally entitled to grant the license stated above.

### Copyright
© 2019 Ottemo, Inc.

