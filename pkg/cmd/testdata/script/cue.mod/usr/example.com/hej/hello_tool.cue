package hej

import "tool/cli"

command: hello: {
    print: cli.Print & {
        text: world 
    }
}
