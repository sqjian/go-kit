set tag=0.0.1

git tag -d %tag%
git push origin :refs/tags/%tag%

git tag  %tag%
git push origin %tag%