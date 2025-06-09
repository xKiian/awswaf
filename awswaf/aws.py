import json
from curl_cffi import requests
from awswaf.verify import CHALLENGE_TYPES
from awswaf.fingerprint import get_fp


class AwsWaf:
    def __init__(self, goku_props, token=None):
        self.session = requests.Session(impersonate="chrome")
        self.session.headers = {
            "host": "fe4385362baa.ead381d8.eu-west-1.token.awswaf.com",
            "connection": "keep-alive",
            "sec-ch-ua-platform": "\"Windows\"",
            "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
            "sec-ch-ua": "\"Chromium\";v=\"136\", \"Google Chrome\";v=\"136\", \"Not.A/Brand\";v=\"99\"",
            "sec-ch-ua-mobile": "?0",
            "accept": "*/*",
            "origin": "https://www.binance.com",
            "sec-fetch-site": "cross-site",
            "sec-fetch-mode": "cors",
            "sec-fetch-dest": "empty",
            "referer": "https://www.binance.com/",
            "accept-encoding": "gzip, deflate, br, zstd",
            "accept-language": "en-US,en;q=0.9"
        }
        self.goku_props = goku_props
        self.token = token
        self.session.get(
            "https://fe4385362baa.ead381d8.eu-west-1.token.awswaf.com/fe4385362baa/306922cde096/8b22eb923d34/challenge.js")

    @staticmethod
    def extract_goku_props(html: str):
        return json.loads(html.split("window.gokuProps = ")[1].split(";")[0])

    def get_inputs(self):
        return self.session.get(
            "https://fe4385362baa.ead381d8.eu-west-1.token.awswaf.com/fe4385362baa/306922cde096/8b22eb923d34/inputs?client=browser").json()

    def build_payload(self, inputs: dict):
        verify = CHALLENGE_TYPES[inputs["challenge_type"]]
        checksum, fp = get_fp()
        return {
            "challenge": inputs["challenge"],
            "checksum": checksum,
            "solution": verify(inputs["challenge"]["input"], checksum, inputs["difficulty"]),
            "signals": [{"name": "KramerAndRio", "value": {"Present": fp}}],
            "existing_token": None,
            "client": "Browser",
            "domain": "www.binance.com",
            "metrics": [
                {
                    "name": "2",
                    "value": 0.4000000022351742,
                    "unit": "2"
                },
                {
                    "name": "100",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "101",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "102",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "103",
                    "value": 8,
                    "unit": "2"
                },
                {
                    "name": "104",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "105",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "106",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "107",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "108",
                    "value": 1,
                    "unit": "2"
                },
                {
                    "name": "undefined",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "110",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "111",
                    "value": 2,
                    "unit": "2"
                },
                {
                    "name": "112",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "undefined",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "3",
                    "value": 4,
                    "unit": "2"
                },
                {
                    "name": "7",
                    "value": 0,
                    "unit": "4"
                },
                {
                    "name": "1",
                    "value": 16.299999997019728,
                    "unit": "2"
                },
                {
                    "name": "4",
                    "value": 36.5,
                    "unit": "2"
                },
                {
                    "name": "5",
                    "value": 0.10000300149011612,
                    "unit": "2"
                },
                {
                    "name": "6",
                    "value": 52.899999918509884,
                    "unit": "2"
                },
                {
                    "name": "0",
                    "value": 136.89999919850988,
                    "unit": "2"
                },
                {
                    "name": "8",
                    "value": 1,
                    "unit": "4"
                }
            ],
            "goku_props": self.goku_props,
        }

    def verify(self, payload):
        self.session.headers = {
            "host": "fe4385362baa.ead381d8.eu-west-1.token.awswaf.com",
            "connection": "keep-alive",
            "sec-ch-ua-platform": "\"Windows\"",
            "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
            "sec-ch-ua": "\"Chromium\";v=\"136\", \"Google Chrome\";v=\"136\", \"Not.A/Brand\";v=\"99\"",
            "content-type": "text/plain;charset=UTF-8",
            "sec-ch-ua-mobile": "?0",
            "accept": "*/*",
            "origin": "https://www.binance.com",
            "sec-fetch-site": "cross-site",
            "sec-fetch-mode": "cors",
            "sec-fetch-dest": "empty",
            "referer": "https://www.binance.com/",
            "accept-encoding": "gzip, deflate, br, zstd",
            "accept-language": "en-US,en;q=0.9"
        }
        res = self.session.post(
            "https://fe4385362baa.ead381d8.eu-west-1.token.awswaf.com/fe4385362baa/306922cde096/8b22eb923d34/verify",
            json=payload).json()
        return res["token"]

    def __call__(self):
        inputs = self.get_inputs()
        payload = self.build_payload(inputs)
        return self.verify(payload)
