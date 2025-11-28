import torch
import torch.nn as nn
from .attention import MultiHeadAttention
# Import tokenizer chỉ để lấy VOCAB_SIZE chuẩn
from .tokenizer import VOCAB_SIZE 

# Cấu hình tối ưu cho RAM 8GB (SLM)
GPT_CONFIG_124M = {
    "vocab_size": VOCAB_SIZE, 
    "context_length": 512, # Giảm context để nhẹ RAM
    "emb_dim": 768,
    "n_heads": 12,
    "n_layers": 6,        # Chỉ dùng 6 lớp thay vì 12
    "drop_rate": 0.1,
    "qkv_bias": False
}

class LayerNorm(nn.Module):
    def __init__(self, emb_dim):
        super().__init__()
        self.eps = 1e-5
        self.scale = nn.Parameter(torch.ones(emb_dim))
        self.shift = nn.Parameter(torch.zeros(emb_dim))
    def forward(self, x):
        mean = x.mean(dim=-1, keepdim=True)
        var = x.var(dim=-1, keepdim=True, unbiased=False)
        norm_x = (x - mean) / torch.sqrt(var + self.eps)
        return self.scale * norm_x + self.shift

class FeedForward(nn.Module):
    def __init__(self, cfg):
        super().__init__()
        self.layers = nn.Sequential(
            nn.Linear(cfg["emb_dim"], 4 * cfg["emb_dim"]),
            nn.GELU(),
            nn.Linear(4 * cfg["emb_dim"], cfg["emb_dim"]),
        )
    def forward(self, x): return self.layers(x)

class TransformerBlock(nn.Module):
    def __init__(self, cfg):
        super().__init__()
        self.att = MultiHeadAttention(
            d_in=cfg["emb_dim"], d_out=cfg["emb_dim"],
            context_length=cfg["context_length"], dropout=cfg["drop_rate"],
            num_heads=cfg["n_heads"], qkv_bias=cfg["qkv_bias"])
        self.ff = FeedForward(cfg)
        self.norm1 = LayerNorm(cfg["emb_dim"])
        self.norm2 = LayerNorm(cfg["emb_dim"])
        self.drop_shortcut = nn.Dropout(cfg["drop_rate"])
    def forward(self, x):
        x = x + self.drop_shortcut(self.att(self.norm1(x)))
        x = x + self.drop_shortcut(self.ff(self.norm2(x)))
        return x

class BlueSentinelGPT(nn.Module):
    def __init__(self, cfg):
        super().__init__()
        self.tok_emb = nn.Embedding(cfg["vocab_size"], cfg["emb_dim"])
        self.pos_emb = nn.Embedding(cfg["context_length"], cfg["emb_dim"])
        self.drop_emb = nn.Dropout(cfg["drop_rate"])
        self.trf_blocks = nn.Sequential(*[TransformerBlock(cfg) for _ in range(cfg["n_layers"])])
        self.final_norm = LayerNorm(cfg["emb_dim"])
        # Head phân loại 2 lớp: [An toàn, Độc hại]
        self.classifier_head = nn.Linear(cfg["emb_dim"], 2, bias=False)

    def forward(self, in_idx):
        batch_size, seq_len = in_idx.shape
        tok_embeds = self.tok_emb(in_idx)
        pos_embeds = self.pos_emb(torch.arange(seq_len, device=in_idx.device))
        x = self.drop_emb(tok_embeds + pos_embeds)
        x = self.trf_blocks(x)
        x = self.final_norm(x)
        # Chỉ lấy token cuối cùng để phân loại câu
        logits = self.classifier_head(x[:, -1, :])
        return logits

if __name__ == "__main__":
    # Sanity check kiến trúc
    model = BlueSentinelGPT(GPT_CONFIG_124M)
    print("Kiến trúc BlueSentinelGPT đã khởi tạo thành công.")
    print(f"Tổng tham số ước tính: {sum(p.numel() for p in model.parameters()) / 1e6:.1f}M")