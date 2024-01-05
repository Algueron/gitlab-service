pipeline {
	agent any

	tools { 
		go 'go-1.21' 
	}

	stages {
		stage("Compile") {
			steps{
				echo 'COMPILE'
				sh 'go build ./...'
			}
		}
		stage("Unit Tests") {
			steps {
				echo 'UNIT TESTS'
				sh 'go test -v ./...'
			}
		}
		stage("Integration Tests") {
			steps {
				echo 'INTEGRATION TESTS'
				sh 'go test -v -tags=integration ./...'
			}
		}
		stage("Code Analysis") {
			steps {
				echo 'CODE ANALYSIS'
				sh 'go test -coverprofile=cover.out -covermode count ./...'
				sh 'go install github.com/t-yuki/gocover-cobertura@latest'
				sh '$HOME/go/bin/gocover-cobertura < cover.out > coverage.xml'
				publishCoverage adapters: [cobertura(coberturaReportFile: 'coverage.xml')], tag: 't'
			}
		}
		stage("Docker build") {
			steps {
				echo 'DOCKER BUILD'
			}
		}
		stage("Docker Publish") {
			steps {
				echo 'DOCKER PUBLISH'
			}
		}
		stage("Deploy") {
			steps {
				echo 'DEPLOY'
			}
		}
	}
}
