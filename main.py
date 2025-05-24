from awswaf.aws import AwsWaf
from curl_cffi import requests

session = requests.Session(impersonate="chrome")

session.headers = headers = {
    'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8',
    'accept-language': 'en-US,en;q=0.5',
    'cache-control': 'no-cache',
    'dnt': '1',
    'pragma': 'no-cache',
    'priority': 'u=0, i',
    'sec-ch-ua': '"Chromium";v="136", "Brave";v="136", "Not.A/Brand";v="99"',
    'sec-ch-ua-mobile': '?0',
    'sec-ch-ua-platform': '"Windows"',
    'sec-fetch-dest': 'document',
    'sec-fetch-mode': 'navigate',
    'sec-fetch-site': 'none',
    'sec-fetch-user': '?1',
    'sec-gpc': '1',
    'upgrade-insecure-requests': '1',
    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36',
    # 'cookie': 'BNC_FV_KEY=3305e06383e0a86291170647f2baf18e94db4fb0; lang=en; language=en; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22196c4cd584dfae-009c68c9b5b17ab-26011f51-3686400-196c4cd584e28f4%22%2C%22first_id%22%3A%22%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_referrer%22%3A%22%22%7D%2C%22identities%22%3A%22eyIkaWRlbnRpdHlfY29va2llX2lkIjoiMTk2YzRjZDU4NGRmYWUtMDA5YzY4YzliNWIxN2FiLTI2MDExZjUxLTM2ODY0MDAtMTk2YzRjZDU4NGUyOGY0In0%3D%22%2C%22history_login_id%22%3A%7B%22name%22%3A%22%22%2C%22value%22%3A%22%22%7D%7D; aws-waf-token=0d08bc50-c13f-4fec-a50b-1cc76ecd9fb6:CQoAfRF8tiSlAAAA:F6FHIaJNi4Ow+/4oj+ossQqeeFikd4jwhhS4AVV/9Z9cSQSX0xd4xzXcyf43tpV9+lUyHdUioKIn7Dg+YtGt/qhLiGjknAQl+OvIEPs1pcruTYFgWBq3U4BmyTKilV50hMqtta1cFbrdOLALDOYHDVz7fqQrdItTxnplCO6QlPW/xWDPm6EYetRRqGtxIfESaqPwkg3HY9L8yg==; bnc-uuid=2282a746-aa05-422c-a06f-cb44962863db; BNC_FV_KEY_T=101-q7dCOkqNNifdlyfviAWoWfRY%2BAhup7KFQNyLIIpw0iNxWTUNjH6R9IpMR8oZJ6Lh5Ine1aTZxDDMhOLwjGW0Nw%3D%3D-DtTG9YH29RHv9APtSAj47A%3D%3D-02; BNC_FV_KEY_EXPIRE=1747871443858; se_gd=hsAVADgEOHICVMHETFwQgZZDBXVoIBUW1Zd5QVE9lZcVwAFNWVgS1; se_gsd=dzYlOz9/LCsiBicsNCUnChAkB1QTBQtSV1RBV11aV1ZaElNT1',
}

response = session.get("https://www.binance.com/")
goku = AwsWaf.extract_goku_props(response.text)

token = AwsWaf(goku)()
session.headers.update({
    "cookie": "aws-waf-token=" + token
})
# print(session.headers)
print(token)
print(session.get("https://www.binance.com/").text)
