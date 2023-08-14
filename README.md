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

[./architecture.png]

# Usage

TODO

# FAQ

TODO

# Licence

MIT