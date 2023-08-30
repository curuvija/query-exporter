pipeline {
    agent { label 'k8s-helm-pod' }

    options{
        buildDiscarder(logRotator(numToKeepStr: '5', daysToKeepStr: '5'))
    }

    environment {
        GITHUB_TOKEN = credentials('github_curuvija_jcasc')
    }

    stages {
        stage('Lint') {
            steps {
                container('helm') {
                    sh 'helm lint query-exporter/'
                }
            }
        }
        // TODO: polaris https://github.com/FairwindsOps/polaris
        stage('Scan for security issues') {
            steps {
                sh 'echo something'
            }
        }
        stage('Test') {
            steps {
                container('helm') {
                    sh 'cd tests && go test ./...'
                }
            }
        }
        stage('Template') {
            steps {
                container('helm') {
                    sh 'helm template query-exporter/'
                }
            }
        }
        stage('Dry run') {
            steps {
                sh 'echo something'
            }
        }
        stage('Static code analysis') {
            steps {
                sh 'echo something'
            }
        }
        stage('Generate docs') {
            steps {
                container('helm') {
                    sh 'helm-docs query-exporter/ && cp query-exporter/README.md .'
                }
            }
        }
        // TODO: git-chglog - https://github.com/git-chglog/git-chglog
        stage('Changelog') {
            steps {
                sh 'echo something'
            }
        }
        stage('Package Helm chart') {
            steps {
                container('helm') {
                    sh 'cr package query-exporter/'
                }
            }
        }
        stage('Upload Helm chart to releases') {
            when {
                allOf {
                    expression {
                        env.BRANCH_NAME == 'master'
                    }
                }
            }
            steps {
                container('helm') {
                    sh 'cr upload -o curuvija --git-repo query-exporter --package-path .cr-release-packages/ --token ${GITHUB_TOKEN_PSW} --release-notes-file CHANGELOG.md'
                    sh 'cr upload -o curuvija --git-repo helm-charts --package-path .cr-release-packages/ --token ${GITHUB_TOKEN_PSW} --release-notes-file CHANGELOG.md'
                }
            }
        }
        // TODO: check https://itnext.io/jenkins-tutorial-part-10-work-with-git-in-pipeline-b5e42f6d124b
        stage('Checkout helm charts repo') {
            when {
                allOf {
                    expression {
                        env.BRANCH_NAME == 'master'
                    }
                }
            }
            steps {
                dir("helm-charts-repo") {
                    git(
                        url: "https://github.com/curuvija/helm-charts.git",
                        branch: "gh-pages",
                        changelog: true,
                        poll: true
                    )
                }
            }
        }
        stage('Generate helm-charts index') {
            when {
                allOf {
                    expression {
                        env.BRANCH_NAME == 'master'
                    }
                }
            }
            steps {
                sh 'cr index --index-path index.yaml --package-path .cr-release-packages/ --owner curuvija --git-repo helm-charts --pr'
                //sh 'cr index --index-path index.yaml --package-path .cr-release-packages/ --owner curuvija --git-repo helm-charts --push'
            }
        }
        stage('Publish index') {
            when {
                allOf {
                    expression {
                        env.BRANCH_NAME == 'master'
                    }
                }
            }
            steps {
                withCredentials([gitUsernamePassword(credentialsId: 'github_curuvija_jcasc', gitToolName: 'Default')]) {
                    //sh "git push -u origin gh-pages"
                }
            }
        }
    }
}
