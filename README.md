# [MeshMap.net](https://meshmap.net/)
A nearly live map of [Meshtastic](https://meshtastic.org/) nodes seen by the official Meshtastic MQTT server

## Features
- Shows all position-reporting nodes heard by Meshtastic's [Public MQTT Server](https://meshtastic.org/docs/software/integrations/mqtt/#public-mqtt-server)
  - Includes nodes self-reporting to MQTT or heard by another node reporting to MQTT
- Node data is updated every minute
- Nodes are removed from the map if their position is not updated after 8 hours
- Search for nodes by name or ID

## FAQs

### How do I get my node on the map?
These are general requirements. Refer to the [official docs](https://meshtastic.org/docs/) or reach out to the fantastic [Meshtastic](https://meshtastic.org/) community for additional support.
- First, make sure you are running a [recent firmware](https://meshtastic.org/downloads/) version
- Use the default primary channel and encryption key
- Enable "OK to MQTT" in LoRa configuration, signaling you want your messages uplinked via MQTT
- Enable position reports from your node
  - This may mean enabling your node's built-in GPS, sharing your phone's location via the app, or setting a fixed position
  - Ensure "Position enabled" is enabled on the primary channel
  - Ensure "Precise location" is disabled on the primary channel and the configured precision is between 23.3 km (14.5 miles) and 364 m (1194 ft), inclusive
    - See [Restrictions on the Public MQTT Server](https://meshtastic.org/docs/software/integrations/mqtt/#restrictions-on-the-public-mqtt-server) for details

If your node can be heard by another node already reporting to MQTT, that's it!

#### To enable MQTT reporting
- Enable the MQTT module, using all default settings, possibly with a custom root topic
  - View nodes around your area on the map to find MQTT topics being used
- Configure your node to connect to wifi or otherwise connect to the internet
- Enable MQTT uplink on your primary channel
  - It is not necessary, and not recommended unless you know what you're doing, to enable MQTT downlink

Note: The "Map reporting" option in the MQTT configuration reports [additional data](https://meshtastic.org/docs/configuration/module/mqtt/#map-reporting-enabled) about your node only. It does not report other nodes heard over LoRa.

### Does the map allow manual/self-reported nodes (not over MQTT)?
No, and that's a feature. The goal of this map is to provide a reasonably up-to-date, reliable data source for node locations.

### What do the lines drawn between nodes represent?
When a node is selected, lines may appear between it and nearby nodes.
These lines indicate communication between nodes over LoRa, but not necessarily a direct connection.
Nodes reporting to MQTT are shown to be communicating with the nodes they report, but the path (hops) is unknown.
[Neighbor Info](https://meshtastic.org/docs/configuration/module/neighbor-info/) reports are parsed to show connections between neighbors.

### Can you add this awesome new feature I just came up with? (Or you'd like to report a bug)
Maybe! Open a GitHub issue and let's discuss it. Pull requests welcome!

### Can I use your code for my own map?
Sure! But please pay attention to the license so we can all benefit from your improvements. :)
