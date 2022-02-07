export const DOCKER_ONELINER = `docker run --interactive --tty \\
    --pull always \\
    --publish 30080:80 --publish 30022:22 \\
    --volume "$HOME/.sturdydata:/var/data" \\
    getsturdy/server`
