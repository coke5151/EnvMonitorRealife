// 定義 L9110 控制腳位
const int IN_A = 16;  // GPIO16
const int IN_B = 17;  // GPIO17

void setup() {
  // 設定腳位為輸出模式
  pinMode(IN_A, OUTPUT);
  pinMode(IN_B, OUTPUT);
}

void loop() {
  digitalWrite(IN_A, LOW);
  digitalWrite(IN_B, HIGH);
} 