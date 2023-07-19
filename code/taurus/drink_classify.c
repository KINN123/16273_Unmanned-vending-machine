/*
 * Copyright (c) 2022 HiSilicon (Shanghai) Technologies CO., LIMITED.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * 该文件提供了基于yolov2的手部检测以及基于resnet18的手势识别，属于两个wk串行推理。
 * 该文件提供了手部检测和手势识别的模型加载、模型卸载、模型推理以及AI flag业务处理的API接口。
 * 若一帧图像中出现多个手，我们通过算法将最大手作为目标手送分类网进行推理，
 * 并将目标手标记为绿色，其他手标记为红色。
 *
 * This file provides hand detection based on yolov2 and gesture recognition based on resnet18,
 * which belongs to two wk serial inferences. This file provides API interfaces for model loading,
 * model unloading, model reasoning, and AI flag business processing for hand detection
 * and gesture recognition. If there are multiple hands in one frame of image,
 * we use the algorithm to use the largest hand as the target hand for inference,
 * and mark the target hand as green and the other hands as red.
 */

#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <errno.h>

#include "sample_comm_nnie.h"
#include "sample_media_ai.h"
#include "ai_infer_process.h"
#include "yolov2_hand_detect.h"
#include "vgs_img.h"
#include "ive_img.h"
#include "misc_util.h"
#include "hisignalling.h"

#ifdef __cplusplus
#if __cplusplus
extern "C" {
#endif
#endif /* End of #ifdef __cplusplus */

#define HAND_FRM_WIDTH     640
#define HAND_FRM_HEIGHT    384
#define DETECT_OBJ_MAX     32
#define RET_NUM_MAX        4
#define DRAW_RETC_THICK    6    // Draw the width of the line
#define WIDTH_LIMIT        32
#define HEIGHT_LIMIT       32
#define IMAGE_WIDTH        224  // The resolution of the model IMAGE sent to the classification is 224*224
#define IMAGE_HEIGHT       224
#define MODEL_FILE_GESTURE    "/userdata/data/temp.wk" // darknet framework wk model

static int biggestBoxIndex;
static IVE_IMAGE_S img;
static DetectObjInfo objs[DETECT_OBJ_MAX] = {0};
static RectBox boxs[DETECT_OBJ_MAX] = {0};
static RectBox objBoxs[DETECT_OBJ_MAX] = {0};
static RectBox remainingBoxs[DETECT_OBJ_MAX] = {0};
static RectBox cnnBoxs[DETECT_OBJ_MAX] = {0}; // Store the results of the classification network
static RecogNumInfo numInfo[RET_NUM_MAX] = {0};
static char objClasses[2 * DETECT_OBJ_MAX + 1] = {0};
static int  clas[DETECT_OBJ_MAX] = {0};
static IVE_IMAGE_S imgIn;
static IVE_IMAGE_S imgDst;
static VIDEO_FRAME_INFO_S frmIn;
static VIDEO_FRAME_INFO_S frmDst;
int uartFd = 0;

/*
 * 加载饮料检测模型 
 * Load hand detect and classify model
 */
HI_S32 Yolo2HandDetectResnetClassifyLoad(uintptr_t* model)
{
    SAMPLE_SVP_NNIE_CFG_S *self = NULL;
    HI_S32 ret;

    ret = CnnCreate(&self, MODEL_FILE_GESTURE);
    *model = ret < 0 ? 0 : (uintptr_t)self;
    HandDetectInit(); // Initialize the hand detection model
    SAMPLE_PRT("Load hand detect classify model success\n");
    /*
     * Uart串口初始化
     * Uart open init
     */
    uartFd = UartOpenInit();
    if (uartFd < 0) {
        printf("uart1 open failed\r\n");
    } else {
        printf("uart1 open successed\r\n");
    }
    return ret;
}

/*
 * 卸载饮料检测和分类模型
 * Unload hand detect and classify model
 */
HI_S32 Yolo2HandDetectResnetClassifyUnload(uintptr_t model)
{
    CnnDestroy((SAMPLE_SVP_NNIE_CFG_S*)model);
    HandDetectExit(); // Uninitialize the hand detection model
    close(uartFd);
    SAMPLE_PRT("Unload hand detect claasify model success\n");

    return 0;
}


/*
 * 商品（饮品）推理
 */
HI_S32 Yolo2HandDetectResnetClassifyCal(uintptr_t model, VIDEO_FRAME_INFO_S *srcFrm, VIDEO_FRAME_INFO_S *dstFrm)
{
    SAMPLE_SVP_NNIE_CFG_S *self = (SAMPLE_SVP_NNIE_CFG_S*)model;
    HI_S32 resLen = 0;
    int objNum;
    int ret;
    int num = 0;
    int class;

    ret = FrmToOrigImg((VIDEO_FRAME_INFO_S*)srcFrm, &img);
    SAMPLE_CHECK_EXPR_RET(ret != HI_SUCCESS, ret, "hand detect for YUV Frm to Img FAIL, ret=%#x\n", ret);

    objNum = HandDetectCal(&img, objs); // Send IMG to the detection net for reasoning
    clas[0] = objNum;
    for (int i = 0; i < objNum; i++) {
        class = objs[i].cls;
        cnnBoxs[i] = objs[i].box;
        RectBox *box = &objs[i].box;
        RectBoxTran(box, HAND_FRM_WIDTH, HAND_FRM_HEIGHT,
            dstFrm->stVFrame.u32Width, dstFrm->stVFrame.u32Height);
        SAMPLE_PRT("yolo2_out: {%d, %d, %d, %d}   class: %d\n", box->xmin, box->ymin, box->xmax, box->ymax, class);
        boxs[i] = *box;
        clas[i+1] = class;
    }

    //将图像检测结果输出为指定数据结构
    int size = sizeof(clas) / sizeof(clas[0]);
    char separator = ',';
    objClasses[0] = '\0'; // 初始化为空字符串
    for (int i = 0; i < size; i++) {
        char numStr[12]; // 存储每个数字的字符串表示（最多 11 位数字）
        sprintf(numStr, "%d", clas[i]); // 将数字转换为字符串
        strcat(objClasses, numStr); // 追加数字字符串到结果中
        if (i < size - 1) {
            strncat(objClasses, &separator, 1); // 追加分隔符到结果中
        }
    }
    // printf("结果: %s\n", objClasses);
    
    /*
     * 将检测信息通过Uart传输
     */ 
    UartSend(uartFd, objClasses, sizeof(objClasses)/sizeof(objClasses[0]));


    /*
     * 当检测到对象时，在DSTFRM中绘制一个矩形
     * When an object is detected, a rectangle is drawn in the DSTFRM
     */
    if (objNum > 0) {
        MppFrmDrawRects(dstFrm, boxs, objNum, RGB888_RED, DRAW_RETC_THICK); // Target hand objnum is equal to 1
    }

    return ret;
}

#ifdef __cplusplus
#if __cplusplus
}
#endif
#endif /* End of #ifdef __cplusplus */
