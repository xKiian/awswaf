import json
from curl_cffi import requests
from awswaf.verify import CHALLENGE_TYPES
from awswaf.fingerprint import get_fp


class AwsWaf:
    def __init__(self, goku_props):
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
            "client": "Browser",
            "domain": "accounts.binance.com",
            "existing_token": None,
            "goku_props": self.goku_props,
            "metrics": [
                {
                    "name": "2",
                    "value": 1,
                    "unit": "2"
                },
                {
                    "name": "undefined",
                    "value": 1,
                    "unit": "2"
                },
                {
                    "name": "101",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "102",
                    "value": 1,
                    "unit": "2"
                },
                {
                    "name": "103",
                    "value": 12,
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
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "undefined",
                    "value": 1,
                    "unit": "2"
                },
                {
                    "name": "110",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "111",
                    "value": 6,
                    "unit": "2"
                },
                {
                    "name": "112",
                    "value": 1,
                    "unit": "2"
                },
                {
                    "name": "undefined",
                    "value": 0,
                    "unit": "2"
                },
                {
                    "name": "3",
                    "value": 4.5,
                    "unit": "2"
                },
                {
                    "name": "7",
                    "value": 0,
                    "unit": "4"
                },
                {
                    "name": "1",
                    "value": 29.800000000745058,
                    "unit": "2"
                },
                {
                    "name": "4",
                    "value": 8.099999999627471,
                    "unit": "2"
                },
                {
                    "name": "5",
                    "value": 0.5,
                    "unit": "2"
                },
                {
                    "name": "6",
                    "value": 38.40000000037253,
                    "unit": "2"
                },
                {
                    "name": "undefined",
                    "value": 79.10000000149012,
                    "unit": "2"
                },
                {
                    "name": "8",
                    "value": 1,
                    "unit": "4"
                }
            ],
            "signals": [{"name": "KramerAndRio", "value": {"Present": fp}}],
            "solution": verify(inputs["challenge"]["input"], checksum, inputs["difficulty"])
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
        return self.session.post(
            "https://fe4385362baa.ead381d8.eu-west-1.token.awswaf.com/fe4385362baa/306922cde096/8b22eb923d34/verify",
            json=payload).json()["token"]

    def __call__(self):
        inputs = {
            "challenge": {
                "input": "eyJ2ZXJzaW9uIjoxLCJ1YmlkIjoiMDQ0MTFkMGEtM2FhMC00MDQ0LTlmNjctNTFjNTQ2ZGNmMTU0IiwiYXR0ZW1wdF9pZCI6Ijc0MjI2ZTU5LTIxYjktNDJkMC05NWMyLTAwNjcxMjBhZDNmNCIsImNyZWF0ZV90aW1lIjoiMjAyNS0wNS0yNVQxODowNDoxNS41ODQxMzY2NTJaIiwiZGlmZmljdWx0eSI6OCwiY2hhbGxlbmdlX3R5cGUiOiJIYXNoY2FzaFNIQTIifQ==",
                "hmac": "ogMP2MV01jAutzc2l5fKH21wNr/PQBjMplij5MPraVk=",
                "region": "eu-west-1"
            },
            "challenge_type": "h72f957df656e80ba55f5d8ce2e8c7ccb59687dba3bfb273d54b08a261b2f3002",
            "difficulty": 8
        }
        payload = self.build_payload(inputs)
        print(payload)
        return self.verify(payload)
