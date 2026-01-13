# 🌍 Travel Planner – AI Powered Trip Planning Platform

An intelligent, full-stack **AI Travel Planner** platform built to help users plan trips in Sri Lanka (and beyond) with smart itineraries, place recommendations, and optimized routes.  
Built using **Go (Gin)** backend, **Nuxt 3 / Vue** frontend, and fully **Dockerized** for easy deployment.

## 🚀 Features

- 🧠 AI-powered itinerary generation  
- 📍 Place discovery & recommendations  
- 🗺️ Google Maps integration  
- 🏨 Hotel & attraction suggestions  
- 🚕 Transport planning (LankaCab ready)  
- 👤 User authentication & profiles  
- 📊 Optimized day-by-day travel plans  
- 🌐 REST API backend (Go + Gin)  
- 🖥️ Modern frontend (Nuxt 3 + Vue)  
- 🐳 Docker + Nginx ready  

## 🛠 Tech Stack

### Backend
- Go (Golang)
- Gin Framework
- SQLite / MySQL
- JWT Authentication
- REST API

### Frontend
- Nuxt 3
- Vue 3
- Tailwind CSS
- Axios

### Infrastructure
- Docker & Docker Compose
- Nginx
- SSL (Certbot)
- Linux (Ubuntu)
- AWS Lightsail compatible

## 📁 Project Structure

travel-planner/
├── trip-planner/          # Go backend (API)
├── trip-planner-web/      # Nuxt 3 frontend
├── nginx/                 # Nginx config
├── certbot/               # SSL configs
├── docker-compose.yml
└── README.md

## ⚙️ Installation & Setup

### Clone the Repository
git clone https://github.com/ashanda/travel-planner.git
cd travel-planner

### Run with Docker
docker-compose up -d --build

Frontend → http://localhost  
API → http://localhost/api

## 🤖 AI Integration

Designed to integrate with OpenAI, Google Gemini, or custom ML models for itinerary generation and smart recommendations.

## 🧾 License

Proprietary / Internal use (add MIT or Apache-2.0 if open sourcing)

## 👨‍💻 Author

Ashanda – Sri Lanka  
Full-Stack Developer | AI & Travel Tech Builder
