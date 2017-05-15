# hive
Unofficial Go wrapper library for the API used by the Hive Hive (by British Gas)
smart devices.

## Usage

```
import "github.com/fstanis/go-hive/hive"
```

Construct a new client and login (it's recommended to keep the credentials in a
variable, as you can reuse them when the token expires):
```
  client := hive.NewClient()
  c := &hive.Credentials{
    Username: "person@example.com",
    Password: "qwerty",
    URL:      "https://beekeeper.hivehome.com/1.0/global/login",
  }
  if err := client.Login(c); err != nil {
    log.Fatalf("Failed to login: %v", err)
  }
```

Example use:

```
  fmt.Println("Listing smart lights in your home...\n")
  for _, d := range client.Devices() {
    if d.IsLight() {
      fmt.Println(d.ID(), d.Color(), d.ColorTemperature(), d.Brightness())

      // If the light bulb is named "Living room", turn it on
      if d.Name() == "Living room" {
        d.Do(hive.NewChange().TurnOn())
      }
    }
  }

```

[GoDoc documentation](https://godoc.org/github.com/fstanis/go-hive/hive)
