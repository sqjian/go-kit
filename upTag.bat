set tag=v0.0.16

git tag -d %tag%
git push origin :refs/tags/%tag%

git tag  %tag%
git push origin %tag%