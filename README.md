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
2. Type `speep + $COMMAND` for your purpose.
    - There are 3 types of commands (`get`, `add` and `delete`)
    - There will be a guide for flags by just typing commands
    ```
    // To retrieve all the shortcuts registered
    speep get -all
     
    // To add a new shortcut (both name and shortcut in string)
    speep add -name Copy -shortcut Ctrl+C
    
    // To know what flags are available for the command
    speep get
    // Result
     -all
        Get full shortcut list
     -keyword string
        Find a shortcut with keyword 
        specify the target shortcut with all or keyword flag
    ```

## Inquiry
Raise an issue on github for any inquiries.

## Contributor
Special thanks to @kinbiko!
