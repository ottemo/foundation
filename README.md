Ottemo 
=========

built for gophers

[![wercker status](https://app.wercker.com/status/97369a2b891e2ff6dd5b37d96301030f/m "wercker status")](https://app.wercker.com/project/bykey/97369a2b891e2ff6dd5b37d96301030f)

[![GoDoc](https://godoc.org/github.com/ottemo/foundation?status.png)](https://godoc.org/github.com/ottemo/foundation)

## Install and Setup coming....

## Contribute to Ottemo development
We use git-flow internally, but if you do not like git-flow you may use [this document](CONTRIBUTE.md) as an alternative.  

Below is a mini quickstart if you are new to git-flow and can't wait to jump into the code. 

### Initialize git-flow

    # fork or clone ottemo like below
    $ git clone https://github.com/ottemo/ottemo-go.git 

    # init git-flow, (git-flow must be installed for your OS locally)
    $ git checkout master
    $ git checkout develop
    $ git flow init -d

### Start a feature branch
    $ git flow feature start <FEATURE-NAME>

### Issue a pull request on github
    $ git push -u origin <FEATURE-BRANCH>
    # if you have git aliased to hub otherwise use the github web interface
    $ git pull-request -b develop

### Delete the local branch
    $ git branch -d <FEATURE-BRANCH>

### How start with Vagrantfile
Clone ottemo/foundation github repo (current bug with ottemo.ini)

    vagrant up
    vagrant ssh
    sudo su -
    export GOPATH=/opt/go
    cd $GOPATH/src/github.com/ottemo/foundation/
    go run main.go

### How to run ottemo/foundation docker container
Pull latest image from docker hub

    docker pull ottemo/foundation

Start the container and access locally (currently set to use sqlite in the container - development use only)

    docker run -d -p 3000:3000 -t ottemo/foundation

### How to proxy foundation with nginx for HTTP/HTTPS traffic
It is suggested to secure Foundation API Server with SSL.  To offload the 
added work to maintain high performance on Foundation, we will use nginx to
 proxy HTTPS.  In production, only allow connections to Foundation over SSL.   
[This gist](https://gist.github.com/vastbinderj/b5e5fa2acfd199d48fa5) explains 
to create the certificate and configure nginx.

## License

[MIT License](http://mit-license.org/) copyright 2014, Ottemo

## Terms and Conditions

All Submissions you make to Ottemo, Inc. (“Ottemo”) through GitHub are subject to the following terms and conditions: 

1. You grant Ottemo a perpetual, worldwide, non-exclusive, no charge, royalty free, irrevocable license under your applicable copyrights and patents to reproduce, prepare derivative works of, display, publically perform, sublicense and distribute any feedback, ideas, code, or other information (“Submission”) you submit through GitHub. 
2. Your Submission is an original work of authorship and you are the owner or are legally entitled to grant the license stated above.
