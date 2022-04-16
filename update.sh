echo "Commit name"
# shellcheck disable=SC2162
read message
echo "Sending update"
git add .
git commit -am "$message"
git push
echo "Sent update"