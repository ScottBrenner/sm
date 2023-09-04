# sm-cli

A command-line interface for interacting with StepMania simfile packs - download new packs and updating existing ones right from your terminal.

## Usage
`sm-cli` creates and interacts with an additional file within a pack's directory, `source.txt`. This file merely contains the source where the pack can be downloaded from.

First, I navigate to a pack's directory and initialize `sm-cli` with the `init` argument:
```
$ cd Games/ITGmania/Songs/
$ mkdir 7guys1pack
$ cd 7guys1pack
$ sm init
Enter download source for pack: https://zenius-i-vanisher.com/v5.2/download.php?type=ddrpack&categoryid=127
```
Please be aware that the source URL here is the direct download location, and not the simfile overview page.

Once initialized I can download the pack with the `download` argument, be patient as this may take some time:
```
$ sm download
Downloading pack...
Unpacking pack...
```

`sm-cli` downloads the pack, extracts it, and deletes the downloaded archive. Already-downloaded packs can be updated with the same method.

## Setting up the development container

### GitHub Codespaces
Follow these steps to open this sample in a Codespace:
1. Click the **Code** drop-down menu.
2. Click on the **Codespaces** tab.
3. Click **Create codespace on main** .

For more info, check out the [GitHub documentation](https://docs.github.com/en/free-pro-team@latest/github/developing-online-with-codespaces/creating-a-codespace#creating-a-codespace).

### VS Code Dev Containers

If you already have VS Code and Docker installed, you can click the badge above or [here](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/microsoft/vscode-remote-try-go) to get started. Clicking these links will cause VS Code to automatically install the Dev Containers extension if needed, clone the source code into a container volume, and spin up a dev container for use.

Follow these steps to open this sample in a container using the VS Code Dev Containers extension:

1. If this is your first time using a development container, please ensure your system meets the pre-reqs (i.e. have Docker installed) in the [getting started steps](https://aka.ms/vscode-remote/containers/getting-started).

2. To use this repository, you can either open the repository in an isolated Docker volume:

    - Press <kbd>F1</kbd> and select the **Dev Containers: Try a Sample...** command.
    - Choose the "Go" sample, wait for the container to start, and try things out!
        > **Note:** Under the hood, this will use the **Dev Containers: Clone Repository in Container Volume...** command to clone the source code in a Docker volume instead of the local filesystem. [Volumes](https://docs.docker.com/storage/volumes/) are the preferred mechanism for persisting container data.

   Or open a locally cloned copy of the code:

   - Clone this repository to your local filesystem.
   - Press <kbd>F1</kbd> and select the **Dev Containers: Open Folder in Container...** command.
   - Select the cloned copy of this folder, wait for the container to start, and try things out!
   
## Things to try

Once you have this sample opened, you'll be able to work with it like you would locally.

Some things to try:

1. **Edit:**
   - Open `server.go`
   - Try adding some code and check out the language features.
   - Make a spelling mistake and notice it is detected. The [Code Spell Checker](https://marketplace.visualstudio.com/items?itemName=streetsidesoftware.code-spell-checker) extension was automatically installed because it is referenced in `.devcontainer/devcontainer.json`.
   - Also notice that utilities like `gopls` and the [Go](https://marketplace.visualstudio.com/items?itemName=golang.Go) extension are installed. Tools are installed in the `mcr.microsoft.com/devcontainers/go` image and Dev Container settings and metadata are automatically picked up from [image labels](https://containers.dev/implementors/reference/#labels).

2. **Terminal:** Press <kbd>ctrl</kbd>+<kbd>shift</kbd>+<kbd>\`</kbd> and type `uname` and other Linux commands from the terminal window.
3. **Build, Run, and Debug:**
   - Open `server.go`
   - Add a breakpoint (e.g. on line 22).
   - Press <kbd>F5</kbd> to launch the app in the container.
   - Once the breakpoint is hit, try hovering over variables, examining locals, and more.   
   - Continue (<kbd>F5</kbd>). You can connect to the server in the container by either: 
        - Clicking on `Open in Browser` in the notification telling you: `Your service running on port 9000 is available`.
        - Clicking the globe icon in the 'Ports' view. The 'Ports' view gives you an organized table of your forwarded ports, and you can get there by clicking on the "1" in the status bar, which means your app has 1 forwarded port.
   - Notice port 9000 in the 'Ports' view is labeled "Hello Remote World." In `devcontainer.json`, you can set `"portsAttributes"`, such as a label for your forwarded ports and the action to be taken when the port is autoforwarded.

   > **Note:** In Dev Containers, you can access your app at `http://localhost:9000` in a local browser. But in a browser-based Codespace, you must click the link from the notification or the `Ports` view so that the service handles port forwarding in the browser and generates the correct URL.
   
4. **Rebuild or update your container:** 

   You may want to make changes to your container, such as installing a different version of a software or forwarding a new port. You'll rebuild your container for your changes to take effect. 

   **Open browser automatically:** As an example change, let's update the `portsAttributes` in the `.devcontainer/devcontainer.json` file to open a browser when our port is automatically forwarded.
   
   - Open the `.devcontainer/devcontainer.json` file.
   - Modify the `"onAutoForward"` attribute in your `portsAttributes` from `"notify"` to `"openBrowser"`.
   - Press <kbd>F1</kbd> and select the **Dev Containers: Rebuild Container** or **Codespaces: Rebuild Container** command so the modifications are picked up.

5. **Install Node.js using a Dev Container Feature:**
   - Press <kbd>F1</kbd> and select the **Dev Containers: Configure Container Features...** or **Codespaces: Configure Container Features...** command.
   - Type "node" in the text box at the top.
   - Check the check box next to "Node.js (via nvm) and yarn" (published by devcontainers) 
   - Click OK
   - Press <kbd>F1</kbd> and select the **Dev Containers: Rebuild Container** or **Codespaces: Rebuild Container** command so the modifications are picked up.

6. **Refactoring - rename:**
    - Open `hello.go`, select method name `Hello` press <kbd>F1</kbd> and run the **Rename Symbol** command.
7. **Refactoring - extract:**
   - Open `hello.go` and select string, press <kbd>F1</kbd> and run the **Go: Extract to variable** command.
   - Open `hello.go` and select line with return statement, press <kbd>F1</kbd> and run the **Go: Extract to function** command.
8. **Generate tests:**
    - Open `hello.go` and press <kbd>F1</kbd> and run the **Go: Generate Unit Tests For File** command.
    - Implement a test case: Open file `hello_test.go` and edit the line with the `TODO` comment: `{"hello without name", "Hello, "},` 
    - You can toggle between implementation file and test file with press <kbd>F1</kbd> and run the **Go: Toggle Test File**
    - Tests can also run as benchmarks: Open file `hello_test.go`, press <kbd>F1</kbd> and run the **Go: Benchmark File**
9. **Stub generation:** ( [details](https://github.com/josharian/impl))
   - define a struct `type mock struct {}`, enter a new line , press <kbd>F1</kbd> and run the **Go: Generate interface stubs** command.
   - edit command `m *mock http.ResponseWriter`
10. **Fill structs:** ([details](https://github.com/davidrjenni/reftools/tree/master/cmd/fillstruct))
   - Open `hello.go` and select `User{}` of variable asignment, press <kbd>F1</kbd> and run the **Go: Fill struct** command.
11. **Add json tags to structs:** ([details](https://github.com/fatih/gomodifytags))
   - Open `hello.go` and go with cursor in to a struct, press <kbd>F1</kbd> and run the **Go: Add Tags To Struct Fields** command.


## Contributing

This project welcomes contributions and suggestions.

## License

Copyright © Scott Brenner All rights reserved.<br />
Licensed under the MIT License. See LICENSE in the project root for license information.
