package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-ini/ini"
	promMW "github.com/labstack/echo-contrib/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var supplyPowerVal = -1.0
var supplyVoltageVal = -1.0

var solarPowerVal = -1.0
var solarVoltageVal = -1.0

var supplyVoltageGauge *prometheus.GaugeVec
var supplyPowerGauge *prometheus.GaugeVec
var supplyCurrentGauge *prometheus.GaugeVec

var solarVoltageGauge *prometheus.GaugeVec
var solarPowerGauge *prometheus.GaugeVec
var solarCurrentGauge *prometheus.GaugeVec

var topicRoot string

func main() {
	cfg, err := ini.Load("../config.ini")
	if err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	mqttSection := cfg.Section("MQTT")
	serverAddress := mqttSection.Key("server_address").String()
	topic := mqttSection.Key("topic").String()

	topicRoot = topic

	e := echo.New()

	p := promMW.NewPrometheus("echo", nil)

	// supply
	supplyVoltageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shelly_em_voltage",
		Help: "Voltage received from a Shelly EM.",
	}, []string{"shelly"})
	if err := prometheus.Register(supplyVoltageGauge); err != nil {
		log.Fatal(err)
	}
	supplyPowerGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shelly_em_power",
		Help: "Power received from a Shelly EM.",
	}, []string{"shelly"})
	if err := prometheus.Register(supplyPowerGauge); err != nil {
		log.Fatal(err)
	}
	supplyCurrentGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shelly_em_current",
		Help: "Current calculated from Voltage and Power received from a Shelly EM.",
	}, []string{"shelly"})
	if err := prometheus.Register(supplyCurrentGauge); err != nil {
		log.Fatal(err)
	}

	// solar
	solarVoltageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shelly_em_solar_voltage",
		Help: "Voltage received from a Shelly EM solar.",
	}, []string{"shelly"})
	if err := prometheus.Register(solarVoltageGauge); err != nil {
		log.Fatal(err)
	}
	solarPowerGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shelly_em_solar_power",
		Help: "Power received from a Shelly EM solar.",
	}, []string{"shelly"})
	if err := prometheus.Register(solarPowerGauge); err != nil {
		log.Fatal(err)
	}
	solarCurrentGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shelly_em_solar_current",
		Help: "Current calculated from Voltage and Power received from a Shelly EM solar.",
	}, []string{"shelly"})
	if err := prometheus.Register(solarCurrentGauge); err != nil {
		log.Fatal(err)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Pow %f, Volt %f", supplyPowerVal, supplyVoltageVal))
	})

	go func() {
		// mqtt.DEBUG = log.New(os.Stdout, "", 0)
		mqtt.ERROR = log.New(os.Stdout, "", 0)
		opts := mqtt.NewClientOptions().AddBroker(
			fmt.Sprintf("tcp://%s:1883", serverAddress),
		).SetClientID(fmt.Sprintf("mqtt_to_prom%d", time.Now().Unix()))

		opts.SetKeepAlive(60 * time.Second)
		// Set the message callback handler
		opts.SetDefaultPublishHandler(f)
		opts.SetPingTimeout(1 * time.Second)

		c := mqtt.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		// Subscribe to a topic
		if token := c.Subscribe(fmt.Sprintf("%s/emeter/#", topic), 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}()

	p.Use(e) // Enable metrics middleware and register route to expose metrics

	e.Logger.Fatal(e.Start(":1323"))
}

func getTopic(source string, measurement string) string {
	return fmt.Sprintf("%s/emeter/%s/%s", topicRoot, source, measurement)
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := msg.Payload()

	// Supply
	if topic == getTopic("1", "voltage") {
		voltNum, err := strconv.ParseFloat(string(payload), 64)
		if err != nil {
			log.Fatalf(err.Error())
			os.Exit(1)
		}
		supplyVoltageGauge.WithLabelValues("em").Set(voltNum)
		supplyVoltageVal = voltNum
		supplyCurrentGauge.WithLabelValues("em").Set(supplyPowerVal / supplyVoltageVal)
	}
	if topic == getTopic("1", "power") {
		powerNum, err := strconv.ParseFloat(string(payload), 64)
		if err != nil {
			log.Fatalf(err.Error())
			os.Exit(1)
		}
		log.Println("Setting supply power to", powerNum)
		supplyPowerGauge.WithLabelValues("em").Set(powerNum)
		supplyPowerVal = powerNum
		supplyCurrentGauge.WithLabelValues("em").Set(supplyPowerVal / supplyVoltageVal)
	}
	// SOlar
	if topic == getTopic("0", "voltage") {
		voltNum, err := strconv.ParseFloat(string(payload), 64)
		if err != nil {
			log.Fatalf(err.Error())
			os.Exit(1)
		}
		solarVoltageGauge.WithLabelValues("em").Set(voltNum)
		solarVoltageVal = voltNum
		solarCurrentGauge.WithLabelValues("em").Set(solarPowerVal / solarVoltageVal)
	}
	if topic == getTopic("0", "power") {
		powerNum, err := strconv.ParseFloat(string(payload), 64)
		if err != nil {
			log.Fatalf(err.Error())
			os.Exit(1)
		}
		log.Println("Setting solar power to", powerNum)
		solarPowerGauge.WithLabelValues("em").Set(powerNum)
		solarPowerVal = powerNum
		solarCurrentGauge.WithLabelValues("em").Set(solarPowerVal / solarVoltageVal)
	}
	// fmt.Printf("TOPIC: %s MSG: %s\n", topic, payload)
}
