# Description
Small utility written in Go that can find and filter pipewire objects based on their properties.  
Uses `pw-cli dump` to get the data. 

This small utility has been motivated by the fact that pipewire does not have a nicer way of finding nodes other than using their ids.

Neither it or wireplumber makes it easy through cli to find the node's property.

# Usage

Get node based on property value:
```
$ pw-info find-node node.name alsa_output.pci-0000_31_00.4.analog-stereo
id: 1
type: PipeWire:Interface:Node/3
```

Get node property based on id:
```
$ pw-info find-node-property 1 node.description
node.description: Headphones
```

# Use-cases

For example, if you know the alsa node name; and you use wireplumber to set its custom description for your riced desktop and would like that custom name appearing instead of default name.

```
# alsa-headphones.lua
rule = {
    matches = {
      {
        { "node.name", "equals", "alsa_output.pci-0000_31_00.4.analog-stereo" },
      },
    },
    apply_properties = {
      ["node.description"] = "Headphones",
    },
  }
  
  table.insert(alsa_monitor.rules,rule)
```

Now, getting the name of the default sink works kinda like this with this utility:

```sh
# this will get you something like: alsa_output.pci-0000_31_00.4.analog-stereo
DEFAULT_SINK_ID=$(pactl get-default-sink)

# this will get you the volume in percentage
VOLUME=$(pamixer --get-volume-human | grep '%' | head -n 1 | cut -d '[' -f 2 | cut -d '%' -f 1)

# this will get you the description you put in wireplumber: Headphones
CUSTOM_SINK_NAME=$(pw-info find-node -s node.name alsa_output.pci-0000_31_00.4.analog-stereo | head -1 | xargs -I{} pw-info find-node-property -s {} node.description)

# we add everything up: 20% Headphone
echo "${VOLUME}% ${CUSTOM_SINK_NAME}"
```