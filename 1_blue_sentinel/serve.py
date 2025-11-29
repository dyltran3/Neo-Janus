# 1_blue_sentinel/serve.py
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import torch
import uvicorn
from src.model import BlueSentinelGPT, GPT_CONFIG_124M
from src.tokenizer import BlueTokenizer

# Khởi tạo FastAPI
app = FastAPI()

# Load Model & Tokenizer (Chỉ load 1 lần khi khởi động)
print("Loading Model...")
device = "cpu" # Hoặc "cuda" nếu có GPU
tokenizer = BlueTokenizer()
model = BlueSentinelGPT(GPT_CONFIG_124M)

# [TODO]: Load weight đã train (bỏ comment khi đã có file .pth)
# model.load_state_dict(torch.load("models/sentinel_v1.pth", map_location=device))
# model = BlueSentinelGPT.optimize_for_cpu(model) # Quantization

model.to(device)
model.eval()
print("Model Ready!")

class AnalyzeRequest(BaseModel):
    text: str

@app.post("/predict")
async def predict(req: AnalyzeRequest):
    try:
        # 1. Tokenize
        input_ids = torch.tensor(tokenizer.encode(req.text)).unsqueeze(0).to(device)
        
        # 2. Inference (Chỉ lấy 128 token đầu để nhanh)
        if input_ids.shape[1] > 128:
            input_ids = input_ids[:, :128]
            
        with torch.no_grad():
            logits = model(input_ids)
            # Giả sử output logits [batch, 2] -> [Safe, Unsafe]
            probs = torch.softmax(logits, dim=-1)
            risk_score = probs[0][1].item() # Lấy xác suất lớp 'Unsafe'
            
        return {
            "risk_score": risk_score,
            "is_blocked": risk_score > 0.85 # Ngưỡng chặn
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8001)