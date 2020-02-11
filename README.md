# goFIND

Command Line Search Utility to search for words in all files inside the specified directory and all it's sub directories.


### Usage:

```
git clone https://github.com/KrishKayc/goFIND.git

go build 

gofind -dir="your search directory path" -search="your search string"

```

### Configurable:

Utility is customizable with a "config.json" file which determines which file extensions to search for and search type.

Sample config.json

```
{
    "excludeDirectories" : ["node_modules","bin","sn-item-bank-picker", "scanit-latest","sn-tinymce"],
    "excludeFiles":[],
    "allowedExtensions" : [".go", ".json", ".txt", ".cs", ".cpp", ".c", ".xml", ".js", ".ts", ".csproj",".html"],
    "matchCase":false,
    "matchFullWord":false
}

```

* excludeDirectories :> Excludes the specified directories from searching
* excludeFiles       :> Excludes the specified files from searching
* allowedExtensions  :> Only searches within the files specified in this
* matchCase          :> Determines whether to do a "case sensitive search"
* matchFullWord      :> Determines whethet to match the "full word" or "partial text"


### Output :

![Sample Output](https://github.com/KrishKayc/goFIND/blob/master/sample_output/gofind_sample_output.jpg)




