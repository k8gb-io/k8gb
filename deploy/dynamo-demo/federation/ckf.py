# SPDX-FileCopyrightText: Copyright (c) 2026 The k8gb Contributors.
# SPDX-License-Identifier: Apache-2.0

"""Bounded Python consumer of Dynamo's CBI1 CKF publication protocol.

Wire contract and addressing are pinned to Dynamo 6b5b6459e21d. This is a
reference implementation for the demo, not the upstream transposed SIMD index.
"""

import struct
import time
from array import array
from dataclasses import dataclass

import xxhash

HEADER = struct.Struct("<4sHHBBHQQQQI")
MAX_BUCKETS = 1 << 20
MAX_PAYLOAD = 4_194_368
ALT_DOMAIN = 0x9E3779B97F4A7C15


class InvalidPublication(ValueError):
    pass


def sequence_hashes(tokens, block_size):
    """DYNAMO_STANDARD_V1 for plain text, without LoRA or a cache namespace."""
    if block_size <= 0 or any(type(t) is not int or not 0 <= t <= 0xFFFFFFFF for t in tokens):
        raise ValueError("invalid token IDs or block size")
    previous = None
    hashes = []
    for start in range(0, len(tokens) - block_size + 1, block_size):
        block = struct.pack(f"<{block_size}I", *tokens[start : start + block_size])
        local = xxhash.xxh3_64_intdigest(block, seed=1337)
        previous = (
            local if previous is None else xxhash.xxh3_64_intdigest(struct.pack("<QQ", previous, local), seed=1337)
        )
        hashes.append(previous)
    return hashes


@dataclass(frozen=True)
class Format:
    version: int
    seed: int
    buckets: int
    fingerprint_bits: int = 16
    slots: int = 4

    def __post_init__(self):
        if (self.version, self.fingerprint_bits, self.slots) != (1, 16, 4):
            raise InvalidPublication("unsupported CKF format")
        if not 2 <= self.buckets <= MAX_BUCKETS or self.buckets & (self.buckets - 1):
            raise InvalidPublication("unsupported CKF bucket count")


class Replica:
    def __init__(self, layout, dc_id):
        self.layout = layout
        self.dc_id = dc_id
        self.clear()

    def clear(self):
        self.words = None
        self.sequence = None
        self.updated = 0.0
        self._assembly = None

    def apply(self, kind, payload, base_sequence, sequence):
        """Invalidate on every malformed frame, including partially applied snapshots."""
        try:
            self._apply(kind, payload, base_sequence, sequence)
            self.updated = time.monotonic()
        except (struct.error, OverflowError, ValueError) as error:
            self.clear()
            raise InvalidPublication(str(error)) from error

    def _apply(self, kind, payload, base_sequence, sequence):
        if kind == 3:
            if payload or self.words is None or sequence != self.sequence:
                raise InvalidPublication("heartbeat without a current complete replica")
            return
        if not HEADER.size <= len(payload) <= MAX_PAYLOAD:
            raise InvalidPublication("invalid CBI1 frame length")
        magic, wire, flags, bits, slots, version, seed, count, dc_id, epoch, checksum = HEADER.unpack_from(payload)
        if magic != b"CBI1" or wire != 1 or flags != kind or kind not in (1, 2):
            raise InvalidPublication("invalid CBI1 header")
        if Format(version, seed, count, bits, slots) != self.layout or dc_id != self.dc_id or epoch != sequence:
            raise InvalidPublication("CBI1 format, pool or sequence drift")
        body = payload[HEADER.size :]
        if xxhash.xxh3_64_intdigest(body) & 0xFFFFFFFF != checksum:
            raise InvalidPublication("CBI1 checksum mismatch")
        if kind == 1:
            if len(body) < 16 or (len(body) - 16) % 8:
                raise InvalidPublication("invalid snapshot body")
            index, chunks, offset = struct.unpack_from("<IIQ", body)
            if not 1 <= chunks <= count or index >= chunks:
                raise InvalidPublication("invalid snapshot chunk index")
            if index == 0:
                self.words = None
                self.sequence = None
                self._assembly = (sequence, chunks, 0, array("Q"))
            if self._assembly is None:
                raise InvalidPublication("snapshot did not start with chunk zero")
            current, total, next_index, words = self._assembly
            if (current, total, next_index, len(words)) != (sequence, chunks, index, offset):
                raise InvalidPublication("non-contiguous snapshot")
            length = (len(body) - 16) // 8
            if not length or offset + length > count:
                raise InvalidPublication("snapshot exceeds bucket bounds")
            words.extend(value[0] for value in struct.iter_unpack("<Q", body[16:]))
            self._assembly = (sequence, chunks, index + 1, words)
            if index + 1 == chunks:
                if len(words) != count:
                    raise InvalidPublication("incomplete snapshot")
                self.words, self.sequence, self._assembly = words, sequence, None
        else:
            if self.words is None or base_sequence != self.sequence or sequence != base_sequence + 1:
                raise InvalidPublication("delta gap or missing snapshot")
            if len(body) < 12:
                raise InvalidPublication("truncated delta")
            base, length = struct.unpack_from("<QI", body)
            if base != base_sequence or len(body) != 12 + 12 * length:
                raise InvalidPublication("invalid delta envelope")
            images = list(struct.iter_unpack("<IQ", body[12:]))
            if any(bucket >= count for bucket, _ in images) or len({b for b, _ in images}) != length:
                raise InvalidPublication("invalid or repeated delta bucket")
            for bucket, word in images:
                self.words[bucket] = word
            self.sequence = sequence

    def prefix_blocks(self, hashes):
        if self.words is None:
            return 0
        mask = self.layout.buckets - 1
        matched = 0
        for sequence_hash in hashes:
            mixed = xxhash.xxh3_64_intdigest(struct.pack("<Q", sequence_hash), seed=self.layout.seed)
            fingerprint = (mixed & 0xFFFF) or 1
            first = (mixed >> 16) & mask
            delta = xxhash.xxh3_64_intdigest(struct.pack("<H", fingerprint), seed=self.layout.seed ^ ALT_DOMAIN) & mask
            second = first ^ (delta or 1)
            if not any(
                (self.words[bucket] >> (16 * slot)) & 0xFFFF == fingerprint
                for bucket in (first, second)
                for slot in range(4)
            ):
                break
            matched += 1
        return matched
