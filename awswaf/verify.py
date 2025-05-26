import hashlib
import binascii
from typing import Union, Callable, Any

import pyscrypt

CHAR_MAP = {
    '0': "0000",
    '1': "0001",
    '2': "0010",
    '3': "0011",
    '4': "0100",
    '5': '0101',
    '6': "0110",
    '7': "0111",
    '8': "1000",
    '9': "1001",
    'A': "1010",
    'B': "1011",
    'C': '1100',
    'D': "1101",
    'E': "1110",
    'F': "1111"
}


def difficulty_met(difficulty, hash_hex):
    bits = difficulty // 4
    mapped = ''
    for i in range(bits):
        c = hash_hex[i].upper()
        mapped += str(CHAR_MAP[c])
    leading = int(mapped[:difficulty])
    return leading == 0


def hashPow(challenge_input, checksum, difficulty):
    combined = challenge_input + checksum
    nonce = 0
    while True:
        data = f"{combined}{nonce}"
        try:
            hash_hex = hashlib.sha256(data.encode('utf-8')).hexdigest()
        except Exception as e:
            print(e)
            hash_hex = "sha2H="

        if difficulty_met(difficulty, hash_hex):
            return str(nonce)

        nonce += 1


def scrypt_func(input_str, salt_str, memory_cost):
    salt = salt_str.encode('utf-8')
    input_bytes = input_str.encode('utf-8')
    result = pyscrypt.hash(password=input_bytes,
                           salt=salt,
                           N=memory_cost,
                           r=8,
                           p=1,
                           dkLen=16)
    return binascii.hexlify(result).decode('ascii')


def compute_scrypt_nonce(challenge_input, checksum, difficulty):
    combined = challenge_input + checksum
    memory = 128
    nonce = 0
    while True:
        try:
            hash_result = scrypt_func(f"{combined}{nonce}", checksum, memory)
        except Exception as e:
            print(e)
            hash_result = "scryptH="

        if difficulty_met(difficulty, hash_result):
            return str(nonce)

        nonce += 1


CHALLENGE_TYPES: dict[str, Union[Callable[[Any, Any, Any], str], str]] = {
    'h72f957df656e80ba55f5d8ce2e8c7ccb59687dba3bfb273d54b08a261b2f3002': compute_scrypt_nonce,
    'h7b0c470f0cfe3a80a9e26526ad185f484f6817d0832712a4a37a908786a6a67f': hashPow,
    'ha9faaffd31b4d5ede2a2e19d2d7fd525f66fee61911511960dcbb52d3c48ce25': "mp_verify"
}

if __name__ == "__main__":
    print(compute_scrypt_nonce(
        "eyJ2ZXJzaW9uIjoxLCJ1YmlkIjoiZDEyMzNjM2EtZGIyNS00ZTJmLThmODQtOWVjZDFkMjk1NDlhIiwiYXR0ZW1wdF9pZCI6ImZhNTUyNTZhLTJmNTctNDM3MS1hYzdjLTY4ZTE0ZjU0ZTVjNCIsImNyZWF0ZV90aW1lIjoiMjAyNS0wNS0yNlQxMjo1ODozMS4yMTQ2ODM2MTJaIiwiZGlmZmljdWx0eSI6NCwiY2hhbGxlbmdlX3R5cGUiOiJIYXNoY2FzaFNjcnlwdCJ9",
        "084CF3DD",
        4))
