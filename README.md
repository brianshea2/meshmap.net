# [MeshMap.net](https://meshmap.net/)
A nearly live map of [Meshtastic](https://meshtastic.org/) nodes seen by the official Meshtastic MQTT server

## Features
- Shows all position-reporting nodes heard by Meshtastic's [official MQTT server](https://meshtastic.org/docs/configuration/module/mqtt/#connect-to-the-default-public-server)
  - Includes nodes self-reporting to MQTT or heard by another node reporting to MQTT
- Node data is updated every minute
- Nodes are removed from the map if their position is not updated after 24 hours
- Search for nodes by name or ID

## FAQs

### How do I get my node on the map?
These are general requirements. Refer to the [official docs](https://meshtastic.org/docs/configuration/) or reach out to the fantastic Meshtastic community
(such as at [r/meshtastic](https://www.reddit.com/r/meshtastic/)) for additional support.
- First, make sure you are running a [recent firmware](https://meshtastic.org/downloads/) version
- Use the default primary channel and encryption key
- Enable "OK to MQTT" in LoRa configuration (signaling you want your messages uplinked via MQTT)
- Enable position reports from your node
  - This may mean enabling your node's built-in GPS, sharing your phone's location via the app, or setting a fixed position
  - Ensure "Position enabled" is enabled on the primary channel
  - Precise locations are filtered (see important update below -- the default precision will work)

If your node can be heard by another node already reporting to MQTT, that's it!

#### Important update as of August, 2024
Meshtastic has [made a change to their MQTT server](https://meshtastic.org/blog/recent-public-mqtt-broker-changes/):

> Only position packets with imprecise location information [10-16 bits] will be passed to the topic, ensuring that sensitive data is not exposed.

The most accurate resolution that conforms to this specification is 364 meters/1194 feet.

#### To enable MQTT reporting
- Enable the MQTT module, using all default settings, possibly with a custom root topic
  - View nodes around your area on the map to find MQTT topics being used
- Configure your node to connect to wifi or otherwise connect to the internet
- Enable MQTT uplink on your primary channel
  - It is not necessary, and not recommended unless you know what you're doing, to enable MQTT downlink

### Does the map allow manual/self-reported nodes (not over MQTT)?
No, and that's a feature. The goal of this map is to provide a reasonably up-to-date, reliable data source for node locations.
This is also why nodes are removed if no position reports are heard after 24 hours.

### Can you add this awesome new feature I just came up with? (Or you'd like to report a bug)
Maybe! Open a GitHub issue and let's discuss it. Pull requests welcome!

### Can I use your code for my own map?
Sure! But please pay attention to the license so we can all benefit from your improvements. :)

### Why do I get an error when trying to build the Docker image?
The included Dockerfile is for building the `meshobserv` program, which is responsible for connecting to the MQTT server and handling node messages.
Meshtastic nodes use [Protocol Buffers](https://protobuf.dev/) to serialize their messages.
The Meshtastic protobuf definitions must be compiled before building `meshobserv`.
See the `scripts` directory for helpful build scripts.
