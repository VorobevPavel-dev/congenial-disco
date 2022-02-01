pipeline {
    agent {
        node {
            label 'centos'
        }
    }
    stages {
        stage("Testing"){
            agent {
                docker {
                    image 'golang:1.17.6-buster'
                    reuseNode true
                }
            }
            steps {
                echo "Start testing for commit"
                sh "find .  -regex \".*_test\\.go\" -type f -printf '%h\\n' | xargs go test"
            }
            when {
                not {
                    branch 'master'
                }
            }
        }
    }
}