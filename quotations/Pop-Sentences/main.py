import os

from flask import Flask, request, jsonify, make_response
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
import json
import random
import html

app = Flask(__name__)
limiter = Limiter(app, key_func=get_remote_address)


def filter_input(input_str):
    # 过滤用户输入，只保留合法字符
    return ''.join(char for char in input_str if char.isalnum() or char in ['_', '-', ' '])


def generate_response_content(file_path='./wyy.json', need_parameter=False, *parameters):
    """
    生成响应内容

    参数：
    file_path -- 语录文件所在路径
    need_parameter -- 是否需要传入参数
    *parameters -- 参数元组

    返回值：
    响应内容
    """
    full_msg = request.full_path  # 获取完整的请求消息
    # 检查是否包含type=json字符串，并提取出来
    response_type = 'string'
    if 'type=json' in full_msg.lower():
        response_type = 'json'

    if need_parameter:
        if len(parameters) == 1:
            if parameters[0] is None:
                error_message = '请提供msg参数作为请求。'
                error_response = {'error': error_message, 'code': 400}
                if response_type == 'json':
                    response = make_response(json.dumps(error_response, ensure_ascii=False).encode('utf-8'), 400)
                    response.headers["Content-Type"] = "application/json;charset=utf-8"
                    return response
                else:
                    return error_message
            # 过滤用户输入，防止恶意输入和木马代码
            filtered_msg = filter_input(parameters[0])
            filtered_msg = filtered_msg.rstrip('typejson')
            # 转义用户输入，防止XSS攻击
            escaped_msg = html.escape(filtered_msg)

            original = '<name>'  # 将此处替换为您想要替换的固定词语
            replacer = escaped_msg

            response_data = {'data': get_pop_sentence(file_path, [original], [replacer]),
                             'code': 200}
            if response_type == 'json':
                response = make_response(json.dumps(response_data, ensure_ascii=False).encode('utf-8'), 200)
                response.headers["Content-Type"] = "application/json;charset=utf-8"
                return response
            else:
                return get_pop_sentence(file_path, [original], [replacer])
        elif len(parameters) == 2:
            # 过滤用户输入，防止恶意输入和木马代码
            filtered_msg_1 = filter_input(parameters[0]).rstrip('typejson')
            filtered_msg_2 = filter_input(parameters[1]).rstrip('typejson')
            # 转义用户输入，防止XSS攻击
            escaped_msg_1 = html.escape(filtered_msg_1)
            escaped_msg_2 = html.escape(filtered_msg_2)

            original_1 = '<name1>'  # 将此处替换为您想要替换的固定词语
            original_2 = '<name2>'  # 将此处替换为您想要替换的固定词语
            replacer_1 = escaped_msg_1
            replacer_2 = escaped_msg_2

            response_data = {'data': get_pop_sentence(file_path, [original_1, original_2], [replacer_1, replacer_2]),
                             'code': 200}
            if response_type == 'json':
                response= make_response(json.dumps(response_data, ensure_ascii=False).encode('utf-8'), 200)
                response.headers["Content-Type"] = "application/json;charset=utf-8"
                return response
            else:
                return get_pop_sentence(file_path, [original_1, original_2], [replacer_1, replacer_2])
    else:
        response_data = {'data': get_pop_sentence(file_path),
                         'code': 200}
        print(response_data)
        if response_type == 'json':
            response = make_response(json.dumps(response_data, ensure_ascii=False).encode('utf-8'), 200)
            response.headers["Content-Type"] = "application/json;charset=utf-8"
            return response
        else:
            return get_pop_sentence(file_path)


def get_pop_sentence(file_path, original=None, replacer=None):
    """
    获取随机语录

    参数：
    file_path -- 语录文件所在路径
    original -- 待替换关键词数组
    replacer -- 替换关键词数组

    返回值：
    随机语录
    """
    if replacer is None:
        replacer = []
    if original is None:
        original = []
    is_file_path = os.path.exists(file_path)
    try:
        if is_file_path:
            with open(file_path, 'r', encoding='utf-8') as file:
                data_list = json.load(file)
                random_item = random.choice(data_list)
                if len(replacer) == 0:
                    return random_item
                else:
                    if len(replacer) == 1:
                        return random_item.replace(original[0], replacer[0])
                    if len(replacer) == 2:
                        return random_item.replace(original[0], replacer[0]).replace(original[1], replacer[1])
        else:
            return make_response('数据文件不存在', 404)
    except Exception as e:
        return make_response(str(e), 500)


@app.errorhandler(429)
def handle_429(error):
    response_data = {'error': 'Too Many Requests', 'code': 429}
    response = make_response(json.dumps(response_data, ensure_ascii=False).encode('utf-8'), 429)
    response.headers["Content-Type"] = "application/json;charset=utf-8"
    return response


@app.route('/', methods=['GET'])
@limiter.limit("10 per minute")  # 使用key_func装饰器设置速率限制
def hello_world():
    return 'blog.gakki.icu'


@app.route('/fd', methods=['GET'])
@limiter.limit("15 per minute")  # 使用key_func装饰器设置速率限制
def process_fd():
    msg = request.args.get('msg')
    return generate_response_content('./psycho.json', True, msg)


@app.route('/tg', methods=['GET'])
@limiter.limit("15 per minute")  # 使用key_func装饰器设置速率限制
def process_tg():
    return generate_response_content('./tg.json')


@app.route('/kfc', methods=['GET'])
@limiter.limit("15 per minute")  # 使用key_func装饰器设置速率限制
def process_kfc():
    return generate_response_content('./kfc.json')


@app.route('/wyy', methods=['GET'])
@limiter.limit("15 per minute")  # 使用key_func装饰器设置速率限制
def process_wyy():
    return generate_response_content()


@app.route('/cp', methods=['GET'])
@limiter.limit("15 per minute")  # 使用key_func装饰器设置速率限制
def process_cp():
    msg_1 = request.args.get('msg')
    msg_2 = request.args.get('msg1')
    return generate_response_content('./cp.json', True, msg_1, msg_2)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
