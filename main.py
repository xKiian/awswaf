from awswaf.aws import AwsWaf
from curl_cffi import requests

session = requests.Session(impersonate="chrome")

session.headers = {
    "host": "www.binance.com",
    "connection": "keep-alive",
    "pragma": "no-cache",
    "cache-control": "no-cache",
    "sec-ch-ua": "\"Chromium\";v=\"136\", \"Google Chrome\";v=\"136\", \"Not.A/Brand\";v=\"99\"",
    "sec-ch-ua-mobile": "?0",
    "sec-ch-ua-platform": "\"Windows\"",
    "upgrade-insecure-requests": "1",
    "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
    "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
    "sec-fetch-site": "none",
    "sec-fetch-mode": "navigate",
    "sec-fetch-user": "?1",
    "sec-fetch-dest": "document",
    "accept-encoding": "gzip, deflate, br, zstd",
    "accept-language": "en-US,en;q=0.9"
}

response = session.get("https://www.binance.com/")

goku = AwsWaf.extract_goku_props(response.text)

token = AwsWaf(goku)()
session.headers.update({
    "cookie": "aws-waf-token=" + token
})
#print(session.headers)
print(token)
#print(session.get("https://www.binance.com/").text)
