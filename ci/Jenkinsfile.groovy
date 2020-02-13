#!/usr/bin/groovy
pipeline {
    agent {
        docker { image 'golang:1.13' }
    }

   stages {
      stage('Build') {
         steps {
            // Get the code from GitHub repository
            git 'https://github.com/jbleduigou/budgetcategorizer.git'
            sh 'mkdir -p $GOPATH/src/github/jbleduigou'
            sh 'ln -s $WORKSPACE $GOPATH/src/github/jbleduigou/budgetcategorizer'

            // Run Maven on the amazon linux agent.
            sh 'export XDG_CACHE_HOME=/tmp/.cache && cd $GOPATH/src/github/jbleduigou/budgetcategorizer && go build -o budgetcategorizer ./cmd/budgetcategorizer'
         }
      }
      stage('Test') {
         steps {
            sh 'export XDG_CACHE_HOME=/tmp/.cache && cd $GOPATH/src/github/jbleduigou/budgetcategorizer && make test'
         }
      }
      stage('Linting') {
         steps {
            sh 'export XDG_CACHE_HOME=/tmp/.cache && go get -u golang.org/x/lint/golint'
            sh 'export XDG_CACHE_HOME=/tmp/.cache && cd $GOPATH/src/github/jbleduigou/budgetcategorizer && go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status'
         }
      }
   }
}
