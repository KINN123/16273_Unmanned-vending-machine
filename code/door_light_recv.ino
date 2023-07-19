const int pinDoor = 3;
const int pinLight = 4;
const int thresh = 1000;
bool door_closed_flag = true;
bool light_in_flag = false;
int light_one_count = 0;
int door_one_count = 0;
void setup()
{
  // put your setup code here, to run once:
  pinMode(pinDoor, INPUT);
  pinMode(pinLight, INPUT); //设置pinRelay脚为输出状态
  Serial.begin(9600);
}

void loop()
{
  int read_light = digitalRead(pinLight);
  int read_door = digitalRead(pinDoor);
  // Serial.println(read_light);  // 没有连线时，输出为 0
  // Serial.println(read_door);   // 没有连线时，输出为 0

  // 检测门状态
  if (!door_closed_flag) {  // 如果门是开的，检测到高电位，计数器计数
    if (read_door == HIGH) {
      door_one_count++;
    } else {
      door_one_count = 0;
    }
  }
  // 如果门是关的，此时检测到低电位，开门，输出'5'
  else if (read_door == LOW) {
    if (read_door == LOW) {
      door_one_count--;
    } else {
      door_one_count = 0;
    }
  }
  if (door_one_count > thresh) {  // 如果计数器大于阈值thresh， 设置门状态为关门，输出 '6'
    door_closed_flag = true;
    door_one_count = 0;
    Serial.println("6");
  } else if (door_one_count < -1 * thresh) { // 如果计数器小于负阈值-thresh， 设置门状态为开门，输出 '5'
    door_one_count = 0;
    door_closed_flag = false;
    Serial.println("5");
  }

  // 检测光栅
  if (!light_in_flag) {  // 如果没有物体进入，光栅电位变为高电位，计数器计数
    if (read_light == HIGH) {
      light_one_count++;
    }
    else {
      light_one_count = 0;
    }
  }
  else {  // 如果有物体进入时，光栅电位变为低电位，表示物体出去，输出 '4'
    if (read_light == LOW) {
      light_one_count--;
    }
    else {
      light_one_count = 0;
    }
  }
  if (light_one_count > thresh) {  // 如果计数器值变为大于阈值thresh，表示物体进入，输出 '3'
    light_in_flag = true;
    light_one_count = 0;
    Serial.println("3");
  } else if (light_one_count < -1 * thresh) {  // 如果计数器值变为小于负阈值-thresh，表示物体进入，输出 '4'
    light_one_count = 0;
    light_in_flag = false;
    Serial.println("4");
  }
}
