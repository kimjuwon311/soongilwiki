from .tool.func import *

async def api_bbs_w_tabom(sub_code = ''):
    other_set = {}
    other_set["sub_code"] = sub_code
    other_set["ip"] = ip_check()

    func_name = sys._getframe().f_code.co_name
    if flask.request.method == 'POST':
        func_name += '_post'

    return flask.Response(response = (await python_to_golang(func_name, other_set)), status = 200, mimetype = 'application/json')