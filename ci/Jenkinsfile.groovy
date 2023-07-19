#!/usr/bin/groovy
pipeline {
    agent {
        docker { image 'golang:1.20' }
    }
    environment {
        XDG_CACHE_HOME='/tmp/.cache'
        GOOS='linux'
        GOARCH='amd64'
    }
    stages {
        stage('Build') {
            steps {
                // Get the code from GitHub repository
                git 'https://github.com/jbleduigou/budgetcategorizer.git'

                // Code is checked out in a separate folder, create symlink to GOPATH
                sh 'mkdir -p $GOPATH/src/github/jbleduigou'
                sh 'ln -s $WORKSPACE $GOPATH/src/github/jbleduigou/budgetcategorizer'

                // Build the go project
                sh 'cd $GOPATH/src/github/jbleduigou/budgetcategorizer && go build -o budgetcategorizer ./cmd/budgetcategorizer'
            }
            post {
                success {
                    archiveArtifacts artifacts: 'budgetcategorizer', fingerprint: true
                }
            }
        }
        stage('Test') {
            steps {
                // Retrieve tool for converting output to junit format
                sh 'go get -u github.com/jstemmer/go-junit-report'
                
                // sh "sed -i -e 's|???|N/A|g' categorizer/categorizer.go"

                // Run unit tests and redirect output to go-junit-report
                sh 'cd $GOPATH/src/github/jbleduigou/budgetcategorizer && go test -v ./... 2>&1 | go-junit-report > report.xml'
            }
            post {
                always {
                  // Publish test results
                  step([$class: 'JUnitResultArchiver', testResults: 'report.xml'])
                }
            }
        }
        stage('Violations') {
            steps {
                // Retrieve golint tool
                sh 'go get -u golang.org/x/lint/golint'
                sh "sed -i -e 's|NewBroker will|this will|g' messaging/message.go"

                // Run golint
                sh 'cd $GOPATH/src/github/jbleduigou/budgetcategorizer && golint ./...'

                // Run go vet 
                sh "sed -i -e 's|objectKey, bucket|objectKey, 1337|g' config/config.go"
                sh 'cd $GOPATH/src/github/jbleduigou/budgetcategorizer && go vet ./... || true'
            }
            post {
              always {
                recordIssues enabledForFailure: true, tool: goLint(), qualityGates: [[threshold: 1, type: 'TOTAL', unstable: true]]
              }
            }
        }
    }
}
