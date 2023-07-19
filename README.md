# 无人智能零售柜

## 作品简介

我们的作品是利用Pegasus，Taurus、五线电插锁等部件搭建的无人智能饮料零售柜。本款无人智能零售柜提供了一种更加便捷、智能、节约物资的饮料购买方式，结合Taurus图像检测技术，可以分析得到消费者拿取的产品，自动结算。

## 作品架构

### 架构

<img src="./pic/comp.png" alt="image-20230719192920859" style="zoom: 67%;" />

项目主要分为：小程序端，门锁端，检测端以及中间服务端。

### 流程

用户通过小程序点击开门按钮，Pegasus通过WiFi与小程序互联，接受到信号后控制电插锁开门。用户取出饮料后，Arduino接收到光电传感器信号。Taurus对饮料进行实时检测并通过串口发送给树莓派，树莓派进行处理后生成购物车传给小程序。



## 代码介绍