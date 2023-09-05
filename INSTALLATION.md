# Installation

## Pre-built Binaries (Cross Platform)

`sm` is available as a pre-built binary from the [releases page](https://github.com/ScottBrenner/sm/releases) for the following:

    macOS (Darwin)
    Windows
    Linux

## Installing Binary on macOS / Linux

Once downloaded, there are two common ways to install. You can place the binary in either of the following directories:

    /usr/local/bin
    $HOME/bin

Using `/usr/local/bin`:

    # Switch directory
    cd /usr/local/bin
    # download the binary
    wget -O sm <URL to binary download>
    # set permissions
    chmod +x sm
    # verify installation
    sm help

Using `$HOME/bin`:

    # create the directory if needed
    mkdir -p $HOME/bin
    # make it the working directory
    cd $HOME/bin
    # download the binary
    wget -O sm <URL to binary download>
    # set permissions
    chmod +x sm
    # verify installation
    ./sm help

To run `sm` without the `./`, you will need to add `$HOME/bin` to your `$PATH`. Here is how you do that:

1. Determine your default shell (zsh or bash).

`echo $SHELL`

2. Edit your shell profile.

If your default shell is zsh:

`nano ~/.zshrc`

If your default shell is bash:

`nano ~/.bash_profile`

In `~/.zshrc` or `~/.bash_profile` add the following line to the file:

`export PATH=$PATH:$HOME/bin`

Save the file by pressing Control-X, then Y.

Close the terminal and open a new terminal to pick up the changes to your profile. Verify the change by running the `sm help` command.