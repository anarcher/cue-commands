package hello

import "tool/cli"

command: print: {
    print: cli.Print & {
        text: hello 
    }
}
