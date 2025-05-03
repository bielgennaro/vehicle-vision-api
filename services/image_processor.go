package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bielgennaro/vehicle-vision-api/configs"
	"github.com/bielgennaro/vehicle-vision-api/models"
)

// PythonServiceResponse representa a resposta do serviço Python
type PythonServiceResponse struct {
	ImageID         int       `json:"image_id"`
	VehicleType     string    `json:"vehicle_type"`
	LicensePlate    string    `json:"license_plate"`
	ConfidenceScore float64   `json:"confidence_score"`
	DamageDetected  bool      `json:"damage_detected"`
	DamageDetails   string    `json:"damage_details"`
	ProcessedAt     time.Time `json:"processed_at"`
}

// EnqueueImageProcessing adiciona uma imagem à fila de processamento
func EnqueueImageProcessing(imageID uint) {
	log.Printf("Imagem %d adicionada à fila de processamento", imageID)

	// Em uma implementação real, você enviaria para uma fila como RabbitMQ ou Kafka
	// Para este exemplo, vamos processar diretamente em uma goroutine

	go ProcessImage(imageID)
}

// ProcessImage processa uma imagem e salva os resultados
func ProcessImage(imageID uint) {
	log.Printf("Processando imagem %d", imageID)

	// Buscar a imagem
	var image models.Image
	result := configs.DB.First(&image, imageID)
	if result.Error != nil {
		log.Printf("Erro ao buscar imagem %d: %v", imageID, result.Error)
		return
	}

	// Verificar se o arquivo existe
	if _, err := os.Stat(image.Path); os.IsNotExist(err) {
		log.Printf("Arquivo de imagem não encontrado: %s", image.Path)
		return
	}

	// Chamar o serviço Python para análise
	response, err := callPythonService(imageID, image.Path)
	if err != nil {
		log.Printf("Erro ao chamar serviço Python: %v", err)

		// Fallback para processamento simulado
		createSimulatedAnalysis(imageID)
		return
	}

	// Criar análise com os resultados do serviço Python
	analysis := models.Analysis{
		ImageID:         imageID,
		VehicleType:     response.VehicleType,
		LicensePlate:    response.LicensePlate,
		ConfidenceScore: response.ConfidenceScore,
		DamageDetected:  response.DamageDetected,
		DamageDetails:   response.DamageDetails,
		ProcessedAt:     response.ProcessedAt,
	}

	// Salvar a análise
	if err := configs.DB.Create(&analysis).Error; err != nil {
		log.Printf("Erro ao salvar análise da imagem %d: %v", imageID, err)
		return
	}

	// Atualizar status da imagem
	image.Processed = true
	configs.DB.Save(&image)

	log.Printf("Imagem %d processada com sucesso", imageID)
}

// callPythonService chama o serviço Python para análise de imagem
func callPythonService(imageID uint, imagePath string) (*PythonServiceResponse, error) {
	// URL do serviço Python
	pythonServiceURL := "http://localhost:8000/analyze/vehicle"

	// Preparar o multipart form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Adicionar o ID da imagem como field
	fw, err := w.CreateFormField("image_id")
	if err != nil {
		return nil, err
	}
	_, err = fw.Write([]byte(strconv.FormatUint(uint64(imageID), 10)))
	if err != nil {
		return nil, err
	}

	// Abrir o arquivo
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Adicionar o arquivo
	fw, err = w.CreateFormFile("file", filepath.Base(imagePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return nil, err
	}

	// Fechar o writer
	w.Close()

	// Criar o request
	req, err := http.NewRequest("POST", pythonServiceURL, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Fazer o request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Verificar o status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}

	// Decodificar a resposta
	var response PythonServiceResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// createSimulatedAnalysis cria uma análise simulada quando o serviço Python falha
func createSimulatedAnalysis(imageID uint) {
	// Buscar o veículo associado à imagem
	var image models.Image
	configs.DB.First(&image, imageID)

	var vehicle models.Vehicle
	configs.DB.First(&vehicle, image.VehicleID)

	// Criar análise simulada
	analysis := models.Analysis{
		ImageID:         imageID,
		VehicleType:     vehicle.Model,
		LicensePlate:    vehicle.LicensePlate,
		ConfidenceScore: 0.7,
		DamageDetected:  false,
		DamageDetails:   "Simulated analysis - failed to call processing service",
		ProcessedAt:     time.Now(),
	}

	// Salva a análise
	if err := configs.DB.Create(&analysis).Error; err != nil {
		log.Printf("Erro ao salvar análise simulada da imagem %d: %v", imageID, err)
		return
	}

	// Atualizar status da imagem
	image.Processed = true
	configs.DB.Save(&image)

	log.Printf("Image %d processed in simulated analysis", imageID)
}
