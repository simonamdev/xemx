# xemx

A proof of concept for a locally hosted electricity consumption and generation monitoring system.

_xemx_ is the maltese word for sun. My original motivation for this project was to be able to track the electrical generation and consumption of my home with minimal hardware required.

# Rationale

[Half of Malta's subsidy spend is on energy](https://www.maltatoday.com.mt/news/national/123554/energy_support_equal_to_half_of_maltas_subsidy_spenda). This meant that, while prices did not increase despite everything happening in the world, the deficit has ballooned to the point that the EU Commission has [recommended using it on that instead](https://timesofmalta.com/articles/view/european-council-tells-malta-wind-energy-support-measures.1033652).

Watching these events unfold and realising the only observability I had on my consumption was my tri-monthly bill, I decided to apply my skills in being more pro-active, giving myself a solid foundation off of which to build electricity bill estimation and tracking the return of investment I made in solar power.

# Setup

## Hardware

It uses a [Shelly EM + 2x 50A Clamp](https://www.shelly.com/en/products/shop/shelly-em-120a/shelly-em-2x-50a) to track both consumption and generation. For solar installations with batteries this will not provide fine-grained information into whether the battery is charging or discharging, however it will still indicate the net performance of the entire system.

## Software

The Shelly EM is setup to use MQTT pointing to a locally hosted MQTT broker. A microservice written in Go listens to this broker, converting it to Prometheus metrics. Prometheus is then configured to scrape these metrics, which are then queried by Grafana for visualisation.

![Architecture of the system](https://github.com/simonamdev/xemx/blob/main/architecture.jpg)

This proof of concept hosts everything using a single docker compose. In my actual homelab I host it in a local Kubernetes cluster.

# Usage

Setup your Shelly EM to publish to the MQTT broker's address set in the config file `config.ini`.

In absence of a Shelly EM, run the python script `pretend.py` to receive mock data to that same broker instead.

Make sure you build the `mqtt-to-prom` service by running `build.sh`. Run the entire stack using `docker-compose up`.

MQTT output can be verified using a tool such as [MQTT Explorer](https://github.com/thomasnordquist/MQTT-Explorer).

![MQTT Explorer output](https://github.com/simonamdev/xemx/blob/main/mqtt.jpg)

Prometheus output can be verified by visting `127.0.0.1:9090` and running a query for `shelly_em`:

![Prometheus output](https://github.com/simonamdev/xemx/blob/main/prom.jpg)

Grafana can be accessed at `127.0.0.1:3000`. The username and password are `admin/admin`. You can see a working dashboard on the `Xemx` dashboard.
Here is a screenshot from my internal dashboard (similar to the one in this repo but with real values):

![Grafana Dashboard](https://github.com/simonamdev/xemx/blob/main/dash.jpg)

# FAQ

* _This code is terrible! Why would you use global variables?_: Global variables are stupid until they work then they're not stupid. In my limited spare time to put something out there, I chose done over perfect. The scope of the code is also very small and inconsequential, I don't condone the use of global variables and wanton duplication of code so easily in my professional work.

* _Why introduce Prometheus when Grafana can read MQTT directly with a plugin?_: This allows me to plug in other devices such as other Shelly EMs or Plus Plugs in the future very easily. I also already had a Prometheus + Grafana setup in my homelab so I leveraged it to get to a working solution quicker.

* _Are you being paid by Shelly?_: Not at all, I'm just a happy customer 🙂.

# Licence

MIT