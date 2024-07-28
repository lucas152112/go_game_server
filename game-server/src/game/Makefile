game.linux:
	env GOOS=linux GOARCH=amd64 go build -o game main.go

game.copy:
	scp game bobi_dev:/mnt/dezhou/script

game.start:
	ssh bobi_dev "cd /mnt/dezhou/script/ ; ./service.sh start"

game.stop:
	ssh bobi_dev "cd /mnt/dezhou/script/ ; ./service.sh stop"

game.product_copy:
	scp game aalive:/home/ec2-user/

game.test_copy:
	scp game paytest:/home/ec2-user/script

git.branch_all:
	git branch -vv

git.branch_main:
	git checkout main

git.branch_dev:
	git checkout dev-auth

git.pull_dev:
	git pull origin dev-auth