import binascii
import hashlib
import itertools
from typing import (
    Any,
    Callable,
    Optional
)


def _check(digest: bytes, difficulty: int) -> bool:
    """
    Return True if digest has at least difficulty leading zero bits.
    """
    full, rem = divmod(difficulty, 8)
    if digest[:full] != b"\x00" * full:
        return False
    if rem and (digest[full] >> (8 - rem)):
        return False
    return True


def hash_pow(challenge: str, salt: str, difficulty: int) -> Optional[str]:
    """
    Find a nonce so that sha256(challenge + salt + nonce) meets difficulty.
    """
    prefix = (challenge + salt).encode()
    for nonce in itertools.count():
        digest = hashlib.sha256(prefix + str(nonce).encode()).digest()
        if _check(digest, difficulty):
            return str(nonce)
    return None


def scrypt_func(input_str: str, salt: str, n: int = 128, r: int = 8, p: int = 1, dklen: int = 16) -> str:
    """
    Compute scrypt hash and return hex-encoded string.
    """
    raw = hashlib.scrypt(password=input_str.encode(), salt=salt.encode(), n=n, r=r, p=p, dklen=dklen)
    return binascii.hexlify(raw).decode()


def compute_scrypt_nonce(
    challenge: str,
    salt: str,
    difficulty: int,
    n: int = 128,
    r: int = 8,
    p: int = 1,
    dklen: int = 16
) -> Optional[str]:
    """
    Find a nonce so that scrypt(challenge + salt + nonce) meets difficulty.
    """
    prefix = challenge + salt
    for nonce in itertools.count():
        digest = hashlib.scrypt(password=f"{prefix}{nonce}".encode(), salt=salt.encode(), n=n, r=r, p=p, dklen=dklen)
        if _check(digest, difficulty):
            return str(nonce)
    return None


CHALLENGE_TYPES: dict[str, Callable[..., Any]] = {
    'h72f957df656e80ba55f5d8ce2e8c7ccb59687dba3bfb273d54b08a261b2f3002': compute_scrypt_nonce,
    'h7b0c470f0cfe3a80a9e26526ad185f484f6817d0832712a4a37a908786a6a67f': hash_pow,
    'ha9faaffd31b4d5ede2a2e19d2d7fd525f66fee61911511960dcbb52d3c48ce25': "mp_verify"
}
