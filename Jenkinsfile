pipeline {
	agent any

	tools { 
		go 'go-1.21' 
	}

	environment {
		imageName = "pulsr-gitlab-service"
		registryCredentials = "adminnexus"
        registry = "docker.algueron.io"
        dockerImage = ''
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
				sh 'go test -v --short ./...'
			}
		}
		stage("Integration Tests") {
			steps {
				echo 'INTEGRATION TESTS'
				sh 'go test -v -coverprofile=cover.out -covermode count ./...'
			}
		}
		stage("Code Analysis") {
			steps {
				echo 'CODE ANALYSIS'
				sh 'go install github.com/t-yuki/gocover-cobertura@latest'
				sh '$HOME/go/bin/gocover-cobertura < cover.out > coverage.xml'
				publishCoverage adapters: [cobertura(coberturaReportFile: 'coverage.xml')], tag: 't'
			}
		}
		stage("Docker build") {
			steps {
				echo 'DOCKER BUILD'
				script {
          			dockerImage = docker.build imageName
        		}
			}
		}
		stage("Docker Publish") {
			steps {
				echo 'DOCKER PUBLISH'
				docker.withRegistry('http://'+registry, registryCredentials) {
                	dockerImage.push('latest')
          		}
			}
		}
		stage("Deploy") {
			steps {
				echo 'DEPLOY'
			}
		}
	}
}
