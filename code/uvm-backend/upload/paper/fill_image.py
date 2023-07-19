import os
from PIL import Image
from PIL import ImageFile

ImageFile.LOAD_TRUNCATED_IMAGES = True


IMAGES_PATH = '../img/'  # 图片集地址
IMAGES_FORMAT = ['.jpg', '.JPG']  # 图片格式


def make_square(im, min_size=256, fill_color=(0, 0, 0)):
    x, y = im.size
    size = max(min_size, x, y)
    new_im = Image.new('RGB', (size, size), fill_color)
    new_im.paste(im, (int((size - x) / 2), int((size - y) / 2)))
    return new_im

def fill(name):
    test_image = Image.open(IMAGES_PATH + name)
    new_image = make_square(test_image, fill_color=(255, 255, 255))
    new_image.save(name)

if __name__ == '__main__':
    image_names = [name for name in os.listdir(IMAGES_PATH) for item in IMAGES_FORMAT if
                   os.path.splitext(name)[1] == item]
    for name in image_names:
        fill(name)
