# xemx
An electricity consumption and generation monitoring system

_xemx_ is the maltese word for sun. My original motivation for this project was to be able to track the electrical generation and consumption of my home with minimal hardware required.

# Rationale

TODO

# Setup

## Hardware

It uses a [Shelly EM + 2x 50A Clamp](https://www.shelly.com/en/products/shop/shelly-em-120a/shelly-em-2x-50a) to track both consumption and generation. For solar installations with batteries this will not provide fine-grained information into whether the battery is charging or discharging, however it will still indicate the net performance of the entire system.

## Software

The Shelly EM is setup to use MQTT pointing to a locally hosted MQTT broker. A microservice written in Go listens to this broker, converting it to Prometheus metrics. Prometheus is then configured to scrape these metrics, which are then queried by Grafana for visualisation.

![Architecture of the system](https://github.com/simonamdev/xemx/blob/main/architecture.jpg)

This proof of concept hosts everything using a single docker compose. In my actual home I host it in a local Kubernetes cluster.

# Usage

Setup your Shelly EM to publish to the MQTT broker's address set in the config file `config.ini`.

In absence of a Shelly EM, run the python script `pretend.py` to receive mock data to that same broker instead.

Run the entire stack using `docker-compose up`.

You can verify the output using a tool such as [MQTT Explorer](https://github.com/thomasnordquist/MQTT-Explorer).

![MQTT Explorer output](https://github.com/simonamdev/xemx/blob/main/mqtt.jpg)


# FAQ

TODO

# Licence

MIT