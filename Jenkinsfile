pipeline {
    agent { label 'k8s-ansible-pod' }

    options{
        buildDiscarder(logRotator(numToKeepStr: '5', daysToKeepStr: '5'))
    }

    stages {
        // TODO: kube-linter -> https://github.com/stackrox/kube-linter (kube-score, kubeconform, kubeeval, datree, kics -> https://kics.io/index.html
        // TODO: check also https://analysis-tools.dev/tag/kubernetes
        // TODO: kube-hunter -> https://aquasecurity.github.io/kube-hunter/
        stage('Lint') {
            steps {
                sh 'echo something'
            }
        }
        // TODO: polaris https://github.com/FairwindsOps/polaris
        stage('Scan for security issues') {
            steps {
                sh 'echo something'
            }
        }
        // TODO: Terratest
        stage('Test') {
            steps {
                sh 'echo something'
            }
        }
        stage('Template') {
            steps {
                sh 'echo something'
            }
        }
        stage('Dry run') {
            steps {
                sh 'echo something'
            }
        }
        stage('Package') {
            steps {
                sh 'echo something'
            }
        }
        // TODO: helm-docs
        stage('Docs') {
            steps {
                sh 'echo something'
            }
        }
        // TODO: git-chglog - https://github.com/git-chglog/git-chglog
        stage('Changelog') {
            steps {
                sh 'echo something'
            }
        }
        stage('Publish') {
            steps {
                sh 'echo something'
            }
        }
        stage('Tag') {
            steps {
                sh 'echo something'
            }
        }
    }
}
