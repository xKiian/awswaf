import hashlib

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
            digest = hashlib.sha256(data.encode('utf-8')).digest()
            hash_hex = ''.join(f"{int.from_bytes(digest[i:i + 4], 'big'):08x}"
                               for i in range(0, len(digest), 4))
        except Exception as e:
            print(e)
            hash_hex = "sha2H="
        if difficulty_met(difficulty, hash_hex):
            return str(nonce)
        nonce += 1


CHALLENGE_TYPES = {
    'h72f957df656e80ba55f5d8ce2e8c7ccb59687dba3bfb273d54b08a261b2f3002': "verify",
    'h7b0c470f0cfe3a80a9e26526ad185f484f6817d0832712a4a37a908786a6a67f': hashPow,
    'ha9faaffd31b4d5ede2a2e19d2d7fd525f66fee61911511960dcbb52d3c48ce25': "mp_verify"
}
