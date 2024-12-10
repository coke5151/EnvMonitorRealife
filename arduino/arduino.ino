#include <WiFi.h>
#include <HTTPClient.h>
#include <ArduinoJson.h>
#include "DHTesp.h"
#include "config.h"

const int serialNumber = 115200;
const int updateDataInterval = 10000;  // 上傳資料間隔
const int checkFanInterval = 1000;     // 檢查風扇狀態間隔
unsigned long lastDataUpdate = 0;
unsigned long lastFanCheck = 0;

// WiFi
const char* ssid = WIFI_SSID;
const char* password = WIFI_PASSWORD;

// DHT22
const int dht22OutPort = 27; // GPIO27
DHTesp dht;

// Fan
const int fanInA = 16;  // GPIO16
const int fanInB = 17;  // GPIO17
String currentFanStatus = "off";

// API
const char* serverInsertURL = INSERT_API_URL;
const char* getFanStatusURL = GET_FAN_STATUS_URL;

void setFanStatus(String status) {
  if (status == "off") {
    digitalWrite(fanInA, LOW);
    digitalWrite(fanInB, LOW);
    currentFanStatus = "off";
  } else if (status == "clockwise") {
    digitalWrite(fanInA, HIGH);
    digitalWrite(fanInB, LOW);
    currentFanStatus = "clockwise";
  } else if (status == "counterclockwise") {
    digitalWrite(fanInA, LOW);
    digitalWrite(fanInB, HIGH);
    currentFanStatus = "counterclockwise";
  }
}

void checkAndUpdateFanStatus() {
    HTTPClient http;
    http.begin(getFanStatusURL);
    
    int httpCode = http.GET();
    if (httpCode <= 0) {
        Serial.println("HTTP 請求失敗：" + String(httpCode));
        http.end();
        return;
    }
    
    String response = http.getString();
    JsonDocument doc;
    DeserializationError error = deserializeJson(doc, response);
    
    if (error) {
        Serial.println("JSON 解析失敗：" + String(error.c_str()));
        http.end();
        return;
    }
    
    String newStatus = doc["status"].as<String>();
    if (newStatus != currentFanStatus) {
        Serial.println("更新風扇狀態：" + newStatus);
        setFanStatus(newStatus);
    }
    
    http.end();
}

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

  // Initialize DHT
  Serial.println("Initializing DHT sensor...");
  dht.setup(dht22OutPort, DHTesp::AM2302);
  delay(2000);  // Wait for dht22 stable after initialized.
  Serial.println("Initialized DHT sensor.");

  // Initialize Fan
  Serial.println("Initializing Fan...");
  pinMode(fanInA, OUTPUT);
  pinMode(fanInB, OUTPUT);
  digitalWrite(fanInA, LOW);
  digitalWrite(fanInB, LOW);
  Serial.println("Initialized Fan");
}

void loop() {
  unsigned long currentMillis = millis();
  
  // 檢查風扇狀態 (每秒)
  if (currentMillis - lastFanCheck >= checkFanInterval) {
    if (WiFi.status() == WL_CONNECTED) {
      checkAndUpdateFanStatus();
    }
    lastFanCheck = currentMillis;
  }
  
  // 上傳溫濕度資料 (每10秒)
  if (currentMillis - lastDataUpdate >= updateDataInterval) {
    if (WiFi.status() == WL_CONNECTED) {
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
        doc["fan_status"] = currentFanStatus;
        
        String jsonString;
        serializeJson(doc, jsonString);

        // Send POST request
        int httpResponseCode = http.POST(jsonString);
        if(httpResponseCode > 0) {
          String response = http.getString();
          Serial.println("HTTP Response Code: " + String(httpResponseCode));
          Serial.println("Response body: " + response);
        } else {
          Serial.println("Error code: " + String(httpResponseCode));
        }
        http.end();
      }
    }
    lastDataUpdate = currentMillis;
  }
}