# golang-bash-completion
A bash-completion shell script for golang package

## To install the files
For most Linux users, run the following commands in the repository:

    ./gen.sh
    make install

The completion files will be installed system-wide, as you may be requested
for administrative privileges. `sudo` command is required in some Linux
distributions.

For installing to current user only, just type the following commands:

    echo ". $(pwd)/go" >> ~/.bash_completion
    echo ". $(pwd)/gofmt" >> ~/.bash_completion

Remember to run a new `bash` to make them available.

The installation have not been run in macOS
(Darwin), any issues or suggestions could be posted
[here](https://github.com/GreenYun/golang-bash-completion/issues/new 'issues').

## License

This project is licensed under the [MIT License](LICENSE). You are encouraged
to embed DNS-over-HTTPS into your other projects, as long as the license
permits.

You are also encouraged to disclose your improvements to the public, so
that others may benefit from your modification, in the same way you receive
benefits from this project.
