## diffloc

Tiny CLI tool to get the number of lines changed in a GitHub PR manner (`+242 -40`). Useful for keeping the PRs small.

### Usage

Download the binary from the "Releases" page.

```sh
diffloc # difference in the working tree

diffloc HEAD~1 # difference from the last commit

diffloc development # difference from the 'development' branch
```