<!--
  Copyright 2026 Samvel Khalatyan. All rights reserved.

  Use of this source code is governed by a BSD-style
  license that can be found in the LICENSE file.
-->

# NAME

**hello** - run Go in Browser using WebAssembly

# DESCRIPTION

Hello demonstrates running Go code in a browser using WebAssembly
([instructions](https://go.dev/wiki/WebAssembly#getting-started)).

1. write a Go program to do work in a `main()` function.
2. build it into WebAssembly (Wasm):
   ```console
   env GOOS=js GOARCH=wasm go build -o hello.wasm .
   ```
3. Create an HTML page with following header:
   ```html
   <html>
     <header>
       <script src="wasm_exec.js"></script>
       <script>
         const go = new Go();
         WebAssembly.instantiateStreaming(
           fetch("hello.wasm"),
           go.importObject,
         ).then((result) => {
           go.run(result.instance);
         });
       </script>
     </header>
   </html>
   ```
4. open the page in a browser and check the console output.

It is important to use `wasm_exec.js` that comes with Go installation:

```console
% ls `go env GOROOT`/lib/wasm
go_js_wasm_exec     go_wasip1_wasm_exec wasm_exec.js        wasm_exec_node.js
```

Alternatively, use make(1):

```console
% make -C ./x/wasm/hello
% go run ./x/serve ./x/wasm/hello
```
