import tiktoken
import torch
from typing import List, Union, Dict, Set

# Định nghĩa các token bảo mật đặc biệt
SPECIAL_TOKENS_MAP: Dict[str, int] = {
    "<|pad|>": 50257,       
    "<|malicious|>": 50258, # Nhãn độc hại
    "<|safe|>": 50259       # Nhãn an toàn
}
# GPT-2 gốc có 50257 token. Tổng cộng: 50260
VOCAB_SIZE = 50257 + len(SPECIAL_TOKENS_MAP) 

class BlueTokenizer:
    def __init__(self, base_encoding: str = "gpt2"):
        try:
            self._base_tokenizer = tiktoken.get_encoding(base_encoding)
        except Exception as e:
            raise RuntimeError(f"Failed to load encoding '{base_encoding}': {e}")

        # Đảm bảo tiktoken không cắt rời các token đặc biệt này
        self._special_token_names: Set[str] = set(SPECIAL_TOKENS_MAP.keys())
        self._special_token_names.add("<|endoftext|>") # Token gốc của GPT-2

    @property
    def vocab_size(self) -> int:
        return VOCAB_SIZE

    def encode(self, text: str, return_tensors: str = "pt") -> Union[List[int], torch.Tensor]:
        # Quan trọng: allowed_special để giữ nguyên token bảo mật
        ids = self._base_tokenizer.encode(text, allowed_special=self._special_token_names)
        if return_tensors == "pt":
            return torch.tensor(ids, dtype=torch.long).clone().detach()
        return ids

    def decode(self, token_ids: Union[List[int], torch.Tensor]) -> str:
        # Hàm decode đơn giản cho mục đích debug
        if isinstance(token_ids, torch.Tensor):
            ids_list = token_ids.tolist()
        else:
            ids_list = token_ids
        # Lưu ý: Decode này chưa xử lý hoàn hảo việc ghép lại special tokens tùy biến,
        # nhưng đủ dùng để kiểm tra văn bản gốc.
        return self._base_tokenizer.decode(ids_list)

if __name__ == "__main__":
    # Sanity Check
    tokenizer = BlueTokenizer()
    text = "Test tiếng việt và teencode h4ck <|malicious|>"
    encoded = tokenizer.encode(text)
    print(f"Input: {text}\nEncoded IDs: {encoded}\nDecoded: {tokenizer.decode(encoded)}")
    assert encoded[-1].item() == SPECIAL_TOKENS_MAP["<|malicious|>"]
    print("Tokenizer sanity check passed!")