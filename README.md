
## How to generate docs

For docs generation I use helm-docs https://github.com/norwoodj/helm-docs.

1. Populate Chart.yaml fields https://helm.sh/docs/topics/charts/
2. You need to download 'helm-docs' executable from this page https://github.com/norwoodj/helm-docs/releases.
3. Create README.md.gotmpl and put it inside Helm chart (ignore this file with .helmignore)
4. Use examples found here https://github.com/norwoodj/helm-docs/tree/master/example-charts to make your template
5. Run ``helm-docs -c <helm chart folder path here>`` for example ``helm-docs -c query-exporter/`` to generate docs
6. Configure pre-commit hook to generate docs on commit https://github.com/norwoodj/helm-docs#pre-commit-hook

