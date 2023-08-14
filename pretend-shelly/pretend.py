import time
import paho.mqtt.client as mqtt
import configparser

config = configparser.ConfigParser()
config.read('../mqtt-to-prom/config.ini')

mqtt_server_address = config.get('MQTT', 'server_address')
mqtt_topic = config.get('MQTT', 'topic')

def on_connect(client, userdata, flags, rc):
    print(f'Connected with result code {str(rc)}')

client = mqtt.Client()
client.on_connect = on_connect

client.connect(mqtt_server_address, 1883, 60)

client.loop_start()

while True:
    client.publish(f'{mqtt_topic}/emeter/0/power', 100)
    client.publish(f'{mqtt_topic}/emeter/0/reactive_power', 100)
    client.publish(f'{mqtt_topic}/emeter/0/pf', -0.98)
    client.publish(f'{mqtt_topic}/emeter/0/voltage', 242.27)
    client.publish(f'{mqtt_topic}/emeter/0/total', 500000)
    client.publish(f'{mqtt_topic}/emeter/0/total_returned', 8000)
    client.publish(f'{mqtt_topic}/emeter/0/energy', 400)
    client.publish(f'{mqtt_topic}/emeter/1/power', 200)
    client.publish(f'{mqtt_topic}/emeter/1/reactive_power', 200)
    client.publish(f'{mqtt_topic}/emeter/1/pf', -0.73)
    client.publish(f'{mqtt_topic}/emeter/1/voltage', 242.27)
    client.publish(f'{mqtt_topic}/emeter/1/total', 600000)
    client.publish(f'{mqtt_topic}/emeter/1/total_returned', 8000)
    client.publish(f'{mqtt_topic}/emeter/1/energy', 600)
    print(f'Published at {int(time.time())}')
    
    time.sleep(5)
