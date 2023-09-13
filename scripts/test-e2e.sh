export GITHUB_ACTIONS=false
export WALLET_PASSWORD=abcde
task generate
task build
task run &
cd api/test
robot robot_tests
