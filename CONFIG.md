# Configuration file documentation

The configuration file is designed to be as customizable as possible, so there may be some complexity coming with it. However, I'll try my best to explain how to set things up accordingly.

First of all, StarBurst is divided into various sub-systems, which are currently:
- Generic functionality (`general`)
- D-Bus Interop (`dbus`)
- OBS Interop (`obs`)
- Twitch API (`twitch`)
- VSeeFace actions (`vseeface`)

Each sub-system has it's own configuration key, which matches it's internal name. Each configuration key contains various settings regarding the sub-systems functionality (like passwords, keys etc.), and a `buttons` key containing button definitions.

### Button definitions
Each button that is visible on the stream deck is represented by a button configuration, which takes the form of the following YAML:

```yaml
action: call_me
params:
    param_1: "so here's my number"
    another_name: 420
button_image: "path/to/awesome/image.png"
button_index: 0
```

The `action` key refers to a hardcoded function inside the particular addon, which are all listed down below. `params` refers to the particular parameters the function takes; those vary and are described in more detail below. `button_image` is just a path to an image that represents the button, and `button_index` declares the position of the button on the stream deck, beginning with 0 starting at the top-left, going down row-majorly.

### Addon descriptions

#### Generic functionality (`general`)

Available actions:

- `set_brightness`: Changes the brightness of the stream deck screen. Brightness is given as an integer between 0 and 100.
    - Parameter `absolute` (`bool`): Defines if the brightness should be set to the value directly (`true`) or be changed relative to the current brightness (`false`)
    - Parameter `value` (`int`): The value to change the brightness to or by.
- `execute`: Executes a given program.
    - Parameter `program` (`string`): The program name or path. If only the name is given, it tries to look up the path itself.
    - Parameter `cmdline` (`string`): The command line given to the program.


#### D-Bus Interop (`dbus`)

Available actions:

- `call`: Calls a D-Bus method on the system bus.
    - Parameter `destination` (`string`): The destination of the D-Bus call.
    - Parameter `path` (`string`): The object path of the D-Bus call.
    - Parameter `method` (`string`): The name of the method to call.
    - Parameter `params` (`array`): A list of any values that are sent as parameters to the called method.


#### OBS Interop (`obs`)

Settings:

- `host` (`string`): The address of the computer where an OBS websocket is listening. By default `localhost`.
- `port` (`int`): The port that an OBS websocket is listening. By default `4455`.
- `password` (`string`): The password used to connect to the OBS websocket.

Available actions:

- `set_scene`: Changes the program scene of OBS to the given scene. Additionally, declaring a button with this action causes it to be marked with a red border whenever the given scene is active.
    - Parameter `scene_name` (`string`): The name of the scene to switch to.

#### Twitch API (`twitch`)

Settings:

- `client_id` (`string`): The client id of the application that is used to send API requests.
- `client_secret` (`string`): The accompanying secret.

Available actions:

- `set_marker`: Places a stream marker at the current time while streaming.
    - Parameter `user_id` (`int`): The user id of the broadcaster to set the marker for.

#### VSeeFace actions (`vseeface`)

Available actions:

- `set_expression`: Sends a key combination via `xdotool` over to the VSeeFace window to set a specific expression. The key combination sent over is `CTRL+SHIFT+<key>`. Additionally, declaring a button with this action will cause it to be highlighted when selected.
    - Parameter `key` (`string`): The key used in the combination that is sent over.