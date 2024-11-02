#include <SimpleDHT.h>

// Pin configuration for DHT22
const int pinDHT22 = 2;
SimpleDHT22 dht22(pinDHT22);

// System configuration
const int MIN_AC_TEMP = 16;
const int MAX_AC_TEMP = 30;
const float TEMP_UPPER_THRESHOLD = 28.0;
const float TEMP_LOWER_THRESHOLD = 26.0;

// Global variables
int sleeping_Condition = 1;  // 0: light sleep, 1: deep sleep
int AC_temperature = 24;     // Default AC temperature
int infared_Sensor[15];      // IR codes for temperatures 16-30
int signal2 = 0;

// Function prototypes
int AC_Monitoring(float temperature);

void setup() {
  Serial.begin(115200);

  // Initialize infrared sensor array (you'll need to add your AC's IR codes)
  for (int i = 0; i < 15; i++) {
    infared_Sensor[i] = 0;  // Replace with actual IR codes
  }
}

void loop() {
  Serial.println("=================================");
  Serial.println("Reading DHT22...");

  float temperature = 0;
  float humidity = 0;
  int err = SimpleDHTErrSuccess;

  // Read temperature and humidity
  if ((err = dht22.read2(&temperature, &humidity, NULL)) != SimpleDHTErrSuccess) {
    Serial.print("Read DHT22 failed, err=");
    Serial.println(err);
    delay(2000);
    return;  // Changed from return 0 to return
  }

  // Validate AC temperature range
  if (AC_temperature < MIN_AC_TEMP || AC_temperature > MAX_AC_TEMP) {
    Serial.println("AC temperature out of range!");
    AC_temperature = 24;  // Reset to default
    delay(2000);
    return;  // Changed from return 0 to return
  }

  // Print current readings
  Serial.print("Reading OK: ");
  Serial.print(temperature, 1);  // Display with 1 decimal place
  Serial.print(" *C, ");
  Serial.print(humidity, 1);
  Serial.println(" RH%");
  Serial.print("Current AC temperature: ");
  Serial.println(AC_temperature);

  // Monitor and adjust AC
  signal2 = AC_Monitoring(temperature);
  if (signal2 != 0) {
    Serial.print("Adjusting AC to: ");
    Serial.println(AC_temperature);
  }

  delay(2500);
}

int AC_Monitoring(float temperature) {
  int index;

  if (temperature > TEMP_UPPER_THRESHOLD) {
    // Decrease temperature
    if (AC_temperature > MIN_AC_TEMP) {
      AC_temperature--;
      index = AC_temperature - MIN_AC_TEMP;
      return infared_Sensor[index];
    }
  } else if (temperature < TEMP_LOWER_THRESHOLD) {
    // Increase temperature
    if (AC_temperature < MAX_AC_TEMP) {
      AC_temperature++;
      index = AC_temperature - MIN_AC_TEMP;
      return infared_Sensor[index];
    }
  }

  return 0;  // No change needed
}