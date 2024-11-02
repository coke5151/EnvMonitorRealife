#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include "DHTesp.h"
#include "config.h"

const int serialNumber = 115200;
const int updatePerMilliseconds = 10000;

// WiFi
const char* ssid = WIFI_SSID;
const char* password = WIFI_PASSWORD;

// DHT22
const int dht22OutPort = 27; // GPIO27
DHTesp dht;

// API
const char* serverInsertURL = INSERT_API_URL;

void setup() {
  Serial.begin(serialNumber);

  // Connect to WiFi
  WiFi.begin(ssid, password);
  Serial.print(String("Connecting to ") + ssid);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.print("\nWiFi status: ");
  WiFi.printDiag(Serial);

  Serial.println("Initializing DHT sensor...");
  dht.setup(dht22OutPort, DHTesp::AM2302);
  // Wait for dht22 stable after initialized.
  delay(2000);
  Serial.println("Initialized DHT sensor.");
}

void loop() {
  // Check if the WiFi is still connected
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("WiFi connection lost");
  } else {
    // Setting http client
    HTTPClient http;
    http.begin(serverInsertURL);
    http.addHeader("Content-Type", "application/json");

    // DHT22
    TempAndHumidity data = dht.getTempAndHumidity();
    if (dht.getStatus() != 0) {
      Serial.println("DHT sensor error status: " + String(dht.getStatusString()));
    } else if (isnan(data.humidity) || isnan(data.temperature)) {
      Serial.println("Data is NaN!");
    } else {
      Serial.print("Humidity: ");
      Serial.print(data.humidity);
      Serial.print("% Temperature: ");
      Serial.println(data.temperature);
      Serial.println("Sending data to server...");
      
      // Create json
      JsonDocument doc;
      doc["temperature"] = data.temperature;
      doc["humidity_percentage"] = data.humidity;
      
      String jsonString;
      serializeJson(doc, jsonString);

      // Send POST request
      int httpResponseCode = http.POST(jsonString);
      if(httpResponseCode > 0) {
        String response = http.getString();
        Serial.println("HTTP Reponse Code: " + String(httpResponseCode));
        Serial.println("Reponse body: " + response);
      } else {
        Serial.println("Error code: " + String(httpResponseCode));
      }
      http.end();
    }
  }
  delay(updatePerMilliseconds);
}