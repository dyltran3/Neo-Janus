"""
PromptFuzzer: Generates adversarial prompts using teencode obfuscation
and injection techniques to bypass content filters
"""

import random
from typing import Dict, List


class PromptFuzzer:
    """Generate fuzzy/obfuscated prompts for adversarial testing"""

    def __init__(self):
        """Initialize fuzzer with mutation rules and injection patterns"""
        # Sensitive keywords that trigger defense mechanisms
        self.sensitive_words = {
            "hack", "kill", "bomb", "password", "admin", "root",
            "xâm nhập", "đánh cắp", "trộm", "tấn công"
        }

        # Teencode mapping for Vietnamese character substitution
        # Based on Vietnamese internet slang for obfuscation
        self.teencode_map: Dict[str, List[str]] = {
            "a": ["4", "@", "a", "á"],
            "e": ["3", "e"],
            "i": ["1", "!", "i", "j"],
            "o": ["0", "()", "o", "oz"],
            "s": ["$", "5", "s", "x"],
            "h": ["k", "h"],
            "ck": ["kc", "k", "q"],
            "ph": ["f", "ph"],
            "ng": ["q", "g", "ng"],
            "qu": ["w", "qu"],
            "kh": ["x", "kh"],
        }

        # Injection payloads to confuse tokenizer
        self.injectors = [
            "\n",
            "\t",
            " --ignore ",
            " ### ",
            " // ",
            "/* */",
            "\u200b",  # Zero-width space
        ]

    def _mutate_word(self, word: str) -> str:
        """
        Mutate word using teencode substitutions
        Prioritizes digraph (2-char) replacements for better obfuscation
        """
        if word.lower() not in self.sensitive_words or random.random() > 0.7:
            return word

        chars: List[str] = []
        i = 0
        word_len = len(word)

        while i < word_len:
            # Try digraph replacement first (2-char sequences)
            if i < word_len - 1:
                digraph = word[i : i + 2].lower()
                if digraph in self.teencode_map:
                    chars.append(random.choice(self.teencode_map[digraph]))
                    i += 2
                    continue

            # Fall back to single character replacement
            char = word[i].lower()
            if char in self.teencode_map:
                chars.append(random.choice(self.teencode_map[char]))
            else:
                chars.append(word[i])

            i += 1

        return "".join(chars)

    def generate_payload(self, base_intent: str) -> str:
        """
        Generate obfuscated payload from base intent
        Uses teencode mutations and injection techniques
        """
        if not base_intent or not isinstance(base_intent, str):
            return base_intent

        # Apply teencode obfuscation to each word
        words = base_intent.split()
        fuzzed_words = [self._mutate_word(w) for w in words]
        payload = " ".join(fuzzed_words)

        # Randomly inject noise to break tokenizer context
        if random.random() > 0.5:
            injector = random.choice(self.injectors)
            injection_pos = random.randint(0, len(payload))
            payload = payload[:injection_pos] + injector + payload[injection_pos:]

        return payload

    def batch_generate(self, base_intents: List[str], count: int = 10) -> List[str]:
        """Generate multiple payloads from base intents"""
        payloads = []
        for _ in range(count):
            base = random.choice(base_intents)
            payloads.append(self.generate_payload(base))
        return payloads