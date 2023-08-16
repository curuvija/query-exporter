pipeline {
    agent { label 'k8s-helm-pod' }

    options{
        buildDiscarder(logRotator(numToKeepStr: '5', daysToKeepStr: '5'))
    }

    environment {
        GITHUB_TOKEN = credentials('github_curuvija_jcasc')
    }


    stages {
        // TODO: kube-linter -> https://github.com/stackrox/kube-linter (kube-score, kubeconform, kubeeval, datree, kics -> https://kics.io/index.html
        // TODO: check also https://analysis-tools.dev/tag/kubernetes
        // TODO: kube-hunter -> https://aquasecurity.github.io/kube-hunter/
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
        // TODO: Terratest (you need to install prometheus crds)
        stage('Test') {
            steps {
                sh 'echo something'
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
                    sh 'cr upload -o curuvija --git-repo query-exporter --package-path .cr-release-packages/ --token ${GITHUB_TOKEN_PSW}'
                }
            }
        }
        stage('Publish') {
            when {
                allOf {
                    expression {
                        env.BRANCH_NAME == 'master'
                    }
                }
            }
            steps {
                sh 'echo something'
            }
        }
        stage('Tag') {
            when {
                allOf {
                    expression {
                        env.BRANCH_NAME == 'master'
                    }
                }
            }
            steps {
                sh 'echo something'
            }
        }
    }
}
