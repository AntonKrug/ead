# Usage

Typical usage excpects a specific folder structure as following:

`ead -standalone_files=fale -source=./src/application/embedded_autogenerated_data/data_root_raw`

```
<PROJECT_ROOT>
    src
        application
            embedded_autogenerated_data       <- output
                data_root_include             <- output_container
                data_root_raw                 <- source

                ead_collection.h              <-
                ead_helpers.c                 <-  Auxiliary files
                ead_helpers.h                 <-
                ead_structures.h              <-
```

It is assumed that the parent of the **output** directory is in project's includes and the **output** is exactly named *embedded_autogenerated_data* by default. If this is not the case then **include_prefix** argument needs to change the default so the generated files have correct paths in the includes.

A current working directory can be used as the **source** folder, but because this could lead to accidental processing of huge amount of files (if executed in a wrong location), therefore it is guarded by the **source_current_folder** flag, use with caution. Add the **-source_current_folder=true** argument when confident that the source directory can be safely the current working directory.

By default the parent of the **source** is used as **output**, the generated files will be saved into **output**/**output_container**. In this example it is the src/application/embedded_autogenerated_data/data_root_include. While the auxiliary files are stored into the **output** path directly.

If the location of the output files needs to be different then change the **-output** and **-output_container** arguments. Then the project's include path or **-include_prefix** might be adjusted as well.

If different auxiliary files are provided or for some other reason they were modified and can't be changed use **-output_auxiliary=false** to disable the auxiliary file generation. Note that the generated files still will produce a metadata structure and will still try to include the structure depend on it to be defined.

To keep the default settings less intrusive the hex-dumps without metadata are produced. Default is **-standalone_files=true** which will produce files which do not depend on any other includes and are fully standalone. But they do not provide any metadata, thus do not provide any virtual filesystem features. When left as true, then the value of **output_auxiliary** is irelevant as no auxiliary files will be generated. The **include_prefix** is irelevant as well because in standalone mode no includes are used in the files.

When using the generated files in a web-server context, then some files can be pre-compressed, use **-compress_web=true** to enable this feature. Use -h to see the list of supported extensions. At the moment only GZIP compression is supported, possibly in the future other compressions might be added.

If the bundled Microchip's copyright is not suitable, then use **-copyright** argument and point to a copyright notice file. The file content will be used as a comment in a C/H files as it is and therefore it needs to have valid syntax. The first line is already indented correctly, but for multiline copyright notices extra attention needs to be made to make sure the syntax and indentation are not broken. 

For all other arguments use **-h** argument to display help.

# Build EAD

- Install and configure golang https://golang.org/dl/
- Install the following packages by typing the following: 

  `go get github.com/logrusorgru/aurora`

  `go get github.com/gabriel-vasile/mimetype`

  `go get github.com/hoisie/mustache`

  `go get github.com/shurcooL/vfsgen`

  `go get github.com/dustin/go-humanize`

- Download this project with `go get github.com/antonkrug/ead` ("undefined: Assets" error is expected as the `go generate` was not run yet)
- Go into the project:
  - On Windows: `cd %GOPATH%/src/github.com/antonkrug/ead`
  - On Linux: `cd $GOPATH/src/github.com/antonkrug/ead`

- Generate vfsdata `go generate`
- Build the project
  - Build the final native binary `go build`
  - To build all other platforms run `bash ./build_all_platforms.sh` (Tested on Linux and on Windows under GitBash command line)

Note: If it is required to cross-build from Windows to Linux only and then copy it into a shared VM folder then the following can be used:
```
GOOS=linux GOARCH=amd64 go build -o release/ead-linux-x86-64
cp release/ead-linux-x86-64 /d/VMs/_shared_folder/ead
```


