package main

import (
	"log"

	"github.com/abdulazizax/udevslab-lesson3/cmd/api"
)

func main() {
	log.Fatal(api.Run())
}

// ListOrdersWithAggregates — Buyurtmalarni sanaga ko'ra agregatlash.
// ListOrdersByCustomer — Mijoz bo'yicha buyurtmalarni olish (join users kolleksiyasi bilan).
// TopSellingProducts — Eng ko'p sotilgan mahsulotlar ro'yxatini olish.
// ListOrdersByDateRange — Vaqt oralig'ida buyurtmalarni qidirish.
