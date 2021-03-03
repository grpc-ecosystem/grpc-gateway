   * Provide couple annotations for replacing values in swagger schema you need specify flag ```atlas_patch=true``` and ```with_custom_annotations=true```
     - ```@example``` annotation can be used for replacing default example with custom one
      support few value types includes float64, string, map[string]interface{}, []map[string]interface{} []float64, []string
       - ```@example 5.0```
       - ```@example "Internal error"```
       - ```@example {"Location": "Tacoma"}```
       - ```@example ["First", "Second"]```
       - ```@example [1, 5, 44]```
       - ```@example [{"Location": "Tacoma"}, {"Group": "Engineering"}]```

     - ```@title``` annotation can be used for replacing default title with custom one
       - ```@title "StringCondition"```

    If you example too long to be presented on one line you could use multiple lines annotation
    ```
       @example <<<EOF
       {
           "Location": "Tacoma"
       }
    ```

    or

   ```
       @example <<<EOF
       {
           "Location": "Tacoma"
       }
       EOF
    ```

    In first case all what presented after line ```@example <<<EOF``` will be rendered as example,
    if you want to manually set boundaries please use ```EOF``` as a closing line
