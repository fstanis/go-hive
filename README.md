# go-hive
Unofficial Go wrapper library compatible with the API used by [Hive](https://www.hivehome.com/)
(by British Gas) smart devices.

## Disclaimer
**This software is in no way endorsed by British Gas**. The underlying REST API
is undocumented and subject to change at any time, which means this code library
may suddenly stop working.

**Use at your own risk.**

In addition, the library only supports a subset of the API - specifically, the
devices I personally own. Most notably, smart heating is not supported.

Contributions are welcome.

## Basic usage

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

## Sending commands

Commands such as turning a light on or off are called *changes* in the library.

To create a new change, use `hive.NewChange()` and then call on it the method
that specifies what kind of change you want to do.

A few examples:

```
  // All methods of the Change object mutate it.
  //
  // This will turn on the light when sent to it.
  ch := hive.NewChange()
  ch.TurnOn()

  // All methods return the change object itself, for the purpose of chaining.
  // The above can thus be simplified like this:
  ch := hive.NewChange().TurnOn()

  // A single change object can contain multiple changes. This will turn on the
  // light and set the brightness to 50 at the same time.
  ch := hive.NewChange().TurnOn().Brightness(50)

  // Some methods may contract / cancel the ones called before.
  ch := hive.NewChange().TurnOn().TurnOff() // same as hive.NewChange().TurnOff()
```

Once the `Change` object is constructed, it can be sent to the device using the
`Do` method of it.

```
  device := client.Device("some-light-device-id")
  ch := hive.NewChange().TurnOn()
  if err := device.Do(ch); err != nil {
    fmt.Printf("Failed to turn on light: %v", err)
  }
```

## Full example

For a full example, check out [hivecli](https://github.com/fstanis/hivecli), a
complete CLI solution for accessing your Hive devices.

## GoDoc

Please see the [GoDoc documentation](https://godoc.org/github.com/fstanis/go-hive/hive)
for more information on the API.
