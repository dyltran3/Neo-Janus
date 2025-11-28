import random
import string

class PromptFuzzer:
    def __init__(self):
        # Mở rộng danh sách từ khóa nhạy cảm
        self.sensitive_words = {"hack", "kill", "bom", "password", "admin", "root", "xâm nhập", "đánh cắp", "trộm"}
        
        # [QUAN TRỌNG] Data Engineering: Teencode & Tiếng lóng (Neo-Janus.xlsx)
        self.teencode_map = {
            "a": ["4", "@", "a", "á"], "e": ["3", "e"], "i": ["1", "!", "i", "j"],
            "o": ["0", "()", "o", "oz"], "s": ["$", "5", "s", "x"], "h": ["k", "h"],
            "ck": ["kc", "k", "q"], "ph": ["f", "ph"], "ng": ["q", "g", "ng"],
            "qu": ["w", "qu"], "kh": ["x", "kh"]
        }
        # Kỹ thuật Injection: Chèn ký tự gây lú tokenizer (Attack Handbook Ch.5)
        self.injectors = ["\n", "\t", " --ignore ", " ### ", " // ", "/* */"]

    def _mutate_word(self, word: str) -> str:
        """Biến đổi từ ngữ dựa trên Teencode."""
        if word.lower() not in self.sensitive_words or random.random() > 0.7:
            return word
        
        # Builder pattern để tối ưu string concatenation
        chars = []
        i = 0
        w_len = len(word)
        while i < w_len:
            # Ưu tiên biến đổi cụm 2 ký tự (digraphs)
            if i < w_len - 1 and word[i:i+2].lower() in self.teencode_map:
                chars.append(random.choice(self.teencode_map[word[i:i+2].lower()]))
                i += 2
            elif word[i].lower() in self.teencode_map:
                chars.append(random.choice(self.teencode_map[word[i].lower()]))
                i += 1
            else:
                chars.append(word[i])
                i += 1
        return "".join(chars)

    def generate_payload(self, base_intent: str) -> str:
        """
        Sinh payload tấn công với cấu trúc ngẫu nhiên (Structured Fuzzing).
        """
        words = base_intent.split()
        fuzzed_words = [self._mutate_word(w) for w in words]
        payload = " ".join(fuzzed_words)

        # [QUAN TRỌNG] Injection Logic: Chèn nhiễu để phá vỡ ngữ cảnh của AI phòng thủ
        if random.random() > 0.5:
            injector = random.choice(self.injectors)
            # Chèn vào vị trí ngẫu nhiên
            pos = random.randint(0, len(payload))
            payload = payload[:pos] + injector + payload[pos:]

        return payload