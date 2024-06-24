# uds-generator

## Test

Generate (OG way/Deprecated):

```bash
# This fails on uds package mutation and lookups
 ./uds-generator generate -c oci://ghcr.io/stefanprodan/charts/podinfo -v 6.6.2 -n podinfo
```



 8727  zarf dev generate podinfo --url https://github.com/stefanprodan/podinfo.git --version 6.4.0 --gitPath charts/podinfo --output-directory upstream
 8728  zarf dev generate podinfo --url https://repo1.dso.mil/big-bang/apps/sandbox/podinfo --version 6.4.0 --gitPath charts --output-directory registry1
 8729  zarf dev generate podinfo --url https://repo1.dso.mil/big-bang/apps/sandbox/podinfo --version 6.4.0 --gitPath chart --output-directory registry1
 8730  zarf dev generate podinfo --url https://repo1.dso.mil/big-bang/apps/sandbox/podinfo.git --version 6.4.0 --gitPath chart --output-directory registry1
 8732  zarf dev generate podinfo --url https://repo1.dso.mil/big-bang/apps/sandbox/podinfo.git --version 6.0.0 --gitPath chart --output-directory registry1
 8733  zarf dev generate podinfo --url https://repo1.dso.mil/big-bang/apps/sandbox/podinfo.git --version 6.0.0-bb.5 --gitPath chart --output-directory registry1
 9139  zarf dev generate 
 9140  zarf dev generate coder --url https://github.com/coder/coder.git --gitPath charts/coder --version v2.10.2
 9141  zarf dev generate coder --url https://github.com/coder/coder.git --gitPath charts/coder --version v2.10.2 --output-dir .
 9142  zarf dev generate coder --url https://github.com/coder/coder.git --gitPath charts/coder --version v2.10.2 --output-directory .
 9143  zarf dev generate coder --url https://github.com/coder/coder.git --gitPath helm/coder --version v2.10.2 --output-directory .
 9443  zarf dev generate
 9445  zarf dev generate gitness --url https://github.com/harness/gitness.git --version v3.0.0-beta.7 --path charts/gitness
 9446  zarf dev generate gitness --url https://github.com/harness/gitness.git --version v3.0.0-beta.7 --gitPath charts/gitness
 9447  zarf dev generate gitness --url https://github.com/harness/gitness.git --version v3.0.0-beta.7 --gitPath charts/gitness --output-directory .
 9506  zarf dev generate athens --url https://github.com/gomods/athens-charts --gitPath charts/athens-proxy --version athens-proxy-0.9.5 --output-dir ./athens/
 9507  zarf dev generate athens --url https://github.com/gomods/athens-charts --gitPath charts/athens-proxy --version athens-proxy-0.9.5 --output-directory ./athens/
 9508  zarf dev generate athens --url https://github.com/gomods/athens-charts.git --gitPath charts/athens-proxy --version athens-proxy-0.9.5 --output-directory ./athens/
 9836  zarf dev generate --url https://charts.verdaccio.org
 9837  zarf dev generate --url https://charts.verdaccio.org --version 4.16.1
 9838  zarf dev generate verdaccio --url https://charts.verdaccio.org --version 4.16.1
 9839  zarf dev generate verdaccio --url https://charts.verdaccio.org --version 4.16.1 --output-directory common
10667  zarf dev generate --help
10668  zarf dev generate --url https://stefanprodan.github.io/podinfo --version 6.6.2
10669  zarf dev generate podinfo --url https://stefanprodan.github.io/podinfo --version 6.6.2
10670  zarf dev generate podinfo --url https://stefanprodan.github.io/podinfo --version 6.6.2 --output-directory tmp/generate
10671  zarf dev generate podinfo --url oci://ghcr.io/stefanprodan/charts/podinfo --version 6.6.2 --output-directory tmp/generate1
10672  zarf dev generate podinfo --url https://github.com/stefanprodan/podinfo.git --version 6.6.2 --chartPath charts/podinfo --output-directory tmp/generate2
10673  zarf dev generate podinfo --url https://github.com/stefanprodan/podinfo.git --version 6.6.2 --gitPath charts/podinfo --output-directory tmp/generate2
```