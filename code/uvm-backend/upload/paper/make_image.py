import PIL.Image as Image
from PIL import ImageFile
import os
 
IMAGES_PATH = '../img/'  # 图片集地址
IMAGES_FORMAT = ['.jpg', '.JPG']  # 图片格式
IMAGE_SIZE = 256  # 每张小图片的大小
IMAGE_ROW = 10  # 图片间隔，也就是合并成一张图后，一共有几行
IMAGE_COLUMN = 10  # 图片间隔，也就是合并成一张图后，一共有几列
IMAGE_SAVE_PATH = './final.jpg'  # 图片转换后的地址
 
ImageFile.LOAD_TRUNCATED_IMAGES = True

# 获取图片集地址下的所有图片名称
# image_names = [name for name in os.listdir(IMAGES_PATH) for item in IMAGES_FORMAT if
#                os.path.splitext(name)[1] == item]
image_names = [
    "BSS.jpg",
    "FANTA.jpg",
    "JDB.jpg",
    "KSFBHC.jpg",
    "VTMSSL.jpg",
    "MD.jpg",
    "MZYGLC.jpg",
    "TYNL.jpg",
    "TYASM.jpg",
    "YYKX.jpg",
    "VITANM.jpg",
    "YLHSNN.jpg",
    "KSFXL.jpg",
    "KSFSMT.jpg",
    "KSFT.jpg",
    "KSFLC.jpg",
    "KSFFY.jpg",
    "KSFMM.jpg",
    "KSFJN.jpg",
    "XMTTLL.jpg",
    "XMLLQN.jpg",
    "YLJH.jpg",
    "MZYRG.jpg",
    "MZYSH.jpg",
    "BYLMZ.jpg",
    "BYBXG.jpg",
    "COCACOLA.jpg",
    "SPIRITNM.jpg",
    "YQSLBS.jpg",
    "YQSLRT.jpg",
    "YQSLRC.jpg",
    "MNDPT.jpg",
    "MNDXG.jpg",
    "NSPEPSI.jpg",
    "JTPEPSI.jpg",
    "DPCP.jpg",
    "HKGN.jpg",
    "KSFWH.jpg",
    "YYMM.jpg",
    "YYLP.jpg",
    "YYHM.jpg",
    "VTMSRS.jpg",
    "CPMW.jpg",
    "BGZMHT.jpg",
    "YYDTGY.jpg",
    "BJKZW.jpg",
    "WHHGWS.jpg",
    "JDLLM.jpg",
    "JDLNM.jpg",
    "DSHC.jpg",
    "BSKL.jpg",
    "COCAXW.jpg",
    "SPIRITXW.jpg",
    "YQYWSD.jpg",
    "YQNMSD.jpg",
    "QCQNSD.jpg",
    "QCMTSD.jpg",
    "NFSQNH.jpg",
    "NFSQYL.jpg",
    "DPWSS.jpg",
    "COCACOF.jpg",
    "KSJCM.jpg",
    "KSJCN.jpg",
    "JJBT.jpg",
    "JJQM.jpg",
    "WTMNM.jpg",
    "KSFCCNM.jpg",
    "COCAYGNM.jpg",
    "CPXM.jpg",
    "COCACH.jpg",
    "COCACL.jpg",
    "COCACP.jpg",
    "QCCCNH.jpg",
    "QCCCTW.jpg",
    "DFSYWL.jpg",
    "NFGYT.jpg",
    "NFGYP.jpg",
    "MZYQO.jpg",
    "MZYZT.jpg",
    "MZYBN.jpg",
    "TYYGXC.jpg",
    "BQYCFS.jpg",
    "YQSLNL.jpg",
    "YQSLYW.jpg",
    "YBDJYS.jpg",
    "YBDJNM.jpg",
    "GYLM.jpg",
    "WTMWX.jpg",
    "YBZCH.jpg",
    "YBZCW.jpg",
    "XWTWQ.jpg",
    "YXCHYQN.jpg",
    "DDNM.jpg",
    "DJKXDW.jpg",
    "YSMTBW.jpg",
    "LLMLH.jpg",
    "DDSL.jpg",
    "QLSL.jpg",
    "DYDPC.jpg",
    "QLHKMB.jpg"
]


# 简单的对于参数的设定和实际图片集的大小进行数量判断
if len(image_names) != IMAGE_ROW * IMAGE_COLUMN:
    raise ValueError("合成图片的参数和要求的数量不能匹配！")

from PIL import Image, ImageDraw, ImageFont

def image_add_text(img, text, left, top, text_color=(0, 0, 0), text_size=40):
    # 创建一个可以在给定图像上绘图的对象
    draw = ImageDraw.Draw(img)
    # 字体的格式 这里的simsun.ttc需要有这个字体
    fontStyle = ImageFont.truetype("simsun.ttc", text_size, encoding="utf-8")
    # 绘制文本
    draw.text((left, top), text, text_color, font=fontStyle)
    return img


# img_path = 'assets/1.jpg'
# im = image_add_text(img_path, '这是一个测试', 50, 100, text_color=(0, 0, 0), text_size=20)
# im.show()


# 定义图像拼接函数
def image_compose():
    to_image = Image.new('RGB', (IMAGE_COLUMN * IMAGE_SIZE, IMAGE_ROW * (IMAGE_SIZE + 40)), color="white") #创建一个新图
    # 循环遍历，把每张图片按顺序粘贴到对应位置上
    for y in range(1, IMAGE_ROW + 1):
        for x in range(1, IMAGE_COLUMN + 1):
            from_image = Image.open(IMAGES_PATH + image_names[IMAGE_COLUMN * (y - 1) + x - 1]).resize((IMAGE_SIZE, IMAGE_SIZE),Image.ANTIALIAS)
            to_image.paste(from_image, ((x - 1) * IMAGE_SIZE, (y - 1) * (IMAGE_SIZE + 40)))
            to_image = image_add_text(to_image, image_names[IMAGE_COLUMN * (y - 1) + x - 1][:-4].center(8), (x - 1) * IMAGE_SIZE + IMAGE_SIZE / 2 - 75, (y - 1) * (IMAGE_SIZE + 40) + IMAGE_SIZE)
    return to_image.save(IMAGE_SAVE_PATH) # 保存新图
image_compose() #调用函数
