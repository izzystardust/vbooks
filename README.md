# vbooks

This is the backend server for the vbooks project - a place for small 
communities to share PDFs, ePubs, and other reading resources with each other.

The repository frontend server can be found at 
[~izzy/vbooks-ui](https://git.sr.ht/~izzy/vbooks-ui).

Planning, feature requests, and bugs live at 
[our issue tracker](https://todo.sr.ht/~izzy/vbooks).


## Backend

### Building

Dependencies:
* go (>=1.12)

    cd $REPO/cmd/vbooks-server && go build

Because the vbooks project is using go modules, the repository should be checked
out in a location outside of GOPATH or compiled with `GO111MODULE=on`.

#### Live Reloading for Development

Live reload of the server while editing code can be accomplished with
[codegangsta/gin](https://github.com/codegangsta/gin) by running

    gin --build cmd/vbooks-server --excludeDir client

in the root of the repository.

## Frontend

The frontend code lives in `client/`

Dependencies:
* npm

Run `npm install` to get all the dependencies.


To run the development frontend, run `npm run dev` from the `client` directory.


## Resources

[Send patches](https://git-send-email.io/) and questions to
[~izzy/vbooks-dev@lists.sr.ht](https://lists.sr.ht/~izzy/vbooks-dev). When 
sending patches, lease use `--subject-prefix PATCH backend` for clarity.
