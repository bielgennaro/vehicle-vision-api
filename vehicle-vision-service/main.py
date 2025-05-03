import os
import shutil
from fastapi import FastAPI, File, UploadFile, Form, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from typing import Optional
import uvicorn
from datetime import datetime
import uuid
from pydantic import BaseModel
import cv2
import numpy as np
from PIL import Image
import io

class AnalysisResponse(BaseModel):
    image_id: int
    vehicle_type: str
    license_plate: Optional[str] = None
    confidence_score: float
    damage_detected: bool
    damage_details: Optional[str] = None
    processed_at: datetime

app = FastAPI(
    title="Vehicle Vision API",
    description="API for analyzing vehicle images",
    version="1.0.0"
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# uploads temporários
os.makedirs("temp", exist_ok=True)

@app.get("/")
async def root():
    return {"message": "Vehicle Vision Python Service"}

@app.post("/analyze/vehicle", response_model=AnalysisResponse)
async def analyze_vehicle(
    image_id: int = Form(...),
    file: UploadFile = File(...),
):
    """
    Analisa uma imagem de veículo e retorna informações detectadas.
    """
    if not file.content_type.startswith("image/"):
        raise HTTPException(
            status_code=400, 
            detail="Arquivo enviado não é uma imagem"
        )
    
    temp_file_path = f"temp/{uuid.uuid4()}{os.path.splitext(file.filename)[1]}"
    with open(temp_file_path, "wb") as buffer:
        shutil.copyfileobj(file.file, buffer)
    
    try:
        results = process_vehicle_image(temp_file_path)
        
        response = AnalysisResponse(
            image_id=image_id,
            vehicle_type=results["vehicle_type"],
            license_plate=results["license_plate"],
            confidence_score=results["confidence_score"],
            damage_detected=results["damage_detected"],
            damage_details=results["damage_details"],
            processed_at=datetime.now()
        )
        
        return response
    
    finally:
        if os.path.exists(temp_file_path):
            os.remove(temp_file_path)

def process_vehicle_image(image_path):
    """
    Como o intuito deste projeto é apenas demonstrar a estrutura,
    e não temos uma IA real, vamos simular o processamento de imagem.
    """
    img = cv2.imread(image_path)
    
    if img is None:
        raise HTTPException(status_code=400, detail="Error reading image")
        
    img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
    height, width, channels = img.shape
    average_color = np.mean(img, axis=(0, 1))
    
    # Simular categorização de veículo baseada em características da imagem
    if width > height:
        vehicle_type = "Sedan"
    else:
        vehicle_type = "SUV"
    
    # Simular detecção de placa
    license_plate = f"ABC{np.random.randint(1000, 9999)}"
    
    # Simular detecção de danos
    brightness = np.mean(average_color)
    damage_detected = brightness < 100  # Exemplo
    
    damage_details = None
    if damage_detected:
        damage_details = "Damage detected on the left side"
    
    return {
        "vehicle_type": vehicle_type,
        "license_plate": license_plate,
        "confidence_score": np.random.uniform(0.8, 0.98),  # Simulado
        "damage_detected": damage_detected,
        "damage_details": damage_details
    }

@app.get("/health")
async def health_check():
    return {
        "status": "Fine",
        "version": "1.0.0",
        "description": "Vehicle Vision Python Service",
        "timestamp": datetime.now().isoformat(),
        "service": "vehicle-vision-python"
    }

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
