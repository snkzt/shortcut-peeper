# shortcut-peeper
A cheat sheet app for shortcut keys!  
It's so confusing to remember shortcuts, especially if you are using several tools that each of them has different shortcut keys to do the same thing.  
If you can relate to this, then this app is for you.  
Nothing complicated, super simple app that support your computer life. 

## How to use
#### Install
1. Install go with `brew install go`.
2. Make sure your local `$PATH` contains `~/go/bin`.
3. Run `go install github.com/snkzt/shortcut-peeper/cmd/speep@latest`.
 

#### Usage
1. Open your terminal.
2. Type `speep <command> <flag1> ... <flag2> ...` for your purpose.
    - There are 3 types of commands (`get`, `add` and `delete`)
    - There will be a guide. Type ```speep help```
    ```
   $ speep help
    Usage of speep: speep <command> <flag1> ... <flag2> ...

    Command Options:
		get [--all] | [--name <name>]
		add --category <category> --name <name> --key <key>
		delete [--all] | [--category <category> --name <name>]

	Flag Options:
		-a, --all retrieve all shortcuts
		-c, --category name of the category of the registered shortcut key: e.g. shell
		-n, --name name of the registered shortcut key: Use "" for more than one word e.g. -name "to the back of the line"
		-k, --key registered shortcut key
    ```

## Inquiry
Raise an issue on github for any inquiries.

## Contributor
Special thanks to @kinbiko for the contribution!
