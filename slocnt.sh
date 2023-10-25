#cloc --by-file --exclude-ext=json,csv,yaml,md  --exclude-dir=test .
cloc --by-file --not-match-f '_test\.go$' --exclude-ext=json,csv,yaml,md  --exclude-dir=test .
