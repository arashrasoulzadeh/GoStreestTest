## GO API Stress Test
Use this appication to stress test your api

Writeen by [Arash Rasoulzadeh]("http://twitter.com/x3n0b1a") , special thanks to [amiraliamhh](https://github.com/amiraliamhh) 

## Example usage

#### Run with a single thread:
`go run main.go single`

#### Run using multiple threads:
*runs using 4 threads*
`go run main.go multiple 4`
#### Scenario 
`go run composer.go`

This command runs the composer , composer reads the configuration from scenario.conf , you may see the example configuration file [here](https://github.com/arashrasoulzadeh/GoStreestTest/blob/master/scenario.yaml)

schema for configation file is :

```yaml
scenarios:
      --
        label: "light"
        config: "config.yaml"
        sleep: 1
      --
        label: "medium"
        config: "config4.yaml"
        sleep : 1
```

scenarios key is an array of 
* label - its just a label
* config - is the config file name
* sleep - sleep in seconds after running this part

---


#### changelog
- readme fixed
- colorized labels in scenario
- scenario added
- fixed bugs in multi-threading

