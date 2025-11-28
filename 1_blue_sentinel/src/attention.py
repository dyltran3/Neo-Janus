import torch
import torch.nn as nn
import torch.nn.functional as F

class MultiHeadAttention(nn.Module):
    def __init__(self, d_in, d_out, context_length, dropout, num_heads, qkv_bias=False):
        super().__init__()
        assert d_out % num_heads == 0, "d_out must be divisible by num_heads"
        self.d_out = d_out
        self.num_heads = num_heads
        self.head_dim = d_out // num_heads
        
        self.W_query = nn.Linear(d_in, d_out, bias=qkv_bias)
        self.W_key = nn.Linear(d_in, d_out, bias=qkv_bias)
        self.W_value = nn.Linear(d_in, d_out, bias=qkv_bias)
        self.out_proj = nn.Linear(d_out, d_out)
        self.dropout_p = dropout

    def forward(self, x):
        b, num_tokens, _ = x.shape
        # Tách heads: (b, num_tokens, num_heads, head_dim) -> (b, num_heads, num_tokens, head_dim)
        keys = self.W_key(x).view(b, num_tokens, self.num_heads, self.head_dim).transpose(1, 2)
        queries = self.W_query(x).view(b, num_tokens, self.num_heads, self.head_dim).transpose(1, 2)
        values = self.W_value(x).view(b, num_tokens, self.num_heads, self.head_dim).transpose(1, 2)
        
        # [QUAN TRỌNG] Tối ưu hóa theo Build LLM From Scratch (Chương 3)
        # Sử dụng Flash Attention (nếu GPU hỗ trợ) hoặc cài đặt tối ưu C++ của PyTorch
        # is_causal=True tự động tạo mask tam giác (Masked Multi-head Attention)
        context_vec = F.scaled_dot_product_attention(
            queries, keys, values, 
            dropout_p=self.dropout_p if self.training else 0.0,
            is_causal=True 
        )
        
        context_vec = context_vec.transpose(1, 2).contiguous().view(b, num_tokens, self.d_out)
        return self.out_proj(context_vec)